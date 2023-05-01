package search

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
)

type Repository interface {
	Get(query string) (models.SearchRes, error)
}
