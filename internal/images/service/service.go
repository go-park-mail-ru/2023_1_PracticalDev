package service

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/images"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
)

type service struct {
	rep images.Repository
}

func NewS3Service(rep images.Repository) images.Service {
	return &service{
		rep: rep,
	}
}

func (serv *service) UploadImage(image *models.Image) (string, error) {
	url, err := serv.rep.UploadImage(image)
	if err != nil {
		return "", err
	}

	return url, nil
}
