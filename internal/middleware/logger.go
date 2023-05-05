package middleware

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

func HandleLogger(handler httprouter.Handle, log *zap.Logger) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		log.Info("New request",
			zap.String("method", r.Method),
			zap.String("url", r.URL.String()),
			zap.String("protocol", r.Proto),
			zap.String("origin", r.Header.Get("Origin")))
		handler(w, r, p)
	}
}

func HandlerFuncLogger(handler http.HandlerFunc, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("New request",
			zap.String("method", r.Method),
			zap.String("url", r.URL.String()),
			zap.String("protocol", r.Proto),
			zap.String("origin", r.Header.Get("Origin")))
		handler(w, r)
	}
}
