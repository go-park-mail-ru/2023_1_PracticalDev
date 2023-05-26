package http

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/xss"
)

//go:generate easyjson -all -snake_case api_models.go

// API responses
type searchResponse struct {
	Pins   []models.Pin     `json:"pins"`
	Boards []models.Board   `json:"boards"`
	Users  []models.Profile `json:"users"`
}

func newSearchResponse(searchRes *models.SearchRes) *searchResponse {
	for i := range searchRes.Pins {
		searchRes.Pins[i].Title = xss.Sanitize(searchRes.Pins[i].Title)
		searchRes.Pins[i].Description = xss.Sanitize(searchRes.Pins[i].Description)
	}

	for i := range searchRes.Boards {
		searchRes.Boards[i].Name = xss.Sanitize(searchRes.Boards[i].Name)
		searchRes.Boards[i].Description = xss.Sanitize(searchRes.Boards[i].Description)
	}

	for i := range searchRes.Users {
		searchRes.Users[i].Name = xss.Sanitize(searchRes.Users[i].Name)
		searchRes.Users[i].Username = xss.Sanitize(searchRes.Users[i].Username)
	}

	return &searchResponse{
		Pins:   searchRes.Pins,
		Boards: searchRes.Boards,
		Users:  searchRes.Users,
	}
}

type suggestionsResponse struct {
	Items []string `json:"items"`
}

func newSuggestionsResponse(suggestions []string) *suggestionsResponse {
	return &suggestionsResponse{Items: suggestions}
}
