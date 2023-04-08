package service

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pins"
)

type service struct {
	rep pins.Repository
}

func NewService(rep pins.Repository) pins.Service {
	return &service{rep}
}

func (serv *service) Create(params *pins.CreateParams) (models.Pin, error) {
	return serv.rep.Create(params)
}

func (serv *service) Get(id int) (models.Pin, error) {
	return serv.rep.Get(id)
}

func (serv *service) ListByUser(userId int, page, limit int) ([]models.Pin, error) {
	return serv.rep.ListByUser(userId, page, limit)
}

func (serv *service) ListByBoard(boardId int, page, limit int) ([]models.Pin, error) {
	return serv.rep.ListByBoard(boardId, page, limit)
}

func (serv *service) List(page, limit int) ([]models.Pin, error) {
	return serv.rep.List(page, limit)
}

func (serv *service) Update(params *models.Pin) (models.Pin, error) {
	return serv.rep.Update(params)
}

func (serv *service) Delete(id int) error {
	return serv.rep.Delete(id)
}

func (serv *service) AddToBoard(boardId, pinId int) error {
	return serv.rep.AddToBoard(boardId, pinId)
}

func (serv *service) RemoveFromBoard(boardId, pinId int) error {
	return serv.rep.RemoveFromBoard(boardId, pinId)
}

func (serv *service) CheckWriteAccess(userId, pinId string) (bool, error) {
	return serv.rep.CheckWriteAccess(userId, pinId)
}

func (serv *service) CheckReadAccess(userId, pinId string) (bool, error) {
	return serv.rep.CheckReadAccess(userId, pinId)
}
