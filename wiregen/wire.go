//go:build wireinject

package wiregen

import (
	"context"
	"github.com/google/wire"
	"log"
	"pb-backend/entities"
	"pb-backend/graph"
	"pb-backend/modifies"
	"pb-backend/services"
)

var serviceSet = wire.NewSet(
	services.NewUserService,
	modifies.ModifiesSet,
)
var dbSet = wire.NewSet(entities.OpenConnectTion, wire.Bind(new(entities.DB), new(*entities.DBConnection)))

type App struct {
	Resolver       *graph.Resolver
	CustomModifies *modifies.MyCustomHttpHandler
}

func InitializeApp(ctx context.Context, log log.Logger) (*App, error) {
	wire.Build(
		dbSet,
		wire.Struct(new(graph.Resolver), "*"),
		wire.Struct(new(App), "*"),
		serviceSet,
	)
	return &App{}, nil
}
