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
	privateKey, address, err := h.getPrivKeyAndAddressFromMnemonic(claims["mnemonic"].(string))

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

	err = h.deleteAddressFromWallet(assetDetails.CreatorAddr)
	if err != nil {
		h.log.WithError(err).Error("Error deleting address from wallet")
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

	privateKey, managerAddr, err := h.getPrivKeyAndAddressFromMnemonic(constants.TestAccountMnemonic)
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

	// Delete current address from wallet
	err = h.deleteAddressFromWallet(managerAddr)
	if err != nil {
		h.log.WithError(err).Error("Error deleting address from wallet")
	}

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

	privateKey, freezeAddr, err := h.getPrivKeyAndAddressFromMnemonic(constants.TestAccountMnemonic)
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

	// Delete current address from wallet
	err = h.deleteAddressFromWallet(freezeAddr)
	if err != nil {
		h.log.WithError(err).Error("Error deleting address from wallet")
	}

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

// Wallet Helper Functions ---- // TODO - MAKE WALLETID A GLOBAL VARIABLE
func (h *ManagerHandler) getPrivKeyAndAddressFromMnemonic(accountMnemonic string) (ed25519.PrivateKey, string, error) {
	// Import Account from Account Mnemonic --------------------------------------
	// Get the list of wallets
	listResponse, err := h.kmd.ListWallets()
	if err != nil {
		h.log.WithError(err).Error("error listing wallets when importing mnemonic")
		return nil, "", err
	}

	// Find our wallet name in the list
	var walletID string
	for _, wallet := range listResponse.Wallets {
		if wallet.Name == constants.TestWalletName {
			h.log.Debugf("Got Wallet '%s' with ID: %s", wallet.Name, wallet.ID)
			walletID = wallet.ID
		}
	}

	// Get a wallet handle
	initResponse, err := h.kmd.InitWalletHandle(walletID, constants.TestWalletPassword)
	if err != nil {
		h.log.WithError(err).Error("Error initializing wallet handle")
		return nil, "", err
	}

	h.log.Debugf("Account Mnemonic: %s", accountMnemonic)
	privateKey, err := mnemonic.ToPrivateKey(accountMnemonic)
	importedAccount, err := h.kmd.ImportKey(initResponse.WalletHandleToken, privateKey)
	h.log.Debugf("Account Successfully Imported: %s", importedAccount)

	return privateKey, importedAccount.Address, nil
}

func (h *ManagerHandler) deleteAddressFromWallet(address string) error {
	listResponse, err := h.kmd.ListWallets()
	if err != nil {
		h.log.WithError(err).Error("error listing wallets when deleting")
		return err
	}

	var walletID string
	for _, wallet := range listResponse.Wallets {
		if wallet.Name == constants.TestWalletName {
			h.log.Debugf("Got Wallet '%s' with ID: %s", wallet.Name, wallet.ID)
			walletID = wallet.ID
		}
	}

	initResponse, err := h.kmd.InitWalletHandle(walletID, constants.TestWalletPassword)
	if err != nil {
		h.log.WithError(err).Error("Error initializing wallet handle")
		return err
	}

	h.kmd.DeleteKey(initResponse.WalletHandleToken, constants.TestWalletPassword, address)
	return nil
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
