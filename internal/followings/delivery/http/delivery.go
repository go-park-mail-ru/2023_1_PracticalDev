package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/followings"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	mw "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
)

func RegisterHandlers(mux *httprouter.Router, logger log.Logger, authorizer mw.Authorizer, serv followings.Service) {
	del := delivery{serv, logger}

	mux.POST("/users/:id/following", mw.HandleLogger(mw.ErrorHandler(mw.Cors(authorizer(del.follow)), logger), logger))
	mux.DELETE("/users/:id/following", mw.HandleLogger(mw.ErrorHandler(mw.Cors(authorizer(del.unfollow)), logger), logger))

	mux.GET("/users/:id/followers", mw.HandleLogger(mw.ErrorHandler(mw.Cors(authorizer(del.getFollowers)), logger), logger))
	mux.GET("/users/:id/followees", mw.HandleLogger(mw.ErrorHandler(mw.Cors(authorizer(del.getFollowees)), logger), logger))
}

type delivery struct {
	serv followings.Service
	log  log.Logger
}

func (del *delivery) follow(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strFollowerId := p.ByName("user-id")
	followerId, err := strconv.Atoi(strFollowerId)
	if err != nil {
		return pkgErrors.ErrInvalidUserIdParam
	}

	strFolloweeId := p.ByName("id")
	followeeId, err := strconv.Atoi(strFolloweeId)
	if err != nil {
		return pkgErrors.ErrInvalidUserIdParam
	}

	err = del.serv.Follow(followerId, followeeId)
	if err != nil {
		return err
	}
	return pkgErrors.ErrNoContent
}

func (del *delivery) unfollow(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strFollowerId := p.ByName("user-id")
	followerId, err := strconv.Atoi(strFollowerId)
	if err != nil {
		return pkgErrors.ErrInvalidUserIdParam
	}

	strFolloweeId := p.ByName("id")
	followeeId, err := strconv.Atoi(strFolloweeId)
	if err != nil {
		return pkgErrors.ErrInvalidUserIdParam
	}

	err = del.serv.Unfollow(followerId, followeeId)
	if err != nil {
		return err
	}
	return pkgErrors.ErrNoContent
}

func (del *delivery) getFollowers(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("id")
	userId, err := strconv.Atoi(strId)
	if err != nil {
		return pkgErrors.ErrInvalidUserIdParam
	}

	followers, err := del.serv.GetFollowers(userId)
	if err != nil {
		return err
	}

	response := followersResponse{
		Followers: followers,
	}
	data, err := json.Marshal(response)
	if err != nil {
		return pkgErrors.ErrCreateResponse
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return pkgErrors.ErrCreateResponse
	}
	return nil
}

func (del *delivery) getFollowees(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("id")
	userId, err := strconv.Atoi(strId)
	if err != nil {
		return pkgErrors.ErrInvalidUserIdParam
	}

	followees, err := del.serv.GetFollowees(userId)
	if err != nil {
		return err
	}

	response := followeesResponse{
		Followees: followees,
	}
	data, err := json.Marshal(response)
	if err != nil {
		return pkgErrors.ErrCreateResponse
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return pkgErrors.ErrCreateResponse
	}
	return nil
}