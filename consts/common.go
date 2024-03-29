package consts

import (
	"context"
	"pb-backend/entities"
	"pb-backend/graph/model"
)

const ERR_USER_NOT_FOUND = "Tài khoản hoặc mật khẩu không chính xác"
const ERR_USER_UN_AUTHORIZATION = "Tài khoản đã hết hạn, vui lòng đăng nhập lại"
const ERR_USER_INVALID_EMAIL = "Email không hợp lệ"
const ERR_USER_INVALID_PASSWORD = "Mật khẩu phải có ít nhất 6 kí tự "
const ERR_USER_DUPPLICATE_EMAIL_ADDRESS = "Email đã tồn tại"
const ERR_USER_LOGIN_REQUIRED = "Bạn phải đăng nhập"
const ERR_DUPLICATE_PRODUCT_KEY = "Mã sản phẩm đã tồn tại"
const ERR_USER_GET_INFORMATION = "Thông tin tài khoản đăng nhập không chính xác"
const ERR_USER_NOT_EXIST = "Tài khoản không tồn tại"
const ERR_USER_INVALID_PERMISSION = "Tài khoản của bạn không thể thực hiện thay đổi này"
const USER_CTX_KEY = "user_context"
const TOKEN_CTX_KEY = "token"
const ROLE_USER_ADMIN = "admin"
const ROLE_USER_SUPER_ADMIN = "super_admin"
const ROLE_STAFF_USER = "staff"

type authString string

func SetCtxKey(key string) authString {
	return authString(key)
}

func CtxValue(ctx context.Context) string {
	raw, err := ctx.Value(authString(TOKEN_CTX_KEY)).(string)
	if !err {
		return ""
	}
	return raw
}

func CtxClaimValue(ctx context.Context) (*entities.MyCustomClaims, error) {
	raw, err := ctx.Value(authString(USER_CTX_KEY)).(entities.MyCustomClaims)
	if !err {
		return nil, &model.MyError{Message: ERR_USER_GET_INFORMATION}
	}
	return &raw, nil
}
