package models

//go:generate easyjson -all -snake_case profile.go

type Profile struct {
	Id           int    `json:"id"`
	Username     string `json:"username"`
	Name         string `json:"name"`
	ProfileImage string `json:"profile_image"`
	WebsiteUrl   string `json:"website_url"`
}
