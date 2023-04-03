package boards

import "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"

type Service interface {
	Create(params *createParams) (models.Board, error)
	List(userId int) ([]models.Board, error)
	Get(id int) (models.Board, error)
	FullUpdate(params *fullUpdateParams) (models.Board, error)
	PartialUpdate(params *partialUpdateParams) (models.Board, error)
	Delete(id int) error

	CheckWriteAccess(userId, boardId string) (bool, error)
	CheckReadAccess(userId, boardId string) (bool, error)
}

func NewService(rep Repository) Service {
	return &service{rep}
}

type service struct {
	rep Repository
}

func (serv *service) Create(params *createParams) (models.Board, error) {
	return serv.rep.Create(params)
}

func (serv *service) List(userId int) ([]models.Board, error) {
	return serv.rep.List(userId)
}

func (serv *service) Get(id int) (models.Board, error) {
	return serv.rep.Get(id)
}

func (serv *service) FullUpdate(params *fullUpdateParams) (models.Board, error) {
	return serv.rep.FullUpdate(params)
}

func (serv *service) PartialUpdate(params *partialUpdateParams) (models.Board, error) {
	return serv.rep.PartialUpdate(params)
}

func (serv *service) Delete(id int) error {
	return serv.rep.Delete(id)
}

func (serv *service) CheckWriteAccess(userId, boardId string) (bool, error) {
	return serv.rep.CheckWriteAccess(userId, boardId)
}

func (serv *service) CheckReadAccess(userId, boardId string) (bool, error) {
	return serv.rep.CheckReadAccess(userId, boardId)
}
