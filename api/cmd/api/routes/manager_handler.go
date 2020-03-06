package routes

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"

	"github.com/go-chi/jwtauth"

	"github.com/algorand/go-algorand-sdk/client/algod"
	"github.com/algorand/go-algorand-sdk/client/kmd"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/mnemonic"
	"github.com/algorand/go-algorand-sdk/transaction"
	"github.com/algorand/go-algorand-sdk/types"
	"github.com/haardikk21/algorand-asset-manager/api/cmd/api/constants"
	"github.com/haardikk21/algorand-asset-manager/api/cmd/api/data"
	"github.com/haardikk21/algorand-asset-manager/api/cmd/api/models"
	"github.com/sirupsen/logrus"
)

// ManagerHandler is the router handler for Mango
type ManagerHandler struct {
	log   *logrus.Logger
	db    *data.DatabaseService
	kmd   *kmd.Client
	algod *algod.Client
	jwt   *jwtauth.JWTAuth
}

type response struct {
	AssetID uint64 `json:"assetId"`
	TXHash  string `json:"txHash"`
}

// NewManagerHandler creates a new instance of ManagerHandler
func NewManagerHandler(log *logrus.Logger, db *data.DatabaseService, kmd *kmd.Client, algod *algod.Client, jwt *jwtauth.JWTAuth) *ManagerHandler {
	return &ManagerHandler{
		log:   log,
		db:    db,
		kmd:   kmd,
		algod: algod,
		jwt:   jwt,
	}
}

