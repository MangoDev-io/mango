package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

type Router struct {
	*chi.Mux
	log *logrus.Logger
}

// Return a new instance of the Router
func NewRouterService(logger *logrus.Logger) *Router {
	router := chi.NewRouter()

	router.Use(middleware.Logger, middleware.RedirectSlashes)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	service := &Router{router, logger}
	return service
}
