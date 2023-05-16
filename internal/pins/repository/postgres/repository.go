package postgres

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg" // импортируем пакет для декодирования JPEG
	_ "image/png"  // импортируем пакет для декодирования JPEG
	"math"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	images "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/images/client"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	pkgPins "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pins"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
)

func NewRepository(db *sql.DB, s3Service images.ImageClient, log *zap.Logger) pkgPins.Repository {
	return &repository{db, log, s3Service}
}

type repository struct {
	db      *sql.DB
	log     *zap.Logger
	imgServ images.ImageClient
}

func bytesToImage(b []byte) (image.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	return img, nil
}

type AvgColor struct {
	Red   uint8
	Green uint8
	Blue  uint8
}

func calcAvgColor(img image.Image) (result AvgColor) {
	imgSize := img.Bounds().Size()

	var redSum float64
	var greenSum float64
	var blueSum float64

	for x := 0; x <= imgSize.X; x++ {
		for y := 0; y <= imgSize.Y; y++ {
			pixel := img.At(x, y)
			col := color.RGBAModel.Convert(pixel).(color.RGBA)

			redSum += float64(col.R)
			greenSum += float64(col.G)
			blueSum += float64(col.B)
		}
	}

	imgArea := float64(imgSize.X * imgSize.Y)

	result.Red = uint8(math.Round(redSum / imgArea))
	result.Green = uint8(math.Round(greenSum / imgArea))
	result.Blue = uint8(math.Round(blueSum / imgArea))

	return
}

const createCmd = `INSERT INTO pins (title, media_source, media_source_color, description, author_id)
				   VALUES ($1, $2, $3, $4, $5)
				   RETURNING id, title, media_source, description, author_id;`

func (repo *repository) Create(params *pkgPins.CreateParams) (models.Pin, error) {
	url, err := repo.imgServ.UploadImage(context.Background(), &params.MediaSource)
	if err != nil {
		return models.Pin{}, errors.Wrap(pkgErrors.ErrImageService, err.Error())
	}

	img, err := bytesToImage(params.MediaSource.Bytes)
	if err != nil {
		return models.Pin{}, err
	}
	avgColor := calcAvgColor(img)
	// (%d, %d, %d) #12F4D4
	avgColorStr := fmt.Sprintf("rgb(%d, %d, %d)", avgColor.Red, avgColor.Green, avgColor.Blue)

	row := repo.db.QueryRow(createCmd,
		params.Title,
		url,
		avgColorStr,
		params.Description,
		params.Author,
	)

	retrievedPin := models.Pin{}
	var title, description, mediaSource sql.NullString
	err = row.Scan(&retrievedPin.Id, &title, &mediaSource, &description, &retrievedPin.Author)
	if err != nil {
		return models.Pin{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	retrievedPin.Title = title.String
	retrievedPin.Description = description.String
	retrievedPin.MediaSource = mediaSource.String
	return retrievedPin, nil
}

func (repo *repository) Get(id int) (models.Pin, error) {
	row := repo.db.QueryRow(getCmd, id)

	pin := models.Pin{}
	var title, description, mediaSource sql.NullString
	err := row.Scan(&pin.Id, &title, &description, &mediaSource, &pin.NumLikes, &pin.Author)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Pin{}, errors.Wrap(pkgErrors.ErrPinNotFound, err.Error())
		} else {
			return models.Pin{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
		}
	}

	pin.Title = title.String
	pin.Description = description.String
	pin.MediaSource = mediaSource.String
	return pin, nil
}

func (repo *repository) ListByAuthor(userId int, page, limit int) ([]models.Pin, error) {
	rows, err := repo.db.Query(listByUserCmd, userId, limit, (page-1)*limit)
	if err != nil {
		return nil, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	var pins []models.Pin
	pin := models.Pin{}
	var title, description, mediaSource sql.NullString

	for rows.Next() {
		err = rows.Scan(&pin.Id, &title, &description, &mediaSource, &pin.NumLikes, &pin.Author)
		if err != nil {
			return nil, errors.Wrap(pkgErrors.ErrDb, err.Error())
		}
		pin.Title = title.String
		pin.Description = description.String
		pin.MediaSource = mediaSource.String
		pins = append(pins, pin)
	}

	return pins, nil
}

const listCmd = `SELECT id, title, description, media_source, media_source_color, n_likes, author_id 
					FROM pins 
					ORDER BY created_at DESC 
					LIMIT $1 OFFSET $2;`

func (repo *repository) List(page, limit int) ([]models.Pin, error) {
	rows, err := repo.db.Query(listCmd, limit, (page-1)*limit)
	if err != nil {
		return nil, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	pins := []models.Pin{}
	pin := models.Pin{}
	var title, description, mediaSourceColor sql.NullString

	for rows.Next() {
		err = rows.Scan(&pin.Id, &title, &description, &pin.MediaSource, &mediaSourceColor, &pin.NumLikes,
			&pin.Author)
		if err != nil {
			return nil, errors.Wrap(pkgErrors.ErrDb, err.Error())
		}
		pin.Title = title.String
		pin.Description = description.String
		pin.MediaSourceColor = mediaSourceColor.String
		pins = append(pins, pin)
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
		if errors.Is(err, sql.ErrNoRows) {
			return models.Pin{}, errors.Wrap(pkgErrors.ErrPinNotFound, err.Error())
		} else {
			return models.Pin{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
		}
	}

	retrievedPin.Title = title.String
	retrievedPin.Description = description.String
	retrievedPin.MediaSource = mediaSource.String
	return retrievedPin, nil
}

func (repo *repository) Delete(id int) error {
	_, err := repo.db.Exec(deleteCmd, id)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrDb, err.Error())
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
		return false, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}
	return liked, nil
}

func (repo *repository) CheckWriteAccess(userId, pinId string) (bool, error) {
	row := repo.db.QueryRow(checkWriteCmd, pinId, userId)

	var access bool
	err := row.Scan(&access)
	if err != nil {
		return false, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}
	return access, err
}

func (repo *repository) CheckReadAccess(userId, pinId string) (bool, error) {
	return true, nil
}
