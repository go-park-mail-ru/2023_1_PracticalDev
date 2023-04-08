package postgres

import (
	"database/sql"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/images"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	_pins "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pins"
)

func NewRepository(db *sql.DB, s3Service images.Service, log log.Logger) _pins.Repository {
	return &repository{db, log, s3Service}
}

type repository struct {
	db        *sql.DB
	log       log.Logger
	s3Service images.Service
}

func (repo *repository) CreatePin(params *models.Pin, image *models.Image) (models.Pin, error) {
	const insertCommand = `INSERT INTO pins (title, media_source, description, author_id)
				      	   VALUES ($1, $2, $3, $4)
						   RETURNING id, title, media_source, description, author_id;`

	url, err := repo.s3Service.UploadImage(image)

	if err != nil {
		return models.Pin{}, err
	}

	row := repo.db.QueryRow(insertCommand,
		params.Title,
		url,
		params.Description,
		params.Author,
	)

	retrievedPin := models.Pin{}
	var title, description, mediaSource sql.NullString

	err = row.Scan(&retrievedPin.Id, &title, &mediaSource, &description, &retrievedPin.Author)

	retrievedPin.Title = title.String
	retrievedPin.Description = description.String
	retrievedPin.MediaSource = mediaSource.String
	return retrievedPin, err
}

func (repo *repository) GetPin(id int) (models.Pin, error) {
	const getCommand = `SELECT id, title, description, media_source, author_id
						FROM pins
						WHERE id = $1`

	row := repo.db.QueryRow(getCommand, id)

	retrievedPin := models.Pin{}
	var title, description, mediaSource sql.NullString

	err := row.Scan(&retrievedPin.Id, &title, &description, &mediaSource, &retrievedPin.Author)

	retrievedPin.Title = title.String
	retrievedPin.Description = description.String
	retrievedPin.MediaSource = mediaSource.String
	return retrievedPin, err
}

func (repo *repository) GetPinsByUser(userId int, page, limit int) ([]models.Pin, error) {
	const getCommand = `SELECT id, title, description, media_source, author_id
						FROM pins 
						WHERE author_id = $1
						ORDER BY created_at DESC 
						LIMIT $2 OFFSET $3;`

	rows, err := repo.db.Query(getCommand, userId, limit, (page-1)*limit)

	if err != nil {
		return []models.Pin{}, err
	}

	pins := []models.Pin{}
	for rows.Next() {

		retrievedPin := models.Pin{}
		var title, description, mediaSource sql.NullString

		err = rows.Scan(&retrievedPin.Id, &title, &description, &mediaSource, &retrievedPin.Author)

		retrievedPin.Title = title.String
		retrievedPin.Description = description.String
		retrievedPin.MediaSource = mediaSource.String
		pins = append(pins, retrievedPin)
	}
	return pins, err
}

func (repo *repository) GetPinsByBoard(boardId int, page, limit int) ([]models.Pin, error) {
	const getCommand = `SELECT pins.id, title, description, media_source, author_id 
						FROM pins 
						JOIN boards_pins AS b
						ON b.board_id = $1 AND b.pin_id = pins.id
						ORDER BY created_at DESC 
						LIMIT $2 OFFSET $3;`

	rows, err := repo.db.Query(getCommand, boardId, limit, (page-1)*limit)

	if err != nil {
		return []models.Pin{}, err
	}

	pins := []models.Pin{}
	for rows.Next() {

		retrievedPin := models.Pin{}
		var title, description, mediaSource sql.NullString

		err = rows.Scan(&retrievedPin.Id, &title, &description, &mediaSource, &retrievedPin.Author)

		retrievedPin.Title = title.String
		retrievedPin.Description = description.String
		retrievedPin.MediaSource = mediaSource.String
		pins = append(pins, retrievedPin)
	}
	return pins, err
}

func (repo *repository) GetPins(page, limit int) ([]models.Pin, error) {
	const getCommand = `SELECT id, title, description, media_source, author_id 
						FROM pins 
						ORDER BY created_at DESC 
						LIMIT $1 OFFSET $2;`

	rows, err := repo.db.Query(getCommand, limit, (page-1)*limit)

	if err != nil {
		return []models.Pin{}, err
	}

	pins := []models.Pin{}
	for rows.Next() {
		retrievedPin := models.Pin{}
		var title, description, mediaSource sql.NullString

		err = rows.Scan(&retrievedPin.Id, &title, &description, &mediaSource, &retrievedPin.Author)

		retrievedPin.Title = title.String
		retrievedPin.Description = description.String
		retrievedPin.MediaSource = mediaSource.String
		pins = append(pins, retrievedPin)
	}
	return pins, err
}

func (repo *repository) UpdatePin(params *models.Pin) (models.Pin, error) {
	const fullUpdateCommand = `UPDATE pins
								SET title = $1::VARCHAR,
								description = $2::TEXT
								WHERE id = $3
								RETURNING id, title, description, media_source, author_id;`

	row := repo.db.QueryRow(fullUpdateCommand,
		params.Title,
		params.Description,
		params.Id,
	)
	retrievedPin := models.Pin{}
	var title, description, mediaSource sql.NullString

	err := row.Scan(&retrievedPin.Id, &title, &description, &mediaSource, &retrievedPin.Author)

	retrievedPin.Title = title.String
	retrievedPin.Description = description.String
	retrievedPin.MediaSource = mediaSource.String
	return retrievedPin, err
}

func (repo *repository) DeletePin(id int) error {
	const deleteCommand = `DELETE FROM pins 
							WHERE id = $1;`

	res, err := repo.db.Exec(deleteCommand, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil || count < 1 {
		return err
	}
	return nil
}

func (repo *repository) AddPinToBoard(boardId, pinId int) error {
	const addCommand = `INSERT INTO boards_pins(pin_id, board_id)
						VALUES($1, $2)`

	res, err := repo.db.Exec(addCommand, pinId, boardId)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil || count < 1 {
		return err
	}
	return nil
}

func (repo *repository) RemovePinFromBoard(boardId, pinId int) error {
	const addCommand = `DELETE FROM boards_pins
						WHERE pin_id = $1 AND board_id = $2`

	res, err := repo.db.Exec(addCommand, pinId, boardId)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil || count < 1 {
		return err
	}
	return nil
}

func (repo *repository) CheckWriteAccess(userId, pinId string) (bool, error) {
	const checkCommand = `SELECT EXISTS(SELECT id
		FROM pins
		 WHERE id = $1 AND author_id = $2);`

	row := repo.db.QueryRow(checkCommand,
		pinId,
		userId,
	)

	var access bool
	err := row.Scan(&access)
	return access, err
}

func (repo *repository) CheckReadAccess(userId, pinId string) (bool, error) {
	return true, nil
}
