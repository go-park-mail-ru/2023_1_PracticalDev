package service

import (
	pkgFollowings "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/followings"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/notifications"
	pkgPins "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pins"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/constants"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
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
	if err := validateTitle(params.Title); err != nil {
		return models.Pin{}, err
	} else if err = validateDescription(params.Description); err != nil {
		return models.Pin{}, err
	}

	pin, err := serv.rep.Create(params)
	if err != nil {
		return models.Pin{}, err
	}

	go func(pin models.Pin) {
		followers, err := serv.followingsRep.GetFollowers(pin.Author.Id)
		if err == nil {
			for _, follower := range followers {
				_ = serv.notificationsServ.Create(follower.Id, constants.NewPin, models.NewPinNotification{
					PinID: pin.Id,
				})
			}
		}
	}(pin)

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

func (serv *service) List(authorized bool, userID int, liked bool, page, limit int) ([]models.Pin, error) {
	if authorized {
		if liked {
			return serv.rep.ListLiked(userID, page, limit)
		}
		return serv.rep.ListWithLikedField(userID, page, limit)
	}
	return serv.rep.List(page, limit)
}

func (serv *service) FullUpdate(params *pkgPins.FullUpdateParams) (models.Pin, error) {
	if err := validateTitle(params.Title); err != nil {
		return models.Pin{}, err
	} else if err = validateDescription(params.Description); err != nil {
		return models.Pin{}, err
	}

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

func validateTitle(title string) error {
	if len(title) > constants.MaxPinTitleLen {
		return pkgErrors.ErrTooLongPinTitle
	}
	return nil
}

func validateDescription(description string) error {
	if len(description) > constants.MaxPinDescriptionLen {
		return pkgErrors.ErrTooLongPinDescription
	}
	return nil
}
