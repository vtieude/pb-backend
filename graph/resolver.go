package graph

import (
	"context"
	"pb-backend/entities"
	"pb-backend/graph/model"
	"pb-backend/services"
)

//go:gen-gql go run github.com/99designs/gqlgen
// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	services.IUserService
	UserResolver
}

func (r *queryResolver) Friend(ctx context.Context, obj *entities.User) (*model.Friend, error) {
	return &model.Friend{City: &obj.Username.String}, nil
}

func (r *Resolver) User() UserResolver { return r.UserResolver }

// Mutation returns graph.MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns graph.QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
