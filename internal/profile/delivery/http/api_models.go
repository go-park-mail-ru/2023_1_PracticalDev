package http

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/xss"
	pkgProfile "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/profile"
)

//go:generate easyjson -all -snake_case api_models.go

// API responses
type getResponse struct {
	Username     string `json:"username"`
	Name         string `json:"name"`
	ProfileImage string `json:"profile_image"`
	WebsiteUrl   string `json:"website_url"`
}

func newGetResponse(profile *pkgProfile.Profile) *getResponse {
	return &getResponse{
		Username:     xss.Sanitize(profile.Username),
		Name:         xss.Sanitize(profile.Name),
		ProfileImage: profile.ProfileImage,
		WebsiteUrl:   xss.Sanitize(profile.WebsiteUrl),
	}
}

type fullUpdateResponse struct {
	Username     string `json:"username"`
	Name         string `json:"name"`
	ProfileImage string `json:"profile_image"`
	WebsiteUrl   string `json:"website_url"`
}

func newFullUpdateResponse(profile *pkgProfile.Profile) *fullUpdateResponse {
	return &fullUpdateResponse{
		Username:     xss.Sanitize(profile.Username),
		Name:         xss.Sanitize(profile.Name),
		ProfileImage: profile.ProfileImage,
		WebsiteUrl:   xss.Sanitize(profile.WebsiteUrl),
	}
}

type partialUpdateResponse struct {
	Username     string `json:"username"`
	Name         string `json:"name"`
	ProfileImage string `json:"profile_image"`
	WebsiteUrl   string `json:"website_url"`
}

func newPartialUpdateResponse(profile *pkgProfile.Profile) *partialUpdateResponse {
	return &partialUpdateResponse{
		Username:     xss.Sanitize(profile.Username),
		Name:         xss.Sanitize(profile.Name),
		ProfileImage: profile.ProfileImage,
		WebsiteUrl:   xss.Sanitize(profile.WebsiteUrl),
	}
}
