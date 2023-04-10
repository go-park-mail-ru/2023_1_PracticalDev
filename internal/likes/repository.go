package likes

import (
	"errors"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
)

var (
	ErrDb                = errors.New("db error")
	ErrLikeNotFound      = errors.New("like not found")
	ErrLikeAlreadyExists = errors.New("like already exists")
)

type Repository interface {
	Create(pinId, authorId int) error
	Delete(pinId, authorId int) error

	ListByAuthor(authorId int) ([]models.Like, error)
	ListByPin(pinId int) ([]models.Like, error)
}
