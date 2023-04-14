package postgres

import (
	"database/sql"
	"github.com/pkg/errors"

	pkgBoards "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/boards"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
)

type repository struct {
	db  *sql.DB
	log log.Logger
}

func NewPostgresRepository(db *sql.DB, log log.Logger) pkgBoards.Repository {
	return &repository{db, log}
}

func (rep *repository) Create(params *pkgBoards.CreateParams) (models.Board, error) {
	const insertCommand = `INSERT INTO boards (name, description, privacy, user_id) 
				      	   VALUES ($1, $2, $3, $4)
						   RETURNING *;`

	row := rep.db.QueryRow(insertCommand,
		params.Name,
		params.Description,
		params.Privacy,
		params.UserId,
	)
	createdBoard := models.Board{}
	var description sql.NullString
	err := row.Scan(&createdBoard.Id, &createdBoard.Name, &description, &createdBoard.Privacy, &createdBoard.UserId)
	createdBoard.Description = description.String
	if err != nil {
		err = pkgBoards.ErrDb
	}
	return createdBoard, err
}

func (rep *repository) List(userId int) ([]models.Board, error) {
	const getBoardsCommand = `SELECT * 
							  FROM boards
							  WHERE user_id = $1;`

	rows, err := rep.db.Query(getBoardsCommand, userId)
	if err != nil {
		return nil, pkgBoards.ErrDb
	}

	var boards []models.Board
	board := models.Board{}
	var description sql.NullString
	for rows.Next() {
		err = rows.Scan(&board.Id, &board.Name, &description, &board.Privacy, &board.UserId)
		if err != nil {
			return nil, pkgBoards.ErrDb
		}
		board.Description = description.String
		boards = append(boards, board)
	}
	return boards, nil
}

func (rep *repository) Get(id int) (models.Board, error) {
	const getBoardCommand = `SELECT *
							 FROM boards
							 WHERE id = $1;`

	row := rep.db.QueryRow(getBoardCommand, id)

	board := models.Board{}
	var description sql.NullString
	err := row.Scan(&board.Id, &board.Name, &description, &board.Privacy, &board.UserId)
	board.Description = description.String
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Board{}, pkgBoards.ErrBoardNotFound
		} else {
			return models.Board{}, pkgBoards.ErrDb
		}
	}
	return board, err
}

func (rep *repository) FullUpdate(params *pkgBoards.FullUpdateParams) (models.Board, error) {
	const fullUpdateBoard = `UPDATE boards
								SET name = $1::VARCHAR,
    							description = $2::TEXT,
    							privacy = $3::privacy
								WHERE id = $4
								RETURNING *;`

	row := rep.db.QueryRow(fullUpdateBoard,
		params.Name,
		params.Description,
		params.Privacy,
		params.Id,
	)
	var updatedBoard models.Board
	var description sql.NullString
	err := row.Scan(&updatedBoard.Id, &updatedBoard.Name, &description, &updatedBoard.Privacy, &updatedBoard.UserId)
	updatedBoard.Description = description.String
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Board{}, pkgBoards.ErrBoardNotFound
		} else {
			return models.Board{}, pkgBoards.ErrDb
		}
	}
	return updatedBoard, err
}

func (rep *repository) PartialUpdate(params *pkgBoards.PartialUpdateParams) (models.Board, error) {
	const partialUpdateBoard = `UPDATE boards
								SET name = CASE WHEN $1::boolean THEN $2::VARCHAR ELSE name END,
    							description = CASE WHEN $3::boolean THEN $4::TEXT ELSE description END,
    							privacy = CASE WHEN $5::boolean THEN $6::privacy ELSE privacy END
								WHERE id = $7
								RETURNING *;`

	row := rep.db.QueryRow(partialUpdateBoard,
		params.UpdateName,
		params.Name,
		params.UpdateDescription,
		params.Description,
		params.UpdatePrivacy,
		params.Privacy,
		params.Id,
	)
	var updatedBoard models.Board
	var description sql.NullString
	err := row.Scan(&updatedBoard.Id, &updatedBoard.Name, &description, &updatedBoard.Privacy, &updatedBoard.UserId)
	updatedBoard.Description = description.String
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Board{}, pkgBoards.ErrBoardNotFound
		} else {
			return models.Board{}, pkgBoards.ErrDb
		}
	}
	return updatedBoard, err
}

