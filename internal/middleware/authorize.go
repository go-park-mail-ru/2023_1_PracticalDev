package middleware

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth/tokens"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/router"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
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
				return ErrUnauthorized
			}

			tmp := strings.Split(sessionCookie.Value, "$")

			if len(tmp) != 2 {
				return ErrUnauthorized
			}

			userId := tmp[0]

			if _, err = serv.CheckAuth(userId, sessionCookie.Value); err != nil {
				return ErrUnauthorized
			}

			csrfToken := r.Header.Get("X-XSRF-TOKEN")
			session := tokens.SessionParams{Token: sessionCookie.Value}
			check, err := token.Check(&session, csrfToken)
			if err != nil || !check {
				log.Warn("Potential CSRF request")
				_, err = w.Write([]byte("{}"))
				return err
			}

			p = append(p, httprouter.Param{Key: "user-id", Value: userId})
			return handler(w, r, p)
		}
	}
}
