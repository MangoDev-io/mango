package routes

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"

	"github.com/algorand/go-algorand-sdk/client/algod"
	algodModels "github.com/algorand/go-algorand-sdk/client/algod/models"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/mnemonic"
	"github.com/algorand/go-algorand-sdk/transaction"
	"github.com/algorand/go-algorand-sdk/types"
	"github.com/mangodev-io/mango/api/cmd/api/models"
	"github.com/sirupsen/logrus"
)

// ManagerHandler is the router handler for Mango
type ManagerHandler struct {
	log          *logrus.Logger
	testnetAlgod *algod.Client
	mainnetAlgod *algod.Client
	jwt          *jwtauth.JWTAuth
}

type response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	AssetID uint64 `json:"assetId"`
	TXHash  string `json:"txHash"`
}

// NewManagerHandler creates a new instance of ManagerHandler
func NewManagerHandler(log *logrus.Logger, testnetAlgod, mainnetAlgod *algod.Client, jwt *jwtauth.JWTAuth) *ManagerHandler {
	return &ManagerHandler{
		log:          log,
		testnetAlgod: testnetAlgod,
		mainnetAlgod: mainnetAlgod,
		jwt:          jwt,
	}
}

func makeResponseJSON(status string, message string, assetID uint64, txHash string) []byte {
	res, _ := json.Marshal(response{
		Status:  status,
		Message: message,
		AssetID: assetID,
		TXHash:  txHash,
	})

	return res
}

// EncodeMnemonic accepts JSON body of the form { "mnemonic" : "abcd" }
// and converts it into a JWT and send it to web
func (h *ManagerHandler) EncodeMnemonic(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)

	type tokenResponse struct {
		Token   string `json:"token"`
		Address string `json:"address"`
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	if err != nil {
		h.log.WithError(err).Error("failed to read body")
		rw.WriteHeader(http.StatusBadRequest)

		rw.Header().Set("Content-Type", "application/json")

		respJSON, _ := json.Marshal(tokenResponse{
			Token:   "",
			Status:  "error",
			Message: "failed to read request body",
		})
		rw.Write(respJSON)
		return
	}

	type mnemonic struct {
		Mnemonic string `json:"mnemonic"`
	}

	var request mnemonic
	err = json.Unmarshal(body, &request)
	if err != nil {
		h.log.WithError(err).Error("failed to unmarshal body")
		rw.WriteHeader(http.StatusBadRequest)

		respJSON, _ := json.Marshal(tokenResponse{
			Token:   "",
			Status:  "error",
			Message: "failed to unmarshal request body",
		})
		rw.Write(respJSON)
		return
	}

	_, address, err := h.recoverAccount(request.Mnemonic)
	if err != nil {
		h.log.WithError(err).Error("failed to recover account from mnemonic")
		rw.WriteHeader(http.StatusBadRequest)
		respJSON, _ := json.Marshal(tokenResponse{
			Token:   "",
			Status:  "error",
			Message: "failed to recover account from mnemonic",
		})
		rw.Write(respJSON)
		return
	}

	ttl := 15 * 60 * time.Second

	_, tokenString, err := h.jwt.Encode(jwt.MapClaims{"mnemonic": request.Mnemonic, "exp": time.Now().UTC().Add(ttl).Unix()})
	if err != nil {
		h.log.WithError(err).Error("failed to encode jwt claims")
		rw.WriteHeader(http.StatusInternalServerError)
		respJSON, _ := json.Marshal(tokenResponse{
			Token:   "",
			Status:  "error",
			Message: "failed to encode jwt claims",
		})
		rw.Write(respJSON)
		return
	}

	var response tokenResponse
	response.Address = address
	response.Token = tokenString
	response.Status = "success"
	response.Message = ""

	rw.Header().Set("Content-Type", "application/json")
	respJSON, err := json.Marshal(response)

	if err != nil {
		h.log.WithError(err).Error("failed to marshal response body")
		rw.WriteHeader(http.StatusInternalServerError)

		respJSON, _ := json.Marshal(tokenResponse{
			Token:   "",
			Status:  "error",
			Message: "failed to marshal token response body",
		})
		rw.Write(respJSON)
		return
	}

	rw.Write(respJSON)
}

