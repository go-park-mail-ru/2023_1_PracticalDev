package middleware

import (
	"net/http"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/julienschmidt/httprouter"
)

func Logger(handler httprouter.Handle, log log.Logger) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		log.Info("New request: ", r.Method, r.URL, r.Proto)
		handler(w, r, p)
	}
}
