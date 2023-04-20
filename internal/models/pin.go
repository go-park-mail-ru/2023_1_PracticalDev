package models

type Pin struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	MediaSource string `json:"media_source"`
	NumLikes    int    `json:"n_likes"`
	Author      int    `json:"author_id"`
}