// GetAssets accepts a JSON body of the form { "address" : "example"  }
// and returns a list of all assets owned by the address (created on Mango)
func (h *ManagerHandler) GetAssets(rw http.ResponseWriter, req *http.Request) {
	_, claims, err := jwtauth.FromContext(req.Context())
	if err != nil {
		h.log.WithError(err).Error("failed to get jwt claims from context")
		rw.WriteHeader(http.StatusBadRequest)

		respJSON := makeResponseJSON("error", "failed to read request body", 0, "")
		rw.Write(respJSON)

		return
	}

	useTestnetNetwork := true
	if network := chi.URLParam(req, "network"); network == "mainnet" {
		useTestnetNetwork = false
	}

	// privateKey, address, err := h.getPrivKeyAndAddressFromMnemonic(claims["mnemonic"].(string))
	_, address, err := h.recoverAccount(claims["mnemonic"].(string))
	if err != nil {
		h.log.WithError(err).Error("failed to get private key from mnemonic")
		rw.WriteHeader(http.StatusBadRequest)

		respJSON := makeResponseJSON("error", "failed to recover account from mnemonic", 0, "")
		rw.Write(respJSON)
		return
	}

	var account algodModels.Account
	if useTestnetNetwork {
		account, err = h.testnetAlgod.AccountInformation(address)
	} else {
		account, err = h.mainnetAlgod.AccountInformation(address)
	}

	var listing []models.AssetListing
	for assetID := range account.Assets {
		listing = append(listing, models.AssetListing{
			AssetID: strconv.FormatUint(assetID, 10),
		})
	}

	jsonResp, err := json.Marshal(listing)
	if err != nil {
		h.log.WithError(err).Error("failed to marshal owned assets")
		rw.WriteHeader(http.StatusInternalServerError)

		respJSON := makeResponseJSON("error", "failed to marshal owned assets to JSON", 0, "")
		rw.Write(respJSON)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(jsonResp)
}

// CreateAsset takes in a JSON body of type AssetCreate and
// creates a new ASA on Algorand
func (h *ManagerHandler) CreateAsset(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		respJSON := makeResponseJSON("error", "failed to read request body", 0, "")
		rw.Write(respJSON)
		return
	}

	useTestnetNetwork := true
	if network := chi.URLParam(req, "network"); network == "mainnet" {
		useTestnetNetwork = false
	}

	var assetDetails models.AssetCreate

	err = json.Unmarshal(body, &assetDetails)

	if err != nil {
		h.log.WithError(err).Error("failed to unmarshal request body")
		rw.WriteHeader(http.StatusBadRequest)

		respJSON := makeResponseJSON("error", "failed to unmarshal request body", 0, "")
		rw.Write(respJSON)
		return
	}

	_, claims, err := jwtauth.FromContext(req.Context())
	if err != nil {
		h.log.WithError(err).Error("failed to get jwt claims from context")
		rw.WriteHeader(http.StatusBadRequest)

		respJSON := makeResponseJSON("error", "failed to get jwt claims from context", 0, "")
		rw.Write(respJSON)
		return
	}

	// privateKey, address, err := h.getPrivKeyAndAddressFromMnemonic(claims["mnemonic"].(string))
	privateKey, address, err := h.recoverAccount(claims["mnemonic"].(string))

	if err != nil {
		h.log.WithError(err).Error("failed to get private key from mnemonic")
		rw.WriteHeader(http.StatusBadRequest)

		respJSON := makeResponseJSON("error", "failed to recover account from mnemonic", 0, "")
		rw.Write(respJSON)
		return
	}

	assetDetails.CreatorAddr = address

	txID, err := h.makeAndSendTxn(assetDetails, "create", privateKey, useTestnetNetwork)
	if err != nil {
		h.log.WithError(err).Error("failed to make and send asset create txn")
		rw.WriteHeader(http.StatusBadRequest)

		respJSON := makeResponseJSON("error", "failed to make and send asset create txn", 0, "")
		rw.Write(respJSON)
		return
	}

	// Retrieve asset ID by grabbing the max asset ID
	// from the creator account's holdings.
	var account algodModels.Account
	if useTestnetNetwork {
		account, err = h.testnetAlgod.AccountInformation(address)
	} else {
		account, err = h.mainnetAlgod.AccountInformation(address)
	}
	if err != nil {
		h.log.WithError(err).Error("failed to get account information")
		rw.WriteHeader(http.StatusInternalServerError)

		respJSON := makeResponseJSON("error", "failed to get account information", 0, "")
		rw.Write(respJSON)
		return
	}
	assetID := uint64(0)
	for i := range account.AssetParams {
		if i > assetID {
			assetID = i
		}
	}
	h.log.Debugf("Asset ID from AssetParams: %d", assetID)
	// Retrieve asset info.
	assetInfo, err := h.testnetAlgod.AssetInformation(assetID)
	h.log.Debugf("Asset info: %#v", assetInfo)

	respJSON := makeResponseJSON("success", "", assetID, txID)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(respJSON)
}

