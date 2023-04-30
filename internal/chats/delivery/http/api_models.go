package http

import (
	"time"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
)

// API requests
type createRequest struct {
	UserID int `json:"user_id"`
}

// API responses
type createResponse struct {
	ID        int       `json:"id"`
	User1ID   int       `json:"user1_id"`
	User2ID   int       `json:"user2_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func newCreateResponse(chat *models.Chat) *createResponse {
	return &createResponse{
		ID:        chat.ID,
		User1ID:   chat.User1ID,
		User2ID:   chat.User2ID,
		CreatedAt: chat.CreatedAt,
		UpdatedAt: chat.UpdatedAt,
	}
}

type listResponse struct {
	Chats []models.Chat `json:"items"`
}

func newListResponse(chats []models.Chat) *listResponse {
	return &listResponse{Chats: chats}
}

type getResponse struct {
	ID        int       `json:"id"`
	User1ID   int       `json:"user1_id"`
	User2ID   int       `json:"user2_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func newGetResponse(chat *models.Chat) *getResponse {
	return &getResponse{
		ID:        chat.ID,
		User1ID:   chat.User1ID,
		User2ID:   chat.User2ID,
		CreatedAt: chat.CreatedAt,
		UpdatedAt: chat.UpdatedAt,
	}
}
