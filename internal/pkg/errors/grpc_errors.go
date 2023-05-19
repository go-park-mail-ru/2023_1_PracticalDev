package errors

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

var grpcCodes = map[error]codes.Code{
	// Common repository
	ErrDb:           codes.Internal,
	ErrImageService: codes.Internal,

	// Invalid Param
	ErrInvalidUserIdParam:  codes.InvalidArgument,
	ErrInvalidBoardIdParam: codes.InvalidArgument,
	ErrInvalidPinIdParam:   codes.InvalidArgument,
	ErrInvalidPageParam:    codes.InvalidArgument,
	ErrInvalidLimitParam:   codes.InvalidArgument,
	ErrInvalidChatIDParam:  codes.InvalidArgument,
	ErrInvalidLinkIDParam:  codes.InvalidArgument,

	ErrBadParams:          codes.InvalidArgument,
	ErrBadRequest:         codes.InvalidArgument,
	ErrBadSessionCookie:   codes.InvalidArgument,
	ErrBadCsrfTokenCookie: codes.InvalidArgument,
	ErrBadTokenTime:       codes.InvalidArgument,
	ErrBadTokenData:       codes.InvalidArgument,
	ErrParseForm:          codes.InvalidArgument,
	ErrParseJson:          codes.InvalidArgument,
	ErrUserAlreadyExists:  codes.InvalidArgument,
	ErrSameUserId:         codes.InvalidArgument,

	// Auth
	ErrWrongLoginOrPassword: codes.NotFound,
	ErrUnauthorized:         codes.Unauthenticated,

	// WebSocket
	ErrUpgradeToWebSocket: codes.InvalidArgument,

	// Not Found
	ErrUserNotFound:    codes.NotFound,
	ErrProfileNotFound: codes.NotFound,
	ErrBoardNotFound:   codes.NotFound,
	ErrPinNotFound:     codes.NotFound,
	ErrChatNotFound:    codes.NotFound,
	ErrLinkNotFound:    codes.NotFound,

	// Profile
	ErrTooShortUsername: codes.InvalidArgument,
	ErrTooLongUsername:  codes.InvalidArgument,
	ErrEmptyName:        codes.InvalidArgument,
	ErrTooLongName:      codes.InvalidArgument,

	ErrNoContent:         codes.DataLoss,
	ErrForbidden:         codes.PermissionDenied,
	ErrTokenExpired:      codes.PermissionDenied,
	ErrLikeNotFound:      codes.AlreadyExists,
	ErrFollowingNotFound: codes.AlreadyExists,

	// Already exists
	ErrLikeAlreadyExists:      codes.AlreadyExists,
	ErrFollowingAlreadyExists: codes.AlreadyExists,
	ErrChatAlreadyExists:      codes.AlreadyExists,
	ErrPinAlreadyAdded:        codes.AlreadyExists,
}

func GetGRPCCodeByError(err error) (codes.Code, bool) {
	httpCode, exist := grpcCodes[err]
	if !exist {
		httpCode = http.StatusInternalServerError
	}
	return httpCode, exist
}
