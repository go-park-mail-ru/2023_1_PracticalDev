package posts

import (
	"database/sql"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
)

type Repository interface {
	GetPosts() ([]models.Pin, error)
}

func NewRepository(db *sql.DB, log log.Logger) Repository {
	return repository{db, log}
}

type repository struct {
	db  *sql.DB
	log log.Logger
}

func (rep repository) GetPosts() ([]models.Pin, error) {
	rows, err := rep.db.Query("SELECT id, title FROM pins")
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
