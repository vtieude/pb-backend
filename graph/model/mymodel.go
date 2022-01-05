package model

type MyError struct {
	Message string
}

func (m *MyError) Error() string {
	return m.Message
}

func (m *MyError) ReturnError() error {
	return m
}

type UserRoleDto struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	UserName string `json:"username"`
	RoleName string `json:"rolename"`
}
