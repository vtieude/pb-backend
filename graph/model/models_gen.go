// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"time"
)

type Friend struct {
	City *string `json:"City"`
}

type NewProduct struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

type NewUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	RoleName string `json:"roleName"`
}

type OverviewUserSaleFilter struct {
	UserName *string    `json:"UserName"`
	DateTime *time.Time `json:"DateTime"`
}

type Pagination struct {
	PerPage *int     `json:"PerPage"`
	Page    *int     `json:"Page"`
	Sort    []string `json:"Sort"`
}

type ProductDto struct {
	ID         int     `json:"id"`
	Name       *string `json:"Name"`
	ProductKey string  `json:"ProductKey"`
}

type UserDto struct {
	ID       int    `json:"id"`
	Token    string `json:"token"`
	Role     string `json:"role"`
	UserName string `json:"userName"`
}
