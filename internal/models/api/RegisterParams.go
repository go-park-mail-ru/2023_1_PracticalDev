package api

type RegisterParams struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}
