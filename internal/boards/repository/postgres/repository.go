package postgres

import (
	"database/sql"

	"github.com/pkg/errors"

	pkgBoards "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/boards"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/log"
)

type repository struct {
	db  *sql.DB
	log log.Logger
}

func NewPostgresRepository(db *sql.DB, log log.Logger) pkgBoards.Repository {
	return &repository{db, log}
}

const insertCommand = `INSERT INTO boards (name, description, privacy, user_id) 
				      	   VALUES ($1, $2, $3, $4)
						   RETURNING *;`

func (rep *repository) Create(params *pkgBoards.CreateParams) (models.Board, error) {
	const fnCreate = "Create"

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
		return models.Board{}, errors.Wrap(pkgErrors.ErrDb,
			pkgErrors.ErrRepositoryQuery{
				Func:   fnCreate,
				Query:  insertCommand,
				Params: []any{params.Name, params.Description, params.Privacy, params.UserId},
				Err:    err,
			}.Error())
	}

	return createdBoard, nil
}

const getBoardsCommand = `SELECT * 
							  FROM boards
							  WHERE user_id = $1;`

func (rep *repository) List(userId int) ([]models.Board, error) {
	const fnList = "List"

	rows, err := rep.db.Query(getBoardsCommand, userId)
	if err != nil {
		return nil, errors.Wrap(pkgErrors.ErrDb,
			pkgErrors.ErrRepositoryQuery{
				Func:   fnList,
				Query:  getBoardsCommand,
				Params: []any{userId},
				Err:    err,
			}.Error())
	}

	var boards []models.Board
	board := models.Board{}
	var description sql.NullString
	for rows.Next() {
		err = rows.Scan(&board.Id, &board.Name, &description, &board.Privacy, &board.UserId)
		if err != nil {
			return nil, errors.Wrap(pkgErrors.ErrDb,
				pkgErrors.ErrRepositoryQuery{
					Func:   fnList,
					Query:  getBoardsCommand,
					Params: []any{userId},
					Err:    err,
				}.Error())
		}

		board.Description = description.String
		boards = append(boards, board)
	}
	return boards, nil
}

const getCmd = `SELECT *
				 FROM boards
				 WHERE id = $1;`

func (rep *repository) Get(id int) (models.Board, error) {
	const fnGet = "Get"

	row := rep.db.QueryRow(getCmd, id)

	board := models.Board{}
	var description sql.NullString

	err := row.Scan(&board.Id, &board.Name, &description, &board.Privacy, &board.UserId)
	board.Description = description.String

	if err != nil {
		errRepo := pkgErrors.ErrRepositoryQuery{
			Func:   fnGet,
			Query:  getCmd,
			Params: []any{id},
			Err:    err,
		}

		if errors.Is(err, sql.ErrNoRows) {
			return models.Board{}, errors.Wrap(pkgErrors.ErrBoardNotFound, errRepo.Error())
		} else {
			return models.Board{}, errors.Wrap(pkgErrors.ErrDb, errRepo.Error())
		}
	}

	return board, nil
}

const fullUpdateCmd = `UPDATE boards
								SET name = $1::VARCHAR,
    							description = $2::TEXT,
    							privacy = $3::privacy
								WHERE id = $4
								RETURNING *;`

func (rep *repository) FullUpdate(params *pkgBoards.FullUpdateParams) (models.Board, error) {
	const fnFullUpdate = "FullUpdate"

	row := rep.db.QueryRow(fullUpdateCmd,
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
		errRepo := pkgErrors.ErrRepositoryQuery{
			Func:   fnFullUpdate,
			Query:  fullUpdateCmd,
			Params: []any{params.Name, params.Description, params.Privacy, params.Id},
			Err:    err,
		}

		if errors.Is(err, sql.ErrNoRows) {
			return models.Board{}, errors.Wrap(pkgErrors.ErrBoardNotFound, errRepo.Error())
		} else {
			return models.Board{}, errors.Wrap(pkgErrors.ErrDb, errRepo.Error())
		}
	}

	return updatedBoard, err
}

const partialUpdateCmd = `UPDATE boards
								SET name = CASE WHEN $1::boolean THEN $2::VARCHAR ELSE name END,
    							description = CASE WHEN $3::boolean THEN $4::TEXT ELSE description END,
    							privacy = CASE WHEN $5::boolean THEN $6::privacy ELSE privacy END
								WHERE id = $7
								RETURNING *;`

func (rep *repository) PartialUpdate(params *pkgBoards.PartialUpdateParams) (models.Board, error) {
	const fnPartialUpdate = "PartialUpdate"

	row := rep.db.QueryRow(partialUpdateCmd,
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
		errRepo := pkgErrors.ErrRepositoryQuery{
			Func:   fnPartialUpdate,
			Query:  partialUpdateCmd,
			Params: []any{params.Name, params.Description, params.Privacy, params.Id},
			Err:    err,
		}

		if errors.Is(err, sql.ErrNoRows) {
			return models.Board{}, errors.Wrap(pkgErrors.ErrBoardNotFound, errRepo.Error())
		} else {
			return models.Board{}, errors.Wrap(pkgErrors.ErrDb, errRepo.Error())
		}
	}

	return updatedBoard, err
}

const deleteCmd = `DELETE FROM boards 
				   WHERE id = $1;`

