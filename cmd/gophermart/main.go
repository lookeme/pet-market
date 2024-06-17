package main

import (
	"pet-market/internal/configuration"
	"pet-market/internal/server"
)

func main() {
	cfg := configuration.New()
	s := server.NewHTTPServer(cfg)
	s.Start()
}
