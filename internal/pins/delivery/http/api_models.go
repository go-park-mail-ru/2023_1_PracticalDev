package http

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/xss"
)

//go:generate easyjson -all -snake_case api_models.go

// API responses
type createResponse struct {
	Id               int            `json:"id"`
	Title            string         `json:"title"`
	Description      string         `json:"description"`
	MediaSource      string         `json:"media_source"`
	MediaSourceColor string         `json:"media_source_color"`
	Author           models.Profile `json:"author"`
}

func newCreateResponse(pin *models.Pin) *createResponse {
	return &createResponse{
		Id:               pin.Id,
		Title:            xss.Sanitize(pin.Title),
		Description:      xss.Sanitize(pin.Description),
		MediaSource:      pin.MediaSource,
		MediaSourceColor: pin.MediaSourceColor,
		Author:           pin.Author,
	}
}

type listResponse struct {
	Pins []models.Pin `json:"pins"`
}

func newListResponse(pins []models.Pin) *listResponse {
	for i := range pins {
		pins[i].Title = xss.Sanitize(pins[i].Title)
		pins[i].Description = xss.Sanitize(pins[i].Description)
	}

	return &listResponse{
		Pins: pins,
	}
}

type getResponse struct {
	Id               int            `json:"id"`
	Title            string         `json:"title"`
	Description      string         `json:"description"`
	MediaSource      string         `json:"media_source"`
	MediaSourceColor string         `json:"media_source_color"`
	NumLikes         int            `json:"n_likes"`
	Liked            bool           `json:"liked"`
	Author           models.Profile `json:"author"`
}

func newGetResponse(pin *models.Pin) *getResponse {
	return &getResponse{
		Id:               pin.Id,
		Title:            xss.Sanitize(pin.Title),
		Description:      xss.Sanitize(pin.Description),
		MediaSource:      pin.MediaSource,
		MediaSourceColor: pin.MediaSourceColor,
		NumLikes:         pin.NumLikes,
		Liked:            pin.Liked,
		Author:           pin.Author,
	}
}

type fullUpdateResponse struct {
	Id               int            `json:"id"`
	Title            string         `json:"title"`
	Description      string         `json:"description"`
	MediaSource      string         `json:"media_source"`
	MediaSourceColor string         `json:"media_source_color"`
	Author           models.Profile `json:"author"`
}

func newFullUpdateResponse(pin *models.Pin) *fullUpdateResponse {
	return &fullUpdateResponse{
		Id:               pin.Id,
		Title:            xss.Sanitize(pin.Title),
		Description:      xss.Sanitize(pin.Description),
		MediaSource:      pin.MediaSource,
		MediaSourceColor: pin.MediaSourceColor,
		Author:           pin.Author,
	}
}
