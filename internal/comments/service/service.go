package service

import (
	pkgComments "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/comments"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/notifications"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pins"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/constants"
)

type service struct {
	rep               pkgComments.Repository
	notificationsServ notifications.Service
	pinsRep           pins.Repository
}

func NewService(rep pkgComments.Repository, notificationsServ notifications.Service, pinsRep pins.Repository) pkgComments.Service {
	return &service{rep: rep, notificationsServ: notificationsServ, pinsRep: pinsRep}
}

func (serv *service) Create(params *pkgComments.CreateParams) (models.Comment, error) {
	comment, err := serv.rep.Create(params)
	if err != nil {
		return models.Comment{}, err
	}

	pin, err := serv.pinsRep.Get(comment.PinID)
	if err == nil && pin.Author != comment.AuthorID {
		_ = serv.notificationsServ.Create(pin.Author, constants.NewComment, models.NewCommentNotification{
			CommentID: comment.ID,
		})
	}

	return comment, nil
}

func (serv *service) List(pinID int) ([]models.Comment, error) {
	return serv.rep.List(pinID)
}
