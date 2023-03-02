package auth

import (
	"context"

	"database/sql"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

var (
	UserAlreadyExistsError    = errors.New("user with this email already exists")
	UserCreationError         = errors.New("failed to create user")
	WrongPasswordOrLoginError = errors.New("wrong password or login")
)

type Repository interface {
	Authenticate(email, hashedPassword string) (models.User, error)
	SetSession(id string, user models.User, expiration time.Duration) error
	CheckAuth(userId, sessionId string) error
	Register(user models.User) error
	DeleteSession(userId, sessionId string) error
}

func NewRepository(db *sql.DB, rdb *redis.Client, ctx context.Context, log log.Logger) Repository {
	return repository{db, rdb, ctx, log}
}

type repository struct {
	db  *sql.DB
	rdb *redis.Client
	ctx context.Context
	log log.Logger
}

func (rep repository) Authenticate(email, password string) (models.User, error) {
	authCommand := "SELECT * FROM users WHERE email = $1"
	row := rep.db.QueryRow(authCommand, email)
	if err := row.Err(); err != nil {
		return models.User{}, WrongPasswordOrLoginError
	}

	user := models.User{}

	if err := row.Scan(&user.Id, &user.Email, &user.Username, &user.HashedPassword); err != nil {
		return models.User{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		return models.User{}, WrongPasswordOrLoginError
	}

	return user, nil
}

func (rep repository) SetSession(sessionId string, user models.User, expiration time.Duration) error {
	tmp, _ := json.Marshal(user)
	err := rep.rdb.HSet(rep.ctx, strconv.Itoa(user.Id), sessionId, tmp).Err()
	return err
}

func (rep repository) CheckAuth(userId, sessionId string) error {
	_, err := rep.rdb.HGet(rep.ctx, userId, sessionId).Result()
	return err
}

func (rep repository) DeleteSession(userId, sessionId string) error {
	return rep.rdb.HDel(rep.ctx, userId, sessionId).Err()
}

func (rep repository) Register(user models.User) error {
	row := rep.db.QueryRow("SELECT * FROM users WHERE email = $1", user.Email)

	if err := row.Err(); err != nil {
		return UserAlreadyExistsError
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(user.HashedPassword), bcrypt.DefaultCost)

	insertCommand := "INSERT INTO users (username, email, hashed_password) VALUES ($1, $2, $3)"

	if _, err := rep.db.Exec(insertCommand, user.Username, user.Email, string(hash)); err != nil {
		return UserCreationError
	}

	return nil
}
