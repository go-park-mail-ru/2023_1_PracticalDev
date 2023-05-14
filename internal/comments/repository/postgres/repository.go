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
	err := row.Scan(&comment.ID, &comment.AuthorID, &comment.PinID, &comment.Text, &comment.CreatedAt)
	if err != nil {
		rep.log.Error(constants.DBScanError, zap.String("sql_query", createCmd),
			zap.Any("create_params", params), zap.Error(err))
		return models.Comment{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	return comment, nil
}

const listCmd = `
		SELECT *
		FROM comments
		WHERE pin_id = $1
		ORDER BY created_at;`

func (rep *repository) List(pinID int) ([]models.Comment, error) {
	rows, err := rep.db.Query(listCmd, pinID)
	if err != nil {
		rep.log.Error(constants.DBQueryError, zap.String("sql_query", listCmd),
			zap.Int("pin_id", pinID), zap.Error(err))
		return nil, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	comments := []models.Comment{}
	comment := models.Comment{}
	for rows.Next() {
		err = rows.Scan(&comment.ID, &comment.AuthorID, &comment.PinID, &comment.Text, &comment.CreatedAt)
		if err != nil {
			rep.log.Error(constants.DBScanError, zap.String("sql_query", listCmd),
				zap.Int("pin_id", pinID), zap.Error(err))
			return nil, errors.Wrap(pkgErrors.ErrDb, err.Error())
		}

		comments = append(comments, comment)
	}

	return comments, nil
}
