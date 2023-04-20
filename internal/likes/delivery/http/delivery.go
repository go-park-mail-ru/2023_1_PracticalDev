package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	pkgLikes "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/likes"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	mw "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
)

func RegisterHandlers(mux *httprouter.Router, logger log.Logger, authorizer mw.Authorizer, serv pkgLikes.Service) {
	del := delivery{serv, logger}

	mux.POST("/pins/:id/like", mw.HandleLogger(mw.ErrorHandler(mw.Cors(authorizer(del.like)), logger), logger))
	mux.DELETE("/pins/:id/like", mw.HandleLogger(mw.ErrorHandler(mw.Cors(authorizer(del.unlike)), logger), logger))

	mux.GET("/pins/:id/likes", mw.HandleLogger(mw.ErrorHandler(mw.Cors(authorizer(del.listByPin)), logger), logger))
	mux.GET("/users/:id/likes", mw.HandleLogger(mw.ErrorHandler(mw.Cors(authorizer(del.listByAuthor)), logger), logger))
}

type delivery struct {
	serv pkgLikes.Service
	log  log.Logger
}

func (del *delivery) like(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("user-id")
	userId, err := strconv.Atoi(strId)
	if err != nil {
		return mw.ErrInvalidUserIdParam
	}

	strId = p.ByName("id")
	pinId, err := strconv.Atoi(strId)
	if err != nil {
		return mw.ErrInvalidPinIdParam
	}

	err = del.serv.Like(pinId, userId)
	if err != nil {
		switch err {
		case pkgLikes.ErrLikeAlreadyExists:
			return mw.ErrLikeAlreadyExists
		case pkgLikes.ErrPinNotFound:
			return mw.ErrPinNotFound
		case pkgLikes.ErrAuthorNotFound:
			return mw.ErrUserNotFound
		default:
			return mw.ErrService
		}
	}
	return mw.ErrNoContent
}

func (del *delivery) unlike(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("user-id")
	userId, err := strconv.Atoi(strId)
	if err != nil {
		return mw.ErrInvalidUserIdParam
	}

	strId = p.ByName("id")
	pinId, err := strconv.Atoi(strId)
	if err != nil {
		return mw.ErrInvalidPinIdParam
	}

	err = del.serv.Unlike(pinId, userId)
	if err != nil {
		switch err {
		case pkgLikes.ErrLikeNotFound:
			return mw.ErrLikeNotFound
		case pkgLikes.ErrPinNotFound:
			return mw.ErrPinNotFound
		case pkgLikes.ErrAuthorNotFound:
			return mw.ErrUserNotFound
		default:
			return mw.ErrService
		}
	}
	return mw.ErrNoContent
}

func (del *delivery) listByPin(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("id")
	pinId, err := strconv.Atoi(strId)
	if err != nil {
		return mw.ErrInvalidPinIdParam
	}

	likes, err := del.serv.ListByPin(pinId)
	if err != nil {
		return err
	}

	response := listByPinResponse{
		Likes: likes,
	}
	data, err := json.Marshal(response)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return mw.ErrCreateResponse
	}
	return nil
}

func (del *delivery) listByAuthor(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("id")
	userId, err := strconv.Atoi(strId)
	if err != nil {
		return mw.ErrInvalidUserIdParam
	}

	likes, err := del.serv.ListByAuthor(userId)
	if err != nil {
		return err
	}

	response := listByAuthorResponse{
		Likes: likes,
	}
	data, err := json.Marshal(response)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return mw.ErrCreateResponse
	}
	return nil
}
