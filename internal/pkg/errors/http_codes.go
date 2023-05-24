package errors

import (
	"net/http"
)

var httpCodes = map[error]int{
	// Common delivery
	ErrMissingFile: http.StatusBadRequest,

	// Common repository
	ErrDb:           http.StatusInternalServerError,
	ErrImageService: http.StatusInternalServerError,

	// Invalid Param
	ErrInvalidUserIdParam:  http.StatusBadRequest,
	ErrInvalidBoardIdParam: http.StatusBadRequest,
	ErrInvalidPinIdParam:   http.StatusBadRequest,
	ErrInvalidPageParam:    http.StatusBadRequest,
	ErrInvalidLimitParam:   http.StatusBadRequest,
	ErrInvalidLikedParam:   http.StatusBadRequest,
	ErrInvalidChatIDParam:  http.StatusBadRequest,
	ErrInvalidLinkIDParam:  http.StatusBadRequest,

	ErrBadParams:          http.StatusBadRequest,
	ErrBadRequest:         http.StatusBadRequest,
	ErrBadSessionCookie:   http.StatusBadRequest,
	ErrBadCsrfTokenCookie: http.StatusBadRequest,
	ErrBadTokenTime:       http.StatusBadRequest,
	ErrBadTokenData:       http.StatusBadRequest,
	ErrParseForm:          http.StatusBadRequest,
	ErrParseJson:          http.StatusBadRequest,
	ErrReadBody:           http.StatusBadRequest,
	ErrUserAlreadyExists:  http.StatusBadRequest,
	ErrSameUserId:         http.StatusBadRequest,

	// Auth
	ErrWrongLoginOrPassword: http.StatusNotFound,
	ErrUnauthorized:         http.StatusUnauthorized,

	// WebSocket
	ErrUpgradeToWebSocket: http.StatusBadRequest,

	// Not Found
	ErrUserNotFound:    http.StatusNotFound,
	ErrProfileNotFound: http.StatusNotFound,
	ErrBoardNotFound:   http.StatusNotFound,
	ErrPinNotFound:     http.StatusNotFound,
	ErrChatNotFound:    http.StatusNotFound,
	ErrLinkNotFound:    http.StatusNotFound,

	// Profile
	ErrTooShortUsername: http.StatusBadRequest,
	ErrTooLongUsername:  http.StatusBadRequest,
	ErrEmptyName:        http.StatusBadRequest,
	ErrTooLongName:      http.StatusBadRequest,

	ErrNoContent:         http.StatusNoContent,
	ErrForbidden:         http.StatusForbidden,
	ErrTokenExpired:      http.StatusForbidden,
	ErrLikeNotFound:      http.StatusConflict,
	ErrFollowingNotFound: http.StatusConflict,

	// Already exists
	ErrLikeAlreadyExists:      http.StatusConflict,
	ErrFollowingAlreadyExists: http.StatusConflict,
	ErrChatAlreadyExists:      http.StatusConflict,
	ErrPinAlreadyAdded:        http.StatusConflict,
}

func GetHTTPCodeByError(err error) (int, bool) {
	httpCode, exist := httpCodes[err]
	if !exist {
		httpCode = http.StatusInternalServerError
	}
	return httpCode, exist
}
