package errors

import (
	"fmt"

	"github.com/pkg/errors"
	"google.golang.org/grpc/status"
)

var (
	// Common delivery
	ErrMissingFile = errors.New("missing file")

	// Common repository
	ErrDb           = errors.New("db error")
	ErrImageService = errors.New("image service error")

	// Profile
	ErrTooShortUsername = errors.New("username must be at least 4 characters")
	ErrTooLongUsername  = errors.New("username must be no more than 30 characters")
	ErrEmptyName        = errors.New("name must not be empty")
	ErrTooLongName      = errors.New("name must be no more than 60 characters")

	// Not found
	ErrUserNotFound      = errors.New("user not found")
	ErrProfileNotFound   = errors.New("profile not found")
	ErrBoardNotFound     = errors.New("board not found")
	ErrPinNotFound       = errors.New("pin not found")
	ErrLikeNotFound      = errors.New("no such like")
	ErrFollowingNotFound = errors.New("no such following")
	ErrChatNotFound      = errors.New("chat not found")
	ErrLinkNotFound      = errors.New("link not found")

	// CSRF
	ErrBadCsrfTokenCookie = errors.New("bad csrf token cookie")
	ErrBadTokenTime       = errors.New("bad token time")
	ErrBadTokenData       = errors.New("bad token data")

	// Auth
	ErrWrongLoginOrPassword = errors.New("wrong login or password")
	ErrUserAlreadyExists    = errors.New("user already exists")

	// Invalid Param
	ErrInvalidUserIdParam  = errors.New("invalid user id param")
	ErrInvalidBoardIdParam = errors.New("invalid board id param")
	ErrInvalidPinIdParam   = errors.New("invalid pin id param")
	ErrInvalidPageParam    = errors.New("invalid page param")
	ErrInvalidLimitParam   = errors.New("invalid limit param")
	ErrInvalidPrivacy      = errors.New("invalid privacy")
	ErrInvalidChatIDParam  = errors.New("invalid chat id param")
	ErrInvalidLinkIDParam  = errors.New("invalid link param")

	// WebSocket
	ErrUpgradeToWebSocket = errors.New("failed to upgrade protocol to websocket")

	ErrBadParams              = errors.New("bad params")
	ErrBadRequest             = errors.New("bad request")
	ErrBadSessionCookie       = errors.New("bad session cookie")
	ErrFileCopy               = errors.New("file copy error")
	ErrParseForm              = errors.New("parse form error")
	ErrParseJson              = errors.New("parse json error")
	ErrSameUserId             = errors.New("same user id")
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
	ErrChatAlreadyExists      = errors.New("chat already exists")
)

var ErrorsByNames = map[string]error{
	// Common delivery
	ErrMissingFile.Error(): ErrMissingFile,

	// Common repository
	ErrDb.Error():           ErrDb,
	ErrImageService.Error(): ErrImageService,

	// Profile
	ErrTooShortUsername.Error(): ErrTooShortUsername,
	ErrTooLongUsername.Error():  ErrTooLongUsername,
	ErrEmptyName.Error():        ErrEmptyName,
	ErrTooLongName.Error():      ErrTooLongName,

	// Not found
	ErrUserNotFound.Error():      ErrUserNotFound,
	ErrProfileNotFound.Error():   ErrProfileNotFound,
	ErrBoardNotFound.Error():     ErrBoardNotFound,
	ErrPinNotFound.Error():       ErrPinNotFound,
	ErrLikeNotFound.Error():      ErrLikeNotFound,
	ErrFollowingNotFound.Error(): ErrFollowingNotFound,
	ErrChatNotFound.Error():      ErrChatNotFound,
	ErrLinkNotFound.Error():      ErrLinkNotFound,

	// CSRF
	ErrBadCsrfTokenCookie.Error(): ErrBadCsrfTokenCookie,
	ErrBadTokenTime.Error():       ErrBadTokenTime,
	ErrBadTokenData.Error():       ErrBadTokenData,

	// Auth
	ErrWrongLoginOrPassword.Error(): ErrWrongLoginOrPassword,
	ErrUserAlreadyExists.Error():    ErrUserAlreadyExists,

	// Invalid Param
	ErrInvalidUserIdParam.Error():  ErrInvalidUserIdParam,
	ErrInvalidBoardIdParam.Error(): ErrInvalidBoardIdParam,
	ErrInvalidPinIdParam.Error():   ErrInvalidPinIdParam,
	ErrInvalidPageParam.Error():    ErrInvalidPageParam,
	ErrInvalidLimitParam.Error():   ErrInvalidLimitParam,
	ErrInvalidPrivacy.Error():      ErrInvalidPrivacy,
	ErrInvalidChatIDParam.Error():  ErrInvalidChatIDParam,
	ErrInvalidLinkIDParam.Error():  ErrInvalidLinkIDParam,

	// WebSocket
	ErrUpgradeToWebSocket.Error(): ErrUpgradeToWebSocket,

	ErrBadParams.Error():              ErrBadParams,
	ErrBadRequest.Error():             ErrBadRequest,
	ErrBadSessionCookie.Error():       ErrBadSessionCookie,
	ErrFileCopy.Error():               ErrFileCopy,
	ErrParseForm.Error():              ErrParseForm,
	ErrParseJson.Error():              ErrParseJson,
	ErrSameUserId.Error():             ErrSameUserId,
	ErrService.Error():                ErrService,
	ErrCreateResponse.Error():         ErrCreateResponse,
	ErrCreateCsrfToken.Error():        ErrCreateCsrfToken,
	ErrUnauthorized.Error():           ErrUnauthorized,
	ErrNoContent.Error():              ErrNoContent,
	ErrForbidden.Error():              ErrForbidden,
	ErrTokenExpired.Error():           ErrTokenExpired,
	ErrLikeAlreadyExists.Error():      ErrLikeAlreadyExists,
	ErrFollowingAlreadyExists.Error(): ErrFollowingAlreadyExists,
	ErrPinAlreadyAdded.Error():        ErrPinAlreadyAdded,
	ErrChatAlreadyExists.Error():      ErrChatAlreadyExists,
}

type ErrRepositoryQuery struct {
	Func   string // the failing function name
	Query  string // query
	Params []any  // query params
	Err    error  // original error
}

func (e ErrRepositoryQuery) Error() string {
	return fmt.Sprintf("%s: query [%s], params %v, error [%s]", e.Func, e.Query, e.Params, e.Err)
}

func GRPCWrapper(err error) error {
	errCause := errors.Cause(err).Error()
	code, _ := GetGRPCCodeByError(err)
	return status.Error(code, errCause)
}

func GRPCUnwrapper(err error) error {
	s, _ := status.FromError(errors.Cause(err))

	return errors.New(s.Message())
}

func RestoreHTTPError(err error) error {
	return ErrorsByNames[err.Error()]
}
