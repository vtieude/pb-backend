package services

import (
	"context"
	"fmt"
	"net/mail"
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
	Login(ctx context.Context, email string, password string) (*model.UserDto, error)
}
type UserService struct {
	DB entities.DB
}

// define provider
var NewUserService = wire.NewSet(wire.Struct(new(UserService), "*"), wire.Bind(new(IUserService), new(*UserService)))

func (u *UserService) CreateUser(ctx context.Context, input model.NewUser) (*entities.User, error) {
	var existUsers entities.User
	if !u.validEmail(input.Name) {
		return nil, &model.MyError{Message: consts.ERR_USER_INVALID_EMAIL_PASSWORD}
	}
	stss := sqrl.Select("email").From("user").Where(sqrl.Eq{"email": input.Name})
	err := u.DB.QueryRowContext(ctx, &existUsers, stss)
	if err != nil {
		return nil, err
	}
	if existUsers.ID != 0 {
		return nil, &model.MyError{Message: consts.ERR_USER_DUPPLICATE_EMAIL_ADDRESS}
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

func (u *UserService) Login(ctx context.Context, email string, password string) (*model.UserDto, error) {
	password = strings.TrimSpace(password)
	email = strings.TrimSpace(email)
	if len(password) == 0 || len(email) == 0 || !u.validEmail(email) {
		return nil, &model.MyError{Message: consts.ERR_USER_INVALID_EMAIL_PASSWORD}
	}
	findUser, err := entities.UserByEmail(ctx, u.DB, email)
	if err != nil {
		return nil, &model.MyError{Message: consts.ERR_USER_NOT_FOUND + string(err.Error())}
	}
	if findUser.ID == 0 {
		return nil, &model.MyError{Message: consts.ERR_USER_NOT_FOUND}
	}
	if !u.checkPasswordHash(password, findUser.Password) {
		return nil, &model.MyError{Message: consts.ERR_USER_NOT_FOUND}
	}
	userResult := model.UserDto{}
	userRoleFiter := sqrl.Select("role.role_name").From("user_role")
	userRoleFiter.Join("role on role.id = user_role.fk_role")
	userRoleFiter.Where(sqrl.Eq{"user_role.fk_user": findUser.ID})
	var roleName []string
	err = u.DB.QueryContext(ctx, &roleName, userRoleFiter)
	if err == nil && len(roleName) > 0 {
		userResult.Role = roleName[0]
	}
	userResult.ID = findUser.ID
	userResult.UserName = findUser.Username
	userResult.Token, err = u.GenerateToken(*findUser)
	return &userResult, err
}

func (u *UserService) validEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
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
