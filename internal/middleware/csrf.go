package middleware

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth/tokens"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/router"
)

type CSRFMiddleware func(h router.Handler) router.Handler

func NewCSRFMiddleware(token *tokens.HashToken, log log.Logger) CSRFMiddleware {
	return func(handler router.Handler) router.Handler {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
			sessionCookie, _ := r.Cookie("JSESSIONID")
			csrfToken := r.Header.Get("X-XSRF-TOKEN")
			session := tokens.SessionParams{Token: sessionCookie.Value}

			check, err := token.Check(&session, csrfToken)
			if err != nil || !check {
				log.Warn(fmt.Sprintf("Potential CSRF request: X-XSRF-TOKEN=%s, err=%v", csrfToken, err))
				_, err = w.Write([]byte("{}"))
				return err
			}

			return handler(w, r, p)
		}
	}
}
