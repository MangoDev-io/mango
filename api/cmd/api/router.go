package main

import (
	"net/http"

	"github.com/haardikk21/algorand-asset-manager/api/cmd/api/data"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

type Router struct {
	*chi.Mux
	log *logrus.Logger
	db  *data.DatabaseService
}

// Return a new instance of the Router
func NewRouterService(logger *logrus.Logger, db *data.DatabaseService) *Router {
	router := chi.NewRouter()

	router.Use(middleware.Logger, middleware.RedirectSlashes)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	service := &Router{router, logger, db}
	return service
}
