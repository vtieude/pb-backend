package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"pb-backend/graph"
	"pb-backend/wiregen"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/rs/cors"
)

const defaultPort = "3000"

func main() {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
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
	port := os.Getenv("PORT")
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
