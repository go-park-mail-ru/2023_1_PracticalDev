package pins

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
)

type Service interface {
	CreatePin(params *models.Pin, image *models.Image) (models.Pin, error)
	GetPin(id int) (models.Pin, error)
	GetPinsByUser(userId int, page, limit int) ([]models.Pin, error)
	GetPinsByBoard(boardId int, page, limit int) ([]models.Pin, error)
	GetPins(page, limit int) ([]models.Pin, error)
	UpdatePin(params *models.Pin) (models.Pin, error)
	DeletePin(id int) error
	AddPinToBoard(boardId, pinId int) error
	RemovePinFromBoard(boardId, pinId int) error

	CheckWriteAccess(userId, pinId string) (bool, error)
	CheckReadAccess(userId, pinId string) (bool, error)
}

func NewService(rep Repository) Service {
	return &service{rep}
}

type service struct {
	rep Repository
}

func (serv service) CreatePin(params *models.Pin, image *models.Image) (models.Pin, error) {
	return serv.rep.CreatePin(params, image)
}

func (serv service) GetPin(id int) (models.Pin, error) {
	return serv.rep.GetPin(id)
}

func (serv service) GetPinsByUser(userId int, page, limit int) ([]models.Pin, error) {
	return serv.rep.GetPinsByUser(userId, page, limit)
}

func (serv service) GetPinsByBoard(boardId int, page, limit int) ([]models.Pin, error) {
	return serv.rep.GetPinsByBoard(boardId, page, limit)
}

func (serv service) GetPins(page, limit int) ([]models.Pin, error) {
	return serv.rep.GetPins(page, limit)
}

func (serv service) UpdatePin(params *models.Pin) (models.Pin, error) {
	return serv.rep.UpdatePin(params)
}

func (serv service) DeletePin(id int) error {
	return serv.rep.DeletePin(id)
}

func (serv service) AddPinToBoard(boardId, pinId int) error {
	return serv.rep.AddPinToBoard(boardId, pinId)
}
func (serv service) RemovePinFromBoard(boardId, pinId int) error {
	return serv.rep.RemovePinFromBoard(boardId, pinId)
}

func (serv service) CheckWriteAccess(userId, pinId string) (bool, error) {
	return serv.rep.CheckWriteAccess(userId, pinId)
}

func (serv service) CheckReadAccess(userId, pinId string) (bool, error) {
	return serv.rep.CheckReadAccess(userId, pinId)
}
