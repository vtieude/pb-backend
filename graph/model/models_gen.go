// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"time"
)

type EditUserModel struct {
	UserID      int     `json:"userId"`
	UserName    string  `json:"userName"`
	RoleName    string  `json:"roleName"`
	PhoneNumber *string `json:"phoneNumber"`
	Password    *string `json:"password"`
}

type NewUser struct {
	UserName    string  `json:"userName"`
	Email       string  `json:"email"`
	Password    string  `json:"password"`
	RoleName    string  `json:"roleName"`
	PhoneNumber *string `json:"phoneNumber"`
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
	ID           int     `json:"id"`
	Name         *string `json:"name"`
	ProductKey   string  `json:"productKey"`
	Category     *string `json:"category"`
	Price        float64 `json:"price"`
	SellingPrice float64 `json:"sellingPrice"`
	Number       int     `json:"number"`
	Description  *string `json:"description"`
	ImageURL     *string `json:"imageUrl"`
}

type ProductInputModel struct {
	ID           *int    `json:"id"`
	Name         string  `json:"name"`
	Key          string  `json:"key"`
	Category     *string `json:"category"`
	Price        float64 `json:"price"`
	SellingPrice float64 `json:"sellingPrice"`
	Number       int     `json:"number"`
	Description  *string `json:"description"`
	ImageBase64  *string `json:"imageBase64"`
	ImagePrefix  *string `json:"imagePrefix"`
}

type ProfileImage struct {
	FileName   *string `json:"fileName"`
	FileBase64 *string `json:"fileBase64"`
}

type UserDto struct {
	ID       int    `json:"id"`
	Token    string `json:"token"`
	Role     string `json:"role"`
	UserName string `json:"userName"`
}
