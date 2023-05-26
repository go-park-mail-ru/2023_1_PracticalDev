package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/constants"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	images "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/images/client"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	pkgPins "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pins"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
	pkgImage "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/image"
)

func NewRepository(db *sql.DB, s3Service images.ImageClient, log *zap.Logger) pkgPins.Repository {
	return &repository{db, log, s3Service}
}

type repository struct {
	db      *sql.DB
	log     *zap.Logger
	imgServ images.ImageClient
}

const createCmd = `
		INSERT INTO pins (title, media_source, media_source_color, description, author_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, title, media_source, media_source_color, description, author_id;`

func (repo *repository) Create(params *pkgPins.CreateParams) (models.Pin, error) {
	url, err := repo.imgServ.UploadImage(context.Background(), &params.MediaSource)
	if err != nil {
		return models.Pin{}, errors.Wrap(pkgErrors.ErrImageService, err.Error())
	}

	avgColor := pkgImage.Color{
		Red:   constants.DefaultRedAvgColor,
		Green: constants.DefaultGreenAvgColor,
		Blue:  constants.DefaultBlueAvgColor,
	}
	img, err := pkgImage.BytesToImage(params.MediaSource.Bytes)
	if err == nil {
		avgColor = pkgImage.CalcAvgColor(img)
	}
	avgColorStr := fmt.Sprintf("rgb(%d, %d, %d)", avgColor.Red, avgColor.Green, avgColor.Blue)

	var pin models.Pin
	var authorID int
	var title, description, mediaSource sql.NullString
	err = repo.db.QueryRow(createCmd, params.Title, url, avgColorStr, params.Description, params.Author).
		Scan(&pin.Id, &title, &mediaSource, &pin.MediaSourceColor, &description, &authorID)
	if err != nil {
		return models.Pin{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}
	pin.Title = title.String
	pin.Description = description.String
	pin.MediaSource = mediaSource.String

	author, err := repo.getAuthor(authorID)
	if err != nil {
		return models.Pin{}, err
	}

	pin.Author = *author
	return pin, nil
}

const getCmd = `
	SELECT p.id,
			title,
			description,
			media_source,
			media_source_color,
			n_likes,
			u.id,
			u.username,
			u.name,
			u.profile_image,
			u.website_url
	FROM pins p
			 JOIN users u ON p.author_id = u.id
	WHERE p.id = $1;`

func (repo *repository) Get(id int) (models.Pin, error) {
	row := repo.db.QueryRow(getCmd, id)

	pin := models.Pin{}
	var title, description, mediaSource, profileImage, websiteUrl sql.NullString
	err := row.Scan(&pin.Id, &title, &description, &mediaSource, &pin.MediaSourceColor, &pin.NumLikes,
		&pin.Author.Id, &pin.Author.Username, &pin.Author.Name, &profileImage, &websiteUrl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Pin{}, errors.Wrap(pkgErrors.ErrPinNotFound, err.Error())
		}

		repo.log.Error(constants.DBScanError, zap.Error(err), zap.String("sql_query", getCmd),
			zap.Int("id", id))
		return models.Pin{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	pin.Title = title.String
	pin.Description = description.String
	pin.MediaSource = mediaSource.String
	pin.Author.ProfileImage = profileImage.String
	pin.Author.WebsiteUrl = websiteUrl.String
	return pin, nil
}

const listByUserCmd = `
	SELECT p.id,
		   title,
		   description,
		   media_source,
		   media_source_color,
		   n_likes,
		   u.id,
		   u.username,
		   u.name,
		   u.profile_image,
		   u.website_url
	FROM pins p
			 JOIN users u ON p.author_id = u.id AND u.id = $1
	ORDER BY created_at DESC
	LIMIT $2 OFFSET $3;`

func (repo *repository) ListByAuthor(userId int, page, limit int) ([]models.Pin, error) {
	rows, err := repo.db.Query(listByUserCmd, userId, limit, (page-1)*limit)
	if err != nil {
		return nil, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	var pins []models.Pin
	pin := models.Pin{}
	var title, description, mediaSource, profileImage, websiteUrl sql.NullString

	for rows.Next() {
		err = rows.Scan(&pin.Id, &title, &description, &mediaSource, &pin.MediaSourceColor, &pin.NumLikes,
			&pin.Author.Id, &pin.Author.Username, &pin.Author.Name, &profileImage, &websiteUrl)
		if err != nil {
			return nil, errors.Wrap(pkgErrors.ErrDb, err.Error())
		}

		pin.Title = title.String
		pin.Description = description.String
		pin.MediaSource = mediaSource.String
		pin.Author.ProfileImage = profileImage.String
		pin.Author.WebsiteUrl = websiteUrl.String
		pins = append(pins, pin)
	}

	return pins, nil
}

const listWithLikedFieldCmd = `
	SELECT p.id,
		   p.title,
		   p.description,
		   p.media_source,
		   p.media_source_color,
		   p.n_likes,
		   CASE WHEN pin_likes.author_id IS NOT NULL THEN true ELSE false END AS liked,
		   u.id,
		   u.username,
		   u.name,
		   u.profile_image,
		   u.website_url
	FROM pins p
			 LEFT JOIN pin_likes ON p.id = pin_likes.pin_id AND pin_likes.author_id = $1
			 JOIN users u ON p.author_id = u.id
	ORDER BY p.created_at DESC
	LIMIT $2 OFFSET $3;`

func (repo *repository) ListWithLikedField(userID int, page, limit int) ([]models.Pin, error) {
	rows, err := repo.db.Query(listWithLikedFieldCmd, userID, limit, (page-1)*limit)
	if err != nil {
		repo.log.Error(constants.DBQueryError, zap.Error(err), zap.String("sql_query", listWithLikedFieldCmd),
			zap.Int("page", page), zap.Int("limit", limit))
		return nil, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			repo.log.Error(constants.FailedCloseQueryRows, zap.Error(err),
				zap.String("sql_query", listWithLikedFieldCmd))
		}
	}()

	pins := []models.Pin{}
	pin := models.Pin{}
	var title, description, mediaSource, profileImage, websiteUrl sql.NullString

	for rows.Next() {
		err = rows.Scan(&pin.Id, &title, &description, &mediaSource, &pin.MediaSourceColor, &pin.NumLikes, &pin.Liked,
			&pin.Author.Id, &pin.Author.Username, &pin.Author.Name, &profileImage, &websiteUrl)
		if err != nil {
			repo.log.Error(constants.DBScanError, zap.Error(err), zap.String("sql_query", listWithLikedFieldCmd),
				zap.Int("page", page), zap.Int("limit", limit))
			return nil, errors.Wrap(pkgErrors.ErrDb, err.Error())
		}

		pin.Title = title.String
		pin.Description = description.String
		pin.MediaSource = mediaSource.String
		pin.Author.ProfileImage = profileImage.String
		pin.Author.WebsiteUrl = websiteUrl.String
		pins = append(pins, pin)
	}

	return pins, nil
}

const listCmd = `
	SELECT p.id,
		   p.title,
		   p.description,
		   p.media_source,
		   p.media_source_color,
		   p.n_likes,
		   false AS liked,
		   u.id,
		   u.username,
		   u.name,
		   u.profile_image,
		   u.website_url
	FROM pins p
			 JOIN users u ON p.author_id = u.id
	ORDER BY created_at DESC
	LIMIT $1 OFFSET $2;`

func (repo *repository) List(page, limit int) ([]models.Pin, error) {
	rows, err := repo.db.Query(listCmd, limit, (page-1)*limit)
	if err != nil {
		repo.log.Error(constants.DBQueryError, zap.Error(err), zap.String("sql_query", listCmd),
			zap.Int("page", page), zap.Int("limit", limit))
		return nil, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			repo.log.Error(constants.FailedCloseQueryRows, zap.Error(err), zap.String("sql_query", listCmd))
		}
	}()

	pins := []models.Pin{}
	pin := models.Pin{}
	var title, description, mediaSource, profileImage, websiteUrl sql.NullString

	for rows.Next() {
		err = rows.Scan(&pin.Id, &title, &description, &mediaSource, &pin.MediaSourceColor, &pin.NumLikes, &pin.Liked,
			&pin.Author.Id, &pin.Author.Username, &pin.Author.Name, &profileImage, &websiteUrl)
		if err != nil {
			repo.log.Error(constants.DBScanError, zap.Error(err), zap.String("sql_query", listCmd),
				zap.Int("page", page), zap.Int("limit", limit))
			return nil, errors.Wrap(pkgErrors.ErrDb, err.Error())
		}

		pin.Title = title.String
		pin.Description = description.String
		pin.MediaSource = mediaSource.String
		pin.Author.ProfileImage = profileImage.String
		pin.Author.WebsiteUrl = websiteUrl.String
		pins = append(pins, pin)
	}

	return pins, nil
}

const listLikedCmd = `
	SELECT p.id,
		   p.title,
		   p.description,
		   p.media_source,
		   p.media_source_color,
		   p.n_likes,
		   true AS liked,
		   u.id,
		   u.username,
		   u.name,
		   u.profile_image,
		   u.website_url
	FROM pins p
			 JOIN pin_likes pl ON p.id = pl.pin_id AND pl.author_id = $1
			 JOIN users u ON p.author_id = u.id
	ORDER BY pl.created_at DESC
	LIMIT $2 OFFSET $3;`

func (repo *repository) ListLiked(userID int, page, limit int) ([]models.Pin, error) {
	rows, err := repo.db.Query(listLikedCmd, userID, limit, (page-1)*limit)
	if err != nil {
		repo.log.Error(constants.DBQueryError, zap.Error(err), zap.String("sql_query", listLikedCmd),
			zap.Int("page", page), zap.Int("limit", limit))
		return nil, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			repo.log.Error(constants.FailedCloseQueryRows, zap.Error(err), zap.String("sql_query", listLikedCmd))
		}
	}()

	pins := []models.Pin{}
	pin := models.Pin{}
	var title, description, mediaSource, profileImage, websiteUrl sql.NullString

	for rows.Next() {
		err = rows.Scan(&pin.Id, &title, &description, &mediaSource, &pin.MediaSourceColor, &pin.NumLikes, &pin.Liked,
			&pin.Author.Id, &pin.Author.Username, &pin.Author.Name, &profileImage, &websiteUrl)
		if err != nil {
			repo.log.Error(constants.DBScanError, zap.Error(err), zap.String("sql_query", listLikedCmd),
				zap.Int("page", page), zap.Int("limit", limit))
			return nil, errors.Wrap(pkgErrors.ErrDb, err.Error())
		}

		pin.Title = title.String
		pin.Description = description.String
		pin.MediaSource = mediaSource.String
		pin.Author.ProfileImage = profileImage.String
		pin.Author.WebsiteUrl = websiteUrl.String
		pins = append(pins, pin)
	}

	return pins, nil
}

