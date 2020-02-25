package main

import (
	"net/http"
	"os"

	"github.com/algorand/go-algorand-sdk/client/algod"

	"github.com/sirupsen/logrus"
)

const algodAddress = "https://testnet-algorand.api.purestake.io/ps1"
const psToken = "FS0ZoE4JAe6MWL1CiJytR9nktogYSVC640C8fgk0"

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

	routerService := NewRouterService(logger)

	logger.Infoln("Starting server on port 5000")
	logger.Fatal(http.ListenAndServe(":5000", routerService))
}
