package auth

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	mw "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models/api"
	"github.com/julienschmidt/httprouter"
)

func RegisterHandlers(mux *httprouter.Router, logger log.Logger, serv Service) {
	del := delivery{serv, logger}
	mux.POST("/auth/login", mw.HandleLogger(mw.ErrorHandler(del.Authenticate, logger), logger))
	mux.DELETE("/auth/logout", mw.HandleLogger(mw.ErrorHandler(del.Logout, logger), logger))
	mux.POST("/auth/signup", mw.HandleLogger(mw.ErrorHandler(del.Register, logger), logger))
	mux.GET("/auth/me", mw.HandleLogger(mw.ErrorHandler(del.CheckAuth, logger), logger))
}

type delivery struct {
	serv Service
	log  log.Logger
}

func parseSessionCookie(c *http.Cookie) (string, string, error) {
	tmp := strings.Split(c.Value, "$")
	if len(tmp) != 2 {
		return "", "", mw.ErrBadSessionCookie
	}

	return tmp[0], c.Value, nil
}

func createSessionCookie(s SessionParams) *http.Cookie {
	return &http.Cookie{
		Name:     "JSESSIONID",
		Value:    s.token,
		Expires:  time.Now().Add(s.livingTime),
		HttpOnly: true,
		Path:     "/",
	}
}

func (del *delivery) Authenticate(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	data := api.LoginParams{}

	if err := decoder.Decode(&data); err != nil {
		return mw.ErrBadRequest
	}

	user, session, err := del.serv.Authenticate(data.Email, data.Password)
	if err != nil {
		if errors.Is(err, WrongPasswordOrLoginError) {
			return mw.ErrUserNotFound
		} else {
			del.log.Error(err)
			return mw.ErrService
		}
	}

	sessionCookie := createSessionCookie(session)

	http.SetCookie(w, sessionCookie)
	encoder := json.NewEncoder(w)

	if err := encoder.Encode(user); err != nil {
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

	newCookie := &http.Cookie{
		Name:     "JSESSIONID",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}

	if err = del.serv.DeleteSession(userId, sessionId); err != nil {
		return mw.ErrUnauthorized
	}

	http.SetCookie(w, newCookie)
	return mw.ErrNoContent
}

func (del *delivery) Register(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	params := api.RegisterParams{}

	if err := decoder.Decode(&params); err != nil {
		return mw.ErrParseJson
	}

	user, sessionParams, err := del.serv.Register(&params)
	if err != nil {
		switch {
		case errors.Is(err, UserAlreadyExistsError):
			return mw.ErrUserAlreadyExists
		case errors.Is(err, UserCreationError):
			return mw.ErrService
		case errors.Is(err, DBConnectionError):
			return mw.ErrService
		}
	}

	cookie := createSessionCookie(sessionParams)
	http.SetCookie(w, cookie)
	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(user); err != nil {
		return mw.ErrCreateResponse
	}

	return nil
}
