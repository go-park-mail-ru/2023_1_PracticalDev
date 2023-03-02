package posts

import "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"

type Service interface {
	GetPosts() ([]models.Pin, error)
}

func NewService(rep Repository) Service {
	return service{rep}
}

type service struct {
	rep Repository
}

func (serv service) GetPosts() ([]models.Pin, error) {
	return serv.rep.GetPosts()
}
