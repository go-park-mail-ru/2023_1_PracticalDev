package pins

import (
	"errors"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
)

type CreateParams struct {
	Title       string
	Description string
	MediaSource models.Image
	Author      int
}

var (
	ErrDb = errors.New("db error")
)

type Repository interface {
	Create(params *CreateParams) (models.Pin, error)
	Get(id int) (models.Pin, error)
	ListByUser(userId int, page, limit int) ([]models.Pin, error)
	ListByBoard(boardId int, page, limit int) ([]models.Pin, error)
	List(page, limit int) ([]models.Pin, error)
	Update(params *models.Pin) (models.Pin, error)
	Delete(id int) error

	AddToBoard(boardId, pinId int) error
	RemoveFromBoard(boardId, pinId int) error

	CheckWriteAccess(userId, pinId string) (bool, error)
	CheckReadAccess(userId, pinId string) (bool, error)
}
