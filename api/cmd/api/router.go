package main

import (
	"github.com/algorand/go-algorand-sdk/client/algod"
	"github.com/algorand/go-algorand-sdk/client/kmd"
	"github.com/go-chi/jwtauth"
	"github.com/haardikk21/algorand-asset-manager/api/cmd/api/data"
	"github.com/haardikk21/algorand-asset-manager/api/cmd/api/routes"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

// Router is the HTTP Router service
type Router struct {
	*chi.Mux
	log   *logrus.Logger
	db    *data.DatabaseService
	kmd   *kmd.Client
	algod *algod.Client
	jwt   *jwtauth.JWTAuth
}

// NewRouterService return a new instance of the Router
func NewRouterService(logger *logrus.Logger, db *data.DatabaseService, kmd *kmd.Client, algod *algod.Client, jwt *jwtauth.JWTAuth) *Router {
	router := chi.NewRouter()

	router.Use(middleware.Logger, middleware.RedirectSlashes)

	managerHandler := routes.NewManagerHandler(logger, db, kmd, algod, jwt)

	// JWT Protected Routes
	router.Group(func(router chi.Router) {
		router.Use(jwtauth.Verifier(jwt))
		router.Use(jwtauth.Authenticator)

		router.Post("/createAsset", managerHandler.CreateAsset)
		router.Get("/assets", managerHandler.GetAssets)
	})

	// Public Routes
	router.Group(func(router chi.Router) {
		router.Post("/encryptMnemonic", managerHandler.EncryptMnemonic)
	})

	service := &Router{router, logger, db, kmd, algod, jwt}
	return service
}
