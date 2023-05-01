package search

import "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"

type Service interface {
	Get(userId int, query string) (models.SearchRes, error)
}
