package services

import (
	"context"
	"database/sql"
	"fmt"
	db_manager "pb-backend/db"
	"pb-backend/graph/model"

	"github.com/elgris/sqrl"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
)

type IUserService interface {
	GetAllUsers(ctx context.Context) ([]*model.User, error)
	CreateUser(ctx context.Context, input model.NewUser) (*model.User, error)
}
type UserService struct {
	DB db_manager.IDb
}

// define provider
var NewUserService = wire.NewSet(wire.Struct(new(UserService), "*"), wire.Bind(new(IUserService), new(*UserService)))

func (u *UserService) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	var user []model.User
	sqlrFilter := sqrl.Expr("Select username from user where id = ?", 1)
	err := u.DB.Select(ctx, &user, sqlrFilter)
	// if there is an error opening the connection, handle it
	if err != nil {
		fmt.Print(err.Error())
		return nil, err
	}
	return &user[0], nil
}
func (u *UserService) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	var users []*model.User
	// Open up our database connection.
	db, err := sql.Open("mysql", "root:qweqwe@tcp(127.0.0.1:3307)/app_db?parseTime=true")

	// if there is an error opening the connection, handle it
	if err != nil {
		fmt.Print(err.Error())
		my := model.MyError{Message: "cannot connect db"}
		return nil, my.ReturnError()
	}
	defer db.Close()
	// Execute the query
	results, err := db.Query("SELECT id, username FROM user")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	for results.Next() {
		var user model.User
		// for each row, scan the result into our tag composite object
		err = results.Scan(&user.ID, &user.Username)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		// and then print out the tag's Name attribute
		users = append(users, &user)
	}
	return users, nil

}
