package middleware

import (
	"net/http"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/router"
	"github.com/julienschmidt/httprouter"
)

var allowedOrigins = map[string]interface{}{
	"http://pickpin.ru": nil,
	"http://localhost":  nil,
}

func OptionsHandler(w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	origin := r.Header.Get("Origin")
	if _, allowed := allowedOrigins[origin]; allowed {
		header.Set("Access-Control-Allow-Origin", origin)
		header.Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	}

	w.WriteHeader(http.StatusNoContent)
}

func CorsChecker(handler router.Handler) router.Handler {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
		origin := r.Header.Get("Origin")
		if _, allowed := allowedOrigins[origin]; allowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		return handler(w, r, p)
	}
}
