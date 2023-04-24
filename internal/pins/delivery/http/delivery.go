package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	mw "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	pkgPins "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pins"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
)

func RegisterHandlers(mux *httprouter.Router, logger log.Logger, authorizer mw.Authorizer, access mw.AccessChecker, serv pkgPins.Service) {
	del := delivery{serv, logger}

	mux.POST("/pins", mw.HandleLogger(mw.ErrorHandler(authorizer(mw.Cors(del.create)), logger), logger))
	mux.GET("/pins", mw.HandleLogger(mw.ErrorHandler(authorizer(mw.Cors(del.list)), logger), logger))
	mux.GET("/pins/:id", mw.HandleLogger(mw.ErrorHandler(authorizer(mw.Cors(del.get)), logger), logger))
	mux.GET("/users/:id/pins", mw.HandleLogger(mw.ErrorHandler(authorizer(mw.Cors(del.listByAuthor)), logger), logger))
	mux.PUT("/pins/:id", mw.HandleLogger(mw.ErrorHandler(authorizer(mw.Cors(access.WriteChecker(del.fullUpdate))), logger), logger))
	mux.DELETE("/pins/:id", mw.HandleLogger(mw.ErrorHandler(authorizer(mw.Cors(access.WriteChecker(del.delete))), logger), logger))
}

type delivery struct {
	serv pkgPins.Service
	log  log.Logger
}

func (del delivery) create(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strUserId := p.ByName("user-id")
	userId, err := strconv.Atoi(strUserId)
	if err != nil {
		return pkgErrors.ErrInvalidUserIdParam
	}

	file, handler, err := r.FormFile("bytes")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			return pkgErrors.ErrMissingFile
		} else {
			return pkgErrors.ErrParseForm
		}
	}
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, file)
	if err != nil {
		return pkgErrors.ErrFileCopy
	}

	image := models.Image{
		ID:    uuid.NewString() + filepath.Ext(handler.Filename),
		Bytes: buf.Bytes(),
	}

	params := pkgPins.CreateParams{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		MediaSource: image,
		Author:      userId,
	}
	pin, err := del.serv.Create(&params)
	if err != nil {
		return err
	}

	response := newCreateResponse(&pin)
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

func (del delivery) get(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		return pkgErrors.ErrInvalidPinIdParam
	}

	strUserId := p.ByName("user-id")
	userId, err := strconv.Atoi(strUserId)
	if err != nil {
		return pkgErrors.ErrInvalidUserIdParam
	}

	pin, err := del.serv.Get(id, userId)
	if err != nil {
		return err
	}

	response := newGetResponse(&pin)
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

func (del delivery) listByAuthor(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strAuthorId := p.ByName("id")
	authorId, err := strconv.Atoi(strAuthorId)
	if err != nil {
		return pkgErrors.ErrInvalidUserIdParam
	}

	strUserId := p.ByName("user-id")
	userId, err := strconv.Atoi(strUserId)
	if err != nil {
		return pkgErrors.ErrInvalidUserIdParam
	}

	queryValues := r.URL.Query()
	page := 1
	strPage := queryValues.Get("page")
	if strPage != "" {
		page, err = strconv.Atoi(strPage)
		if err != nil || page < 1 {
			return pkgErrors.ErrInvalidPageParam
		}
	}

	limit := 30
	strLimit := queryValues.Get("limit")
	if strLimit != "" {
		limit, err = strconv.Atoi(strLimit)
		if err != nil || limit < 0 {
			return pkgErrors.ErrInvalidLimitParam
		}
	}

	pins, err := del.serv.ListByAuthor(authorId, userId, page, limit)
	if err != nil {
		return err
	}

	response := newListResponse(pins)
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

func (del delivery) list(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strUserId := p.ByName("user-id")
	userId, err := strconv.Atoi(strUserId)
	if err != nil {
		return pkgErrors.ErrInvalidUserIdParam
	}

	queryValues := r.URL.Query()
	page := 1
	strPage := queryValues.Get("page")
	if strPage != "" {
		page, err = strconv.Atoi(strPage)
		if err != nil || page < 1 {
			return pkgErrors.ErrInvalidPageParam
		}
	}

	limit := 30
	strLimit := queryValues.Get("limit")
	if strLimit != "" {
		limit, err = strconv.Atoi(strLimit)
		if err != nil || limit < 0 {
			return pkgErrors.ErrInvalidLimitParam
		}
	}

	pins, err := del.serv.List(userId, page, limit)
	if err != nil {
		return err
	}

	response := newListResponse(pins)
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

func (del delivery) fullUpdate(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		return pkgErrors.ErrInvalidPinIdParam
	}

	params := pkgPins.FullUpdateParams{
		Id:          id,
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
	}
	pin, err := del.serv.FullUpdate(&params)
	if err != nil {
		return err
	}

	response := newFullUpdateResponse(&pin)
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

func (del delivery) delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		return pkgErrors.ErrInvalidPinIdParam
	}

	err = del.serv.Delete(id)
	if err != nil {
		return err
	}
	return pkgErrors.ErrNoContent
}
