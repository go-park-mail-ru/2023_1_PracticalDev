package middleware

import (
	"net/http"
	"strings"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/router"
	"github.com/julienschmidt/httprouter"
)

type (
	Authorizer func(h router.Handler) router.Handler

	AuthService interface {
		CheckAuth(userId, sessionId string) (models.User, error)
	}
)

func NewAuthorizer(serv AuthService) func(h router.Handler) router.Handler {
	return func(handler router.Handler) router.Handler {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
			sessionCookie, err := r.Cookie("JSESSIONID")
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return err
			}

			tmp := strings.Split(sessionCookie.Value, "$")

			if len(tmp) != 2 {
				w.WriteHeader(http.StatusUnauthorized)
				return http.ErrNoCookie
			}

			userId := tmp[0]

			if _, err = serv.CheckAuth(userId, sessionCookie.Value); err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return err
			}

			return handler(w, r, p)
		}
	}
}
