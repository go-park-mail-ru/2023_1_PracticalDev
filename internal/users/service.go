package users

import "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"

type Service interface {
	GetUser(id int) (models.User, error)
}

func NewService(rep Repository) Service {
	return service{rep}
}

type service struct {
	rep Repository
}

func (serv service) GetUser(id int) (models.User, error) {
	return serv.rep.GetUser(id)
}
