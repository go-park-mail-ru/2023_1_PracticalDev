package http

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
)

//go:generate easyjson -all -snake_case api_models.go

// API responses
type listByAuthorResponse struct {
	Likes []models.Like `json:"likes"`
}

type listByPinResponse struct {
	Likes []models.Like `json:"likes"`
}
