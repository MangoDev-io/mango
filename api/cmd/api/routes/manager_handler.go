package routes

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"

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

type ManagerHandler struct {
	log *logrus.Logger
	db  *data.DatabaseService

	kmd   *kmd.Client
	algod *algod.Client
}

func NewManagerHandler(log *logrus.Logger, db *data.DatabaseService, kmd *kmd.Client, algod *algod.Client) *ManagerHandler {
	return &ManagerHandler{
		log:   log,
		db:    db,
		kmd:   kmd,
		algod: algod,
	}
}

func (h *ManagerHandler) GetHello(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte("{ message: 'it works' }"))
}

func (h *ManagerHandler) GetAssets(rw http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)

	type getAssetReq struct {
		Address string `json:"address"`
	}

	var request getAssetReq
	_ = json.Unmarshal(body, &request)

	ownedassets, err := h.db.SelectAllAssetsForAddress(request.Address)
	if err != nil {
		h.log.WithError(err).Error("cabnnot select")
	}
	h.log.Info(ownedassets.AssetIds)

	jsonResp, _ := json.Marshal(ownedassets)

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(jsonResp)
}

func (h *ManagerHandler) CreateAsset(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	var assetDetails models.AssetCreate

	err = json.Unmarshal(body, &assetDetails)

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	privateKey, err := h.getPrivateKeyFromMnemonic(constants.TestAccountMnemonic)
	if err != nil {
		h.log.WithError(err).Error("failed to get private key from mnemonic")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Create CreateAsset Transaction
	txnParams, err := h.algod.SuggestedParams()
	note := []byte(nil)

	gHash := base64.StdEncoding.EncodeToString(txnParams.GenesisHash)

	txn, err := transaction.MakeAssetCreateTxn(assetDetails.CreatorAddr, txnParams.Fee, txnParams.LastRound, txnParams.LastRound+1000, note, txnParams.GenesisID, gHash, assetDetails.TotalIssuance, assetDetails.Decimals, assetDetails.DefaultFrozen, assetDetails.ManagerAddr, assetDetails.ReserveAddr, assetDetails.FreezeAddr, assetDetails.ClawbackAddr, assetDetails.UnitName, assetDetails.AssetName, assetDetails.URL, assetDetails.MetaDataHash)

	if err != nil {
		h.log.WithError(err).Error("Failed to make asset")
		return
	}
	h.log.Debugf("Asset created AssetName: %s", txn.AssetConfigTxnFields.AssetParams.AssetName)

	txid, stx, err := crypto.SignTransaction(privateKey, txn)
	if err != nil {
		h.log.WithError(err).Error("Failed to sign transaction")
		return
	}
	h.log.Debugf("Transaction ID: %s", txid)
	// Broadcast the transaction to the network
	sendResponse, err := h.algod.SendRawTransaction(stx, &algod.Header{Key: "Content-Type", Value: "application/x-binary"})
	if err != nil {
		h.log.WithError(err).Error("failed to send transaction")
		return
	}

	// Wait for transaction to be confirmed
	h.waitForConfirmation(h.algod, sendResponse.TxID)

	// Retrieve asset ID by grabbing the max asset ID
	// from the creator account's holdings.
	act, err := h.algod.AccountInformation(constants.TestAccountPublicKey)
	if err != nil {
		h.log.WithError(err).Error("failed to get account information")
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

	type response struct {
		AssetID uint64 `json:"assetId"`
	}

	resp := response{AssetID: assetID}
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

func (h *ManagerHandler) getPrivateKeyFromMnemonic(accountMnemonic string) (ed25519.PrivateKey, error) {
	// Import Account from Account Mnemonic --------------------------------------
	// Get the list of wallets
	listResponse, err := h.kmd.ListWallets()
	if err != nil {
		h.log.WithError(err).Error("error listing wallets")
		return nil, err
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
		return nil, err
	}

	h.log.Debugf("Account Mnemonic: %s", accountMnemonic)
	privateKey, err := mnemonic.ToPrivateKey(accountMnemonic)
	importedAccount, err := h.kmd.ImportKey(initResponse.WalletHandleToken, privateKey)
	h.log.Debugf("Account Successfully Imported: %s", importedAccount)

	return privateKey, nil
}
