package modifies

import (
	"context"
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

const tokenKey = "token"

type IMyCustomHttpHandler interface {
	Authorization(ctx context.Context, token string) (entities.MyCustomClaims, error)
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

		ctx := context.WithValue(r.Context(), authString(tokenKey), token)
		if token != "" {
			log.Println("user token call")
			customClaim, _ := services.ParseToken(token)
			ctx = context.WithValue(ctx, authString(consts.USER_CTX_KEY), customClaim)
		}

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

func (handler *MyCustomHttpHandler) AuthGraphql(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	tokenData := CtxValue(ctx)
	claims, err := handler.Authorization(ctx, tokenData)
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
	tokenData := CtxValue(ctx)
	if tokenData == "" {
		return nil, &gqlerror.Error{
			Message: consts.ERR_USER_LOGIN_REQUIRED,
		}
	}
	claims, err := handler.Authorization(ctx, tokenData)
	if err != nil {
		return nil, err
	}
	if claims.ID == 0 {
		return nil, &gqlerror.Error{
			Message: "Access Denied",
		}
	}
	userFilter := sqrl.Select("count(*)").From("user u")
	userFilter.Join("user_role ur on ur.fk_user = u.id")
	userFilter.Join("role r on ur.fk_role = r.id").Where(sqrl.Eq{"u.id": claims.ID})
	userFilter.Where(sqrl.Eq{"r.role_name": "admin"})

	var roleCount = 0
	err = handler.DB.QueryRowContext(ctx, &roleCount, userFilter)
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
func CtxValue(ctx context.Context) string {
	raw, err := ctx.Value(authString(tokenKey)).(string)
	if !err {
		return ""
	}
	return raw
}
