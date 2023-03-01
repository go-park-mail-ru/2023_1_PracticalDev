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
	rows, err := rep.db.Query("SELECT id, username, email FROM users WHERE id = $1", id)
	if err != nil {
		return models.User{}, err
	}

	user := models.User{}
	rows.Next()
	err = rows.Scan(&user.Id, &user.Username, &user.Email, &user.HashedPassword)
	return user, err
}
