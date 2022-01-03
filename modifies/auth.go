package modifies

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"pb-backend/consts"
	"pb-backend/entities"
	"pb-backend/graph/model"
	"pb-backend/services"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/wire"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type IMyCustomHttpHandler interface {
	Authorization(ctx context.Context, token string) (entities.MyCustomClaims, error)
	LoggingHandler(next http.Handler) http.Handler
}

type MyCustomHttpHandler struct {
	DB entities.DB
}

var ModifiesSet = wire.NewSet(wire.Struct(new(MyCustomHttpHandler), "*"), wire.Bind(new(IMyCustomHttpHandler), new(*MyCustomHttpHandler)))

func (errorHandler *MyCustomHttpHandler) LoggingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			next.ServeHTTP(w, r)
			return
		}
		bearer := "Bearer "
		if len(auth) <= len(bearer) {
			next.ServeHTTP(w, r)
			return
		}
		defer func() {
			if err := recover(); err != nil {
				log.Println("panic occurred:", err)
				return
			}
		}()
		auth = auth[len(bearer):]
		claims, err := errorHandler.Authorization(r.Context(), auth)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			resp := make(map[string]string)
			resp["message"] = err.Error()
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			}
			w.Write(jsonResp)
			return
		}
		ctx := context.WithValue(r.Context(), authString("auth"), claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// I return nil for the sake of example.
func (e *MyCustomHttpHandler) Authorization(ctx context.Context, token string) (entities.MyCustomClaims, error) {
	customClaim, err := services.ParseToken(token)
	if err != nil || customClaim.ID == 0 {
		return customClaim, &model.MyError{Message: consts.ERR_USER_UN_AUTHORIZATION}
	}
	currentUser, err := entities.UserByID(ctx, e.DB, customClaim.ID)
	if err != nil {
		return customClaim, &model.MyError{Message: consts.ERR_USER_UN_AUTHORIZATION}
	}
	if currentUser.Username == currentUser.Username {
		return customClaim, nil
	}
	return customClaim, &model.MyError{Message: consts.ERR_USER_UN_AUTHORIZATION}
}

// // I return nil for the sake of example.
// func (e *MyCustomHttpHandler) ServeHTTP() error {
// 	return nil
// }

type authString string

func Auth(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	tokenData := CtxValue(ctx)
	if tokenData == nil || tokenData.ID == 0 {
		return nil, &gqlerror.Error{
			Message: "Access Denied",
		}
	}

	return next(ctx)
}
func CtxValue(ctx context.Context) *entities.MyCustomClaims {
	raw, err := ctx.Value(authString("auth")).(entities.MyCustomClaims)
	if !err {
		return nil
	}
	return &raw
}
