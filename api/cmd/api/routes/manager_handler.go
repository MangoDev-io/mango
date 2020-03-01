package routes

import (
	"encoding/json"
	"fmt"
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

	// Import Account from Account Mnemonic --------------------------------------
	// Get the list of wallets
	listResponse, err := h.kmd.ListWallets()
	if err != nil {
		h.log.WithError(err).Error("error listing wallets: %s\n")
		return
	}

	// Find our wallet name in the list
	var walletID string
	for _, wallet := range listResponse.Wallets {
		if wallet.Name == constants.TestWalletName {
			h.log.Info("Got Wallet '%s' with ID: %s\n", wallet.Name, wallet.ID)
			walletID = wallet.ID
		}
	}

	// Get a wallet handle
	initResponse, err := h.kmd.InitWalletHandle(walletID, constants.TestWalletPassword)
	if err != nil {
		h.log.WithError(err).Error("Error initializing wallet handle")
		return
	}

	h.log.Info("Account Mnemonic: %s\n", constants.TestAccountMnemonic)
	privateKey, err := mnemonic.ToPrivateKey(constants.TestAccountMnemonic)
	importedAccount, err := h.kmd.ImportKey(initResponse.WalletHandleToken, privateKey)
	h.log.Info("Account Sucessfully Imported: ", importedAccount)

	// Create CreateAsset Transaction
	txnParams, err := h.algod.SuggestedParams()
	note := []byte(nil)

	txn, err := transaction.MakeAssetCreateTxn(assetDetails.CreatorAddr, txnParams.Fee, txnParams.LastRound, txnParams.LastRound+1000, note, txnParams.GenesisID, string(txnParams.GenesisHash), assetDetails.TotalIssuance, assetDetails.Decimals, assetDetails.DefaultFrozen, assetDetails.ManagerAddr, assetDetails.ReserveAddr, assetDetails.FreezeAddr, assetDetails.ClawbackAddr, assetDetails.UnitName, assetDetails.AssetName, assetDetails.URL, assetDetails.MetaDataHash)

	if err != nil {
		h.log.WithError(err).Error("Failed to make asset")
		return
	}
	h.log.Info("Asset created AssetName: %s\n", txn.AssetConfigTxnFields.AssetParams.AssetName)

	txid, stx, err := crypto.SignTransaction(privateKey, txn)
	if err != nil {
		h.log.WithError(err).Error("Failed to sign transaction")
		return
	}
	h.log.Info("Transaction ID: %s\n", txid)
	// Broadcast the transaction to the network
	sendResponse, err := h.algod.SendRawTransaction(stx)
	if err != nil {
		h.log.WithError(err).Error("failed to send transaction")
		return
	}

	// Wait for transaction to be confirmed
	waitForConfirmation(h.algod, sendResponse.TxID)

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
	h.log.Infof("Asset ID from AssetParams: %d\n", assetID)
	// Retrieve asset info.
	assetInfo, err := h.algod.AssetInformation(assetID)
	h.log.Infof("Asset info: %#v\n", assetInfo)
}

func waitForConfirmation(algodClient *algod.Client, txID string) {
	for {
		pt, err := algodClient.PendingTransactionInformation(txID)
		if err != nil {
			fmt.Printf("waiting for confirmation... (pool error, if any): %s\n", err)
			continue
		}
		if pt.ConfirmedRound > 0 {
			fmt.Printf("Transaction "+pt.TxID+" confirmed in round %d\n", pt.ConfirmedRound)
			break
		}
		nodeStatus, err := algodClient.Status()
		if err != nil {
			fmt.Printf("error getting algod status: %s\n", err)
			return
		}
		algodClient.StatusAfterBlock(nodeStatus.LastRound + 1)
	}
}
