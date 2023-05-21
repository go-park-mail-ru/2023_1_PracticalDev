package models

import "time"

//go:generate easyjson -all -snake_case like.go

type Like struct {
	PinId     int       `json:"pin_id"`
	AuthorId  int       `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
}
