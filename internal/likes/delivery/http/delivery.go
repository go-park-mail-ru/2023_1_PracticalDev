package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	_likes "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/likes"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	"github.com/julienschmidt/httprouter"
)

var (
	ErrInvalidPinId  = errors.New("invalid pin id")
	ErrInvalidUserId = errors.New("invalid user id")
	ErrPinNotFound   = errors.New("pin not found")
	ErrService       = errors.New("service error")
)

func RegisterHandlers(mux *httprouter.Router, logger log.Logger, authorizer middleware.Authorizer, serv _likes.Service) {
	del := delivery{serv, logger}

	mux.POST("/pins/:id/like", middleware.HandleLogger(middleware.ErrorHandler(middleware.CorsChecker(authorizer(del.like)), logger), logger))
	mux.DELETE("/pins/:id/like", middleware.HandleLogger(middleware.ErrorHandler(middleware.CorsChecker(authorizer(del.unlike)), logger), logger))

	mux.GET("/pins/:id/likes", middleware.HandleLogger(middleware.ErrorHandler(middleware.CorsChecker(authorizer(del.listByPin)), logger), logger))
	mux.GET("/users/:id/likes", middleware.HandleLogger(middleware.ErrorHandler(middleware.CorsChecker(authorizer(del.listByAuthor)), logger), logger))
}

type delivery struct {
	serv _likes.Service
	log  log.Logger
}

func (del *delivery) like(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("user-id")
	userId, err := strconv.Atoi(strId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return ErrInvalidUserId
	}

	strId = p.ByName("id")
	pinId, err := strconv.Atoi(strId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return ErrInvalidPinId
	}

	err = del.serv.Like(pinId, userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	return err
}

func (del *delivery) unlike(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("user-id")
	userId, err := strconv.Atoi(strId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return ErrInvalidUserId
	}

	strId = p.ByName("id")
	pinId, err := strconv.Atoi(strId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return ErrInvalidPinId
	}

	err = del.serv.Unlike(pinId, userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	return err
}

func (del *delivery) listByPin(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("id")
	pinId, err := strconv.Atoi(strId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return ErrInvalidPinId
	}

	likes, err := del.serv.ListByPin(pinId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	response := listByPinResponse{
		Likes: likes,
	}
	data, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	return err
}

func (del *delivery) listByAuthor(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("id")
	userId, err := strconv.Atoi(strId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return ErrInvalidUserId
	}

	likes, err := del.serv.ListByAuthor(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	response := listByAuthorResponse{
		Likes: likes,
	}
	data, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	return err
}
