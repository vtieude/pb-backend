package services

import (
	"context"
	"database/sql"
	"fmt"
	"pb-backend/graph/model"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
)

type IUserService interface {
	GetAllUsers(ctx context.Context) ([]*model.User, error)
	CreateUser(ctx context.Context, input model.NewUser) (*model.User, error)
}
type UserService struct {
}

// define provider
var NewUserService = wire.NewSet(wire.Struct(new(UserService), "*"), wire.Bind(new(IUserService), new(*UserService)))

func (u *UserService) OpenConnectTion() {
	db, err := sql.Open("mysql", "user:password@/dbname?parseTime=true")
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
}
func (u *UserService) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	return &model.User{}, nil
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
		err = results.Scan(&user.ID, &user.Name)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		// and then print out the tag's Name attribute
		users = append(users, &user)
	}
	return users, nil

}
