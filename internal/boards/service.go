package boards

import "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"

type Service interface {
	CreateBoard(board *models.Board) (models.Board, error)
	GetBoards(userId int) ([]models.Board, error)
	GetBoard(id int) (models.Board, error)
	DeleteBoard(id int) error
}

func NewService(rep Repository) Service {
	return &service{rep}
}

type service struct {
	rep Repository
}

func (serv *service) CreateBoard(board *models.Board) (models.Board, error) {
	return serv.rep.CreateBoard(board)
}

func (serv *service) GetBoards(userId int) ([]models.Board, error) {
	return serv.rep.GetBoards(userId)
}

func (serv *service) GetBoard(id int) (models.Board, error) {
	return serv.rep.GetBoard(id)
}

func (serv *service) DeleteBoard(id int) error {
	return serv.rep.DeleteBoard(id)
}
