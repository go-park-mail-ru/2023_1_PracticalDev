package middleware

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/router"
)

func SetUserID(handler router.Handler) router.Handler {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
		sessionCookie, err := r.Cookie("JSESSIONID")
		if err != nil {
			return handler(w, r, p)
		}

		tmp := strings.Split(sessionCookie.Value, "$")
		if len(tmp) != 2 {
			return handler(w, r, p)
		}

		p = append(p, httprouter.Param{Key: "user-id", Value: tmp[0]})
		return handler(w, r, p)
	}
}
