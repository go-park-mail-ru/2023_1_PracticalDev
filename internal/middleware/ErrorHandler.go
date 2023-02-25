package middleware

import (
	"net/http"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/julienschmidt/httprouter"
)

func ErrorHandler(handler func(w http.ResponseWriter, r *http.Request, p httprouter.Params) error, log log.Logger) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		defer func() {
			if err := recover(); err != nil {
				log.Error(err)
			}
		}()

		err := handler(w, r, p)
		if err != nil {
			log.Error(err)
			w.Write([]byte(err.Error()))
		}
	}
}
