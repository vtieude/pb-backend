package resolvers

import (
	"context"
	"pb-backend/entities"
)

type UserResolver struct {
}

// var NewUserSet = wire.NewSet(wire.Struct(new(UserResolver), "*"), wire.Bind(new(graph.UserResolver), new(*UserResolver)))

func (us *UserResolver) Username(ctx context.Context, obj *entities.User) (string, error) {
	return obj.Username, nil
}
