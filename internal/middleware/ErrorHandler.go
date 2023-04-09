package middleware

import (
	"github.com/pkg/errors"
	"net/http"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/julienschmidt/httprouter"
)

var (
	ErrMissingFile         = errors.New("missing file")
	ErrInvalidUserIdParam  = errors.New("invalid user id param")
	ErrInvalidBoardIdParam = errors.New("invalid board id param")
	ErrInvalidPinIdParam   = errors.New("invalid pin id param")
	ErrInvalidPageParam    = errors.New("invalid page param")
	ErrInvalidLimitParam   = errors.New("invalid limit param")
	ErrProfileNotFound     = errors.New("profile not found")
	ErrBoardNotFound       = errors.New("board not found")
	ErrPinNotFound         = errors.New("pin not found")
	ErrUserNotFound        = errors.New("user not found")
	ErrBadParams           = errors.New("bad params")
	ErrBadRequest          = errors.New("bad request")
	ErrBadSessionCookie    = errors.New("bad session cookie")
	ErrFileCopy            = errors.New("file copy error")
	ErrParseForm           = errors.New("parse form error")
	ErrParseJson           = errors.New("parse json error")
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrService             = errors.New("service error")
	ErrCreateResponse      = errors.New("create response error")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrNoContent           = errors.New("no content")
	ErrForbidden           = errors.New("access denied")
)

var statusCodes = map[error]int{
	ErrMissingFile:         http.StatusBadRequest,
	ErrInvalidUserIdParam:  http.StatusBadRequest,
	ErrInvalidBoardIdParam: http.StatusBadRequest,
	ErrInvalidPinIdParam:   http.StatusBadRequest,
	ErrInvalidPageParam:    http.StatusBadRequest,
	ErrInvalidLimitParam:   http.StatusBadRequest,
	ErrBadParams:           http.StatusBadRequest,
	ErrBadRequest:          http.StatusBadRequest,
	ErrBadSessionCookie:    http.StatusBadRequest,
	ErrParseForm:           http.StatusBadRequest,
	ErrParseJson:           http.StatusBadRequest,
	ErrUserAlreadyExists:   http.StatusBadRequest,
	ErrProfileNotFound:     http.StatusNotFound,
	ErrBoardNotFound:       http.StatusNotFound,
	ErrPinNotFound:         http.StatusNotFound,
	ErrUserNotFound:        http.StatusNotFound,
	ErrUnauthorized:        http.StatusUnauthorized,
	ErrNoContent:           http.StatusNoContent,
	ErrForbidden:           http.StatusForbidden,
}

func ErrorHandler(handler func(w http.ResponseWriter, r *http.Request, p httprouter.Params) error, log log.Logger) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		defer func() {
			if err := recover(); err != nil {
				log.Error(err)
			}
		}()

		err := handler(w, r, p)
		if err != nil {
			code, ok := statusCodes[err]
			if !ok {
				w.WriteHeader(http.StatusInternalServerError)
				log.Error(err.Error())
			} else {
				w.WriteHeader(code)
			}

			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				log.Error(err.Error())
			}
		}
	}
}
