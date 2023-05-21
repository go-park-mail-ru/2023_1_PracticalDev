package models

//go:generate easyjson -all -snake_case board.go

type Board struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Privacy     string `json:"privacy"`
	UserId      int    `json:"user_id"`
}
