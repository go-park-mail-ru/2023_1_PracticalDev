package service

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/likes"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/notifications"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pins"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/constants"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
	"go.uber.org/zap"
)

type service struct {
	rep               likes.Repository
	notificationsServ notifications.Service
	pinsRep           pins.Repository
	log               *zap.Logger
}

func NewService(rep likes.Repository, notificationsServ notifications.Service, pinsRep pins.Repository, log *zap.Logger) likes.Service {
	return &service{rep: rep, notificationsServ: notificationsServ, pinsRep: pinsRep, log: log}
}

func (serv *service) Like(pinID, authorID int) error {
	exists, err := serv.rep.PinExists(pinID)
	if err != nil {
		return err
	}
	if !exists {
		return pkgErrors.ErrPinNotFound
	}

	exists, err = serv.rep.UserExists(authorID)
	if err != nil {
		return err
	}
	if !exists {
		return pkgErrors.ErrUserNotFound
	}

	exists, err = serv.rep.LikeExists(pinID, authorID)
	if err != nil {
		return err
	}
	if exists {
		return pkgErrors.ErrLikeAlreadyExists
	}

	err = serv.rep.Create(pinID, authorID)
	if err != nil {
		return err
	}

	pin, err := serv.pinsRep.Get(pinID)
	if err == nil && pin.Author != authorID {
		_ = serv.notificationsServ.Create(pin.Author, constants.NewLike, models.NewLikeNotification{
			PinID:    pinID,
			AuthorID: authorID})
	}

	return nil
}

func (serv *service) Unlike(pinId, authorId int) error {
	exists, err := serv.rep.PinExists(pinId)
	if err != nil {
		return err
	}
	if !exists {
		return pkgErrors.ErrPinNotFound
	}

	exists, err = serv.rep.UserExists(authorId)
	if err != nil {
		return err
	}
	if !exists {
		return pkgErrors.ErrUserNotFound
	}

	exists, err = serv.rep.LikeExists(pinId, authorId)
	if err != nil {
		return err
	}
	if !exists {
		return pkgErrors.ErrLikeNotFound
	}

	return serv.rep.Delete(pinId, authorId)
}

func (serv *service) ListByAuthor(authorId int) ([]models.Like, error) {
	return serv.rep.ListByAuthor(authorId)
}

func (serv *service) ListByPin(pinId int) ([]models.Like, error) {
	return serv.rep.ListByPin(pinId)
}
