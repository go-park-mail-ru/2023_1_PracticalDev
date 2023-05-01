package models

type SearchRes struct {
	Pins   []Pin     `json:"pins"`
	Boards []Board   `json:"boards"`
	Users  []Profile `json:"users"`
}
