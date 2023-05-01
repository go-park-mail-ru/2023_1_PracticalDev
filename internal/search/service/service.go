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

func (serv service) Get(userId int, query string) (models.SearchRes, error) {
	res, err := serv.rep.Get(query)
	if err != nil {
		return res, err
	}

	for i := range res.Pins {
		pin, err := serv.pinServ.Get(res.Pins[i].Id, userId)
		if err != nil {
			return res, err
		}

		res.Pins[i].Liked = pin.Liked
	}
	return res, nil
}
