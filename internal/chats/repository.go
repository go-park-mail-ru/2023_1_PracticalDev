package chats

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
)

type CreateParams struct {
	User1ID int
	User2ID int
}

type SendMessageParams struct {
	AuthorID int
	ChatID   int
	Text     string
}

type Repository interface {
	Create(params *CreateParams) (models.Chat, error)
	ListByUser(userId int) ([]models.Chat, error)
	Get(id int) (models.Chat, error)

	SendMessage(params *SendMessageParams) (*models.Message, error)

	ChatExists(user1ID, user2ID int) (bool, error)
	GetByUsers(user1ID, user2ID int) (models.Chat, error)
}
