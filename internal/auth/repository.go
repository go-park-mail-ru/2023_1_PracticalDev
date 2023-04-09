package auth

import (
	"github.com/pkg/errors"
	"time"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
)

var (
	UserAlreadyExistsError    = errors.New("user with this email already exists")
	UserCreationError         = errors.New("failed to create user")
	WrongPasswordOrLoginError = errors.New("wrong password or login")
	DBConnectionError         = errors.New("failed to connnect to db")
)

type Repository interface {
	Authenticate(email, hashedPassword string) (models.User, error)
	SetSession(id string, session *models.Session, expiration time.Duration) error
	CheckAuth(userId, sessionId string) (models.User, error)
	Register(user *models.User) error
	DeleteSession(userId, sessionId string) error
}
