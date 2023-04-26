package middleware

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth/tokens"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
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

func NewAuthorizer(serv AuthService, token *tokens.HashToken, log log.Logger) func(h router.Handler) router.Handler {
	return func(handler router.Handler) router.Handler {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
			sessionCookie, err := r.Cookie("JSESSIONID")
			if err != nil {
				return errors.Wrap(pkgErrors.ErrUnauthorized, err.Error())
			}

			tmp := strings.Split(sessionCookie.Value, "$")

			if len(tmp) != 2 {
				return errors.Wrap(pkgErrors.ErrUnauthorized, "invalid cookie")
			}

			userId := tmp[0]

			if _, err = serv.CheckAuth(userId, sessionCookie.Value); err != nil {
				return errors.Wrap(pkgErrors.ErrUnauthorized, err.Error())
			}

			csrfToken := r.Header.Get("X-XSRF-TOKEN")
			session := tokens.SessionParams{Token: sessionCookie.Value}
			check, err := token.Check(&session, csrfToken)
			if err != nil || !check {
				log.Warn("Potential CSRF request. X-XSRF-TOKEN:", "\""+csrfToken+"\"")
				_, err = w.Write([]byte("{}"))
				return err
			}

			p = append(p, httprouter.Param{Key: "user-id", Value: userId})
			return handler(w, r, p)
		}
	}
}
