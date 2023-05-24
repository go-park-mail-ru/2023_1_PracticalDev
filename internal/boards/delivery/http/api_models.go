package http

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/xss"
)

//go:generate easyjson -all -snake_case api_models.go

// API requests
type createRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Privacy     *string `json:"privacy"`
}

type fullUpdateRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Privacy     *string `json:"privacy"`
}

type partialUpdateRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Privacy     *string `json:"privacy"`
}

// API responses
type listResponse struct {
	Boards []models.Board `json:"boards"`
}

func newListResponse(boards []models.Board) *listResponse {
	for i := range boards {
		boards[i].Name = xss.Sanitize(boards[i].Name)
		boards[i].Description = xss.Sanitize(boards[i].Description)
	}

	return &listResponse{
		Boards: boards,
	}
}

type pinListResponse struct {
	Pins []models.Pin `json:"pins"`
}

func newPinsListResponse(pins []models.Pin) *pinListResponse {
	for i := range pins {
		pins[i].Title = xss.Sanitize(pins[i].Title)
		pins[i].Description = xss.Sanitize(pins[i].Description)
	}

	return &pinListResponse{
		Pins: pins,
	}
}

type createResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Privacy     string `json:"privacy"`
	UserId      int    `json:"user_id"`
}

func newCreateResponse(board *models.Board) *createResponse {
	return &createResponse{
		ID:          board.Id,
		Name:        xss.Sanitize(board.Name),
		Description: xss.Sanitize(board.Description),
		Privacy:     board.Privacy,
		UserId:      board.UserId,
	}
}

type getResponse struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Privacy     string `json:"privacy"`
	UserId      int    `json:"user_id"`
}

func newGetResponse(board *models.Board) *getResponse {
	return &getResponse{
		Id:          board.Id,
		Name:        xss.Sanitize(board.Name),
		Description: xss.Sanitize(board.Description),
		Privacy:     board.Privacy,
		UserId:      board.UserId,
	}
}

type fullUpdateResponse struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Privacy     string `json:"privacy"`
	UserId      int    `json:"user_id"`
}

func newFullUpdateResponse(board *models.Board) *fullUpdateResponse {
	return &fullUpdateResponse{
		Id:          board.Id,
		Name:        xss.Sanitize(board.Name),
		Description: xss.Sanitize(board.Description),
		Privacy:     board.Privacy,
		UserId:      board.UserId,
	}
}

type partialUpdateResponse struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Privacy     string `json:"privacy"`
	UserId      int    `json:"user_id"`
}

func newPartialUpdateResponse(board *models.Board) *partialUpdateResponse {
	return &partialUpdateResponse{
		Id:          board.Id,
		Name:        xss.Sanitize(board.Name),
		Description: xss.Sanitize(board.Description),
		Privacy:     board.Privacy,
		UserId:      board.UserId,
	}
}
