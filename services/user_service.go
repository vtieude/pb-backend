package services

import (
	"context"
	"fmt"
	db_manager "pb-backend/db"
	"pb-backend/entities"
	"pb-backend/graph/model"

	"github.com/elgris/sqrl"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
)

type IUserService interface {
	GetAllUsers(ctx context.Context, pagination *model.Pagination) ([]*entities.User, error)
	CreateUser(ctx context.Context, input model.NewUser) (*entities.User, error)
}
type UserService struct {
	DB db_manager.IDb
}

// define provider
var NewUserService = wire.NewSet(wire.Struct(new(UserService), "*"), wire.Bind(new(IUserService), new(*UserService)))

func (u *UserService) CreateUser(ctx context.Context, input model.NewUser) (*entities.User, error) {
	var user []entities.User
	stss := sqrl.Select("username").From("user")
	u.DB.AddPagination(stss, nil)
	err := u.DB.Select(ctx, &user, stss)
	// if there is an error opening the connection, handle it
	if err != nil {
		fmt.Print(err.Error())
		return nil, err
	}
	return &user[0], nil
}
func (u *UserService) GetAllUsers(ctx context.Context, pagination *model.Pagination) ([]*entities.User, error) {
	var users []*entities.User
	stss := sqrl.Select("username, id").From("user")
	u.DB.AddPagination(stss, pagination)
	err := u.DB.Select(ctx, &users, stss)
	// if there is an error opening the connection, handle it
	if err != nil {
		fmt.Print(err.Error())
		return nil, err
	}
	return users, nil

}
