package handlerTodo

import (
	"errors"
	"net/http"
	"pb-backend/todo"

	"github.com/go-chi/render"
)

// $ curl -X DELETE http://localhost:3333/articles/1
// {"id":"1","title":"Hi"}
//
// $ curl http://localhost:3333/articles/1
// "Not Found"
//
// $ curl -X POST -d '{"id":"will-be-omitted","title":"awesomeness"}' http://localhost:3333/articles
// {"id":"97","title":"awesomeness"}

func ListTodo(w http.ResponseWriter, r *http.Request) {
	if err := render.RenderList(w, r, NewTodoListResponse()); err != nil {
		render.Render(w, r, nil)
		return
	}
}

func NewTodoListResponse() []render.Renderer {
	todoList := todo.Get()
	list := []render.Renderer{}
	for _, item := range todoList {
		list = append(list, NewTodoResponse(&item))
	}
	return list
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	item := r.Context().Value("todo").(todo.Todo)

	_ = todo.Complete(item.ID)

	render.Status(r, http.StatusCreated)
	render.RenderList(w, r, NewTodoListResponse())
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	data := &TodoRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	article := data.Message
	_ = todo.Add(*article)

	render.Status(r, http.StatusCreated)
	render.RenderList(w, r, NewTodoListResponse())
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	item := r.Context().Value("todo").(todo.Todo)

	err := todo.Delete(item.ID)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Render(w, r, NewTodoResponse(&item))
}

func NewTodoResponse(article *todo.Todo) *TodoResponse {
	resp := &TodoResponse{Todo: *article}

	return resp
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

type TodoResponse struct {
	todo.Todo

	// We add an additional field to the response here.. such as this
	// elapsed computed property
	Elapsed int64 `json:"elapsed"`
}

func (rd *TodoResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	rd.Elapsed = 10
	return nil
}

type TodoRequest struct {
	Message *string
	Id      *string
}

func (a *TodoRequest) Bind(r *http.Request) error {
	// a.Article is nil if no Article fields are sent in the request. Return an
	// error to avoid a nil pointer dereference.
	if a.Message == nil {
		return errors.New("missing required todo fields.")
	}

	// a.User is nil if no Userpayload fields are sent in the request. In this app
	// this won't cause a panic, but checks in this Bind method may be required if
	// a.User or futher nested fields like a.User.Name are accessed elsewhere.

	// just a post-process after a decode..
	// a.ProtectedID = "" // unset the protected ID
	// a.Article.Title = strings.ToLower(a.Article.Title) // as an example, we down-case
	return nil
}

func newStringResponse(message string) render.Renderer {
	return &StringResponse{Message: message}
}

type StringResponse struct {
	Message string
}

func (rd *StringResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}
