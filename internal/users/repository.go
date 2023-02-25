package users

import (
	"database/sql"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
)

type Repository interface {
	GetUser(id int) (models.User, error)
}

func NewRepository(db *sql.DB, log log.Logger) Repository {
	return repository{db, log}
}

type repository struct {
	db  *sql.DB
	log log.Logger
}

func (rep repository) GetUser(id int) (models.User, error) {
	rows, err := rep.db.Query("SELECT user_id, username, email FROM users WHERE user_id = $1", id)
	if err != nil {
		return models.User{}, err
	}

	user := models.User{}
	rows.Next()
	err = rows.Scan(&user.Id, &user.Email, &user.Username)
	return user, err
}
