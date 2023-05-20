package service

import (
	"fmt"
	"os"

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

func (serv *service) CreatePinLink(id int) (string, error) {
	if os.Getenv("SHORT_HOST") == "localhost:8091" {
		return serv.Create(fmt.Sprintf("http://localhost/pin/%d", id))
	}
	return serv.Create(fmt.Sprintf("https://pickpin.ru/pin/%d", id))
}
