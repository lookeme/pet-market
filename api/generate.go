package api

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest -package api -generate "types,chi-server,spec" -o shop.gen.go  openapi.yaml
