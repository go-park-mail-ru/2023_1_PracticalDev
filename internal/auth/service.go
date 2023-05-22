package auth

import (
	"time"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
)

type SessionParams struct {
	Token      string
	LivingTime time.Duration
}

type RegisterParams struct {
	Username string
	Email    string
	Name     string
	Password string
}

type Service interface {
	Authenticate(login, hashedPassword string) (models.User, SessionParams, error)
	Register(user *RegisterParams) (models.User, SessionParams, error)
	SetSession(id string, session *models.Session, expiration time.Duration) error
	CheckAuth(userId, sessionId string) (models.User, error)
	DeleteSession(userId, sessionId string) error
	CreateSession(userId int) SessionParams
}
