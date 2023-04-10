package auth

import (
	"time"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models/api"
)

type SessionParams struct {
	Token      string
	LivingTime time.Duration
}

type Service interface {
	Authenticate(login, hashedPassword string) (models.User, SessionParams, error)
	Register(user *api.RegisterParams) (models.User, SessionParams, error)
	SetSession(id string, session *models.Session, expiration time.Duration) error
	CheckAuth(userId, sessionId string) (models.User, error)
	DeleteSession(userId, sessionId string) error
	CreateSession(userId int) SessionParams
}
