package service

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth"
	hasherPkg "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth/hasher"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models/api"
	"github.com/google/uuid"
	"strconv"
	"time"
)

const defaultAvatar = "https://pickpin.hb.bizmrg.com/default-user-icon-8-4024862977"

func NewService(rep auth.Repository) auth.Service {
	return &service{rep}
}

type service struct {
	rep auth.Repository
}

func (serv *service) Authenticate(email, password string) (models.User, auth.SessionParams, error) {
	user, err := serv.rep.Authenticate(email, password)
	session := auth.SessionParams{}
	if err != nil {
		return user, session, err
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

func (serv *service) Register(user *api.RegisterParams) (models.User, auth.SessionParams, error) {
	hasher := hasherPkg.NewHasher()
	hash, _ := hasher.GetHashedPassword(user.Password)

	tmp := models.User{
		Name:           user.Name,
		Username:       user.Username,
		Email:          user.Email,
		ProfileImage:   defaultAvatar,
		WebsiteUrl:     "https://pickpin.hb.bizmrg.com/default-user-icon-8-4024862977",
		AccountType:    "personal",
		HashedPassword: hash,
	}

	err := serv.rep.Register(&tmp)
	if err != nil {
		return tmp, auth.SessionParams{}, err
	}

	return serv.Authenticate(user.Email, user.Password)
}
