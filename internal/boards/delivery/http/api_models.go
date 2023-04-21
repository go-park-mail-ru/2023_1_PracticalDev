package http

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/xss"
)

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

type pinListResponse struct {
	Pins []models.Pin `json:"pins"`
}

type GetResponse struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Privacy     string `json:"privacy"`
	UserId      int    `json:"user_id"`
}

func NewGetResponse(board *models.Board) *GetResponse {
	return &GetResponse{
		Id:          board.Id,
		Name:        xss.Sanitize(board.Name),
		Description: xss.Sanitize(board.Description),
		Privacy:     board.Privacy,
		UserId:      board.UserId,
	}
}

type FullUpdateResponse struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Privacy     string `json:"privacy"`
	UserId      int    `json:"user_id"`
}

func NewFullUpdateResponse(board *models.Board) *FullUpdateResponse {
	return &FullUpdateResponse{
		Id:          board.Id,
		Name:        xss.Sanitize(board.Name),
		Description: xss.Sanitize(board.Description),
		Privacy:     board.Privacy,
		UserId:      board.UserId,
	}
}

type PartialUpdateResponse struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Privacy     string `json:"privacy"`
	UserId      int    `json:"user_id"`
}

func NewPartialUpdateResponse(board *models.Board) *PartialUpdateResponse {
	return &PartialUpdateResponse{
		Id:          board.Id,
		Name:        xss.Sanitize(board.Name),
		Description: xss.Sanitize(board.Description),
		Privacy:     board.Privacy,
		UserId:      board.UserId,
	}
}
