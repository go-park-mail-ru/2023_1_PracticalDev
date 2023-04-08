package http

import "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"

// API responses
type createResponse struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	MediaSource string `json:"media_source"`
	Author      int    `json:"author_id"`
}

type listResponse struct {
	Pins []models.Pin `json:"pins"`
}

type getResponse struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	MediaSource string `json:"media_source"`
	Author      int    `json:"author_id"`
}

type fullUpdateResponse struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	MediaSource string `json:"media_source"`
	Author      int    `json:"author_id"`
}
