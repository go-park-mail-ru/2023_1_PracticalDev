package middleware

import (
	"github.com/pkg/errors"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
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
			httpCode, exist := pkgErrors.GetHTTPCodeByError(err)
			if !exist {
				err = errors.Wrap(err, "undefined error")
			}
			log.Error(err.Error())
			w.WriteHeader(httpCode)

			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				log.Error(err.Error())
			}
		}
	}
}
