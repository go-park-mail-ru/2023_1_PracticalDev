package postgres

import (
	"database/sql"

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

const getPinsCmd = `SELECT id, title, description, media_source, n_likes, author_id FROM pins
                        WHERE to_tsquery($1) @@ to_tsvector(pins.title || pins.description);`

const getBoardsCmd = `SELECT * FROM boards
                        WHERE to_tsquery($1) @@ to_tsvector(boards.name);`

const getUsersCmd = `SELECT id, username, name, profile_image, website_url 
					    FROM users
                        WHERE to_tsquery($1) @@ to_tsvector(users.username)`

func (rep repository) Get(query string) (models.SearchRes, error) {
	rows, err := rep.db.Query(getPinsCmd, query)
	if err != nil {
		return models.SearchRes{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}
	var pins []models.Pin
	pin := models.Pin{}
	var title, description, mediaSource sql.NullString

	for rows.Next() {
		err = rows.Scan(&pin.Id, &title, &description, &mediaSource, &pin.NumLikes, &pin.Author)
		if err != nil {
			return models.SearchRes{}, errors.Wrap(pkgErrors.ErrDb,
				pkgErrors.ErrRepositoryQuery{
					Func:   "get",
					Query:  getPinsCmd,
					Params: []any{query},
					Err:    err,
				}.Error())
		}
		pin.Title = title.String
		pin.Description = description.String
		pin.MediaSource = mediaSource.String

		pins = append(pins, pin)
	}

	rows, err = rep.db.Query(getUsersCmd, query)
	if err != nil {
		return models.SearchRes{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	var users []models.Profile
	user := models.Profile{}
	var profileImage, websiteUrl sql.NullString

	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Username, &user.Name, &profileImage, &websiteUrl)
		if err != nil {
			return models.SearchRes{}, errors.Wrap(pkgErrors.ErrDb,
				pkgErrors.ErrRepositoryQuery{
					Func:   "get",
					Query:  getUsersCmd,
					Params: []any{query},
					Err:    err,
				}.Error())
		}
		user.ProfileImage = profileImage.String
		user.WebsiteUrl = websiteUrl.String

		users = append(users, user)
	}

	rows, err = rep.db.Query(getBoardsCmd, query)

	if err != nil {
		return models.SearchRes{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	var boards []models.Board
	board := models.Board{}

	for rows.Next() {
		err = rows.Scan(&board.Id, &board.Name, &description, &board.Privacy, &board.UserId)
		if err != nil {
			return models.SearchRes{}, errors.Wrap(pkgErrors.ErrDb,
				pkgErrors.ErrRepositoryQuery{
					Func:   "get",
					Query:  getUsersCmd,
					Params: []any{query},
					Err:    err,
				}.Error())
		}
		board.Description = description.String

		boards = append(boards, board)
	}
	return models.SearchRes{Pins: pins, Boards: boards, Users: users}, nil
}
