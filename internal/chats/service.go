package chats

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
)

type Service interface {
	Create(params *CreateParams) (models.Chat, error)
	ListByUser(userID int) ([]models.Chat, error)
	Get(id int) (models.Chat, error)
}
