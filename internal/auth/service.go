package auth

import (
	"time"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
)

type SessionParams struct {
	Token      string
	LivingTime time.Duration
}

type LoginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterParams struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Service interface {
	Authenticate(login, hashedPassword string) (models.User, SessionParams, error)
	Register(user *RegisterParams) (models.User, SessionParams, error)
	SetSession(id string, session *models.Session, expiration time.Duration) error
	CheckAuth(userId, sessionId string) (models.User, error)
	DeleteSession(userId, sessionId string) error
	CreateSession(userId int) SessionParams
}
