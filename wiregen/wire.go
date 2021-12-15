//go:build wireinject
package wiregen

import (
	"context"
	"pb-backend/graph"
	"pb-backend/services"

	"github.com/google/wire"
)

var serviceSet = wire.NewSet(
	services.NewUserService,
)

type App struct {
	Resolver *graph.Resolver
}

func InitializeApp(ctx context.Context) (*App, error) {
	wire.Build(
		wire.Struct(new(graph.Resolver), "*"),
		wire.Struct(new(App), "*"),
		serviceSet,
	)
	return &App{}, nil
}
