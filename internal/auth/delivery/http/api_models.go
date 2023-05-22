package http

//go:generate easyjson -all -snake_case api_models.go

// API requests
type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}
