package models

//go:generate easyjson -all -snake_case search-res.go

type SearchRes struct {
	Pins   []Pin     `json:"pins"`
	Boards []Board   `json:"boards"`
	Users  []Profile `json:"users"`
}
