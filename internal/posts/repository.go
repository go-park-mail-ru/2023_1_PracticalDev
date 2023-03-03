package posts

import (
	"database/sql"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
)

type Repository interface {
	GetPosts(limit, offset int) ([]models.Pin, error)
}

func NewRepository(db *sql.DB, log log.Logger) Repository {
	return repository{db, log}
}

type repository struct {
	db  *sql.DB
	log log.Logger
}

func (rep repository) GetPosts(limit, offset int) ([]models.Pin, error) {
	rows, err := rep.db.Query("SELECT id, title FROM pins ORDER BY created_at DESC LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return []models.Pin{}, err
	}

	pins := []models.Pin{}
	pin := models.Pin{}
	for rows.Next() {
		err = rows.Scan(&pin.Id, &pin.Title)
		pins = append(pins, pin)
	}
	return pins, err
}