// ModifyAsset is used to modify the mutable addresses linked to an asset
func (h *ManagerHandler) ModifyAsset(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		h.log.WithError(err).Error("unable to read request body")
		rw.WriteHeader(http.StatusBadRequest)

		respJSON := makeResponseJSON("error", "failed to read request body", 0, "")
		rw.Write(respJSON)
		return
	}

	useTestnetNetwork := true
	if network := chi.URLParam(req, "network"); network == "mainnet" {
		useTestnetNetwork = false
	}

	var assetDetails models.AssetModify

	err = json.Unmarshal(body, &assetDetails)

	if err != nil {
		h.log.WithError(err).Error("unable to unmarshal request into JSON")
		rw.WriteHeader(http.StatusBadRequest)

		respJSON := makeResponseJSON("error", "failed to unmarshal request body", 0, "")
		rw.Write(respJSON)
		return
	}

	_, claims, err := jwtauth.FromContext(req.Context())
	if err != nil {
		h.log.WithError(err).Error("failed to get jwt claims from context")
		rw.WriteHeader(http.StatusBadRequest)

		respJSON := makeResponseJSON("error", "failed to get jwt claims from context", 0, "")
		rw.Write(respJSON)
		return
	}

	privateKey, managerAddr, err := h.recoverAccount(claims["mnemonic"].(string))
	if err != nil {
		h.log.WithError(err).Error("failed to get private key from mnemonic")
		rw.WriteHeader(http.StatusBadRequest)

		respJSON := makeResponseJSON("error", "failed to recover account from mnemonic", 0, "")
		rw.Write(respJSON)
		return
	}

	assetDetails.CurrManagerAddr = managerAddr

	txID, err := h.makeAndSendTxn(assetDetails, "modify", privateKey, useTestnetNetwork)
	if err != nil {
		h.log.WithError(err).Error("failed to make and send asset modify txn")
		rw.WriteHeader(http.StatusBadRequest)

		respJSON := makeResponseJSON("error", "failed to make and send asset modify txn", 0, "")
		rw.Write(respJSON)
		return
	}
	h.log.Debug("Transaction ID: ", txID)

	respJSON := makeResponseJSON("success", "", assetDetails.AssetID, txID)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(respJSON)
}

