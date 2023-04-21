package http

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth/tokens"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	mw "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	"github.com/julienschmidt/httprouter"
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
		return "", "", mw.ErrBadSessionCookie
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
	data := LoginParams{}
	if err := decoder.Decode(&data); err != nil {
		return mw.ErrBadRequest
	}

	user, session, err := del.serv.Authenticate(data.Email, data.Password)
	if err != nil {
		if errors.Is(err, auth.WrongPasswordOrLoginError) {
			return mw.ErrUserNotFound
		} else {
			return mw.ErrService
		}
	}

	sessionCookie := createSessionCookie(session)
	http.SetCookie(w, sessionCookie)

	token, err := del.token.Create(&tokens.SessionParams{Token: session.Token}, time.Now().Add(session.LivingTime).Unix())
	if err != nil {
		del.log.Error("csrf token creation error:", err)
		return mw.ErrCreateCsrfToken
	}

	csrfCookie := createCsrfTokenCookie(token)
	http.SetCookie(w, csrfCookie)

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if err = encoder.Encode(user); err != nil {
		return mw.ErrCreateResponse
	}
	return nil
}

func (del *delivery) CheckAuth(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	sessionCookie, err := r.Cookie("JSESSIONID")
	if err != nil {
		return mw.ErrBadSessionCookie
	}

	userId, sessionId, err := parseSessionCookie(sessionCookie)
	if err != nil {
		return mw.ErrBadRequest
	}

	user, err := del.serv.CheckAuth(userId, sessionId)
	if err != nil {
		return mw.ErrUnauthorized
	}

	decoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	if err = decoder.Encode(user); err != nil {
		return mw.ErrCreateResponse
	}

	return nil
}

func (del *delivery) Logout(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	sessionCookie, err := r.Cookie("JSESSIONID")
	if err != nil {
		return mw.ErrUserNotFound
	}

	userId, sessionId, err := parseSessionCookie(sessionCookie)
	if err != nil {
		return mw.ErrBadRequest
	}

	if err = del.serv.DeleteSession(userId, sessionId); err != nil {
		return mw.ErrUnauthorized
	}

	newCookie := &http.Cookie{
		Name:     "JSESSIONID",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
	http.SetCookie(w, newCookie)
	return mw.ErrNoContent
}

func (del *delivery) Register(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	decoder := json.NewDecoder(r.Body)
	defer func() {
		err := r.Body.Close()
		if err != nil {
			del.log.Error(err)
		}
	}()
	params := RegisterParams{}

	if err := decoder.Decode(&params); err != nil {
		return mw.ErrParseJson
	}

	user, sessionParams, err := del.serv.Register(&params)
	if err != nil {
		switch {
		case errors.Is(err, auth.UserAlreadyExistsError):
			return mw.ErrUserAlreadyExists
		case errors.Is(err, auth.UserCreationError):
			return mw.ErrService
		case errors.Is(err, auth.DBConnectionError):
			return mw.ErrService
		}
	}

	cookie := createSessionCookie(sessionParams)
	http.SetCookie(w, cookie)
	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	if err = encoder.Encode(user); err != nil {
		return mw.ErrCreateResponse
	}

	return nil
}
