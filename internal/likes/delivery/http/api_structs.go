package http

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
)

// API responses
type listByAuthorResponse struct {
	Likes []models.Like `json:"likes"`
}

type listByPinResponse struct {
	Likes []models.Like `json:"likes"`
}
