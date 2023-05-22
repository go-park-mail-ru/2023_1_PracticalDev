package http

//go:generate easyjson -all -snake_case api_models.go

type url struct {
	URL string `json:"url"`
}
