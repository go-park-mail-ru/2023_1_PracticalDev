package service

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/shortener"
)

type service struct {
	rep shortener.ShortenerRepository
}

func NewShortenerService(rep shortener.ShortenerRepository) shortener.ShortenerService {
	return &service{
		rep: rep,
	}
}

func (serv *service) Get(hash string) (string, error) {
	return serv.rep.Get(hash)
}

func (serv *service) Create(url string) (string, error) {
	return serv.rep.Create(url)
}
