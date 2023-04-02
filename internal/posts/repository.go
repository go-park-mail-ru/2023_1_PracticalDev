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
	return &repository{db, log}
}

type repository struct {
	db  *sql.DB
	log log.Logger
}

func (rep *repository) GetPosts(limit, offset int) ([]models.Pin, error) {
	rows, err := rep.db.Query(`SELECT id, title, description, media_source
									FROM pins 
									ORDER BY created_at DESC 
									LIMIT $1 OFFSET $2;`, limit, offset)
	if err != nil {
		return []models.Pin{}, err
	}

	pins := []models.Pin{}
	pin := models.Pin{}
	var title, description, mediaSource sql.NullString
	for rows.Next() {
		err = rows.Scan(&pin.Id, &title, &description, &mediaSource)
		if err != nil {
			break
		}
		pin.Title = title.String
		pin.Description = description.String
		pin.MediaSource = mediaSource.String
		pins = append(pins, pin)
	}
	return pins, err
}
