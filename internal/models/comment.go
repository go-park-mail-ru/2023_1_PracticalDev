package models

import "time"

type Comment struct {
	ID        int       `json:"id"`
	AuthorID  int       `json:"author_id"`
	PinID     int       `json:"pin_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}
