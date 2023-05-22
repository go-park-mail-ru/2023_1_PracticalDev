package auth

import (
	"github.com/pkg/errors"
	"time"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
)

var (
	WrongPasswordOrLoginError = errors.New("wrong password or login")
)

type Repository interface {
	Authenticate(email, password string) (models.User, error)
	SetSession(id string, session *models.Session, expiration time.Duration) error
	CheckAuth(userId, sessionId string) (models.User, error)
	Register(user *models.User) error
	DeleteSession(userId, sessionId string) error
}
