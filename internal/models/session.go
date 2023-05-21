package models

//go:generate easyjson -all -snake_case session.go

type Session struct {
	UserId    int    `json:"id"`
	UserEmail string `json:"email"`
}
