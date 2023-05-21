package http

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/xss"
	"time"
)

//go:generate easyjson -all -snake_case api_models.go

// API requests
type createRequest struct {
	Text string `json:"text"`
}

// API responses
type createResponse struct {
	ID        int       `json:"id"`
	AuthorID  int       `json:"author_id"`
	PinID     int       `json:"pin_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

func newCreateResponse(comment *models.Comment) *createResponse {
	return &createResponse{
		ID:        comment.ID,
		AuthorID:  comment.AuthorID,
		PinID:     comment.PinID,
		Text:      xss.Sanitize(comment.Text),
		CreatedAt: comment.CreatedAt,
	}
}

type listResponse struct {
	Items []models.Comment `json:"items"`
}

func newListResponse(comments []models.Comment) *listResponse {
	for i := range comments {
		comments[i].Text = xss.Sanitize(comments[i].Text)
	}

	return &listResponse{
		Items: comments,
	}
}
