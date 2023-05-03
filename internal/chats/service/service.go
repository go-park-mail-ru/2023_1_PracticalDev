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

func (serv *service) MessagesList(chatID int) ([]models.Message, error) {
	return serv.repo.MessagesList(chatID)
}

func (serv *service) SendMessage(params *pkgChats.SendMessageParams) (*models.Message, error) {
	return serv.repo.SendMessage(params)
}

func (serv *service) ChatExists(user1ID, user2ID int) (bool, error) {
	return serv.repo.ChatExists(user1ID, user2ID)
}

func (serv *service) GetByUsers(user1ID, user2ID int) (models.Chat, error) {
	return serv.repo.GetByUsers(user1ID, user2ID)
}
