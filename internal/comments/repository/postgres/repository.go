package postgres

import (
	"database/sql"
	pkgComments "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/comments"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/constants"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type repository struct {
	db  *sql.DB
	log *zap.Logger
}

func NewRepository(db *sql.DB, log *zap.Logger) pkgComments.Repository {
	return &repository{db, log}
}

const createCmd = `
		INSERT INTO comments (author_id, pin_id, text)
		VALUES ($1, $2, $3)
		RETURNING *;`

func (rep *repository) Create(params *pkgComments.CreateParams) (models.Comment, error) {
	row := rep.db.QueryRow(createCmd, params.AuthorID, params.PinID, params.Text)

	comment := models.Comment{}
	var authorID int
	err := row.Scan(&comment.ID, &authorID, &comment.PinID, &comment.Text, &comment.CreatedAt)
	if err != nil {
		rep.log.Error(constants.DBScanError, zap.Error(err), zap.String("sql_query", createCmd),
			zap.Any("create_params", params))
		return models.Comment{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	author, err := rep.getAuthor(authorID)
	if err != nil {
		return models.Comment{}, err
	}

	comment.Author = *author
	return comment, nil
}

const listCmd = `
	SELECT c.id,
		   c.pin_id,
		   c.text,
		   c.created_at,
		   u.id,
		   u.username,
		   u.name,
		   u.profile_image,
		   u.website_url
	FROM comments c
    	     JOIN users u on u.id = c.author_id
	WHERE pin_id = $1
	ORDER BY created_at;`

func (rep *repository) List(pinID int) ([]models.Comment, error) {
	rows, err := rep.db.Query(listCmd, pinID)
	if err != nil {
		rep.log.Error(constants.DBQueryError, zap.Error(err), zap.String("sql_query", listCmd),
			zap.Int("pin_id", pinID))
		return nil, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	comments := []models.Comment{}
	comment := models.Comment{}
	var profileImage, websiteUrl sql.NullString
	for rows.Next() {
		err = rows.Scan(&comment.ID, &comment.PinID, &comment.Text, &comment.CreatedAt,
			&comment.Author.Id, &comment.Author.Username, &comment.Author.Name, &profileImage, &websiteUrl)
		if err != nil {
			rep.log.Error(constants.DBScanError, zap.Error(err), zap.String("sql_query", listCmd),
				zap.Int("pin_id", pinID))
			return nil, errors.Wrap(pkgErrors.ErrDb, err.Error())
		}

		comment.Author.ProfileImage = profileImage.String
		comment.Author.WebsiteUrl = websiteUrl.String
		comments = append(comments, comment)
	}

	return comments, nil
}

const getAuthorCmd = `
	SELECT id, username, name, profile_image, website_url
	FROM users
	WHERE id = $1;`

func (rep *repository) getAuthor(authorID int) (*models.Profile, error) {
	author := &models.Profile{}
	var profileImage, websiteUrl sql.NullString

	err := rep.db.QueryRow(getAuthorCmd, authorID).
		Scan(&author.Id, &author.Username, &author.Name, &profileImage, &websiteUrl)
	if err != nil {
		rep.log.Error(constants.DBScanError, zap.Error(err), zap.String("sql_query", getAuthorCmd),
			zap.Any("author_id", authorID))
		return nil, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	author.ProfileImage = profileImage.String
	author.WebsiteUrl = websiteUrl.String
	return author, nil
}
