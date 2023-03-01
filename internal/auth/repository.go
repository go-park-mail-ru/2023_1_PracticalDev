package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	"github.com/redis/go-redis/v9"
)

var (
	UserAlreadyExistsError = errors.New("user with this email already exists")
	UserCreationError      = errors.New("Failed to create user")
)

type Repository interface {
	Authenticate(email, hashedPassword string) (models.User, error)
	SetSession(id string, user models.User, expiration time.Duration) error
	CheckAuth(id string) error
	Register(user models.User) error
	DeleteSession(id string) error
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

func (rep repository) Authenticate(email, hashedPassword string) (models.User, error) {
	authCommand := "SELECT * FROM users WHERE email = $1 AND hashed_password = $2"
	row := rep.db.QueryRow(authCommand, email, hashedPassword)
	user := models.User{}

	if err := row.Scan(&user.Id, &user.Email, &user.Username, &user.HashedPassword); err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (rep repository) SetSession(id string, user models.User, expiration time.Duration) error {
	tmp, _ := json.Marshal(user)
	err := rep.rdb.Set(rep.ctx, id, string(tmp), expiration).Err()
	return err
}

func (rep repository) CheckAuth(id string) error {
	_, err := rep.rdb.Get(rep.ctx, id).Result()
	return err
}

func (rep repository) DeleteSession(id string) error {
	return rep.rdb.Del(rep.ctx, id).Err()
}

func (rep repository) Register(user models.User) error {
	row := rep.db.QueryRow("SELECT * FROM users WHERE email = $1", user.Email)

	if err := row.Err(); err != nil {
		return UserAlreadyExistsError
	}

	insertCommand := "INSERT INTO users (username, email, hashed_password) VALUES ($1, $2, $3)"
	if _, err := rep.db.Exec(insertCommand, user.Username, user.Email, user.HashedPassword); err != nil {
		return UserCreationError
	}

	return nil
}