func (rep *repository) Delete(id int) error {
	const fnDelete = "Delete"

	_, err := rep.db.Exec(deleteCmd, id)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrDb,
			pkgErrors.ErrRepositoryQuery{
				Func:   fnDelete,
				Query:  deleteCmd,
				Params: []any{id},
				Err:    err,
			}.Error())
	}

	return nil
}

const AddPinCmd = `INSERT INTO boards_pins(pin_id, board_id)
						VALUES($1, $2);`

func (rep *repository) AddPin(boardId, pinId int) error {
	const fnAddPin = "AddPin"

	_, err := rep.db.Exec(AddPinCmd, pinId, boardId)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrDb,
			pkgErrors.ErrRepositoryQuery{
				Func:   fnAddPin,
				Query:  AddPinCmd,
				Params: []any{pinId, boardId},
				Err:    err,
			}.Error())
	}

	return nil
}

const pinsListCmd = `SELECT pins.id, title, description, media_source, author_id 
						FROM pins 
						JOIN boards_pins AS b
						ON b.board_id = $1 AND b.pin_id = pins.id
						ORDER BY created_at DESC 
						LIMIT $2 OFFSET $3;`

func (rep *repository) PinsList(boardId int, page, limit int) ([]models.Pin, error) {
	const fnPinsList = "PinsList"

	rows, err := rep.db.Query(pinsListCmd, boardId, limit, (page-1)*limit)
	if err != nil {
		return nil, errors.Wrap(pkgErrors.ErrDb,
			pkgErrors.ErrRepositoryQuery{
				Func:   fnPinsList,
				Query:  pinsListCmd,
				Params: []any{boardId, limit, (page - 1) * limit},
				Err:    err,
			}.Error())
	}

	var pins []models.Pin
	retrievedPin := models.Pin{}
	var title, description, mediaSource sql.NullString

	for rows.Next() {
		err = rows.Scan(&retrievedPin.Id, &title, &description, &mediaSource, &retrievedPin.Author)
		if err != nil {
			return nil, errors.Wrap(pkgErrors.ErrDb,
				pkgErrors.ErrRepositoryQuery{
					Func:   fnPinsList,
					Query:  pinsListCmd,
					Params: []any{boardId, limit, (page - 1) * limit},
					Err:    err,
				}.Error())
		}
		retrievedPin.Title = title.String
		retrievedPin.Description = description.String
		retrievedPin.MediaSource = mediaSource.String
		pins = append(pins, retrievedPin)
	}

	return pins, nil
}

const RemovePinCmd = `DELETE FROM boards_pins
						WHERE pin_id = $1 AND board_id = $2;`

func (rep *repository) RemovePin(boardId, pinId int) error {
	const fnRemovePin = "RemovePin"

	_, err := rep.db.Exec(RemovePinCmd, pinId, boardId)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrDb,
			pkgErrors.ErrRepositoryQuery{
				Func:   fnRemovePin,
				Query:  RemovePinCmd,
				Params: []any{pinId, boardId},
				Err:    err,
			}.Error())
	}

	return nil
}

const hasPinCmd = `SELECT EXISTS(SELECT pin_id
									FROM boards_pins
									WHERE board_id = $1 AND pin_id = $2) AS has;`

func (rep *repository) HasPin(boardId, pinId int) (bool, error) {
	const fnHasPin = "HasPin"

	row := rep.db.QueryRow(hasPinCmd, boardId, pinId)

	var has bool
	err := row.Scan(&has)
	if err != nil {
		return false, errors.Wrap(pkgErrors.ErrDb,
			pkgErrors.ErrRepositoryQuery{
				Func:   fnHasPin,
				Query:  hasPinCmd,
				Params: []any{boardId, pinId},
				Err:    err,
			}.Error())
	}

	return has, nil
}

const checkWriteCommand = `SELECT EXISTS(SELECT id
     			          				FROM boards
              			  				WHERE id = $1 AND user_id = $2) AS access;`

func (rep *repository) CheckWriteAccess(userId, boardId string) (bool, error) {
	const fnCheckWriteAccess = "CheckWriteAccess"

	row := rep.db.QueryRow(checkWriteCommand, boardId, userId)

	var access bool
	err := row.Scan(&access)
	if err != nil {
		return false, errors.Wrap(pkgErrors.ErrDb,
			pkgErrors.ErrRepositoryQuery{
				Func:   fnCheckWriteAccess,
				Query:  checkWriteCommand,
				Params: []any{boardId, userId},
				Err:    err,
			}.Error())
	}

	return access, nil
}

const checkReadCommand = `SELECT EXISTS(SELECT
              							FROM boards
              							WHERE id = $1 AND (privacy = 'public' OR user_id = $2)) AS access;`

func (rep *repository) CheckReadAccess(userId, boardId string) (bool, error) {
	const fnCheckReadAccess = "CheckReadAccess"

	row := rep.db.QueryRow(checkReadCommand, boardId, userId)

	var access bool
	err := row.Scan(&access)
	if err != nil {
		return false, errors.Wrap(pkgErrors.ErrDb,
			pkgErrors.ErrRepositoryQuery{
				Func:   fnCheckReadAccess,
				Query:  checkReadCommand,
				Params: []any{boardId, userId},
				Err:    err,
			}.Error())
	}

	return access, nil
}
