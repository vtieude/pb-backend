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
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/elgris/sqrl"
	"github.com/google/wire"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type IMyCustomHttpHandler interface {
	Authorization(ctx context.Context) (entities.MyCustomClaims, error)
	LoggingHandler(next http.Handler) http.Handler
	AuthGraphql(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error)
	AdminValidate(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error)
}

type MyCustomHttpHandler struct {
	DB entities.DB
}

var ModifiesSet = wire.NewSet(wire.Struct(new(MyCustomHttpHandler), "*"), wire.Bind(new(IMyCustomHttpHandler), new(*MyCustomHttpHandler)))

func (errorHandler *MyCustomHttpHandler) LoggingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("panic occurred:", err)
				return
			}
		}()
		token := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", -1)

		ctx := context.WithValue(r.Context(), consts.SetCtxKey(consts.TOKEN_CTX_KEY), token)
		if token != "" {
			customClaim, err := services.ParseToken(token)
			if err != nil {
				resp := make(map[string]string)
				resp["message"] = "Unauthorized"
				w.WriteHeader(401)
				jsonResp, err := json.Marshal(resp)
				if err != nil {
					log.Fatalf("Error happened in JSON marshal. Err: %s", err)
				}
				w.Write(jsonResp)
				return
			}
			ctx = context.WithValue(ctx, consts.SetCtxKey(consts.USER_CTX_KEY), customClaim)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// I return nil for the sake of example.
func (e *MyCustomHttpHandler) Authorization(ctx context.Context) (entities.MyCustomClaims, error) {
	customClaim, errParse := consts.CtxClaimValue(ctx)
	if errParse != nil || customClaim.ID == 0 {
		return *customClaim, &model.MyError{Message: consts.ERR_USER_UN_AUTHORIZATION}
	}
	currentUser, err := entities.UserByID(ctx, e.DB, customClaim.ID)
	if err != nil {
		return *customClaim, &model.MyError{Message: consts.ERR_USER_UN_AUTHORIZATION}
	}
	if currentUser.Username == currentUser.Username {
		return *customClaim, nil
	}
	return *customClaim, &model.MyError{Message: consts.ERR_USER_UN_AUTHORIZATION}
}

// // I return nil for the sake of example.
// func (e *MyCustomHttpHandler) ServeHTTP() error {
// 	return nil
// }

func (handler *MyCustomHttpHandler) AuthGraphql(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	claims, err := handler.Authorization(ctx)
	if err != nil {
		return nil, err
	}
	if claims.ID == 0 {
		return nil, &gqlerror.Error{
			Message: "Access Denied",
		}
	}

	return next(ctx)
}
func (handler *MyCustomHttpHandler) AdminValidate(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	tokenData := consts.CtxValue(ctx)
	if tokenData == "" {
		return nil, &gqlerror.Error{
			Message: consts.ERR_USER_LOGIN_REQUIRED,
		}
	}
	claims, errParse := consts.CtxClaimValue(ctx)
	if errParse != nil || claims.ID == 0 {
		return claims, &model.MyError{Message: consts.ERR_USER_UN_AUTHORIZATION}
	}
	adminRole := [2]string{consts.ROLE_USER_ADMIN, consts.ROLE_USER_SUPER_ADMIN}
	userFilter := sqrl.Select("count(*)").From("user u").
		Where(sqrl.Eq{"u.id": claims.ID}).Where(sqrl.Eq{"u.role": adminRole})

	var roleCount = 0
	err := handler.DB.QueryRowContext(ctx, &roleCount, userFilter)
	if err != nil {
		return nil, &gqlerror.Error{
			Message: "Access Denied",
		}
	}
	if roleCount == 0 {
		return nil, &gqlerror.Error{
			Message: "Permission Denied",
		}
	}
	return next(ctx)
}
