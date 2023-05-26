package postgres

import (
	"database/sql"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/constants"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/utils"

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
		SELECT p.id,
			   p.title,
			   p.description,
			   p.media_source,
			   p.media_source_color,
			   p.n_likes,
			   u.id,
			   u.username,
			   u.name,
			   u.profile_image,
			   u.website_url
		FROM pins p
			JOIN users u ON p.author_id = u.id
		WHERE websearch_to_tsquery('russian', $1) @@ to_tsvector('russian', p.title)
		   OR websearch_to_tsquery($1) @@ to_tsvector(p.title)
		   OR lower(p.title) LIKE lower('%' || $1 || '%')
		ORDER BY ts_rank(to_tsvector('russian', p.title), websearch_to_tsquery('russian', $1)),
				 ts_rank(to_tsvector(p.title), websearch_to_tsquery($1)) DESC;`

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

func (rep repository) Search(query string) (models.SearchRes, error) {
	rows, err := rep.db.Query(getPinsCmd, query)
	if err != nil {
		rep.log.Error(constants.DBQueryError, zap.String("sql_query", getPinsCmd),
			zap.String("search_query", query), zap.Error(err))
		return models.SearchRes{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	var pins []models.Pin
	pin := models.Pin{}
	var title, description, mediaSource, profileImage, websiteUrl sql.NullString
	for rows.Next() {
		err = rows.Scan(&pin.Id, &title, &description, &mediaSource, &pin.MediaSourceColor, &pin.NumLikes,
			&pin.Author.Id, &pin.Author.Username, &pin.Author.Name, &profileImage, &websiteUrl)
		if err != nil {
			rep.log.Error(constants.DBScanError, zap.String("sql_query", getPinsCmd),
				zap.String("search_query", query), zap.Error(err))
			return models.SearchRes{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
		}

		pin.Title = title.String
		pin.Description = description.String
		pin.MediaSource = mediaSource.String
		pin.Author.ProfileImage = profileImage.String
		pin.Author.WebsiteUrl = websiteUrl.String
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

const pinsSuggestionsCmd = `
	SELECT title
	FROM pins
	WHERE websearch_to_tsquery('russian', $1) @@ to_tsvector('russian', title)
	   OR websearch_to_tsquery($1) @@ to_tsvector(title)
	   OR lower(title) LIKE lower('%' || $1 || '%')
	ORDER BY ts_rank(to_tsvector('russian', title), websearch_to_tsquery('russian', $1)),
			 ts_rank(to_tsvector(title), websearch_to_tsquery($1)) DESC;`

const boardsSuggestionsCmd = `
		SELECT name
		FROM boards
		WHERE websearch_to_tsquery('russian', $1) @@ to_tsvector('russian', name)
		   OR websearch_to_tsquery($1) @@ to_tsvector(name)
		   OR lower(name) LIKE lower('%' || $1 || '%')
		ORDER BY ts_rank(to_tsvector('russian', name), websearch_to_tsquery('russian', $1)),
				 ts_rank(to_tsvector(name), websearch_to_tsquery($1)) DESC;`

const usersSuggestionsCmd = `
		SELECT username, name
		FROM users
		WHERE websearch_to_tsquery('russian', $1) @@ (to_tsvector('russian', username) || to_tsvector('russian', name))
		   OR websearch_to_tsquery($1) @@ (to_tsvector(username) || to_tsvector(name))
		   OR lower(username) || lower(name) LIKE lower('%' || $1 || '%')
		ORDER BY ts_rank(to_tsvector('russian', username) || to_tsvector('russian', name),
						 websearch_to_tsquery('russian', $1)),
				 ts_rank(to_tsvector(username) || to_tsvector(name), websearch_to_tsquery($1)) DESC;`

func (rep repository) Suggestions(query string) ([]string, error) {
	rows, err := rep.db.Query(pinsSuggestionsCmd, query)
	if err != nil {
		rep.log.Error(constants.DBQueryError, zap.String("sql_query", pinsSuggestionsCmd),
			zap.String("search_query", query), zap.Error(err))
		return nil, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	suggestions := []string{}
	var title sql.NullString
	for rows.Next() {
		err = rows.Scan(&title)
		if err != nil {
			rep.log.Error(constants.DBScanError, zap.String("sql_query", pinsSuggestionsCmd),
				zap.String("search_query", query), zap.Error(err))
			return nil, errors.Wrap(pkgErrors.ErrDb, err.Error())
		}

		suggestions = append(suggestions, title.String)
	}

	rows, err = rep.db.Query(usersSuggestionsCmd, query)
	if err != nil {
		rep.log.Error(constants.DBQueryError, zap.String("sql_query", usersSuggestionsCmd),
			zap.String("search_query", query), zap.Error(err))
		return nil, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	var username, name string
	for rows.Next() {
		err = rows.Scan(&username, &name)
		if err != nil {
			rep.log.Error(constants.DBScanError, zap.String("sql_query", usersSuggestionsCmd),
				zap.String("search_query", query), zap.Error(err))
			return nil, errors.Wrap(pkgErrors.ErrDb, err.Error())
		}

		suggestions = append(suggestions, username, name)
	}

	rows, err = rep.db.Query(boardsSuggestionsCmd, query)
	if err != nil {
		rep.log.Error(constants.DBQueryError, zap.String("sql_query", boardsSuggestionsCmd),
			zap.String("search_query", query), zap.Error(err))
		return nil, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	var boardName string
	for rows.Next() {
		err = rows.Scan(&boardName)
		if err != nil {
			rep.log.Error(constants.DBScanError, zap.String("sql_query", boardsSuggestionsCmd),
				zap.String("search_query", query), zap.Error(err))
			return nil, errors.Wrap(pkgErrors.ErrDb, err.Error())
		}

		suggestions = append(suggestions, boardName)
	}

	suggestions = utils.RemoveDuplicates(suggestions)
	if len(suggestions) < 10 {
		return suggestions, nil
	}
	return suggestions[:10], nil
}
