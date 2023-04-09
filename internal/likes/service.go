package likes

import "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"

type Service interface {
	Like(pinId, authorId int) error
	Unlike(pinId, authorId int) error

	ListByAuthor(authorId int) ([]models.Like, error)
	ListByPin(pinId int) ([]models.Like, error)
}