const fullUpdateCmd = `
		UPDATE pins
		SET title = $1::VARCHAR,
		description = $2::TEXT
		WHERE id = $3
		RETURNING id, title, description, media_source, media_source_color, author_id;`

func (repo *repository) FullUpdate(params *pkgPins.FullUpdateParams) (models.Pin, error) {
	var pin models.Pin
	var authorID int
	var title, description, mediaSource sql.NullString
	err := repo.db.QueryRow(fullUpdateCmd, params.Title, params.Description, params.Id).
		Scan(&pin.Id, &title, &description, &mediaSource, &pin.MediaSourceColor, &authorID)
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

	author, err := repo.getAuthor(authorID)
	if err != nil {
		return models.Pin{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	pin.Author = *author
	return pin, nil
}

const deleteCmd = `
		DELETE FROM pins 
		WHERE id = $1;`

func (repo *repository) Delete(id int) error {
	_, err := repo.db.Exec(deleteCmd, id)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	return nil
}

const isLikedByUserCmd = `
		SELECT EXISTS(SELECT pin_id
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

const checkWriteCmd = `
		SELECT EXISTS(SELECT id
		FROM pins
		WHERE id = $1 AND author_id = $2);`

func (repo *repository) CheckWriteAccess(userId, pinId string) (bool, error) {
	row := repo.db.QueryRow(checkWriteCmd, pinId, userId)

	var access bool
	err := row.Scan(&access)
	if err != nil {
		return false, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}
	return access, err
}

func (repo *repository) CheckReadAccess(_, _ string) (bool, error) {
	return true, nil
}

const getAuthorCmd = `
	SELECT id, username, name, profile_image, website_url
	FROM users
	WHERE id = $1;`

func (repo *repository) getAuthor(authorID int) (*models.Profile, error) {
	author := &models.Profile{}
	var profileImage, websiteUrl sql.NullString

	err := repo.db.QueryRow(getAuthorCmd, authorID).
		Scan(&author.Id, &author.Username, &author.Name, &profileImage, &websiteUrl)
	if err != nil {
		repo.log.Error(constants.DBScanError, zap.Error(err), zap.String("sql_query", getAuthorCmd),
			zap.Any("author_id", authorID))
		return nil, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	author.ProfileImage = profileImage.String
	author.WebsiteUrl = websiteUrl.String
	return author, nil
}
