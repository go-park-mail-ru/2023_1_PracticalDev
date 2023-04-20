package service

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	pkgPins "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pins"
)

type service struct {
	rep pkgPins.Repository
}

func NewService(rep pkgPins.Repository) pkgPins.Service {
	return &service{rep}
}

func (serv *service) Create(params *pkgPins.CreateParams) (models.Pin, error) {
	return serv.rep.Create(params)
}

func (serv *service) Get(id, userId int) (pkgPins.Pin, error) {
	pin, err := serv.rep.Get(id)
	if err != nil {
		return pkgPins.Pin{}, err
	}
	likedPin, err := serv.addLikedField(&pin, userId)
	if err != nil {
		return pkgPins.Pin{}, err
	}
	return likedPin, nil
}

func (serv *service) ListByAuthor(authorId, userId, page, limit int) ([]pkgPins.Pin, error) {
	pins, err := serv.rep.ListByAuthor(authorId, page, limit)
	if err != nil {
		return []pkgPins.Pin{}, err
	}
	likedPins := []pkgPins.Pin{}
	for _, pin := range pins {
		likedPin, err := serv.addLikedField(&pin, userId)
		if err != nil {
			return []pkgPins.Pin{}, err
		}
		likedPins = append(likedPins, likedPin)
	}
	return likedPins, nil
}

func (serv *service) List(userId, page, limit int) ([]pkgPins.Pin, error) {
	pins, err := serv.rep.List(page, limit)
	if err != nil {
		return []pkgPins.Pin{}, err
	}
	likedPins := []pkgPins.Pin{}
	for _, pin := range pins {
		likedPin, err := serv.addLikedField(&pin, userId)
		if err != nil {
			return []pkgPins.Pin{}, err
		}
		likedPins = append(likedPins, likedPin)
	}
	return likedPins, nil
}

func (serv *service) FullUpdate(params *pkgPins.FullUpdateParams) (models.Pin, error) {
	return serv.rep.FullUpdate(params)
}

func (serv *service) Delete(id int) error {
	return serv.rep.Delete(id)
}

func (serv *service) CheckWriteAccess(userId, pinId string) (bool, error) {
	return serv.rep.CheckWriteAccess(userId, pinId)
}

func (serv *service) CheckReadAccess(userId, pinId string) (bool, error) {
	return serv.rep.CheckReadAccess(userId, pinId)
}

func (serv *service) addLikedField(pin *models.Pin, userId int) (pkgPins.Pin, error) {
	liked, err := serv.rep.IsLikedByUser(pin.Id, userId)
	if err != nil {
		return pkgPins.Pin{}, err
	}
	return pkgPins.Pin{
		Id:          pin.Id,
		Title:       pin.Title,
		Description: pin.Description,
		MediaSource: pin.MediaSource,
		NumLikes:    pin.NumLikes,
		Liked:       liked,
		Author:      pin.Author,
	}, nil
}