// DestroyAsset broadcasts a destroy asset transaction
func (h *ManagerHandler) DestroyAsset(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		h.log.WithError(err).Error("unable to read request body")
		rw.WriteHeader(http.StatusBadRequest)

		respJSON := makeResponseJSON("error", "failed to read request body", 0, "")
		rw.Write(respJSON)
		return
	}

	useTestnetNetwork := true
	if network := chi.URLParam(req, "network"); network == "mainnet" {
		useTestnetNetwork = false
	}

	var assetDetails models.AssetDestroy

	err = json.Unmarshal(body, &assetDetails)

	if err != nil {
		h.log.WithError(err).Error("unable to unmarshal request into JSON")
		rw.WriteHeader(http.StatusBadRequest)

		respJSON := makeResponseJSON("error", "failed to unmarshal request body", 0, "")
		rw.Write(respJSON)
		return
	}

	_, claims, err := jwtauth.FromContext(req.Context())
	if err != nil {
		h.log.WithError(err).Error("failed to get jwt claims from context")
		rw.WriteHeader(http.StatusBadRequest)

		respJSON := makeResponseJSON("error", "failed to get jwt claims from context", 0, "")
		rw.Write(respJSON)
		return
	}

	privateKey, managerAddr, err := h.recoverAccount(claims["mnemonic"].(string))
	if err != nil {
		h.log.WithError(err).Error("failed to get private key from mnemonic")
		rw.WriteHeader(http.StatusBadRequest)

		respJSON := makeResponseJSON("error", "failed to recover account from mnemonic", 0, "")
		rw.Write(respJSON)
		return
	}

	assetDetails.ManagerAddr = managerAddr

	txID, err := h.makeAndSendTxn(assetDetails, "destory", privateKey, useTestnetNetwork)
	if err != nil {
		h.log.WithError(err).Error("failed to make and send asset destroy txn")
		rw.WriteHeader(http.StatusInternalServerError)

		respJSON := makeResponseJSON("error", "failed to make and send asset destroy txn", 0, "")
		rw.Write(respJSON)
		return
	}
	h.log.Debug("Transaction ID: ", txID)

	respJSON := makeResponseJSON("success", "", assetDetails.AssetID, txID)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(respJSON)
}

// FreezeAsset is used to broadcast a freeze asset txn for some address
func (h *ManagerHandler) FreezeAsset(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		h.log.WithError(err).Error("unable to read request body")
		rw.WriteHeader(http.StatusBadRequest)

		respJSON := makeResponseJSON("error", "failed to read request body", 0, "")
		rw.Write(respJSON)
		return
	}

	useTestnetNetwork := true
	if network := chi.URLParam(req, "network"); network == "mainnet" {
		useTestnetNetwork = false
	}

	var assetDetails models.AssetFreeze

	err = json.Unmarshal(body, &assetDetails)

	if err != nil {
		h.log.WithError(err).Error("unable to unmarshal request into JSON")
		rw.WriteHeader(http.StatusBadRequest)

		respJSON := makeResponseJSON("error", "failed to unmarshal request body", 0, "")
		rw.Write(respJSON)
		return
	}

	_, claims, err := jwtauth.FromContext(req.Context())
	if err != nil {
		h.log.WithError(err).Error("failed to get jwt claims from context")
		rw.WriteHeader(http.StatusBadRequest)

		respJSON := makeResponseJSON("error", "failed to get jwt claims from context", 0, "")
		rw.Write(respJSON)
		return
	}

	privateKey, freezeAddr, err := h.recoverAccount(claims["mnemonic"].(string))
	if err != nil {
		h.log.WithError(err).Error("failed to get private key from mnemonic")
		rw.WriteHeader(http.StatusBadRequest)

		respJSON := makeResponseJSON("error", "failed to recover account from mnemonic", 0, "")
		rw.Write(respJSON)
		return
	}

	assetDetails.FreezeAddr = freezeAddr

	txID, err := h.makeAndSendTxn(assetDetails, "freeze", privateKey, useTestnetNetwork)
	if err != nil {
		h.log.WithError(err).Error("failed to make and send asset freeze txn")
		rw.WriteHeader(http.StatusBadRequest)

		respJSON := makeResponseJSON("error", "failed to make and send asset freeze txn", 0, "")
		rw.Write(respJSON)
		return
	}
	h.log.Debug("Transaction ID: ", txID)

	respJSON := makeResponseJSON("success", "", assetDetails.AssetID, txID)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(respJSON)
}

