package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"pb-backend/dataloader"
	"pb-backend/entities"
	"pb-backend/graph"
	"pb-backend/modifies"
	"pb-backend/wiregen"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/rs/cors"
	"gopkg.in/yaml.v2"
)

var (
	flags = flag.NewFlagSet("migrate", flag.ExitOnError)
	dir   = flags.String("dir", "/pb-backend/db/migrations", "directory with migration files")
)

const defaultPort = "3000"

func main() {
	r := chi.NewRouter()
	cfg, err := GetConfig()
	if err != nil {
		panic(err)
	}
	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(dataloader.Middleware)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
		Debug:            false,
	}).Handler)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))
	baseCtx := context.WithValue(context.Background(), "", "")
	baseCtx = context.WithValue(baseCtx, entities.ConfigKey, cfg)
	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.AppPort
		if port == "" {
			port = defaultPort
		}
	}
	app, err := wiregen.InitializeApp(baseCtx, *log.Default())
	if err != nil {
		log.Fatal(err)
	}
	config := graph.Config{
		Resolvers:  app.Resolver,
		Directives: graph.DirectiveRoot{},
		Complexity: graph.ComplexityRoot{},
	}
	config.Directives.Auth = modifies.Auth
	r.Use(app.CustomModifies.LoggingHandler)
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(config))
	srv.SetRecoverFunc(func(ctx context.Context, err interface{}) error {
		log.Println(fmt.Sprint(err))
		return errors.New("Internal server error! : ")
	})
	r.Handle("/", playground.Handler("GraphQL playground", "/query"))
	r.Handle("/query", srv)
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func GetConfig() (entities.PbConfig, error) {
	var config entities.PbConfig
	source, err := ioutil.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(source, &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return config, nil
}
