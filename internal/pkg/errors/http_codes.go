package errors

import (
	"net/http"
)

var httpCodes = map[error]int{
	ErrMissingFile:            http.StatusBadRequest,
	ErrInvalidUserIdParam:     http.StatusBadRequest,
	ErrInvalidBoardIdParam:    http.StatusBadRequest,
	ErrInvalidPinIdParam:      http.StatusBadRequest,
	ErrInvalidPageParam:       http.StatusBadRequest,
	ErrInvalidLimitParam:      http.StatusBadRequest,
	ErrBadParams:              http.StatusBadRequest,
	ErrBadRequest:             http.StatusBadRequest,
	ErrBadSessionCookie:       http.StatusBadRequest,
	ErrBadCsrfTokenCookie:     http.StatusBadRequest,
	ErrBadTokenTime:           http.StatusBadRequest,
	ErrBadTokenData:           http.StatusBadRequest,
	ErrParseForm:              http.StatusBadRequest,
	ErrParseJson:              http.StatusBadRequest,
	ErrUserAlreadyExists:      http.StatusBadRequest,
	ErrSameUserId:             http.StatusBadRequest,
	ErrProfileNotFound:        http.StatusNotFound,
	ErrBoardNotFound:          http.StatusNotFound,
	ErrPinNotFound:            http.StatusNotFound,
	ErrUserNotFound:           http.StatusNotFound,
	ErrUnauthorized:           http.StatusUnauthorized,
	ErrNoContent:              http.StatusNoContent,
	ErrForbidden:              http.StatusForbidden,
	ErrTokenExpired:           http.StatusForbidden,
	ErrLikeAlreadyExists:      http.StatusConflict,
	ErrFollowingAlreadyExists: http.StatusConflict,
	ErrLikeNotFound:           http.StatusConflict,
	ErrFollowingNotFound:      http.StatusConflict,
	ErrPinAlreadyAdded:        http.StatusConflict,
}

func GetHTTPCodeByError(err error) (int, bool) {
	httpCode, exist := httpCodes[err]
	if !exist {
		httpCode = http.StatusInternalServerError
	}
	return httpCode, exist
}
