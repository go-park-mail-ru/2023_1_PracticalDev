package http

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	pkgPins "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pins"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/xss"
)

// API responses
type createResponse struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	MediaSource string `json:"media_source"`
	Author      int    `json:"author_id"`
}

func newCreateResponse(pin *models.Pin) *createResponse {
	return &createResponse{
		Id:          pin.Id,
		Title:       xss.Sanitize(pin.Title),
		Description: xss.Sanitize(pin.Description),
		MediaSource: pin.MediaSource,
		Author:      pin.Author,
	}
}

type listResponse struct {
	Pins []pkgPins.Pin `json:"pins"`
}

func newListResponse(pins []pkgPins.Pin) *listResponse {
	for i := range pins {
		pins[i].Title = xss.Sanitize(pins[i].Title)
		pins[i].Description = xss.Sanitize(pins[i].Description)
	}

	return &listResponse{
		Pins: pins,
	}
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

func newGetResponse(pin *pkgPins.Pin) *getResponse {
	return &getResponse{
		Id:          pin.Id,
		Title:       xss.Sanitize(pin.Title),
		Description: xss.Sanitize(pin.Description),
		MediaSource: pin.MediaSource,
		NumLikes:    pin.NumLikes,
		Liked:       pin.Liked,
		Author:      pin.Author,
	}
}

type fullUpdateResponse struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	MediaSource string `json:"media_source"`
	Author      int    `json:"author_id"`
}

func newFullUpdateResponse(pin *models.Pin) *fullUpdateResponse {
	return &fullUpdateResponse{
		Id:          pin.Id,
		Title:       xss.Sanitize(pin.Title),
		Description: xss.Sanitize(pin.Description),
		MediaSource: pin.MediaSource,
		Author:      pin.Author,
	}
}
