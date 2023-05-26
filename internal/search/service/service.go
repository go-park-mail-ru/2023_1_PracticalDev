package service

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	pkgPins "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pins"
	pkgSearch "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/search"
)

type service struct {
	rep     pkgSearch.Repository
	pinServ pkgPins.Service
}

func NewService(rep pkgSearch.Repository, serv pkgPins.Service) pkgSearch.Service {
	return &service{rep, serv}
}

func (serv *service) Search(userId int, query string) (models.SearchRes, error) {
	res, err := serv.rep.Search(query)
	if err != nil {
		return res, err
	}

	for i := range res.Pins {
		err := serv.pinServ.SetLikedField(&res.Pins[i], userId)
		if err != nil {
			return res, err
		}
	}

	return res, nil
}

func (serv *service) Suggestions(query string) ([]string, error) {
	return serv.rep.Suggestions(query)
}
