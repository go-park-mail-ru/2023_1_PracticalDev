package service

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/boards"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
)

type boardsService struct {
	repo boards.Repository
}

func NewBoardsService(repo boards.Repository) boards.Service {
	return &boardsService{repo: repo}
}

func (serv *boardsService) Create(params *boards.CreateParams) (models.Board, error) {
	return serv.repo.Create(params)
}

func (serv *boardsService) List(userId int) ([]models.Board, error) {
	return serv.repo.List(userId)
}

func (serv *boardsService) Get(id int) (models.Board, error) {
	return serv.repo.Get(id)
}

func (serv *boardsService) FullUpdate(params *boards.FullUpdateParams) (models.Board, error) {
	return serv.repo.FullUpdate(params)
}

func (serv *boardsService) PartialUpdate(params *boards.PartialUpdateParams) (models.Board, error) {
	return serv.repo.PartialUpdate(params)
}

func (serv *boardsService) Delete(id int) error {
	return serv.repo.Delete(id)
}

func (serv *boardsService) CheckWriteAccess(userId, boardId string) (bool, error) {
	return serv.repo.CheckWriteAccess(userId, boardId)
}

func (serv *boardsService) CheckReadAccess(userId, boardId string) (bool, error) {
	return serv.repo.CheckReadAccess(userId, boardId)
}
