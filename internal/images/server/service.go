package images

import "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"

type Service interface {
	UploadImage(image *models.Image) (string, error)
}
