package main

import (
	"github.com/algorand/go-algorand-sdk/client/algod"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth"
	"github.com/mangodev-io/mango/api/cmd/api/routes"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

// Router is the HTTP Router service
type Router struct {
	*chi.Mux
	log   *logrus.Logger
	algod *algod.Client
	jwt   *jwtauth.JWTAuth
}

// NewRouterService return a new instance of the Router
func NewRouterService(logger *logrus.Logger, algod *algod.Client, jwt *jwtauth.JWTAuth) *Router {
	router := chi.NewRouter()

	cors := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://mangodev-io", "https://www.mangodev-io"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	router.Use(cors.Handler, middleware.Logger, middleware.RedirectSlashes)

	managerHandler := routes.NewManagerHandler(logger, algod, jwt)

	// JWT Protected Routes
	router.Group(func(router chi.Router) {
		router.Use(jwtauth.Verifier(jwt))
		router.Use(jwtauth.Authenticator)

		router.Post("/createAsset", managerHandler.CreateAsset)
		router.Post("/modifyAsset", managerHandler.ModifyAsset)
		router.Post("/destroyAsset", managerHandler.DestroyAsset)
		router.Post("/freezeAsset", managerHandler.FreezeAsset)
		router.Post("/revokeAsset", managerHandler.RevokeAsset)
		router.Get("/assets", managerHandler.GetAssets)
	})

	// Public Routes
	router.Group(func(router chi.Router) {
		router.Post("/encodeMnemonic", managerHandler.EncodeMnemonic)
	})

	service := &Router{router, logger, algod, jwt}
	return service
}
