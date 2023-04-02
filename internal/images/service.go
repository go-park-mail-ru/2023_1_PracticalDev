package images

import "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"

type Service interface {
	UploadImage(image *models.Image) (string, error)
}

func NewS3Service(rep Repository) Service {
	return &service{
		rep: rep,
	}
}

type service struct {
	rep Repository
}

func (serv *service) UploadImage(image *models.Image) (string, error) {
	url, err := serv.rep.UploadImage(image)
	if err != nil {
		return "", err
	}

	return url, nil
}
