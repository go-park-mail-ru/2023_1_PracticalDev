package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models/api"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

var BadRequestError = errors.New("bad request")

func RegisterHandlers(mux *httprouter.Router, logger log.Logger, serv Service) {
	del := delivery{serv, logger}
	mux.POST("/auth/login", middleware.HandleLogger(middleware.ErrorHandler(del.Authenticate, logger), logger))
	mux.DELETE("/auth/logout", middleware.HandleLogger(middleware.ErrorHandler(del.Logout, logger), logger))
	mux.POST("/auth/signup", middleware.HandleLogger(middleware.ErrorHandler(del.Register, logger), logger))
}

type delivery struct {
	serv Service
	log  log.Logger
}

func (del *delivery) Authenticate(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	data := api.LoginParams{}

	if err := decoder.Decode(&data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return BadRequestError
	}

	user, err := del.serv.Authenticate(data.Email, data.Password)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return WrongPasswordOrLoginError
	}

	token := uuid.New().String()
	livingTime := 3 * time.Hour
	expiration := time.Now().Add(livingTime)
	sessionCookie := http.Cookie{Name: "JSESSIONID", Value: strconv.Itoa(user.Id) + "$" + token, Expires: expiration, HttpOnly: true, Path: "/"}

	session := models.Session{
		UserId:    user.Id,
		UserEmail: user.Email,
	}

	if err := del.serv.SetSession(token, &session, livingTime); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	http.SetCookie(w, &sessionCookie)
	usr, _ := json.Marshal(user)

	_, err = w.Write(usr)
	return err
}

func (del *delivery) Logout(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	sessionCookie, err := r.Cookie("JSESSIONID")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return err
	}

	tmp := strings.Split(sessionCookie.Value, "$")
	if len(tmp) != 2 {
		w.WriteHeader(http.StatusBadRequest)
		return BadRequestError
	}

	userId, sessionId := tmp[0], tmp[1]

	newCookie := http.Cookie{
		Name:   sessionCookie.Name,
		MaxAge: -1,
	}

	if err = del.serv.DeleteSession(userId, sessionId); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return err
	}

	http.SetCookie(w, &newCookie)
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (del *delivery) Register(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	user := api.RegisterParams{}

	if err := decoder.Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	if err := del.serv.Register(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}
