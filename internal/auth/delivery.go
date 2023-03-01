package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

func RegisterHandlers(mux *httprouter.Router, logger log.Logger, serv Service) {
	del := delivery{serv, logger}
	mux.POST("/auth/login", middleware.Logger(middleware.ErrorHandler(del.Authenticate, logger), logger))
	mux.DELETE("/auth/logout", middleware.Logger(middleware.ErrorHandler(del.Logout, logger), logger))
	mux.POST("/auth/signup", middleware.Logger(middleware.ErrorHandler(del.Register, logger), logger))
	mux.POST("/check", middleware.Logger(middleware.ErrorHandler(del.CheckAuth, logger), logger))
}

type delivery struct {
	serv Service
	log  log.Logger
}

func (del delivery) Authenticate(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	data := struct {
		Email          string `json:"email"`
		HashedPassword string `json:"hashed_password"`
	}{}

	if err := decoder.Decode(&data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return errors.New("bad request")
	}

	user, err := del.serv.Authenticate(data.Email, data.HashedPassword)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return errors.New("wrong login or password")
	}

	token := uuid.New().String()
	livingTime := 3 * time.Hour
	expiration := time.Now().Add(livingTime)
	cookie := http.Cookie{Name: "JSESSIONID", Value: token, Expires: expiration, HttpOnly: true}
	del.serv.SetSession(token, user, livingTime)
	http.SetCookie(w, &cookie)
	return nil
}

func (del delivery) Logout(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	cookie, err := r.Cookie("JSESSIONID")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return errors.New("cookie not found")
	}

	tmp := http.Cookie{
		Name:   cookie.Name,
		MaxAge: -1,
	}

	if err = del.serv.CheckAuth(cookie.Value); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return err
	}

	del.serv.DeleteSession(cookie.Value)
	http.SetCookie(w, &tmp)
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (del delivery) CheckAuth(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	cookie, err := r.Cookie("JSESSIONID")

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return err
	}

	if err = del.serv.CheckAuth(cookie.Value); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return err
	}

	return nil

}

func (del delivery) Register(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	user := models.User{}

	if err := decoder.Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	if err := del.serv.Register(user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	w.WriteHeader(http.StatusCreated)
	return nil
}
