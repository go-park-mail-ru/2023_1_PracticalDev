package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models/api"
	"github.com/julienschmidt/httprouter"
)

var (
	BadRequestError  = errors.New("bad request")
	BadSessionCookie = errors.New("bad session cookie")
)

func RegisterHandlers(mux *httprouter.Router, logger log.Logger, serv Service) {
	del := delivery{serv, logger}
	mux.POST("/auth/login", middleware.HandleLogger(middleware.ErrorHandler(del.Authenticate, logger), logger))
	mux.DELETE("/auth/logout", middleware.HandleLogger(middleware.ErrorHandler(del.Logout, logger), logger))
	mux.POST("/auth/signup", middleware.HandleLogger(middleware.ErrorHandler(del.Register, logger), logger))
	mux.GET("/auth/me", middleware.HandleLogger(middleware.ErrorHandler(del.CheckAuth, logger), logger))
}

type delivery struct {
	serv Service
	log  log.Logger
}

func parseSessionCookie(c *http.Cookie) (string, string, error) {
	tmp := strings.Split(c.Value, "$")
	if len(tmp) != 2 {
		return "", "", BadSessionCookie
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
		w.WriteHeader(http.StatusBadRequest)
		return BadRequestError
	}

	user, session, err := del.serv.Authenticate(data.Email, data.Password)
	if err != nil {
		if errors.Is(err, WrongPasswordOrLoginError) {
			w.WriteHeader(http.StatusNotFound)
			return nil
		} else {
			log.New().Error(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return err
	}

	sessionCookie := createSessionCookie(session)

	http.SetCookie(w, sessionCookie)
	encoder := json.NewEncoder(w)

	if err := encoder.Encode(user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		del.log.Error(err)
		return err
	}

	return nil
}

func (del *delivery) CheckAuth(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	sessionCookie, err := r.Cookie("JSESSIONID")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return err
	}

	userId, sessionId, err := parseSessionCookie(sessionCookie)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	user, err := del.serv.CheckAuth(userId, sessionId)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	}

	decoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	if err = decoder.Encode(user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

func (del *delivery) Logout(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	sessionCookie, err := r.Cookie("JSESSIONID")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return err
	}

	userId, sessionId, err := parseSessionCookie(sessionCookie)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return BadRequestError
	}

	newCookie := &http.Cookie{
		Name:     "JSESSIONID",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}

	if err = del.serv.DeleteSession(userId, sessionId); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return err
	}

	http.SetCookie(w, newCookie)
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (del *delivery) Register(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	params := api.RegisterParams{}

	if err := decoder.Decode(&params); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	user, sessionParams, err := del.serv.Register(&params)
	if err != nil {
		switch {
		case errors.Is(err, UserAlreadyExistsError):
			w.WriteHeader(http.StatusBadRequest)
			return nil
		case errors.Is(err, UserCreationError):
			w.WriteHeader(http.StatusInternalServerError)
			return err
		case errors.Is(err, DBConnectionError):
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}
	}

	cookie := createSessionCookie(sessionParams)
	http.SetCookie(w, cookie)
	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		del.log.Error(err)
	}

	return nil
}
