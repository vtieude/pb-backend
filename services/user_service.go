package services

import (
	"context"
	"fmt"
	"net/mail"
	"pb-backend/consts"
	"pb-backend/entities"
	"pb-backend/graph/model"
	"pb-backend/helper"
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
	EditUser(ctx context.Context, input model.EditUserModel) (*entities.User, error)
	DeleteUser(ctx context.Context, userId int) (bool, error)
	Login(ctx context.Context, email string, password string) (*model.UserDto, error)
	Me(ctx context.Context) (*model.UserDto, error)
}
type UserService struct {
	DB entities.DB
}

// define provider
var NewUserService = wire.NewSet(wire.Struct(new(UserService), "*"), wire.Bind(new(IUserService), new(*UserService)))

func (u *UserService) DeleteUser(ctx context.Context, userId int) (bool, error) {
	editUser, err := entities.UserByID(ctx, u.DB, userId)
	if err != nil {
		return false, err
	}
	if editUser == nil {
		return false, &model.MyError{Message: consts.ERR_USER_NOT_EXIST}
	}
	claims, _ := consts.CtxClaimValue(ctx)
	if claims.ID != userId {
		if !u.validUserPermissionAction(ctx, editUser.Permission) {
			return false, &model.MyError{Message: consts.ERR_USER_INVALID_PERMISSION}
		}
	}
	err = editUser.Delete(ctx, u.DB)
	return true, err
}

func (u *UserService) EditUser(ctx context.Context, input model.EditUserModel) (*entities.User, error) {

	editUser, err := entities.UserByID(ctx, u.DB, input.UserID)
	if err != nil {
		return nil, err
	}
	if editUser == nil {
		return nil, &model.MyError{Message: consts.ERR_USER_NOT_EXIST}
	}
	claims, _ := consts.CtxClaimValue(ctx)
	if claims.ID != input.UserID {
		if !u.validUserPermissionAction(ctx, editUser.Permission) {
			return nil, &model.MyError{Message: consts.ERR_USER_INVALID_PERMISSION}
		}
	}

	u.validUserPermissionAction(ctx, editUser.Permission)
	if input.RoleName == "" {
		input.RoleName = "user"
	}
	editUser.Username = input.UserName
	editUser.PhoneNumber = helper.ConvertToNullPointSqlString(input.PhoneNumber)
	editUser.Role = input.RoleName
	editUser.RoleLabel = helper.GetRoleLabelByRole(input.RoleName)
	editUser.Permission = helper.GetPermissionByRole(input.RoleName)

	err = editUser.Update(ctx, u.DB)
	// if there is an error opening the connection, handle it
	if err != nil {
		return nil, err
	}
	return editUser, nil
}

func (u *UserService) CreateUser(ctx context.Context, input model.NewUser) (*entities.User, error) {
	if !u.validEmail(input.Email) || len(strings.TrimSpace(input.Password)) < 6 {
		return nil, &model.MyError{Message: consts.ERR_USER_INVALID_EMAIL_PASSWORD}
	}
	email := strings.TrimSpace(input.Email)
	stss := sqrl.Select("count(*)").From("user").Where(sqrl.Eq{"email": email})
	var existUsers int
	err := u.DB.QueryRowContext(ctx, &existUsers, stss)
	if err != nil {
		return nil, err
	}
	if existUsers != 0 {
		return nil, &model.MyError{Message: consts.ERR_USER_DUPPLICATE_EMAIL_ADDRESS}
	}
	hsPwd, err := u.hashPassword(input.Password)
	if err != nil {
		return nil, err
	}
	if input.RoleName == "" {
		input.RoleName = "user"
	}

	newUsers := entities.User{
		Username:    input.UserName,
		Email:       input.Email,
		Password:    hsPwd,
		Role:        input.RoleName,
		RoleLabel:   helper.GetRoleLabelByRole(input.RoleName),
		Permission:  helper.GetPermissionByRole(input.RoleName),
		PhoneNumber: helper.ConvertToNullPointSqlString(input.PhoneNumber),
		Active:      true,
	}
	err = newUsers.Insert(ctx, u.DB)
	// if there is an error opening the connection, handle it
	if err != nil {
		return nil, err
	}
	return &newUsers, nil
}