// RevokeAsset broadcasts a revoke (clawback) txn for some address
func (h *ManagerHandler) RevokeAsset(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		h.log.WithError(err).Error("unable to read request body")
		rw.WriteHeader(http.StatusBadRequest)

		respJSON := makeResponseJSON("error", "failed to read request body", 0, "")
		rw.Write(respJSON)
		return
	}

	useTestnetNetwork := true
	if network := chi.URLParam(req, "network"); network == "mainnet" {
		useTestnetNetwork = false
	}

	var assetDetails models.AssetRevoke

	err = json.Unmarshal(body, &assetDetails)

	if err != nil {
		h.log.WithError(err).Error("unable to unmarshal request into JSON")
		rw.WriteHeader(http.StatusBadRequest)

		respJSON := makeResponseJSON("error", "failed to unmarshal request body", 0, "")
		rw.Write(respJSON)
		return
	}

	_, claims, err := jwtauth.FromContext(req.Context())
	if err != nil {
		h.log.WithError(err).Error("failed to get jwt claims from context")
		rw.WriteHeader(http.StatusBadRequest)

		respJSON := makeResponseJSON("error", "failed to get jwt claims from context", 0, "")
		rw.Write(respJSON)
		return
	}

	privateKey, clawbackAddr, err := h.recoverAccount(claims["mnemonic"].(string))
	if err != nil {
		h.log.WithError(err).Error("failed to get private key from mnemonic")
		rw.WriteHeader(http.StatusBadRequest)

		respJSON := makeResponseJSON("error", "failed to recover account from mnemonic", 0, "")
		rw.Write(respJSON)
		return
	}

	assetDetails.ClawbackAddr = clawbackAddr

	txID, err := h.makeAndSendTxn(assetDetails, "revoke", privateKey, useTestnetNetwork)
	if err != nil {
		h.log.WithError(err).Error("failed to make and send asset revoke txn")
		rw.WriteHeader(http.StatusBadRequest)

		respJSON := makeResponseJSON("error", "failed to make and send asset revoke txn", 0, "")
		rw.Write(respJSON)
		return
	}
	h.log.Debug("Transaction ID: ", txID)

	respJSON := makeResponseJSON("success", "", assetDetails.AssetID, txID)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(respJSON)
}

func (h *ManagerHandler) waitForConfirmation(algodClient *algod.Client, txID string) error {
	for {
		pt, err := algodClient.PendingTransactionInformation(txID)
		if err != nil {
			h.log.WithError(err).Error("waiting for confirmation... (pool error, if any)")
			continue
		}
		if pt.ConfirmedRound > 0 {
			h.log.Debugf("Transaction "+pt.TxID+" confirmed in round %d", pt.ConfirmedRound)
			break
		}
		nodeStatus, err := algodClient.Status()
		if err != nil {
			h.log.WithError(err).Error("error getting algod status")
			return err
		}
		algodClient.StatusAfterBlock(nodeStatus.LastRound + 1)
	}

	return nil
}

// Jason Weathersby's function utility function to recover account and return sk and address :)
func (h *ManagerHandler) recoverAccount(userMnemonic string) (ed25519.PrivateKey, string, error) {
	sk, err := mnemonic.ToPrivateKey(userMnemonic)
	if err != nil {
		h.log.WithError(err).Error("error recovering account")
		return nil, "", err
	}
	pk := sk.Public()
	var a types.Address
	cpk := pk.(ed25519.PublicKey)
	copy(a[:], cpk[:])
	h.log.Debugf("Address: %s\n", a.String())
	address := a.String()
	return sk, address, nil
}

