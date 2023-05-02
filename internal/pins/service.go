package pins

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
)

type Service interface {
	Create(params *CreateParams) (models.Pin, error)
	Get(id, userId int) (models.Pin, error)
	ListByAuthor(authorId, userId, page, limit int) ([]models.Pin, error)
	List(userId, page, limit int) ([]models.Pin, error)
	FullUpdate(params *FullUpdateParams) (models.Pin, error)
	Delete(id int) error

	SetLikedField(pin *models.Pin, userId int) error

	CheckWriteAccess(userId, pinId string) (bool, error)
	CheckReadAccess(userId, pinId string) (bool, error)
}