// EncodeMnemonic accepts JSON body of the form { "mnemonic" : "abcd" }
// and converts it into a JWT and send it to web
func (h *ManagerHandler) EncodeMnemonic(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		h.log.WithError(err).Error("failed to read body")
		rw.WriteHeader(http.StatusBadRequest)
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
		return
	}

	_, tokenString, err := h.jwt.Encode(jwt.MapClaims{"mnemonic": request.Mnemonic})
	if err != nil {
		h.log.WithError(err).Error("failed to encode jwt claims")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	type token struct {
		Token string `json:"token"`
	}

	var response token
	response.Token = tokenString

	respJSON, err := json.Marshal(response)
	if err != nil {
		h.log.WithError(err).Error("failed to marshal body")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(respJSON)
}

// GetAssets accepts a JSON body of the form { "address" : "example"  }
// and returns a list of all assets owned by the address (created on Mango)
func (h *ManagerHandler) GetAssets(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		h.log.WithError(err).Error("failed to read body")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	type getAssetReq struct {
		Address string `json:"address"`
	}

	var request getAssetReq
	err = json.Unmarshal(body, &request)
	if err != nil {
		h.log.WithError(err).Error("failed to unmarshal body")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	ownedAssets, err := h.db.SelectAllAssetsForAddress(request.Address)
	if err != nil {
		h.log.WithError(err).Error("failed to select rows")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	jsonResp, err := json.Marshal(ownedAssets)
	if err != nil {
		h.log.WithError(err).Error("failed to marshal owned assets")
		rw.WriteHeader(http.StatusInternalServerError)
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
		return
	}

	var assetDetails models.AssetCreate

	err = json.Unmarshal(body, &assetDetails)

	if err != nil {
		h.log.WithError(err).Error("failed to unmarshal body")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	_, claims, err := jwtauth.FromContext(req.Context())
	if err != nil {
		h.log.WithError(err).Error("failed to get jwt claims from context")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// privateKey, address, err := h.getPrivKeyAndAddressFromMnemonic(claims["mnemonic"].(string))
	privateKey, address := h.recoverAccount(claims["mnemonic"].(string))

	if err != nil {
		h.log.WithError(err).Error("failed to get private key from mnemonic")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	assetDetails.CreatorAddr = address

	txID, err := h.makeAndSendAssetCreateTxn(assetDetails, privateKey)
	if err != nil {
		h.log.WithError(err).Error("failed to make and send asset create txn")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Retrieve asset ID by grabbing the max asset ID
	// from the creator account's holdings.
	act, err := h.algod.AccountInformation(constants.TestAccountPublicKey)
	if err != nil {
		h.log.WithError(err).Error("failed to get account information")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	assetID := uint64(0)
	for i := range act.AssetParams {
		if i > assetID {
			assetID = i
		}
	}
	h.log.Debugf("Asset ID from AssetParams: %d", assetID)
	// Retrieve asset info.
	assetInfo, err := h.algod.AssetInformation(assetID)
	h.log.Debugf("Asset info: %#v", assetInfo)

	err = h.db.InsertNewAsset(assetDetails.CreatorAddr, strconv.FormatUint(assetID, 10))
	if err != nil {
		h.log.WithError(err).Error("failed to insert new asset in database")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := response{AssetID: assetID, TXHash: txID}
	respJSON, err := json.Marshal(resp)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(respJSON)
}

func (h *ManagerHandler) DestroyAsset(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		h.log.WithError(err).Error("unable to read request body")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	var assetDetails models.AssetDestroy

	err = json.Unmarshal(body, &assetDetails)

	if err != nil {
		h.log.WithError(err).Error("unable to unmarshal request into JSON")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	privateKey, managerAddr := h.recoverAccount(constants.TestAccountMnemonic)
	if err != nil {
		h.log.WithError(err).Error("failed to get private key from mnemonic")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	assetDetails.ManagerAddr = managerAddr

	txID, err := h.makeAndSendAssetDestroyTxn(assetDetails, privateKey)
	if err != nil {
		h.log.WithError(err).Error("failed to make and send asset destroy txn")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	h.log.Debug("Transaction ID: ", txID)

	resp := response{AssetID: 0, TXHash: txID}
	respJSON, err := json.Marshal(resp)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(respJSON)
}

func (h *ManagerHandler) FreezeAsset(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		h.log.WithError(err).Error("unable to read request body")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	var assetDetails models.AssetFreeze

	err = json.Unmarshal(body, &assetDetails)

	if err != nil {
		h.log.WithError(err).Error("unable to unmarshal request into JSON")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	privateKey, freezeAddr := h.recoverAccount(constants.TestAccountMnemonic)
	if err != nil {
		h.log.WithError(err).Error("failed to get private key from mnemonic")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	assetDetails.FreezeAddr = freezeAddr

	txID, err := h.makeAndSendAssetFreezeTxn(assetDetails, privateKey)
	if err != nil {
		h.log.WithError(err).Error("failed to make and send asset freeze txn")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	h.log.Debug("Transaction ID: ", txID)

	resp := response{AssetID: assetDetails.AssetID, TXHash: txID}
	respJSON, err := json.Marshal(resp)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(respJSON)
}

func (h *ManagerHandler) waitForConfirmation(algodClient *algod.Client, txID string) {
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
			return
		}
		algodClient.StatusAfterBlock(nodeStatus.LastRound + 1)
	}
}

// Jason Weathersby's function utility function to recover account and return sk and address :)
func (h *ManagerHandler) recoverAccount(userMnemonic string) (ed25519.PrivateKey, string) {
	sk, err := mnemonic.ToPrivateKey(userMnemonic)
	if err != nil {
		h.log.WithError(err).Error("error recovering account")
		return nil, ""
	}
	pk := sk.Public()
	var a types.Address
	cpk := pk.(ed25519.PublicKey)
	copy(a[:], cpk[:])
	h.log.Debugf("Address: %s\n", a.String())
	address := a.String()
	return sk, address
}

// Making and Sending Transactions:
func (h *ManagerHandler) makeAndSendAssetCreateTxn(assetDetails models.AssetCreate, privateKey ed25519.PrivateKey) (string, error) {

	txnParams, err := h.algod.SuggestedParams()
	note := []byte(nil)
	gHash := base64.StdEncoding.EncodeToString(txnParams.GenesisHash)

	txn, err := transaction.MakeAssetCreateTxn(assetDetails.CreatorAddr, txnParams.Fee, txnParams.LastRound, txnParams.LastRound+1000, note, txnParams.GenesisID, gHash, assetDetails.TotalIssuance, assetDetails.Decimals, assetDetails.DefaultFrozen, assetDetails.ManagerAddr, assetDetails.ReserveAddr, assetDetails.FreezeAddr, assetDetails.ClawbackAddr, assetDetails.UnitName, assetDetails.AssetName, assetDetails.URL, assetDetails.MetaDataHash)

	if err != nil {
		h.log.WithError(err).Error("Failed to make asset")
		return "", err
	}
	h.log.Debugf("Asset created AssetName: %s", txn.AssetConfigTxnFields.AssetParams.AssetName)

	txid, stx, err := crypto.SignTransaction(privateKey, txn)
	if err != nil {
		h.log.WithError(err).Error("Failed to sign transaction")
		return "", err
	}
	h.log.Debugf("Transaction ID: %s", txid)
	// Broadcast the transaction to the network
	sendResponse, err := h.algod.SendRawTransaction(stx, &algod.Header{Key: "Content-Type", Value: "application/x-binary"})
	if err != nil {
		h.log.WithError(err).Error("failed to send transaction")
		return "", err
	}

	// Wait for transaction to be confirmed
	h.waitForConfirmation(h.algod, sendResponse.TxID)

	return sendResponse.TxID, nil
}

func (h *ManagerHandler) makeAndSendAssetDestroyTxn(assetDetails models.AssetDestroy, privateKey ed25519.PrivateKey) (string, error) {
	txnParams, err := h.algod.SuggestedParams()
	note := []byte(nil)
	gHash := base64.StdEncoding.EncodeToString(txnParams.GenesisHash)

	txn, err := transaction.MakeAssetDestroyTxn(assetDetails.ManagerAddr, txnParams.Fee,
		txnParams.LastRound, txnParams.LastRound+1000, note, txnParams.GenesisID, gHash, assetDetails.AssetID)

	if err != nil {
		h.log.WithError(err).Error("failed to send txn")
		return "", err
	}

	txid, stx, err := crypto.SignTransaction(privateKey, txn)
	if err != nil {
		h.log.WithError(err).Error("Failed to sign transaction")
		return "", err
	}
	h.log.Debugf("Transaction ID: %s", txid)
	// Broadcast the transaction to the network
	sendResponse, err := h.algod.SendRawTransaction(stx, &algod.Header{Key: "Content-Type", Value: "application/x-binary"})
	if err != nil {
		h.log.WithError(err).Error("failed to send transaction")
		return "", err
	}
	h.log.Infof("Transaction ID raw: %s", sendResponse.TxID)
	// Wait for transaction to be confirmed
	h.waitForConfirmation(h.algod, sendResponse.TxID)

	return sendResponse.TxID, nil
}

func (h *ManagerHandler) makeAndSendAssetFreezeTxn(assetDetails models.AssetFreeze, privateKey ed25519.PrivateKey) (string, error) {
	txnParams, err := h.algod.SuggestedParams()
	note := []byte(nil)
	gHash := base64.StdEncoding.EncodeToString(txnParams.GenesisHash)

	txn, err := transaction.MakeAssetFreezeTxn(assetDetails.FreezeAddr, txnParams.Fee,
		txnParams.LastRound, txnParams.LastRound+1000, note, txnParams.GenesisID, gHash, assetDetails.AssetID, assetDetails.TargetAddr, assetDetails.FreezeSetting)

	if err != nil {
		h.log.WithError(err).Error("failed to send txn")
		return "", err
	}

	txid, stx, err := crypto.SignTransaction(privateKey, txn)
	if err != nil {
		h.log.WithError(err).Error("Failed to sign transaction")
		return "", err
	}
	h.log.Debugf("Transaction ID: %s", txid)
	// Broadcast the transaction to the network
	sendResponse, err := h.algod.SendRawTransaction(stx, &algod.Header{Key: "Content-Type", Value: "application/x-binary"})
	if err != nil {
		h.log.WithError(err).Error("failed to send transaction")
		return "", err
	}
	h.log.Infof("Transaction ID raw: %s", sendResponse.TxID)
	// Wait for transaction to be confirmed
	h.waitForConfirmation(h.algod, sendResponse.TxID)

	return sendResponse.TxID, nil
}
