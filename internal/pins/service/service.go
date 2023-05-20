package service

import (
	pkgFollowings "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/followings"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/notifications"
	pkgPins "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pins"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/constants"
)

type service struct {
	rep               pkgPins.Repository
	notificationsServ notifications.Service
	followingsRep     pkgFollowings.Repository
}

func NewService(rep pkgPins.Repository, notificationsServ notifications.Service, followingsRep pkgFollowings.Repository) pkgPins.Service {
	return &service{rep: rep, notificationsServ: notificationsServ, followingsRep: followingsRep}
}

func (serv *service) Create(params *pkgPins.CreateParams) (models.Pin, error) {
	pin, err := serv.rep.Create(params)
	if err != nil {
		return models.Pin{}, err
	}

	go func() {
		followers, err := serv.followingsRep.GetFollowers(pin.Author)
		if err == nil {
			for _, follower := range followers {
				_ = serv.notificationsServ.Create(follower.Id, constants.NewPin, models.NewPinNotification{
					PinID: pin.Id,
				})
			}
		}
	}()

	return pin, err
}

func (serv *service) Get(id, userId int) (models.Pin, error) {
	pin, err := serv.rep.Get(id)
	if err != nil {
		return models.Pin{}, err
	}

	err = serv.SetLikedField(&pin, userId)
	if err != nil {
		return models.Pin{}, err
	}

	return pin, nil
}

func (serv *service) ListByAuthor(authorId, userId, page, limit int) ([]models.Pin, error) {
	pins, err := serv.rep.ListByAuthor(authorId, page, limit)
	if err != nil {
		return []models.Pin{}, err
	}

	for i := range pins {
		err = serv.SetLikedField(&pins[i], userId)
		if err != nil {
			return []models.Pin{}, err
		}
	}

	return pins, nil
}

func (serv *service) List(userId, page, limit int) ([]models.Pin, error) {
	pins, err := serv.rep.List(page, limit)
	if err != nil {
		return []models.Pin{}, err
	}

	for i := range pins {
		err = serv.SetLikedField(&pins[i], userId)
		if err != nil {
			return []models.Pin{}, err
		}
	}

	return pins, nil
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

func (serv *service) SetLikedField(pin *models.Pin, userId int) error {
	liked, err := serv.rep.IsLikedByUser(pin.Id, userId)
	if err != nil {
		return err
	}
	pin.Liked = liked
	return nil
}
