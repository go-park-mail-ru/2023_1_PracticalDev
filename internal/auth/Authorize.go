package auth

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type Handler func(w http.ResponseWriter, r *http.Request, p httprouter.Params) error
type Authorizer func(h Handler) Handler

func NewAuthorizer(serv Service) func(h Handler) Handler {
	return func(handler Handler) Handler {
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

			userId, sessionId := tmp[0], tmp[1]

			if err = serv.CheckAuth(userId, sessionId); err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return err
			}

			return handler(w, r, p)
		}
	}
}
