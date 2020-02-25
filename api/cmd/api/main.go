package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/algorand/go-algorand-sdk/mnemonic"
	"github.com/algorand/go-algorand-sdk/types"

	"github.com/algorand/go-algorand-sdk/client/kmd"

	"github.com/algorand/go-algorand-sdk/client/algod"

	"github.com/sirupsen/logrus"
)

const algodAddress = "https://testnet-algorand.api.purestake.io/ps1"
const psToken = "FS0ZoE4JAe6MWL1CiJytR9nktogYSVC640C8fgk0"

const kmdAddress = "http://299799fb.ngrok.io"
const kmdToken = "ff590ebc0ac6793eb075dcbcb48df407a17008514612a99ad154be5c8f49eb9e"

func main() {

	logger := logrus.New()

	logrus.SetFormatter(&logrus.JSONFormatter{})

	logrus.SetOutput(os.Stdout)

	logrus.SetLevel(logrus.DebugLevel)

	retCode := 0

	// Exit with return code
	defer func() {
		logger.Exit(retCode)
		os.Exit(retCode)
	}()

	// Catch panics
	defer func() {
		if err := recover(); err != nil {
			logger.Panic(err.(string))
			retCode = 1
			return
		}
	}()

	var headers []*algod.Header
	headers = append(headers, &algod.Header{"X-API-Key", psToken})
	algodClient, err := algod.MakeClientWithHeaders(algodAddress, "", headers)
	if err != nil {
		logger.WithError(err).Fatal("failed to make algod client")
	}

	nodeStatus, err := algodClient.Status()
	if err != nil {
		logger.WithError(err).Error("Error getting algod status")
		return
	}

	logger.Printf("algod last round: %d\n", nodeStatus.LastRound)
	logger.Printf("algod time since last round: %d\n", nodeStatus.TimeSinceLastRound)
	logger.Printf("algod catchup: %d\n", nodeStatus.CatchupTime)
	logger.Printf("algod latest version: %s\n", nodeStatus.LastVersion)

	kmdClient, err := kmd.MakeClient(kmdAddress, kmdToken)
	if err != nil {
		logger.WithError(err).Fatal("failed to make kmd client")
		return
	}

	backupPhrase := "fire enlist diesel stamp nuclear chunk student stumble call snow flock brush example slab guide choice option recall south kangaroo hundred matrix school above zero"
	keyBytes, err := mnemonic.ToKey(backupPhrase)
	if err != nil {
		fmt.Printf("failed to get key: %s\n", err)
		return
	}

	var mdk types.MasterDerivationKey
	copy(mdk[:], keyBytes)
	cwResponse, err := kmdClient.CreateWallet("testwallet", "testpassword", kmd.DefaultWalletDriver, mdk)
	if err != nil {
		fmt.Printf("error creating wallet: %s\n", err)
		return
	}
	fmt.Printf("Created wallet '%s' with ID: %s\n", cwResponse.Wallet.Name, cwResponse.Wallet.ID)

	routerService := NewRouterService(logger)

	logger.Infoln("Starting server on port 5000")
	logger.Fatal(http.ListenAndServe(":5000", routerService))
}