func (u *UserService) Me(ctx context.Context) (*model.UserDto, error) {
	claims, errParse := consts.CtxClaimValue(ctx)
	if errParse != nil {
		return nil, &model.MyError{Message: consts.ERR_USER_LOGIN_REQUIRED}
	}
	userRoleFilter := sqrl.Select("u.id, u.email, u.password, u.username, u.role as rolename").From("user u")
	userRoleFilter.Where(sqrl.Eq{"u.id": claims.ID})
	var findUsers []model.UserRoleDto
	err := u.DB.QueryContext(ctx, &findUsers, userRoleFilter)
	if err != nil {
		return nil, &model.MyError{Message: consts.ERR_USER_GET_INFORMATION + string(err.Error())}
	}
	if len(findUsers) == 0 {
		return nil, &model.MyError{Message: consts.ERR_USER_GET_INFORMATION}
	}
	userLogin := findUsers[0]
	userResult := model.UserDto{}
	userResult.Role = userLogin.RoleName
	userResult.ID = userLogin.ID
	userResult.UserName = userLogin.UserName
	return &userResult, nil
}

func (u *UserService) GetAllUsers(ctx context.Context, pagination *model.Pagination) ([]*entities.User, error) {
	claims, errParse := consts.CtxClaimValue(ctx)
	if errParse != nil {
		return nil, &model.MyError{Message: consts.ERR_USER_LOGIN_REQUIRED}
	}
	currentUse, err := entities.UserByID(ctx, u.DB, claims.ID)
	if err != nil {
		return nil, err
	}
	if currentUse == nil {
		return nil, nil

	}
	var users []*entities.User
	stss := sqrl.Select("id, username, email, role, role_label, active, phone_number").From("user").Where(sqrl.LtOrEq{"permission": currentUse.Permission})
	u.DB.AddPagination(stss, pagination)
	err = u.DB.QueryContext(ctx, &users, stss)
	// if there is an error opening the connection, handle it
	if err != nil {
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
	userRoleFilter := sqrl.Select("u.id, u.email, u.password, u.username, u.role as rolename").From("user u")
	userRoleFilter.Where(sqrl.Eq{"u.email": email})
	var findUsers []model.UserRoleDto
	err := u.DB.QueryContext(ctx, &findUsers, userRoleFilter)
	if err != nil {
		return nil, &model.MyError{Message: consts.ERR_USER_NOT_FOUND}
	}
	if len(findUsers) == 0 {
		return nil, &model.MyError{Message: consts.ERR_USER_NOT_FOUND}
	}
	userLogin := findUsers[0]
	if !u.checkPasswordHash(password, userLogin.Password) {
		return nil, &model.MyError{Message: consts.ERR_USER_NOT_FOUND}
	}
	userResult := model.UserDto{}
	userResult.Role = userLogin.RoleName
	userResult.ID = userLogin.ID
	userResult.UserName = userLogin.UserName
	userResult.Token, err = u.GenerateToken(userLogin)
	return &userResult, err
}

func (u *UserService) validEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func (u *UserService) validUserPermissionAction(ctx context.Context, changePermission int) bool {
	claims, errParse := consts.CtxClaimValue(ctx)
	if errParse != nil {
		return false
	}
	userLogin, err := entities.UserByID(ctx, u.DB, claims.ID)
	if err != nil {
		return false
	}
	if userLogin.Role == consts.ROLE_USER_SUPER_ADMIN {
		return true
	}
	return userLogin.Permission > changePermission
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
func (u *UserService) GenerateToken(userLogin model.UserRoleDto) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, entities.MyCustomClaims{
		Email: userLogin.Email,
		ID:    userLogin.ID,
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
