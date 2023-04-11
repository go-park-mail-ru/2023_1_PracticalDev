package service

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/likes"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
)

type service struct {
	rep likes.Repository
}

func NewService(rep likes.Repository) likes.Service {
	return &service{rep}
}

func (serv *service) Like(pinId, authorId int) error {
	exists, err := serv.rep.PinExists(pinId)
	if err != nil {
		return err
	}
	if !exists {
		return likes.ErrPinNotFound
	}

	exists, err = serv.rep.UserExists(authorId)
	if err != nil {
		return err
	}
	if !exists {
		return likes.ErrAuthorNotFound
	}

	exists, err = serv.rep.LikeExists(pinId, authorId)
	if err != nil {
		return err
	}
	if exists {
		return likes.ErrLikeAlreadyExists
	}

	return serv.rep.Create(pinId, authorId)
}

func (serv *service) Unlike(pinId, authorId int) error {
	exists, err := serv.rep.PinExists(pinId)
	if err != nil {
		return err
	}
	if !exists {
		return likes.ErrPinNotFound
	}

	exists, err = serv.rep.UserExists(authorId)
	if err != nil {
		return err
	}
	if !exists {
		return likes.ErrAuthorNotFound
	}

	exists, err = serv.rep.LikeExists(pinId, authorId)
	if err != nil {
		return err
	}
	if !exists {
		return likes.ErrLikeNotFound
	}

	return serv.rep.Delete(pinId, authorId)
}

func (serv *service) ListByAuthor(authorId int) ([]models.Like, error) {
	return serv.rep.ListByAuthor(authorId)
}

func (serv *service) ListByPin(pinId int) ([]models.Like, error) {
	return serv.rep.ListByPin(pinId)
}
