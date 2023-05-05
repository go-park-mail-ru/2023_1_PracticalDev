package middleware

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
)

func ErrorHandler(handler func(w http.ResponseWriter, r *http.Request, p httprouter.Params) error, log *zap.Logger) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		defer func() {
			if err := recover(); err != nil {
				log.Error("after recover: ", zap.Any("error", err))
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

			if httpCode >= 500 {
				log.Error("Internal Server Error", zap.Int("http_code", httpCode), zap.String("error", err.Error()))
			} else {
				log.Info("Response", zap.Int("http_code", httpCode), zap.String("message", errCause.Error()))
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
