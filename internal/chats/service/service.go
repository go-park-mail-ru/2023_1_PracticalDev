package service

import (
	pkgChats "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/chats"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
)

type service struct {
	repo pkgChats.Repository
}

func NewService(repo pkgChats.Repository) pkgChats.Service {
	return &service{repo: repo}
}

func (serv *service) Create(params *pkgChats.CreateParams) (models.Chat, error) {
	if params.User1ID == params.User2ID {
		return models.Chat{}, errors.ErrSameUserId
	}

	return serv.repo.Create(params)
}

func (serv *service) ListByUser(userId int) ([]models.Chat, error) {
	return serv.repo.ListByUser(userId)
}

func (serv *service) Get(id int) (models.Chat, error) {
	return serv.repo.Get(id)
}