// Making and Sending Transactions:
func (h *ManagerHandler) makeAndSendTxn(request interface{}, txnType string, privateKey ed25519.PrivateKey, useTestnetNetwork bool) (string, error) {
	var txnParams algodModels.TransactionParams
	if useTestnetNetwork {
		txnParams, _ = h.testnetAlgod.SuggestedParams()
	} else {
		txnParams, _ = h.mainnetAlgod.SuggestedParams()
	}
	note := []byte(nil)
	gHash := base64.StdEncoding.EncodeToString(txnParams.GenesisHash)

	var txn types.Transaction

	switch assetDetails := request.(type) {
	case models.AssetCreate:
		txn, err := transaction.MakeAssetCreateTxn(assetDetails.CreatorAddr, txnParams.Fee, txnParams.LastRound, txnParams.LastRound+1000, note, txnParams.GenesisID, gHash, assetDetails.TotalIssuance, assetDetails.Decimals, assetDetails.DefaultFrozen, assetDetails.ManagerAddr, assetDetails.ReserveAddr, assetDetails.FreezeAddr, assetDetails.ClawbackAddr, assetDetails.UnitName, assetDetails.AssetName, assetDetails.URL, assetDetails.MetaDataHash)
		if err != nil {
			h.log.WithError(err).Error("Failed to make asset")
			return "", err
		}
		h.log.Debugf("Asset created AssetName: %s", txn.AssetConfigTxnFields.AssetParams.AssetName)
		break

	case models.AssetModify:
		_, err := transaction.MakeAssetConfigTxn(assetDetails.CurrManagerAddr, txnParams.Fee,
			txnParams.LastRound, txnParams.LastRound+1000, note, txnParams.GenesisID, gHash, assetDetails.AssetID, assetDetails.NewManagerAddr, assetDetails.NewReserveAddr, assetDetails.NewFreezeAddr, assetDetails.NewClawbackAddr, false)
		if err != nil {
			h.log.WithError(err).Error("Failed to modify asset")
			return "", err
		}
		h.log.Debugf("Asset modified AssetID: %s", assetDetails.AssetID)
		break

	case models.AssetDestroy:
		_, err := transaction.MakeAssetDestroyTxn(assetDetails.ManagerAddr, txnParams.Fee,
			txnParams.LastRound, txnParams.LastRound+1000, note, txnParams.GenesisID, gHash, assetDetails.AssetID)
		if err != nil {
			h.log.WithError(err).Error("Failed to destroy asset")
			return "", err
		}
		h.log.Debugf("Asset destroyed AssetID: %s", assetDetails.AssetID)
		break

	case models.AssetFreeze:
		_, err := transaction.MakeAssetFreezeTxn(assetDetails.FreezeAddr, txnParams.Fee,
			txnParams.LastRound, txnParams.LastRound+1000, note, txnParams.GenesisID, gHash, assetDetails.AssetID, assetDetails.TargetAddr, assetDetails.FreezeSetting)
		if err != nil {
			h.log.WithError(err).Error("Failed to freeze asset")
			return "", err
		}
		h.log.Debugf("Asset freezed AssetID: %s", assetDetails.AssetID)
		break

	case models.AssetRevoke:
		_, err := transaction.MakeAssetRevocationTxn(assetDetails.ClawbackAddr, assetDetails.TargetAddr, assetDetails.RecipientAddr, assetDetails.Amount, txnParams.Fee,
			txnParams.LastRound, txnParams.LastRound+1000, note, txnParams.GenesisID, gHash, assetDetails.AssetID)
		if err != nil {
			h.log.WithError(err).Error("Failed to revoke asset")
			return "", err
		}
		h.log.Debugf("Asset revoked AssetID: %s", assetDetails.AssetID)
		break

	default:
		return "", errors.New("invalid txn type")
	}

	txid, stx, err := crypto.SignTransaction(privateKey, txn)
	if err != nil {
		h.log.WithError(err).Error("Failed to sign transaction")
		return "", err
	}
	h.log.Debugf("Transaction ID: %s", txid)
	// Broadcast the transaction to the network
	var sendResponse algodModels.TransactionID
	if useTestnetNetwork {
		sendResponse, err = h.testnetAlgod.SendRawTransaction(stx, &algod.Header{Key: "Content-Type", Value: "application/x-binary"})
	} else {
		sendResponse, err = h.mainnetAlgod.SendRawTransaction(stx, &algod.Header{Key: "Content-Type", Value: "application/x-binary"})
	}
	if err != nil {
		h.log.WithError(err).Error("failed to send transaction")
		return "", err
	}

	// Wait for transaction to be confirmed
	if useTestnetNetwork {
		err = h.waitForConfirmation(h.testnetAlgod, sendResponse.TxID)
	} else {
		err = h.waitForConfirmation(h.mainnetAlgod, sendResponse.TxID)
	}
	if err != nil {
		return "", err
	}

	return sendResponse.TxID, nil
}
