package main

// import (
// 	"context"
// 	"crypto/md5"
// 	"fmt"
// 	"html/template"
// 	"io"
// 	"net/http"
// 	handler "pb-backend/httpHandler"
// 	"pb-backend/todo"
// 	"strconv"
// 	"strings"
// 	"time"

// 	"github.com/go-chi/chi/v5"
// 	"github.com/go-chi/chi/v5/middleware"
// 	"github.com/go-chi/jwtauth/v5"
// 	"github.com/go-chi/render"
// 	"github.com/rs/cors"
// )

// var tokenAuth *jwtauth.JWTAuth

// func init() {
// 	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)

// 	// For debugging/example purposes, we generate and print
// 	// a sample jwt token with claims `user_id:123` here:
// 	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"user_id": 123})
// 	fmt.Printf("DEBUG: a sample jwt is %s\n\n", tokenString)
// }

// func main22() {
// 	r := chi.NewRouter()

// 	// A good base middleware stack
// 	r.Use(middleware.RequestID)
// 	r.Use(middleware.RealIP)
// 	r.Use(middleware.Logger)
// 	r.Use(middleware.Recoverer)
// 	r.Use(render.SetContentType(render.ContentTypeJSON))
// 	r.Use(cors.New(cors.Options{
// 		AllowedOrigins:   []string{"https://*", "http://*"},
// 		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
// 		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
// 		AllowCredentials: true,
// 		Debug:            false,
// 	}).Handler)

// 	// Set a timeout value on the request context (ctx), that will signal
// 	// through ctx.Done() that the request has timed out and further
// 	// processing should be stopped.
// 	r.Use(middleware.Timeout(60 * time.Second))

// 	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
// 		w.Write([]byte("hi"))
// 	})

// 	// RESTy routes for "articles" resource
// 	r.Route("/todo", func(r chi.Router) {
// 		r.Get("/", handler.ListTodo)
// 		r.Post("/", handler.CreateTodo)
// 		// Subrouters:
// 		r.Route("/{todoId}", func(r chi.Router) {
// 			r.Use(TodoCtx)
// 			r.Put("/", handler.UpdateTodo)    // PUT /articles/123
// 			r.Delete("/", handler.DeleteTodo) // DELETE /articles/123
// 		})
// 	})
// 	// Protected routes
// 	r.Group(func(r chi.Router) {
// 		// Seek, verify and validate JWT tokens
// 		r.Use(jwtauth.Verifier(tokenAuth))

// 		// Handle valid / invalid tokens. In this example, we use
// 		// the provided authenticator middleware, but you can write your
// 		// own very easily, look at the Authenticator method in jwtauth.go
// 		// and tweak it, its not scary.
// 		r.Use(jwtauth.Authenticator)

// 		r.Get("/alltodo", handler.ListTodo)
// 		r.Get("/admin", func(w http.ResponseWriter, r *http.Request) {
// 			_, claims, _ := jwtauth.FromContext(r.Context())
// 			w.Write([]byte(fmt.Sprintf("protected area. hi %v", claims["user_id"])))
// 		})
// 	})
// 	// Public routes
// 	// r.Group(func(r chi.Router) {
// 	// 	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
// 	// 		w.Write([]byte("welcome anonymous"))
// 	// 	})
// 	// })

// 	http.ListenAndServe(":3000", r)
// }

// // ArticleCtx middleware is used to load an Article object from
// // the URL parameters passed through as the request. In case
// // the Article could not be found, we stop here and return a 404.
// func TodoCtx(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		var todoItem todo.Todo
// 		var err error

// 		if todoId := chi.URLParam(r, "todoId"); todoId != "" {
// 			todoItem, err = todo.FindTodo(todoId)
// 		} else {
// 			render.Render(w, r, handler.ErrRender(err))
// 			return
// 		}
// 		if err != nil {
// 			render.Render(w, r, handler.ErrRender(err))
// 			return
// 		}
// 		fmt.Println(todoItem)

// 		ctx := context.WithValue(r.Context(), "todo", todoItem)
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }

// func sayhelloName(w http.ResponseWriter, r *http.Request) {
// 	r.ParseForm()       // parse arguments, you have to call this by yourself
// 	fmt.Println(r.Form) // print form information in server side
// 	fmt.Println("path", r.URL.Path)
// 	fmt.Println("scheme", r.URL.Scheme)
// 	fmt.Println(r.Form["url_long"])
// 	for k, v := range r.Form {
// 		fmt.Println("key:", k)
// 		fmt.Println("val:", strings.Join(v, ""))
// 	}
// 	fmt.Fprintf(w, "Hello astaxie!") // send data to client side
// }
// func login(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("method:", r.Method) // get request method
// 	if r.Method == "GET" {
// 		crutime := time.Now().Unix()
// 		h := md5.New()
// 		io.WriteString(h, strconv.FormatInt(crutime, 10))
// 		token := fmt.Sprintf("%x", h.Sum(nil))
// 		fmt.Println("token serve sid length:", token)
// 		t, _ := template.ParseFiles("login.gtpl")
// 		t.Execute(w, token)
// 	} else {
// 		// log in request
// 		r.ParseForm()
// 		token := r.Form.Get("token")
// 		fmt.Println("token client sid length:", token)
// 		if token != "" {
// 			// check token validity
// 		} else {
// 			// give error if no token
// 		}
// 		fmt.Println("username length:", len(r.Form["username"][0]))
// 		fmt.Println("username:", template.HTMLEscapeString(r.Form.Get("username"))) // print in server side
// 		fmt.Println("password:", template.HTMLEscapeString(r.Form.Get("password")))
// 		template.HTMLEscape(w, []byte(r.Form.Get("username"))) // respond to client
// 	}
// }

// func main() {
// 	r := chi.NewRouter()

// 	// A good base middleware stack
// 	r.Use(middleware.RequestID)
// 	r.Use(middleware.RealIP)
// 	r.Use(middleware.Logger)
// 	r.Use(middleware.Recoverer)
// 	r.Use(render.SetContentType(render.ContentTypeJSON))
// 	r.Use(cors.New(cors.Options{
// 		AllowedOrigins:   []string{"https://*", "http://*"},
// 		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
// 		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
// 		AllowCredentials: true,
// 		Debug:            false,
// 	}).Handler)

// 	// Set a timeout value on the request context (ctx), that will signal
// 	// through ctx.Done() that the request has timed out and further
// 	// processing should be stopped.
// 	r.Use(middleware.Timeout(60 * time.Second))

// 	port := os.Getenv("PORT")
// 	if port == "" {
// 		port = defaultPort
// 	}

// 	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

// 	r.Handle("/", playground.Handler("GraphQL playground", "/query"))
// 	r.Handle("/query", srv)

// 	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
// 	log.Fatal(http.ListenAndServe(":"+port, r))
// }
