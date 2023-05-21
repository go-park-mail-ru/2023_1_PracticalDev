package models

import "time"

//go:generate easyjson -all -snake_case message.go

type Message struct {
	ID        int       `json:"id"`
	AuthorID  int       `json:"author_id"`
	ChatID    int       `json:"chat_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}
