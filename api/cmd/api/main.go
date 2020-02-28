package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/haardikk21/algorand-asset-manager/api/cmd/api/config"

	"github.com/algorand/go-algorand-sdk/mnemonic"
	"github.com/algorand/go-algorand-sdk/types"

	"github.com/algorand/go-algorand-sdk/client/kmd"

	"github.com/algorand/go-algorand-sdk/client/algod"

	"github.com/sirupsen/logrus"
)

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

	config, err := config.New()
	if err != nil {
		logger.WithError(err).Panic("Missing configuration")
	}

	var headers []*algod.Header
	headers = append(headers, &algod.Header{"X-API-Key", config.PSToken})
	algodClient, err := algod.MakeClientWithHeaders(config.AlgodAddress, "", headers)
	if err != nil {
		logger.WithError(err).Panic("failed to make algod client")
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

	kmdClient, err := kmd.MakeClient(config.KMDAddress, config.KMDToken)
	if err != nil {
		logger.WithError(err).Panic("failed to make kmd client")
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
	cwResponse, err := kmdClient.CreateWallet("sdk,jfgsdlkufjgasdfbsaldikf", "testpassword", kmd.DefaultWalletDriver, mdk)
	if err != nil {
		fmt.Printf("error creating wallet: %s\n", err)
		return
	}
	fmt.Printf("Created wallet '%s' with ID: %s\n", cwResponse.Wallet.Name, cwResponse.Wallet.ID)

	routerService := NewRouterService(logger)

	logger.Infoln("Starting server on port 5000")
	logger.Fatal(http.ListenAndServe(":5000", routerService))
}
