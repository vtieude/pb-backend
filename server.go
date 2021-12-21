package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"pb-backend/dataloader"
	"pb-backend/entities"
	"pb-backend/graph"
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
	r.Use(middleware.Recoverer)
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
	port := cfg.AppPort
	if port == "" {
		port = defaultPort
	}
	app, err := wiregen.InitializeApp(baseCtx)
	if err != nil {
		log.Fatal(err)
	}
	config := graph.Config{
		Resolvers:  app.Resolver,
		Directives: graph.DirectiveRoot{},
		Complexity: graph.ComplexityRoot{},
	}
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(config))
	srv.SetRecoverFunc(func(ctx context.Context, err interface{}) error {
		return errors.New("Internal server error! : " + fmt.Sprint(err))
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
