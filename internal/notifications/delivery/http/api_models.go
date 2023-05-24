package http

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/constants"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/xss"
)

//go:generate easyjson -all -snake_case api_models.go

// API responses
type listResponse struct {
	Notifications []models.Notification `json:"items"`
}

func newListResponse(notifications []models.Notification) *listResponse {
	for i := range notifications {
		if notifications[i].Type == constants.NewComment {
			nc := notifications[i].Data.(models.NewCommentNotification)
			nc.Text = xss.Sanitize(nc.Text)
			notifications[i].Data = nc
		}
	}

	return &listResponse{Notifications: notifications}
}
