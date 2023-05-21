package middleware

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/router"
)

type (
	Authorizer func(h router.Handler) router.Handler

	AuthService interface {
		CheckAuth(userId, sessionId string) (models.User, error)
	}
)

func NewAuthorizer(serv AuthService, log *zap.Logger) Authorizer {
	return func(handler router.Handler) router.Handler {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
			sessionCookie, err := r.Cookie("JSESSIONID")
			if err != nil {
				log.Debug("Failed to get session cookie", zap.Error(err))
				return errors.Wrap(pkgErrors.ErrUnauthorized, err.Error())
			}

			tmp := strings.Split(sessionCookie.Value, "$")
			if len(tmp) != 2 {
				log.Debug("Failed to split session cookie", zap.String("cookie", sessionCookie.String()))
				return errors.Wrap(pkgErrors.ErrUnauthorized, "invalid cookie")
			}

			userId := tmp[0]
			if _, err = serv.CheckAuth(userId, sessionCookie.Value); err != nil {
				return errors.Wrap(pkgErrors.ErrUnauthorized, err.Error())
			}

			p = append(p, httprouter.Param{Key: "user-id", Value: userId})
			return handler(w, r, p)
		}
	}
}
