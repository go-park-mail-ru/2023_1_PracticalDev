package auth

import (
	"time"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models/api"
)

type Service interface {
	Authenticate(login, hashed_password string) (models.User, error)
	Register(user *api.RegisterParams) error
	SetSession(id string, user *models.User, expiration time.Duration) error
	CheckAuth(userId, sessionId string) error
	DeleteSession(userId, sessionId string) error
}

func NewService(rep Repository) Service {
	return &service{rep}
}

type service struct {
	rep Repository
}

func (serv *service) Authenticate(email, hashed_password string) (models.User, error) {
	return serv.rep.Authenticate(email, hashed_password)
}

func (serv *service) SetSession(id string, user *models.User, expiration time.Duration) error {
	return serv.rep.SetSession(id, user, expiration)
}

func (serv *service) CheckAuth(userId, sessionId string) error {
	return serv.rep.CheckAuth(userId, sessionId)
}

func (serv *service) DeleteSession(userId, sessionId string) error {
	return serv.rep.DeleteSession(userId, sessionId)
}

func (serv *service) Register(user *api.RegisterParams) error {
	return serv.rep.Register(user)
}
