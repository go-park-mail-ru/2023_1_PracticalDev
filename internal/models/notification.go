package models

import "time"

type Notification struct {
	ID        int         `json:"id"`
	UserID    int         `json:"user_id"`
	Message   string      `json:"message"`
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
	CommentID int `json:"comment_id"`
}
