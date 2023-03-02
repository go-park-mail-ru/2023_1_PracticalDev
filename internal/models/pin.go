package models

type Pin struct {
	Id          int    `json:"id"`
	Link        string `json:"link"`
	Title       string `json:"title"`
	Description string `json:"description"`
	MediaSource string `json:"media_source"`
	BoardId     int    `json:"board_id"`
}
