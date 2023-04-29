package middleware

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/log"
)

func HandleLogger(handler httprouter.Handle, log log.Logger) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		log.Info("New request: ", r.Method, r.URL, r.Proto, "Origin="+r.Header.Get("Origin"))
		handler(w, r, p)
	}
}

func HandlerFuncLogger(handler http.HandlerFunc, log log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("New request: ", r.Method, r.URL, r.Proto, "Origin="+r.Header.Get("Origin"))
		handler(w, r)
	}
}
