package models

//go:generate easyjson -all -snake_case pin.go

type Pin struct {
	Id               int     `json:"id"`
	Title            string  `json:"title"`
	Description      string  `json:"description"`
	MediaSource      string  `json:"media_source"`
	MediaSourceColor string  `json:"media_source_color"`
	NumLikes         int     `json:"n_likes"`
	Liked            bool    `json:"liked"`
	Author           Profile `json:"author"`
}
