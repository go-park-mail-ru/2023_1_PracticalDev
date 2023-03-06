package models

type Session struct {
	UserId    int    `json:"id"`
	UserEmail string `json:"email"`
}
