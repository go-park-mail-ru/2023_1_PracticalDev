package auth

import (
	"strconv"
	"time"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth/hasher"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models/api"
	"github.com/google/uuid"
)

type SessionParams struct {
	token      string
	livingTime time.Duration
}

type Service interface {
	Authenticate(login, hashed_password string) (models.User, SessionParams, error)
	Register(user *api.RegisterParams) (models.User, SessionParams, error)
	SetSession(id string, session *models.Session, expiration time.Duration) error
	CheckAuth(userId, sessionId string) (models.User, error)
	DeleteSession(userId, sessionId string) error
	CreateSession(userId int) SessionParams
}

func NewService(rep Repository) Service {
	return &service{rep}
}

type service struct {
	rep Repository
}

func (serv *service) Authenticate(email, hashed_password string) (models.User, SessionParams, error) {
	user, err := serv.rep.Authenticate(email, hashed_password)
	session := SessionParams{}

	if err != nil {
		return user, session, err
	}

	sessionParams := serv.CreateSession(user.Id)
	sessionData := models.Session{
		UserId:    user.Id,
		UserEmail: user.Email,
	}

	err = serv.SetSession(sessionParams.token, &sessionData, sessionParams.livingTime)

	return user, sessionParams, err
}

func (serv *service) CreateSession(userId int) SessionParams {
	token := strconv.Itoa(userId) + "$" + uuid.New().String()
	livingTime := 5 * time.Hour
	return SessionParams{token, livingTime}
}

func (serv *service) SetSession(token string, session *models.Session, expiration time.Duration) error {
	return serv.rep.SetSession(token, session, expiration)
}

func (serv *service) CheckAuth(userId, sessionId string) (models.User, error) {
	return serv.rep.CheckAuth(userId, sessionId)
}

func (serv *service) DeleteSession(userId, sessionId string) error {
	return serv.rep.DeleteSession(userId, sessionId)
}

func (serv *service) Register(user *api.RegisterParams) (models.User, SessionParams, error) {
	hasher := hasher.NewHasher()
	hash, _ := hasher.GetHashedPassword(user.Password)

	tmp := models.User{
		Name:           user.Name,
		Username:       user.Username,
		Email:          user.Email,
		ProfileImage:   "",
		WebsiteUrl:     "",
		AccountType:    "personal",
		HashedPassword: hash,
	}

	err := serv.rep.Register(&tmp)
	if err != nil {
		return tmp, SessionParams{}, err
	}

	return serv.Authenticate(user.Email, user.Password)
}
