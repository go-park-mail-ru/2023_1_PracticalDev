package search

import "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"

type Service interface {
	Search(userId int, query string) (models.SearchRes, error)
	Suggestions(query string) ([]string, error)
}
