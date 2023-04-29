package middleware

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"

	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/log"
)

func ErrorHandler(handler func(w http.ResponseWriter, r *http.Request, p httprouter.Params) error, log log.Logger) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		defer func() {
			if err := recover(); err != nil {
				log.Error("after recover: ", err)
			}
		}()

		err := handler(w, r, p)
		if err != nil {
			errCause := errors.Cause(err)

			httpCode, exist := pkgErrors.GetHTTPCodeByError(errCause)
			if !exist {
				errCause = errors.Wrap(errCause, "undefined error")
			}
			w.WriteHeader(httpCode)

			if 200 <= httpCode && httpCode <= 399 {
				log.Info("Response Code:", httpCode, "Message:", "\""+err.Error()+"\"")
			} else {
				log.Error("Response Code:", httpCode, "Error Message:", "\""+err.Error()+"\"")
			}

			if httpCode != http.StatusNoContent {
				_, err = w.Write([]byte(errCause.Error()))
				if err != nil {
					log.Error(err.Error())
				}
			}
		}
	}
}
