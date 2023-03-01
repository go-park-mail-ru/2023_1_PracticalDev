package models

type User struct {
	Id             int    `json:"id"`
	Email          string `json:"email"`
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
}
