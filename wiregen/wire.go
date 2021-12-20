//go:build wireinject
package wiregen

import (
	"context"
	"pb-backend/db"
	"pb-backend/graph"
	"pb-backend/graph/resolvers"
	"pb-backend/services"

	"github.com/google/wire"
)

var serviceSet = wire.NewSet(
	services.NewUserService,
	resolvers.NewUserSet,
)
var dbSet = wire.NewSet(db_manager.OpenConnectTion, wire.Bind(new(db_manager.IDb), new(*db_manager.DB)))

type App struct {
	Resolver *graph.Resolver
}

func InitializeApp(ctx context.Context) (*App, error) {
	wire.Build(
		dbSet,
		wire.Struct(new(graph.Resolver), "*"),
		wire.Struct(new(App), "*"),
		serviceSet,
	)
	return &App{}, nil
}
