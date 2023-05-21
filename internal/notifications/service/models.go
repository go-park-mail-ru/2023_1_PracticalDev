package service

//go:generate easyjson -all -snake_case models.go

type Request struct {
	ID int `json:"id"`
}

type Message struct {
	Type    string      `json:"type"` // response, notification
	Content interface{} `json:"content"`
}

// Code
// 20: ok,
// 40: bad request,
// 50: internal error
type ResponseContent struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}
