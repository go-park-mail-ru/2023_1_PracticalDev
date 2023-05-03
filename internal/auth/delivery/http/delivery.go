package http

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth/tokens"
	mw "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/log"
)

func RegisterHandlers(mux *httprouter.Router, logger log.Logger, serv auth.Service, token *tokens.HashToken) {
	del := delivery{serv, logger, token}
	mux.POST("/auth/login", mw.HandleLogger(mw.ErrorHandler(del.Authenticate, logger), logger))
	mux.DELETE("/auth/logout", mw.HandleLogger(mw.ErrorHandler(del.Logout, logger), logger))
	mux.POST("/auth/signup", mw.HandleLogger(mw.ErrorHandler(del.Register, logger), logger))
	mux.GET("/auth/me", mw.HandleLogger(mw.ErrorHandler(del.CheckAuth, logger), logger))
}

type delivery struct {
	serv  auth.Service
	log   log.Logger
	token *tokens.HashToken
}

func parseSessionCookie(c *http.Cookie) (string, string, error) {
	tmp := strings.Split(c.Value, "$")
	if len(tmp) != 2 {
		return "", "", pkgErrors.ErrBadSessionCookie
	}

	return tmp[0], c.Value, nil
}

func createSessionCookie(s auth.SessionParams) *http.Cookie {
	return &http.Cookie{
		Name:     "JSESSIONID",
		Value:    s.Token,
		Expires:  time.Now().Add(s.LivingTime),
		HttpOnly: true,
		Path:     "/",
	}
}

func createCsrfTokenCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:     "XSRF-TOKEN",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 5),
		HttpOnly: false,
		Path:     "/",
	}
}

func (del *delivery) Authenticate(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	data := auth.LoginParams{}
	if err := decoder.Decode(&data); err != nil {
		return pkgErrors.ErrBadRequest
	}

	user, session, err := del.serv.Authenticate(data.Email, data.Password)
	if err != nil {
		return err
	}

	sessionCookie := createSessionCookie(session)
	http.SetCookie(w, sessionCookie)

	token, err := del.token.Create(&tokens.SessionParams{Token: session.Token}, time.Now().Add(session.LivingTime).Unix())
	if err != nil {
		del.log.Error("csrf token creation error:", err)
		return pkgErrors.ErrCreateCsrfToken
	}

	csrfCookie := createCsrfTokenCookie(token)
	http.SetCookie(w, csrfCookie)

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if err = encoder.Encode(user); err != nil {
		return pkgErrors.ErrCreateResponse
	}
	return nil
}

func (del *delivery) CheckAuth(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	sessionCookie, err := r.Cookie("JSESSIONID")
	if err != nil {
		return pkgErrors.ErrBadSessionCookie
	}

	userId, sessionId, err := parseSessionCookie(sessionCookie)
	if err != nil {
		return pkgErrors.ErrBadRequest
	}

	user, err := del.serv.CheckAuth(userId, sessionId)
	if err != nil {
		return pkgErrors.ErrUnauthorized
	}

	decoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	if err = decoder.Encode(user); err != nil {
		return pkgErrors.ErrCreateResponse
	}

	return nil
}

func (del *delivery) Logout(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	sessionCookie, err := r.Cookie("JSESSIONID")
	if err != nil {
		return pkgErrors.ErrUserNotFound
	}

	userId, sessionId, err := parseSessionCookie(sessionCookie)
	if err != nil {
		return pkgErrors.ErrBadRequest
	}

	if err = del.serv.DeleteSession(userId, sessionId); err != nil {
		return pkgErrors.ErrUnauthorized
	}

	newCookie := &http.Cookie{
		Name:     "JSESSIONID",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
	http.SetCookie(w, newCookie)
	return pkgErrors.ErrNoContent
}

func (del *delivery) Register(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	decoder := json.NewDecoder(r.Body)
	defer func() {
		err := r.Body.Close()
		if err != nil {
			del.log.Error(err)
		}
	}()
	params := auth.RegisterParams{}

	if err := decoder.Decode(&params); err != nil {
		return pkgErrors.ErrParseJson
	}

	user, sessionParams, err := del.serv.Register(&params)
	if err != nil {
		return err
	}

	cookie := createSessionCookie(sessionParams)
	http.SetCookie(w, cookie)

	token, err := del.token.Create(&tokens.SessionParams{Token: sessionParams.Token}, time.Now().Add(sessionParams.LivingTime).Unix())
	if err != nil {
		del.log.Error("csrf token creation error:", err)
		return pkgErrors.ErrCreateCsrfToken
	}

	csrfCookie := createCsrfTokenCookie(token)
	http.SetCookie(w, csrfCookie)

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if err = encoder.Encode(user); err != nil {
		return pkgErrors.ErrCreateResponse
	}

	return nil
}
