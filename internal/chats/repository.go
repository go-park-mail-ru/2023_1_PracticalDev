package chats

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
)

type CreateParams struct {
	User1ID int
	User2ID int
}

type Repository interface {
	Create(params *CreateParams) (models.Chat, error)
	ListByUser(userId int) ([]models.Chat, error)
	Get(id int) (models.Chat, error)
}
