package resolvers

import (
	"context"
	"pb-backend/entities"
	"pb-backend/graph"

	"github.com/google/wire"
)

type UserResolver struct {
}

var NewUserSet = wire.NewSet(wire.Struct(new(UserResolver), "*"), wire.Bind(new(graph.UserResolver), new(*UserResolver)))

func (us *UserResolver) Username(ctx context.Context, obj *entities.User) (string, error) {
	return obj.Username, nil
}

func (us *UserResolver) PhoneNumber(ctx context.Context, obj *entities.User) (*string, error) {
	if obj.PhoneNumber.Valid {
		return &obj.PhoneNumber.String, nil
	}
	return nil, nil
}
