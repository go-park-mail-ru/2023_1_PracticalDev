package http

import (
	"time"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
)

// API requests

// API responses
type listResponse struct {
	Chats []models.Chat `json:"items"`
}

func newListResponse(chats []models.Chat) *listResponse {
	return &listResponse{Chats: chats}
}

type messagesListResponse struct {
	Messages []models.Message `json:"items"`
}

func newMessagesListResponse(messages []models.Message) *messagesListResponse {
	return &messagesListResponse{Messages: messages}
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
