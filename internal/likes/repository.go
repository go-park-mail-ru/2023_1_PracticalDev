package likes

import (
	"errors"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
)

var (
	ErrDb           = errors.New("db error")
	ErrPinNotFound  = errors.New("pin not found")
	ErrUserNotFound = errors.New("user not found")
)

type Repository interface {
	Create(pinId, authorId int) error
	Delete(pinId, authorId int) error

	ListByAuthor(authorId int) ([]models.Like, error)
	ListByPin(pinId int) ([]models.Like, error)
}
