package models

type User struct {
	Id             int    `json:"id"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	HashedPassword string `json:"-"`
	Name           string `json:"name"`
	Profile_image  string `json:"profile_image"`
	Website_url    string `json:"website_url"`
	Account_type   string `json:"account_type"`
}
