package models

import "time"

//go:generate easyjson -all -snake_case notification.go

type Notification struct {
	ID        int         `json:"id"`
	UserID    int         `json:"user_id"`
	IsRead    bool        `json:"is_read"`
	CreatedAt time.Time   `json:"created_at"`
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
}

type NewPinNotification struct {
	PinID int `json:"pin_id"`
}

type NewLikeNotification struct {
	PinID    int `json:"pin_id"`
	AuthorID int `json:"author_id"`
}

type NewCommentNotification struct {
	PinID    int    `json:"pin_id"`
	AuthorID int    `json:"author_id"`
	Text     string `json:"text"`
}

type NewFollowerNotification struct {
	FollowerID int `json:"follower_id"`
}
