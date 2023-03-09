package models

type User struct {
	Id             int    `json:"id"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	HashedPassword string `json:"-"`
	Name           string `json:"name"`
	ProfileImage   string `json:"profile_image"`
	WebsiteUrl     string `json:"website_url"`
	AccountType    string `json:"account_type"`
}
