package postgres

import (
	"database/sql"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/images"
	pkgLikes "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/likes"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	pkgPins "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pins"
)

func NewRepository(db *sql.DB, s3Service images.Service, log log.Logger) pkgPins.Repository {
	return &repository{db, log, s3Service}
}

type repository struct {
	db      *sql.DB
	log     log.Logger
	imgServ images.Service
}

func (repo *repository) Create(params *pkgPins.CreateParams) (models.Pin, error) {
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
		err = pkgPins.ErrDb
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
		err = pkgPins.ErrDb
	}
	retrievedPin.Title = title.String
	retrievedPin.Description = description.String
	retrievedPin.MediaSource = mediaSource.String
	return retrievedPin, err
}

func (repo *repository) ListByAuthor(userId int, page, limit int) ([]models.Pin, error) {
	rows, err := repo.db.Query(listByUserCmd, userId, limit, (page-1)*limit)
	if err != nil {
		return nil, pkgPins.ErrDb
	}

	var pins []models.Pin
	retrievedPin := models.Pin{}
	var title, description, mediaSource sql.NullString
	for rows.Next() {
		err = rows.Scan(&retrievedPin.Id, &title, &description, &mediaSource, &retrievedPin.Author)
		if err != nil {
			return nil, pkgPins.ErrDb
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
		return nil, pkgPins.ErrDb
	}

	pins := []models.Pin{}
	retrievedPin := models.Pin{}
	var title, description, mediaSource sql.NullString
	for rows.Next() {
		err = rows.Scan(&retrievedPin.Id, &title, &description, &mediaSource, &retrievedPin.Author)
		if err != nil {
			return nil, pkgPins.ErrDb
		}
		retrievedPin.Title = title.String
		retrievedPin.Description = description.String
		retrievedPin.MediaSource = mediaSource.String
		pins = append(pins, retrievedPin)
	}
	return pins, nil
}

func (repo *repository) FullUpdate(params *pkgPins.FullUpdateParams) (models.Pin, error) {
	row := repo.db.QueryRow(fullUpdateCmd,
		params.Title,
		params.Description,
		params.Id,
	)

	retrievedPin := models.Pin{}
	var title, description, mediaSource sql.NullString
	err := row.Scan(&retrievedPin.Id, &title, &description, &mediaSource, &retrievedPin.Author)
	if err != nil {
		err = pkgPins.ErrDb
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

const isLikedByUserCmd = `SELECT EXISTS(SELECT pin_id
										FROM pin_likes
										WHERE pin_id = $1 AND author_id = $2) AS liked;`

func (repo *repository) IsLikedByUser(pinId, userId int) (bool, error) {
	row := repo.db.QueryRow(isLikedByUserCmd, pinId, userId)

	var liked bool
	err := row.Scan(&liked)
	if err != nil {
		return false, pkgLikes.ErrDb
	}
	return liked, nil
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
