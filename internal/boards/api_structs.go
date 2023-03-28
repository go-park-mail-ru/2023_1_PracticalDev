package boards

import "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"

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

type getResponse struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Privacy     string `json:"privacy"`
	UserId      int    `json:"user_id"`
}

type fullUpdateResponse struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Privacy     string `json:"privacy"`
	UserId      int    `json:"user_id"`
}

type partialUpdateResponse struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Privacy     string `json:"privacy"`
	UserId      int    `json:"user_id"`
}
