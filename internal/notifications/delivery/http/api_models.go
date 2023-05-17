package http

import "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"

// API responses
type listResponse struct {
	Notifications []models.Notification `json:"items"`
}

func newListResponse(notifications []models.Notification) *listResponse {
	return &listResponse{Notifications: notifications}
}
