package service

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/boards"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
)

type service struct {
	repo boards.Repository
}

func NewBoardsService(repo boards.Repository) boards.Service {
	return &service{repo: repo}
}

func validatePrivacy(privacy string) error {
	if privacy != "secret" && privacy != "public" {
		return boards.ErrInvalidPrivacy
	}
	return nil
}

func (serv *service) Create(params *boards.CreateParams) (models.Board, error) {
	if err := validatePrivacy(params.Privacy); err != nil {
		return models.Board{}, err
	}
	return serv.repo.Create(params)
}

func (serv *service) List(userId int) ([]models.Board, error) {
	return serv.repo.List(userId)
}

func (serv *service) Get(id int) (models.Board, error) {
	return serv.repo.Get(id)
}

func (serv *service) FullUpdate(params *boards.FullUpdateParams) (models.Board, error) {
	return serv.repo.FullUpdate(params)
}

func (serv *service) PartialUpdate(params *boards.PartialUpdateParams) (models.Board, error) {
	return serv.repo.PartialUpdate(params)
}

func (serv *service) Delete(id int) error {
	return serv.repo.Delete(id)
}

func (serv *service) AddPin(boardId, pinId int) error {
	return serv.repo.AddPin(boardId, pinId)
}

func (serv *service) PinsList(boardId int, page, limit int) ([]models.Pin, error) {
	return serv.repo.PinsList(boardId, page, limit)
}

func (serv *service) RemovePin(boardId, pinId int) error {
	return serv.repo.RemovePin(boardId, pinId)
}

func (serv *service) CheckWriteAccess(userId, boardId string) (bool, error) {
	return serv.repo.CheckWriteAccess(userId, boardId)
}

func (serv *service) CheckReadAccess(userId, boardId string) (bool, error) {
	return serv.repo.CheckReadAccess(userId, boardId)
}
