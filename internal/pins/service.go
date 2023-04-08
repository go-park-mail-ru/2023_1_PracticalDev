package pins

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
)

type Service interface {
	CreatePin(params *models.Pin, image *models.Image) (models.Pin, error)
	GetPin(id int) (models.Pin, error)
	GetPinsByUser(userId int, page, limit int) ([]models.Pin, error)
	GetPinsByBoard(boardId int, page, limit int) ([]models.Pin, error)
	GetPins(page, limit int) ([]models.Pin, error)
	UpdatePin(params *models.Pin) (models.Pin, error)
	DeletePin(id int) error

	AddPinToBoard(boardId, pinId int) error
	RemovePinFromBoard(boardId, pinId int) error

	CheckWriteAccess(userId, pinId string) (bool, error)
	CheckReadAccess(userId, pinId string) (bool, error)
}
