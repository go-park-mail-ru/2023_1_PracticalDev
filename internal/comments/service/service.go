package service

import (
	pkgComments "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/comments"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
)

type service struct {
	rep pkgComments.Repository
}

func NewService(rep pkgComments.Repository) pkgComments.Service {
	return &service{rep}
}

func (serv *service) Create(params *pkgComments.CreateParams) (models.Comment, error) {
	return serv.rep.Create(params)
}

func (serv *service) List(pinID int) ([]models.Comment, error) {
	return serv.rep.List(pinID)
}
