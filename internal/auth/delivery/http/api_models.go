package http

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/xss"
)

//go:generate easyjson -all -snake_case api_models.go

// API requests
type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

// API responses
type authenticateResponse struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	ProfileImage string `json:"profile_image"`
	WebsiteUrl   string `json:"website_url"`
	AccountType  string `json:"account_type"`
}

func newAuthenticateResponse(user *models.User) *authenticateResponse {
	return &authenticateResponse{
		ID:           user.Id,
		Username:     xss.Sanitize(user.Username),
		Email:        user.Email,
		Name:         xss.Sanitize(user.Name),
		ProfileImage: user.ProfileImage,
		WebsiteUrl:   user.WebsiteUrl,
		AccountType:  user.AccountType,
	}
}

type registerResponse struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	ProfileImage string `json:"profile_image"`
	WebsiteUrl   string `json:"website_url"`
	AccountType  string `json:"account_type"`
}

func newRegisterResponse(user *models.User) *registerResponse {
	return &registerResponse{
		ID:           user.Id,
		Username:     xss.Sanitize(user.Username),
		Email:        user.Email,
		Name:         xss.Sanitize(user.Name),
		ProfileImage: user.ProfileImage,
		WebsiteUrl:   user.WebsiteUrl,
		AccountType:  user.AccountType,
	}
}

type checkAuthResponse struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	ProfileImage string `json:"profile_image"`
	WebsiteUrl   string `json:"website_url"`
	AccountType  string `json:"account_type"`
}

func newCheckAuthResponse(user *models.User) *checkAuthResponse {
	return &checkAuthResponse{
		ID:           user.Id,
		Username:     xss.Sanitize(user.Username),
		Email:        user.Email,
		Name:         xss.Sanitize(user.Name),
		ProfileImage: user.ProfileImage,
		WebsiteUrl:   user.WebsiteUrl,
		AccountType:  user.AccountType,
	}
}
