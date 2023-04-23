package errors

import "github.com/pkg/errors"

var (
	// Common delivery
	ErrMissingFile = errors.New("missing file")

	// Common repository
	ErrDb = errors.New("db error")

	// Not found
	ErrUserNotFound      = errors.New("user not found")
	ErrProfileNotFound   = errors.New("profile not found")
	ErrBoardNotFound     = errors.New("board not found")
	ErrPinNotFound       = errors.New("pin not found")
	ErrLikeNotFound      = errors.New("no such like")
	ErrFollowingNotFound = errors.New("no such following")

	// CSRF
	ErrBadCsrfTokenCookie = errors.New("bad csrf token cookie")
	ErrBadTokenTime       = errors.New("bad token time")
	ErrBadTokenData       = errors.New("bad token data")

	ErrInvalidUserIdParam     = errors.New("invalid user id param")
	ErrInvalidBoardIdParam    = errors.New("invalid board id param")
	ErrInvalidPinIdParam      = errors.New("invalid pin id param")
	ErrInvalidPageParam       = errors.New("invalid page param")
	ErrInvalidLimitParam      = errors.New("invalid limit param")
	ErrBadParams              = errors.New("bad params")
	ErrBadRequest             = errors.New("bad request")
	ErrBadSessionCookie       = errors.New("bad session cookie")
	ErrFileCopy               = errors.New("file copy error")
	ErrParseForm              = errors.New("parse form error")
	ErrParseJson              = errors.New("parse json error")
	ErrUserAlreadyExists      = errors.New("user already exists")
	ErrSameUserId             = errors.New("same user id: user cannot follow himself")
	ErrService                = errors.New("service error")
	ErrCreateResponse         = errors.New("create response error")
	ErrCreateCsrfToken        = errors.New("create csrf token error")
	ErrUnauthorized           = errors.New("unauthorized")
	ErrNoContent              = errors.New("no content")
	ErrForbidden              = errors.New("access denied")
	ErrTokenExpired           = errors.New("token expired")
	ErrLikeAlreadyExists      = errors.New("like already exists")
	ErrFollowingAlreadyExists = errors.New("following already exists")
	ErrPinAlreadyAdded        = errors.New("pin already added")
)
