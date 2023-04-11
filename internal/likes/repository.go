package likes

import (
	"errors"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
)

var (
	ErrDb                = errors.New("db error")
	ErrPinNotFound       = errors.New("no such pin")
	ErrAuthorNotFound    = errors.New("no such author")
	ErrLikeNotFound      = errors.New("no such like")
	ErrLikeAlreadyExists = errors.New("like already exists")
)

type Repository interface {
	Create(pinId, authorId int) error
	Delete(pinId, authorId int) error

	ListByAuthor(authorId int) ([]models.Like, error)
	ListByPin(pinId int) ([]models.Like, error)

	PinExists(pinId int) (bool, error)
	UserExists(userId int) (bool, error)
	LikeExists(pinId, authorId int) (bool, error)
}
