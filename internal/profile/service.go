package profile

import (
	"errors"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
)

var (
	ErrTooShortUsername = errors.New("username must be at least 4 characters")
	ErrTooLongUsername  = errors.New("username must be no more than 30 characters")
	ErrEmptyName        = errors.New("name must not be empty")
	ErrTooLongName      = errors.New("name must be no more than 60 characters")
)

type Service interface {
	FullUpdate(params *FullUpdateParams, image *models.Image) (Profile, error)
	PartialUpdate(params *PartialUpdateParams) (Profile, error)
}
