package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth/hasher"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	"github.com/redis/go-redis/v9"
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

func NewRepository(db *sql.DB, rdb *redis.Client, ctx context.Context, log log.Logger) Repository {
	return &repository{db, rdb, ctx, log}
}

type repository struct {
	db  *sql.DB
	rdb *redis.Client
	ctx context.Context
	log log.Logger
}

func scanUser(user *models.User, row *sql.Row) error {
	var profile_image, website_url sql.NullString
	err := row.Scan(&user.Id, &user.Username, &user.Email, &user.HashedPassword, &user.Name, &profile_image, &website_url, &user.AccountType)

	user.WebsiteUrl = website_url.String
	user.ProfileImage = profile_image.String
	return err
}

func (rep *repository) Authenticate(email, password string) (models.User, error) {
	authCommand := "SELECT * FROM users WHERE email = $1"
	row := rep.db.QueryRow(authCommand, email)
	user := models.User{}
	hasher := hasher.NewHasher()

	err := scanUser(&user, row)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return models.User{}, WrongPasswordOrLoginError
		} else {
			return models.User{}, DBConnectionError
		}
	}

	if err := hasher.CompareHashAndPassword(user.HashedPassword, password); err != nil {
		return models.User{}, WrongPasswordOrLoginError
	}

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

func (rep *repository) CheckAuth(userId, sessionId string) (models.User, error) {
	err := rep.rdb.HGet(rep.ctx, userId, sessionId).Err()
	user := models.User{}

	if err != nil {
		return user, err
	}

	row := rep.db.QueryRow("SELECT * FROM users WHERE id = $1", userId)
	err = scanUser(&user, row)

	return user, err
}

func (rep *repository) DeleteSession(userId, sessionId string) error {
	if err := rep.rdb.HGet(rep.ctx, userId, sessionId).Err(); err != nil {
		return err
	}
	rep.rdb.HDel(rep.ctx, userId, sessionId)
	return nil
}

func (rep *repository) Register(user *models.User) error {
	row := rep.db.QueryRow("SELECT email FROM users WHERE email = $1", user.Email)
	tmp := ""
	err := row.Scan(&tmp)

	if err == nil {
		return UserAlreadyExistsError
	}

	if err.Error() != "sql: no rows in result set" {
		return DBConnectionError
	}

	insertCommand := "INSERT INTO users (username, name, email, hashed_password, account_type, profile_image, website_url) VALUES ($1, $2, $3, $4, $5, $6, $7)"

	if _, err := rep.db.Exec(insertCommand, user.Username, user.Name, user.Email, user.HashedPassword, user.AccountType, user.ProfileImage, user.WebsiteUrl); err != nil {
		return UserCreationError
	}

	return nil
}