func (rep *repository) Delete(id int) error {
	const deleteCommand = `DELETE FROM boards 
						   WHERE id = $1;`

	res, err := rep.db.Exec(deleteCommand, id)
	if err != nil {
		return pkgBoards.ErrDb
	}
	count, err := res.RowsAffected()
	if err != nil || count < 1 {
		return pkgBoards.ErrBoardNotFound
	}
	return nil
}

const addToBoardCmd = `INSERT INTO boards_pins(pin_id, board_id)
						VALUES($1, $2);`

func (rep *repository) AddPin(boardId, pinId int) error {
	res, err := rep.db.Exec(addToBoardCmd, pinId, boardId)
	if err != nil {
		return pkgBoards.ErrDb
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected < 1 {
		return pkgBoards.ErrDb
	}
	return nil
}

const listByBoardCmd = `SELECT pins.id, title, description, media_source, author_id 
						FROM pins 
						JOIN boards_pins AS b
						ON b.board_id = $1 AND b.pin_id = pins.id
						ORDER BY created_at DESC 
						LIMIT $2 OFFSET $3;`

func (rep *repository) PinsList(boardId int, page, limit int) ([]models.Pin, error) {
	rows, err := rep.db.Query(listByBoardCmd, boardId, limit, (page-1)*limit)
	if err != nil {
		return nil, pkgBoards.ErrDb
	}

	var pins []models.Pin
	retrievedPin := models.Pin{}
	var title, description, mediaSource sql.NullString
	for rows.Next() {
		err = rows.Scan(&retrievedPin.Id, &title, &description, &mediaSource, &retrievedPin.Author)
		if err != nil {
			return nil, pkgBoards.ErrDb
		}
		retrievedPin.Title = title.String
		retrievedPin.Description = description.String
		retrievedPin.MediaSource = mediaSource.String
		pins = append(pins, retrievedPin)
	}
	return pins, nil
}

const deleteFromBoardCmd = `DELETE FROM boards_pins
							WHERE pin_id = $1 AND board_id = $2;`

func (rep *repository) RemovePin(boardId, pinId int) error {
	res, err := rep.db.Exec(deleteFromBoardCmd, pinId, boardId)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil || count < 1 {
		return err
	}
	return nil
}

const hasPinCmd = `SELECT EXISTS(SELECT pin_id
									FROM boards_pins
									WHERE board_id = $1 AND pin_id = $2) AS has;`

func (rep *repository) HasPin(boardId, pinId int) (bool, error) {
	row := rep.db.QueryRow(hasPinCmd, boardId, pinId)

	var has bool
	err := row.Scan(&has)
	if err != nil {
		return false, pkgBoards.ErrDb
	}
	return has, nil
}

const checkWriteCommand = `SELECT EXISTS(SELECT id
     			          				FROM boards
              			  				WHERE id = $1 AND user_id = $2) AS access;`

func (rep *repository) CheckWriteAccess(userId, boardId string) (bool, error) {
	row := rep.db.QueryRow(checkWriteCommand,
		boardId,
		userId,
	)

	var access bool
	err := row.Scan(&access)
	if err != nil {
		err = pkgBoards.ErrDb
	}
	return access, err
}

const checkReadCommand = `SELECT EXISTS(SELECT
              							FROM boards
              							WHERE id = $1 AND (privacy = 'public' OR user_id = $2)) AS access;`

func (rep *repository) CheckReadAccess(userId, boardId string) (bool, error) {
	row := rep.db.QueryRow(checkReadCommand, boardId, userId)

	var access bool
	err := row.Scan(&access)
	if err != nil {
		err = pkgBoards.ErrDb
	}
	return access, err
}
