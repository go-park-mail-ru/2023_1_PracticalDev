package postgres

import (
	"database/sql"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/constants"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/search"
)

func NewRepository(db *sql.DB, log *zap.Logger) search.Repository {
	return &repository{db, log}
}

type repository struct {
	db  *sql.DB
	log *zap.Logger
}

const getPinsCmd = `
		SELECT id, title, description, media_source, n_likes, author_id
		FROM pins
		WHERE websearch_to_tsquery('russian', $1) @@ to_tsvector('russian', title)
		   OR websearch_to_tsquery($1) @@ to_tsvector(title)
		   OR lower(title) LIKE lower('%' || $1 || '%')
		ORDER BY ts_rank(to_tsvector('russian', title), websearch_to_tsquery('russian', $1)),
				 ts_rank(to_tsvector(title), websearch_to_tsquery($1)) DESC;`

const getBoardsCmd = `
		SELECT *
		FROM boards
		WHERE websearch_to_tsquery('russian', $1) @@ to_tsvector('russian', name)
		   OR websearch_to_tsquery($1) @@ to_tsvector(name)
		   OR lower(name) LIKE lower('%' || $1 || '%')
		ORDER BY ts_rank(to_tsvector('russian', name), websearch_to_tsquery('russian', $1)),
				 ts_rank(to_tsvector(name), websearch_to_tsquery($1)) DESC;`

const getUsersCmd = `
		SELECT id, username, name, profile_image, website_url
		FROM users
		WHERE websearch_to_tsquery('russian', $1) @@ (to_tsvector('russian', username) || to_tsvector('russian', name))
		   OR websearch_to_tsquery($1) @@ (to_tsvector(username) || to_tsvector(name))
		   OR lower(username) || lower(name) LIKE lower('%' || $1 || '%')
		ORDER BY ts_rank(to_tsvector('russian', username) || to_tsvector('russian', name),
						 websearch_to_tsquery('russian', $1)),
				 ts_rank(to_tsvector(username) || to_tsvector(name), websearch_to_tsquery($1)) DESC;`

func (rep repository) Get(query string) (models.SearchRes, error) {
	rows, err := rep.db.Query(getPinsCmd, query)
	if err != nil {
		rep.log.Error(constants.DBQueryError, zap.String("sql_query", getPinsCmd),
			zap.String("search_query", query), zap.Error(err))
		return models.SearchRes{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	var pins []models.Pin
	pin := models.Pin{}
	var title, description, mediaSource sql.NullString
	for rows.Next() {
		err = rows.Scan(&pin.Id, &title, &description, &mediaSource, &pin.NumLikes, &pin.Author)
		if err != nil {
			rep.log.Error(constants.DBScanError, zap.String("sql_query", getPinsCmd),
				zap.String("search_query", query), zap.Error(err))
			return models.SearchRes{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
		}

		pin.Title = title.String
		pin.Description = description.String
		pin.MediaSource = mediaSource.String
		pins = append(pins, pin)
	}

	rows, err = rep.db.Query(getUsersCmd, query)
	if err != nil {
		rep.log.Error(constants.DBQueryError, zap.String("sql_query", getUsersCmd),
			zap.String("search_query", query), zap.Error(err))
		return models.SearchRes{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	var users []models.Profile
	user := models.Profile{}
	var profileImage, websiteUrl sql.NullString
	for rows.Next() {
		err = rows.Scan(&user.Id, &user.Username, &user.Name, &profileImage, &websiteUrl)
		if err != nil {
			rep.log.Error(constants.DBScanError, zap.String("sql_query", getUsersCmd),
				zap.String("search_query", query), zap.Error(err))
			return models.SearchRes{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
		}

		user.ProfileImage = profileImage.String
		user.WebsiteUrl = websiteUrl.String
		users = append(users, user)
	}

	rows, err = rep.db.Query(getBoardsCmd, query)
	if err != nil {
		rep.log.Error(constants.DBQueryError, zap.String("sql_query", getBoardsCmd),
			zap.String("search_query", query), zap.Error(err))
		return models.SearchRes{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	var boards []models.Board
	board := models.Board{}
	for rows.Next() {
		err = rows.Scan(&board.Id, &board.Name, &description, &board.Privacy, &board.UserId)
		if err != nil {
			rep.log.Error(constants.DBScanError, zap.String("sql_query", getBoardsCmd),
				zap.String("search_query", query), zap.Error(err))
			return models.SearchRes{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
		}

		board.Description = description.String
		boards = append(boards, board)
	}

	return models.SearchRes{Pins: pins, Boards: boards, Users: users}, nil
}
