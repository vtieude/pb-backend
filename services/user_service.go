package services

import (
	"context"
	"fmt"
	"pb-backend/consts"
	"pb-backend/entities"
	"pb-backend/graph/model"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/elgris/sqrl"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
	"golang.org/x/crypto/bcrypt"
)

const keyHassPwd = "wilson-pb-app"

type IUserService interface {
	GetAllUsers(ctx context.Context, pagination *model.Pagination) ([]*entities.User, error)
	CreateUser(ctx context.Context, input model.NewUser) (*entities.User, error)
	Login(ctx context.Context, email string, password string) (string, error)
}
type UserService struct {
	DB entities.DB
}

// define provider
var NewUserService = wire.NewSet(wire.Struct(new(UserService), "*"), wire.Bind(new(IUserService), new(*UserService)))

func (u *UserService) CreateUser(ctx context.Context, input model.NewUser) (*entities.User, error) {
	var existUsers entities.User
	stss := sqrl.Select("email").From("user").Where(sqrl.Eq{"email": input.Name})
	err := u.DB.QueryRowContext(ctx, &existUsers, stss)
	if err != nil {
		return nil, err
	}
	if existUsers.ID != 0 {
		return nil, &model.MyError{"Email already exist"}
	}
	newUsers := entities.User{
		Username: input.Name,
	}
	err = newUsers.Insert(ctx, u.DB)
	// if there is an error opening the connection, handle it
	if err != nil {
		fmt.Print(err.Error())
		return nil, err
	}
	return &newUsers, nil
}
func (u *UserService) GetAllUsers(ctx context.Context, pagination *model.Pagination) ([]*entities.User, error) {
	var users []*entities.User
	stss := sqrl.Select("username, id").From("user")
	u.DB.AddPagination(stss, pagination)
	err := u.DB.QueryContext(ctx, &users, stss)
	// if there is an error opening the connection, handle it
	if err != nil {
		fmt.Print(err.Error())
		return nil, err
	}
	return users, nil

}

func (u *UserService) Login(ctx context.Context, email string, password string) (string, error) {
	password = strings.TrimSpace(password)
	email = strings.TrimSpace(email)
	if len(password) == 0 || len(email) == 0 {
		return "", &model.MyError{Message: "Invalid email or password "}
	}
	findUser, err := entities.UserByEmail(ctx, u.DB, email)
	if err != nil {
		return "", &model.MyError{Message: "User not found, " + string(err.Error())}
	}
	if findUser.ID == 0 {
		return "", &model.MyError{Message: consts.ERR_USER_NOT_FOUND}
	}
	if !u.checkPasswordHash(password, findUser.Password) {
		return "", &model.MyError{Message: consts.ERR_USER_NOT_FOUND}
	}
	return u.GenerateToken(*findUser)
}

func (u *UserService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password+keyHassPwd), bcrypt.DefaultCost)
	return string(bytes), err
}

func (u *UserService) checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password+keyHassPwd))
	return err == nil
}

// ParseToken parses a jwt token and returns the username in it's claims
func (u *UserService) GenerateToken(userLogin entities.User) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, entities.MyCustomClaims{
		Username: userLogin.Username,
		ID:       userLogin.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().AddDate(0, 3, 0).Unix(),
			// ExpiresAt: time.Now().Add(time.Second * 10).Unix(),
		},
	})
	return claims.SignedString([]byte(keyHassPwd))
}

// ParseToken parses a jwt token and returns the username in it's claims
func ParseToken(tokenStr string) (entities.MyCustomClaims, error) {
	dataClaims, err := jwt.ParseWithClaims(tokenStr, &entities.MyCustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there's a problem with the signing method")
		}
		return []byte(keyHassPwd), nil
	})
	if err != nil || dataClaims == nil {
		return entities.MyCustomClaims{}, err
	}
	result := dataClaims.Claims.(*entities.MyCustomClaims)
	return *result, err
}
