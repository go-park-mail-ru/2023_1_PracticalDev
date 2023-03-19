package boards

import (
	"database/sql"
	"errors"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
)

var (
	ErrDeleteBoard   = errors.New("failed to delete board")
	ErrBoardNotFound = errors.New("board not found")
)

type Repository interface {
	CreateBoard(board *models.Board) (models.Board, error)
	GetBoards(userId int) ([]models.Board, error)
	GetBoard(id int) (models.Board, error)
	DeleteBoard(id int) error
}

func NewRepository(db *sql.DB, log log.Logger) Repository {
	return &repository{db, log}
}

type repository struct {
	db  *sql.DB
	log log.Logger
}

func (rep *repository) CreateBoard(board *models.Board) (models.Board, error) {
	const insertCommand = `INSERT INTO boards (name, description, privacy, user_id) 
				      	   VALUES ($1, $2, $3, $4)
						   RETURNING *;`

	row := rep.db.QueryRow(insertCommand, board.Name, board.Description, board.Privacy, board.UserId)
	createdBoard := models.Board{}
	var description sql.NullString
	err := row.Scan(&createdBoard.Id, &createdBoard.Name, &description, &createdBoard.Privacy, &createdBoard.UserId)
	createdBoard.Description = description.String
	return createdBoard, err
}

func (rep *repository) GetBoards(userId int) ([]models.Board, error) {
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

func (rep *repository) GetBoard(id int) (models.Board, error) {
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

func (rep *repository) DeleteBoard(id int) error {
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
