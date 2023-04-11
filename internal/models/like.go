package models

import "time"

type Like struct {
	PinId     int       `json:"pin_id"`
	AuthorId  int       `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
}
