package auth

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	"time"
)

type Service interface {
	Authenticate(login, hashed_password string) (models.User, error)
	Register(user models.User) error
	SetSession(id string, user models.User, expiration time.Duration) error
	CheckAuth(id string) error
	DeleteSession(id string) error
}

func NewService(rep Repository) Service {
	return service{rep}
}

type service struct {
	rep Repository
}

func (serv service) Authenticate(email, hashed_password string) (models.User, error) {
	return serv.rep.Authenticate(email, hashed_password)
}

func (serv service) SetSession(id string, user models.User, expiration time.Duration) error {
	return serv.rep.SetSession(id, user, expiration)
}

func (serv service) CheckAuth(id string) error {
	return serv.rep.CheckAuth(id)
}

func (serv service) DeleteSession(id string) error {
	return serv.rep.DeleteSession(id)
}

func (serv service) Register(user models.User) error {
	return serv.rep.Register(user)
}
