package middleware

import (
	"net/http"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/router"
	"github.com/julienschmidt/httprouter"
)

const mainOrigin = "http://pickpin.ru"

var allowedOrigins = map[string]struct{}{
	mainOrigin:                {},
	"http://localhost":        {},
	"http://127.0.0.1":        {},
	"https://park.vk.company": {},
}

func OptionsHandler(w http.ResponseWriter, r *http.Request) {
	origin := mainOrigin
	gotOrigin := r.Header.Get("Origin")
	if _, allowed := allowedOrigins[gotOrigin]; allowed {
		origin = gotOrigin
	}
	wHeader := w.Header()
	wHeader.Set("Access-Control-Allow-Origin", origin)
	wHeader.Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	wHeader.Set("Access-Control-Allow-Headers", "Content-Type")
	wHeader.Set("Access-Control-Max-Age", "86400")
	wHeader.Set("Vary", "Origin")
	w.WriteHeader(http.StatusNoContent)
}

func CorsChecker(handler router.Handler) router.Handler {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
		origin := mainOrigin
		gotOrigin := r.Header.Get("Origin")
		if _, allowed := allowedOrigins[gotOrigin]; allowed {
			origin = gotOrigin
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
		return handler(w, r, p)
	}
}
