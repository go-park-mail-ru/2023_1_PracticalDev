package auth

import (
	"context"

	"database/sql"
	"encoding/json"
	"errors"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth/hasher"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models/api"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
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
	CheckAuth(userId, sessionId string) error
	Register(user *api.RegisterParams) error
	DeleteSession(userId, sessionId string) error
}

func NewRepository(db *sql.DB, rdb *redis.Client, ctx context.Context, log log.Logger) Repository {
	return &repository{db, rdb, ctx, log}
}

type repository struct {
	db  *sql.DB
	rdb *redis.Client
	ctx context.Context
	log log.Logger
}

func (rep *repository) Authenticate(email, password string) (models.User, error) {
	authCommand := "SELECT * FROM users WHERE email = $1"
	row := rep.db.QueryRow(authCommand, email)
	user := models.User{}
	hasher := hasher.NewHasher()

	var profile_image, website_url sql.NullString
	err := row.Scan(&user.Id, &user.Username, &user.Email, &user.HashedPassword, &user.Name, &profile_image, &website_url, &user.Account_type)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return models.User{}, WrongPasswordOrLoginError
		} else {
			return models.User{}, DBConnectionError
		}
	}

	if err := hasher.CompareHashAndPassword(user.HashedPassword, password); err != nil {
		return models.User{}, WrongPasswordOrLoginError
	}

	user.Website_url = website_url.String
	user.Profile_image = profile_image.String

	return user, nil
}

func (rep *repository) SetSession(sessionId string, session *models.Session, expiration time.Duration) error {
	tmp, _ := json.Marshal(session)

	err := rep.rdb.HSet(rep.ctx, strconv.Itoa(session.UserId), sessionId, tmp).Err()

	if err != nil {
		rep.rdb.Expire(rep.ctx, strconv.Itoa(session.UserId), expiration)
	}

	return err
}

func (rep *repository) CheckAuth(userId, sessionId string) error {
	err := rep.rdb.HGet(rep.ctx, userId, sessionId).Err()
	return err
}

func (rep *repository) DeleteSession(userId, sessionId string) error {
	return rep.rdb.HDel(rep.ctx, userId, sessionId).Err()
}

func (rep *repository) Register(user *api.RegisterParams) error {
	row := rep.db.QueryRow("SELECT email FROM users WHERE email = $1", user.Email)
	hasher := hasher.NewHasher()
	tmp := ""
	err := row.Scan(&tmp)

	if err == nil {
		return UserAlreadyExistsError
	}

	if (err != nil) && (err.Error() == "no rows in result set") {
		return DBConnectionError
	}

	hash, _ := hasher.GetHashedPassword(user.Password)

	insertCommand := "INSERT INTO users (username, name, email, hashed_password, account_type) VALUES ($1, $2, $3, $4, $5)"

	if _, err := rep.db.Exec(insertCommand, user.Username, user.Name, user.Email, string(hash), "personal"); err != nil {
		rep.log.Info(err)
		return UserCreationError
	}

	return nil
}
