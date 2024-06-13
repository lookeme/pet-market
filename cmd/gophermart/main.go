package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"pet-market/api"
	"pet-market/internal/configuration"
	"pet-market/internal/controller"
	"pet-market/internal/integration"
	"pet-market/internal/logger"
	"pet-market/internal/repository"
	"pet-market/internal/security"
	"pet-market/internal/service"

	"github.com/getkin/kin-openapi/openapi3filter"
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
	ctx := context.Background()
	postgres, err := repository.New(ctx, log, cfg.Storage)
	if err != nil {
		log.Log.Fatal(err.Error())
	}
	swagger, err := api.GetSwagger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
		os.Exit(1)
	}
	swagger.Servers = nil
	auth := security.New(log)
	usrRepository := repository.NewUsrRepository(postgres)
	orderRepo := repository.NewOrderRepository(postgres)
	balanceRepo := repository.NewBalanceRepository(postgres)
	withdrawRepo := repository.NewWithdrawRepository(postgres)
	accural := integration.New(cfg.Network.AccuralAddress)
	orderService := service.NewOrderService(accural, orderRepo)
	usrService := service.NewUserService(usrRepository, auth, log)
	balanceService := service.NewBalanceService(balanceRepo, withdrawRepo)

	httpServer := controller.NewController(
		auth,
		usrService,
		orderService,
		balanceService,
		log)
	r := chi.NewRouter()
	r.Use(middleware.OapiRequestValidatorWithOptions(swagger, &middleware.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: openapi3filter.NoopAuthenticationFunc,
		},
	}))
	api.HandlerFromMux(httpServer, r)
	s := &http.Server{
		Handler: r,
		Addr:    cfg.Network.ServerAddress,
	}
	log.Log.Fatal(s.ListenAndServe().Error())
}
