package main

import (
	"fmt"
	"net/http"
	"os"
	"pet-market/api"
	"pet-market/internal/configuration"
	"pet-market/internal/controller"
	"pet-market/internal/logger"

	"github.com/go-chi/chi/v5"
	middleware "github.com/oapi-codegen/nethttp-middleware"
)

func main() {
	cfg := configuration.New()
	log, err := logger.CreateLogger(cfg.Logger)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing logger\n: %s", err)
		os.Exit(1)
	}
	swagger, err := api.GetSwagger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
		os.Exit(1)
	}
	swagger.Servers = nil
	httpServer := controller.NewController()
	r := chi.NewRouter()
	r.Use(middleware.OapiRequestValidator(swagger))
	r.Use(log.Middleware)
	api.HandlerFromMux(httpServer, r)
	s := &http.Server{
		Handler: r,
		Addr:    cfg.Network.ServerAddress,
	}
	log.Log.Fatal(s.ListenAndServe().Error())
}
