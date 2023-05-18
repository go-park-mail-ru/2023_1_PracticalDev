package pins

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
)

type Service interface {
	Create(params *CreateParams) (models.Comment, error)
	List(pinID int) ([]models.Comment, error)
}
