package http

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pins"
)

// API responses
type createResponse struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	MediaSource string `json:"media_source"`
	Author      int    `json:"author_id"`
}

type listResponse struct {
	Pins []pins.Pin `json:"pins"`
}

type getResponse struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	MediaSource string `json:"media_source"`
	NumLikes    int    `json:"n_likes"`
	Liked       bool   `json:"liked"`
	Author      int    `json:"author_id"`
}

type fullUpdateResponse struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	MediaSource string `json:"media_source"`
	Author      int    `json:"author_id"`
}
