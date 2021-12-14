package graph

import "pb-backend/services"

//go:gen-gql go run github.com/99designs/gqlgen
// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	services.IUserService
}
