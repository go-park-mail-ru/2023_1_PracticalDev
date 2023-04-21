package http

type LoginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterParams struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}
