package boards

import (
	"database/sql"
	"errors"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
)

type createParams struct {
	Name        string
	Description string
	Privacy     string
	UserId      int
}

type partialUpdateParams struct {
	Id                int
	Name              string
	UpdateName        bool
	Description       string
	UpdateDescription bool
	Privacy           string
	UpdatePrivacy     bool
}

type fullUpdateParams struct {
	Id          int
	Name        string
	Description string
	Privacy     string
}

var (
	ErrDeleteBoard   = errors.New("failed to delete board")
	ErrBoardNotFound = errors.New("board not found")
)

type Repository interface {
	Create(params *createParams) (models.Board, error)
	List(userId int) ([]models.Board, error)
	Get(id int) (models.Board, error)
	FullUpdate(params *fullUpdateParams) (models.Board, error)
	PartialUpdate(params *partialUpdateParams) (models.Board, error)
	Delete(id int) error
	CheckWriteAccess(userId, boardId string) (bool, error)
	CheckReadAccess(userId, boardId string) (bool, error)
}

func NewRepository(db *sql.DB, log log.Logger) Repository {
	return &repository{db, log}
}

type repository struct {
	db  *sql.DB
	log log.Logger
}

func (rep *repository) Create(params *createParams) (models.Board, error) {
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
	return createdBoard, err
}

func (rep *repository) List(userId int) ([]models.Board, error) {
	const getBoardsCommand = `SELECT * 
							  FROM boards
							  WHERE user_id = $1;`

	boards := []models.Board{}
	rows, err := rep.db.Query(getBoardsCommand, userId)
	if err != nil {
		return boards, err
	}

	board := models.Board{}
	var description sql.NullString
	for rows.Next() {
		err = rows.Scan(&board.Id, &board.Name, &description, &board.Privacy, &board.UserId)
		if err != nil {
			break
		}
		board.Description = description.String
		boards = append(boards, board)
	}
	return boards, err
}

func (rep *repository) Get(id int) (models.Board, error) {
	const getBoardCommand = `SELECT *
							 FROM boards
							 WHERE id = $1;`

	board := models.Board{}
	rows, err := rep.db.Query(getBoardCommand, id)
	if err != nil {
		return board, err
	}

	var description sql.NullString
	rows.Next()
	err = rows.Scan(&board.Id, &board.Name, &description, &board.Privacy, &board.UserId)
	board.Description = description.String
	return board, err
}

func (rep *repository) FullUpdate(params *fullUpdateParams) (models.Board, error) {
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
	return updatedBoard, err
}

func (rep *repository) PartialUpdate(params *partialUpdateParams) (models.Board, error) {
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

func (rep *repository) Delete(id int) error {
	const deleteCommand = `DELETE FROM boards 
						   WHERE id = $1;`

	res, err := rep.db.Exec(deleteCommand, id)
	if err != nil {
		return ErrDeleteBoard
	}
	count, err := res.RowsAffected()
	if err != nil || count < 1 {
		return ErrBoardNotFound
	}
	return nil
}

func (rep *repository) CheckWriteAccess(userId, boardId string) (bool, error) {
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

func (rep *repository) CheckReadAccess(userId, boardId string) (bool, error) {
	const checkCommand = `SELECT EXISTS(SELECT
              							FROM boards
              							WHERE id = $1 AND (privacy = 'public' OR user_id = $2));`

	row := rep.db.QueryRow(checkCommand, boardId, userId)

	var access bool
	err := row.Scan(&access)
	return access, err
}
