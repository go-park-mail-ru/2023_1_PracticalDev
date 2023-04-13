package pins

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	"github.com/pkg/errors"
)

type CreateParams struct {
	Title       string
	Description string
	MediaSource models.Image
	Author      int
}

type FullUpdateParams struct {
	Id          int
	Title       string
	Description string
}

var (
	ErrDb          = errors.New("db error")
	ErrPinNotFound = errors.New("pin not found")
)

type Repository interface {
	Create(params *CreateParams) (models.Pin, error)
	Get(id int) (models.Pin, error)
	ListByUser(userId int, page, limit int) ([]models.Pin, error)
	List(page, limit int) ([]models.Pin, error)
	FullUpdate(params *FullUpdateParams) (models.Pin, error)
	Delete(id int) error

	CheckWriteAccess(userId, pinId string) (bool, error)
	CheckReadAccess(userId, pinId string) (bool, error)
}
