package main

import (
	"net/http"
	"os"

	"github.com/go-chi/jwtauth"

	"github.com/mangodev-io/mango/api/cmd/api/config"

	"github.com/algorand/go-algorand-sdk/client/algod"

	"github.com/sirupsen/logrus"
)

func main() {

	logger := logrus.New()

	logrus.SetFormatter(&logrus.JSONFormatter{})

	logrus.SetOutput(os.Stdout)

	logrus.SetLevel(logrus.TraceLevel)

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

	// Load Configuration
	config, err := config.New()
	if err != nil {
		logger.WithError(err).Panic("Missing configuration")
	}
	logger.Info("Loaded configuration...")

	// Setup Algod Client
	var headers []*algod.Header
	headers = append(headers, &algod.Header{"X-API-Key", config.PSToken})
	algodClient, err := algod.MakeClientWithHeaders(config.AlgodAddress, "", headers)
	if err != nil {
		logger.WithError(err).Panic("failed to make algod client")
	}

	// Setup JWT Auth
	tokenAuth := jwtauth.New("HS256", []byte(config.TokenAuthPassword), nil)

	// Setup Router
	routerService := NewRouterService(logger, &algodClient, tokenAuth)

	// Serve
	logger.Info("Starting server on port 5000")
	logger.Fatal(http.ListenAndServe(":5000", routerService))
}
