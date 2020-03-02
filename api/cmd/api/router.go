package main

import (
	"github.com/algorand/go-algorand-sdk/client/algod"
	"github.com/algorand/go-algorand-sdk/client/kmd"
	"github.com/haardikk21/algorand-asset-manager/api/cmd/api/data"
	"github.com/haardikk21/algorand-asset-manager/api/cmd/api/routes"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

type Router struct {
	*chi.Mux
	log   *logrus.Logger
	db    *data.DatabaseService
	kmd   *kmd.Client
	algod *algod.Client
}

// Return a new instance of the Router
func NewRouterService(logger *logrus.Logger, db *data.DatabaseService, kmd *kmd.Client, algod *algod.Client) *Router {
	router := chi.NewRouter()

	router.Use(middleware.Logger, middleware.RedirectSlashes)

	managerHandler := routes.NewManagerHandler(logger, db, kmd, algod)

	router.Get("/", managerHandler.GetHello)
	router.Post("/createAsset", managerHandler.CreateAsset)

	service := &Router{router, logger, db, kmd, algod}
	return service
}
