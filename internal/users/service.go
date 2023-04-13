package users

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
)

type Service interface {
	Get(id int) (models.User, error)
}
