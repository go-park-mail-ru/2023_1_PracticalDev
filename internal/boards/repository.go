package boards

import (
	"database/sql"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
)

type Repository interface {
	GetBoards(userId int) ([]models.Board, error)
}

func NewRepository(db *sql.DB, log log.Logger) Repository {
	return &repository{db, log}
}

type repository struct {
	db  *sql.DB
	log log.Logger
}

func (rep *repository) GetBoards(userId int) ([]models.Board, error) {
	const getBoardsCommand = `SELECT id, name, description, privacy, user_id 
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
