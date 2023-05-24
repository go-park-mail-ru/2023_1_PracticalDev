package http

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/xss"
)

//go:generate easyjson -all -snake_case api_models.go

// API responses
type getResponse struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	ProfileImage string `json:"profile_image"`
	WebsiteUrl   string `json:"website_url"`
	AccountType  string `json:"account_type"`
}

func newGetResponse(user *models.User) *getResponse {
	return &getResponse{
		ID:           user.Id,
		Username:     xss.Sanitize(user.Username),
		Email:        user.Email,
		Name:         xss.Sanitize(user.Name),
		ProfileImage: user.ProfileImage,
		WebsiteUrl:   user.WebsiteUrl,
		AccountType:  user.AccountType,
	}
}
