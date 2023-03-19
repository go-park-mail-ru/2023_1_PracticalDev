package boards

import "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"

type Service interface {
	GetBoards(userId int) ([]models.Board, error)
}

func NewService(rep Repository) Service {
	return &service{rep}
}

type service struct {
	rep Repository
}

func (serv *service) GetBoards(userId int) ([]models.Board, error) {
	return serv.rep.GetBoards(userId)
}
