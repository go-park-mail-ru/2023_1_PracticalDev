package pins

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
)

type Service interface {
	Create(params *CreateParams) (models.Pin, error)
	Get(id int) (models.Pin, error)
	ListByUser(userId int, page, limit int) ([]models.Pin, error)
	ListByBoard(boardId int, page, limit int) ([]models.Pin, error)
	List(page, limit int) ([]models.Pin, error)
	FullUpdate(params *FullUpdateParams) (models.Pin, error)
	Delete(id int) error

	AddToBoard(boardId, pinId int) error
	RemoveFromBoard(boardId, pinId int) error

	CheckWriteAccess(userId, pinId string) (bool, error)
	CheckReadAccess(userId, pinId string) (bool, error)
}
