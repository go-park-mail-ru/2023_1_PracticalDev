package service

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth"
	hasherPkg "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth/hasher"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/constants"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

func NewService(rep auth.Repository) auth.Service {
	return &service{rep}
}

type service struct {
	rep auth.Repository
}

func (serv *service) Authenticate(email, password string) (models.User, auth.SessionParams, error) {
	user, err := serv.rep.Authenticate(email, password)
	if err != nil {
		return user, auth.SessionParams{}, errors.Wrap(err, "Authenticate")
	}

	sessionParams := serv.CreateSession(user.Id)
	sessionData := models.Session{
		UserId:    user.Id,
		UserEmail: user.Email,
	}
	err = serv.SetSession(sessionParams.Token, &sessionData, sessionParams.LivingTime)
	return user, sessionParams, err
}

func (serv *service) CreateSession(userId int) auth.SessionParams {
	token := strconv.Itoa(userId) + "$" + uuid.New().String()
	livingTime := 5 * time.Hour
	return auth.SessionParams{Token: token, LivingTime: livingTime}
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

func (serv *service) Register(user *auth.RegisterParams) (models.User, auth.SessionParams, error) {
	hasher := hasherPkg.NewHasher()
	hash, _ := hasher.GetHashedPassword(user.Password)

	tmp := models.User{
		Name:           user.Name,
		Username:       user.Username,
		Email:          user.Email,
		ProfileImage:   constants.DefaultAvatar,
		WebsiteUrl:     constants.DefaultWebsiteUrl,
		AccountType:    constants.DefaultAccountType,
		HashedPassword: hash,
	}

	err := serv.rep.Register(&tmp)
	if err != nil {
		return models.User{}, auth.SessionParams{}, errors.Wrap(err, "Register")
	}

	return serv.Authenticate(user.Email, user.Password)
}
