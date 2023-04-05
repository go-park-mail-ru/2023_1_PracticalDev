package postgres

import (
	"database/sql"
	_boards "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/boards"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
)

type postgresRepository struct {
	db  *sql.DB
	log log.Logger
}

func NewPostgresRepository(db *sql.DB, log log.Logger) _boards.Repository {
	return &postgresRepository{db, log}
}

func (rep *postgresRepository) Create(params *_boards.CreateParams) (models.Board, error) {
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
		err = _boards.ErrDb
	}
	return createdBoard, err
}

func (rep *postgresRepository) List(userId int) ([]models.Board, error) {
	const getBoardsCommand = `SELECT * 
							  FROM boards
							  WHERE user_id = $1;`

	rows, err := rep.db.Query(getBoardsCommand, userId)
	if err != nil {
		return nil, _boards.ErrDb
	}

	var boards []models.Board
	board := models.Board{}
	var description sql.NullString
	for rows.Next() {
		err = rows.Scan(&board.Id, &board.Name, &description, &board.Privacy, &board.UserId)
		if err != nil {
			return nil, _boards.ErrDb
		}
		board.Description = description.String
		boards = append(boards, board)
	}
	return boards, nil
}

func (rep *postgresRepository) Get(id int) (models.Board, error) {
	const getBoardCommand = `SELECT *
							 FROM boards
							 WHERE id = $1;`

	row := rep.db.QueryRow(getBoardCommand, id)

	board := models.Board{}
	var description sql.NullString
	err := row.Scan(&board.Id, &board.Name, &description, &board.Privacy, &board.UserId)
	board.Description = description.String
	if err != nil {
		err = _boards.ErrDb
	}
	return board, err
}

func (rep *postgresRepository) FullUpdate(params *_boards.FullUpdateParams) (models.Board, error) {
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
		err = _boards.ErrDb
	}
	return updatedBoard, err
}

func (rep *postgresRepository) PartialUpdate(params *_boards.PartialUpdateParams) (models.Board, error) {
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
	return updatedBoard, err
}

func (rep *postgresRepository) Delete(id int) error {
	const deleteCommand = `DELETE FROM boards 
						   WHERE id = $1;`

	res, err := rep.db.Exec(deleteCommand, id)
	if err != nil {
		return _boards.ErrDeleteBoard
	}
	count, err := res.RowsAffected()
	if err != nil || count < 1 {
		return _boards.ErrBoardNotFound
	}
	return nil
}

func (rep *postgresRepository) CheckWriteAccess(userId, boardId string) (bool, error) {
	const checkCommand = `SELECT EXISTS(SELECT id
     			          				FROM boards
              			  				WHERE id = $1 AND user_id = $2);`

	row := rep.db.QueryRow(checkCommand,
		boardId,
		userId,
	)

	var access bool
	err := row.Scan(&access)
	return access, err
}

func (rep *postgresRepository) CheckReadAccess(userId, boardId string) (bool, error) {
	const checkCommand = `SELECT EXISTS(SELECT
              							FROM boards
              							WHERE id = $1 AND (privacy = 'public' OR user_id = $2));`

	row := rep.db.QueryRow(checkCommand, boardId, userId)

	var access bool
	err := row.Scan(&access)
	return access, err
}
