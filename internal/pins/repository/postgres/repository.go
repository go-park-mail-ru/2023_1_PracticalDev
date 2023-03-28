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
	db      *sql.DB
	log     log.Logger
	imgServ images.Service
}

func (repo *repository) Create(params *_pins.CreateParams) (models.Pin, error) {
	url, err := repo.imgServ.UploadImage(&params.MediaSource)
	if err != nil {
		return models.Pin{}, err
	}

	row := repo.db.QueryRow(createCmd,
		params.Title,
		url,
		params.Description,
		params.Author,
	)

	retrievedPin := models.Pin{}
	var title, description, mediaSource sql.NullString
	err = row.Scan(&retrievedPin.Id, &title, &mediaSource, &description, &retrievedPin.Author)
	if err != nil {
		err = _pins.ErrDb
	}
	retrievedPin.Title = title.String
	retrievedPin.Description = description.String
	retrievedPin.MediaSource = mediaSource.String
	return retrievedPin, err
}

func (repo *repository) Get(id int) (models.Pin, error) {
	row := repo.db.QueryRow(getCmd, id)

	retrievedPin := models.Pin{}
	var title, description, mediaSource sql.NullString
	err := row.Scan(&retrievedPin.Id, &title, &description, &mediaSource, &retrievedPin.Author)
	if err != nil {
		err = _pins.ErrDb
	}
	retrievedPin.Title = title.String
	retrievedPin.Description = description.String
	retrievedPin.MediaSource = mediaSource.String
	return retrievedPin, err
}

func (repo *repository) ListByUser(userId int, page, limit int) ([]models.Pin, error) {
	rows, err := repo.db.Query(listByUserCmd, userId, limit, (page-1)*limit)
	if err != nil {
		return nil, _pins.ErrDb
	}

	var pins []models.Pin
	retrievedPin := models.Pin{}
	var title, description, mediaSource sql.NullString
	for rows.Next() {
		err = rows.Scan(&retrievedPin.Id, &title, &description, &mediaSource, &retrievedPin.Author)
		if err != nil {
			return nil, _pins.ErrDb
		}
		retrievedPin.Title = title.String
		retrievedPin.Description = description.String
		retrievedPin.MediaSource = mediaSource.String
		pins = append(pins, retrievedPin)
	}
	return pins, nil
}

func (repo *repository) ListByBoard(boardId int, page, limit int) ([]models.Pin, error) {
	rows, err := repo.db.Query(listByBoardCmd, boardId, limit, (page-1)*limit)
	if err != nil {
		return nil, _pins.ErrDb
	}

	var pins []models.Pin
	retrievedPin := models.Pin{}
	var title, description, mediaSource sql.NullString
	for rows.Next() {
		err = rows.Scan(&retrievedPin.Id, &title, &description, &mediaSource, &retrievedPin.Author)
		if err != nil {
			return nil, _pins.ErrDb
		}
		retrievedPin.Title = title.String
		retrievedPin.Description = description.String
		retrievedPin.MediaSource = mediaSource.String
		pins = append(pins, retrievedPin)
	}
	return pins, nil
}

func (repo *repository) List(page, limit int) ([]models.Pin, error) {
	rows, err := repo.db.Query(listCmd, limit, (page-1)*limit)
	if err != nil {
		return nil, _pins.ErrDb
	}

	var pins []models.Pin
	retrievedPin := models.Pin{}
	var title, description, mediaSource sql.NullString
	for rows.Next() {
		err = rows.Scan(&retrievedPin.Id, &title, &description, &mediaSource, &retrievedPin.Author)
		if err != nil {
			return nil, _pins.ErrDb
		}
		retrievedPin.Title = title.String
		retrievedPin.Description = description.String
		retrievedPin.MediaSource = mediaSource.String
		pins = append(pins, retrievedPin)
	}
	return pins, nil
}

func (repo *repository) FullUpdate(params *_pins.FullUpdateParams) (models.Pin, error) {
	url, err := repo.imgServ.UploadImage(&params.MediaSource)
	if err != nil {
		return models.Pin{}, err
	}

	row := repo.db.QueryRow(fullUpdateCmd,
		params.Title,
		params.Description,
		url,
		params.Id,
	)

	retrievedPin := models.Pin{}
	var title, description, mediaSource sql.NullString
	err = row.Scan(&retrievedPin.Id, &title, &description, &mediaSource, &retrievedPin.Author)
	if err != nil {
		err = _pins.ErrDb
	}
	retrievedPin.Title = title.String
	retrievedPin.Description = description.String
	retrievedPin.MediaSource = mediaSource.String
	return retrievedPin, err
}

func (repo *repository) Delete(id int) error {
	res, err := repo.db.Exec(deleteCmd, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil || count < 1 {
		return err
	}
	return nil
}

func (repo *repository) AddToBoard(boardId, pinId int) error {
	res, err := repo.db.Exec(addToBoardCmd, pinId, boardId)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil || count < 1 {
		return err
	}
	return nil
}

func (repo *repository) RemoveFromBoard(boardId, pinId int) error {
	res, err := repo.db.Exec(deleteFromBoardCmd, pinId, boardId)
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
	row := repo.db.QueryRow(checkWriteCmd,
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