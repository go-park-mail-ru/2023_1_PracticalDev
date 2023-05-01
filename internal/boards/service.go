package boards

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
)

type Service interface {
	Create(params *CreateParams) (models.Board, error)
	List(userId int) ([]models.Board, error)
	Get(id int) (models.Board, error)
	FullUpdate(params *FullUpdateParams) (models.Board, error)
	PartialUpdate(params *PartialUpdateParams) (models.Board, error)
	Delete(id int) error

	AddPin(boardId, pinId int) error
	PinsList(userId, boardId int, page, limit int) ([]models.Pin, error)
	RemovePin(boardId, pinId int) error

	CheckWriteAccess(userId, boardId string) (bool, error)
	CheckReadAccess(userId, boardId string) (bool, error)
}
