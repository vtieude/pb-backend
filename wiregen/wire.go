//go:build wireinject

package wiregen

import (
	"context"
	"log"
	"pb-backend/entities"
	"pb-backend/graph"
	"pb-backend/graph/resolvers"
	"pb-backend/modifies"
	"pb-backend/services"

	"github.com/google/wire"
)

var serviceSet = wire.NewSet(
	services.NewUserService,
	services.NewProductService,
	services.NewSaleService,
	modifies.ModifiesSet,
	services.NewGoogleService,
)
var dbSet = wire.NewSet(entities.OpenConnection, wire.Bind(new(entities.DB), new(*entities.DBConnection)))

type App struct {
	Resolver       *graph.Resolver
	CustomModifies *modifies.MyCustomHttpHandler
}

func InitializeApp(ctx context.Context, log log.Logger) (*App, error) {
	wire.Build(
		dbSet,
		resolvers.NewUserSet,
		wire.Struct(new(graph.Resolver), "*"),
		wire.Struct(new(App), "*"),
		serviceSet,
	)
	return &App{}, nil
}
