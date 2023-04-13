package service

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/users"
)

type service struct {
	rep users.Repository
}

func NewService(rep users.Repository) users.Service {
	return service{rep}
}

func (serv service) Get(id int) (models.User, error) {
	return serv.rep.Get(id)
}
