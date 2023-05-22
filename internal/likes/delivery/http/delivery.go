package http

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	pkgLikes "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/likes"
	mw "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
)

func RegisterHandlers(mux *httprouter.Router, logger *zap.Logger, authorizer mw.Authorizer, csrf mw.CSRFMiddleware, serv pkgLikes.Service, m *mw.HttpMetricsMiddleware) {
	del := delivery{serv, logger}

	mux.POST("/pins/:id/like", mw.HandleLogger(mw.ErrorHandler(m.MetricsMiddleware(mw.Cors(authorizer(csrf(del.like))), logger), logger), logger))
	mux.DELETE("/pins/:id/like", mw.HandleLogger(mw.ErrorHandler(m.MetricsMiddleware(mw.Cors(authorizer(csrf(del.unlike))), logger), logger), logger))

	mux.GET("/pins/:id/likes", mw.HandleLogger(mw.ErrorHandler(m.MetricsMiddleware(mw.Cors(authorizer(csrf(del.listByPin))), logger), logger), logger))
	mux.GET("/users/:id/likes", mw.HandleLogger(mw.ErrorHandler(m.MetricsMiddleware(mw.Cors(authorizer(csrf(del.listByAuthor))), logger), logger), logger))
}

type delivery struct {
	serv pkgLikes.Service
	log  *zap.Logger
}

func (del *delivery) like(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("user-id")
	userId, err := strconv.Atoi(strId)
	if err != nil {
		return pkgErrors.ErrInvalidUserIdParam
	}

	strId = p.ByName("id")
	pinId, err := strconv.Atoi(strId)
	if err != nil {
		return pkgErrors.ErrInvalidPinIdParam
	}

	err = del.serv.Like(pinId, userId)
	if err != nil {
		return err
	}
	return pkgErrors.ErrNoContent
}

func (del *delivery) unlike(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("user-id")
	userId, err := strconv.Atoi(strId)
	if err != nil {
		return pkgErrors.ErrInvalidUserIdParam
	}

	strId = p.ByName("id")
	pinId, err := strconv.Atoi(strId)
	if err != nil {
		return pkgErrors.ErrInvalidPinIdParam
	}

	err = del.serv.Unlike(pinId, userId)
	if err != nil {
		return err
	}
	return pkgErrors.ErrNoContent
}

func (del *delivery) listByPin(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("id")
	pinId, err := strconv.Atoi(strId)
	if err != nil {
		return pkgErrors.ErrInvalidPinIdParam
	}

	likes, err := del.serv.ListByPin(pinId)
	if err != nil {
		return err
	}

	response := listByPinResponse{Likes: likes}
	data, err := response.MarshalJSON()
	if err != nil {
		return errors.Wrap(pkgErrors.ErrCreateResponse, err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrCreateResponse, err.Error())
	}
	return nil
}

func (del *delivery) listByAuthor(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("id")
	userId, err := strconv.Atoi(strId)
	if err != nil {
		return pkgErrors.ErrInvalidUserIdParam
	}

	likes, err := del.serv.ListByAuthor(userId)
	if err != nil {
		return err
	}

	response := listByAuthorResponse{Likes: likes}
	data, err := response.MarshalJSON()
	if err != nil {
		return errors.Wrap(pkgErrors.ErrCreateResponse, err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrCreateResponse, err.Error())
	}
	return nil
}
