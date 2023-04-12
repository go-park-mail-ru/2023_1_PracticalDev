package http

import (
	"bytes"
	"encoding/json"
	mw "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	_pins "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pins"
)

func RegisterHandlers(mux *httprouter.Router, logger log.Logger, authorizer mw.Authorizer, access mw.AccessChecker, serv _pins.Service) {
	del := delivery{serv, logger}

	mux.POST("/pins", mw.HandleLogger(mw.ErrorHandler(authorizer(mw.CorsChecker(del.create)), logger), logger))
	mux.GET("/pins", mw.HandleLogger(mw.ErrorHandler(authorizer(mw.CorsChecker(del.list)), logger), logger))
	mux.GET("/pins/:id", mw.HandleLogger(mw.ErrorHandler(authorizer(mw.CorsChecker(del.get)), logger), logger))
	mux.GET("/users/:id/pins", mw.HandleLogger(mw.ErrorHandler(authorizer(mw.CorsChecker(del.listByUser)), logger), logger))
	mux.PUT("/pins/:id", mw.HandleLogger(mw.ErrorHandler(authorizer(mw.CorsChecker(access.WriteChecker(del.fullUpdate))), logger), logger))
	mux.DELETE("/pins/:id", mw.HandleLogger(mw.ErrorHandler(authorizer(mw.CorsChecker(access.WriteChecker(del.delete))), logger), logger))
}

type delivery struct {
	serv _pins.Service
	log  log.Logger
}

func (del delivery) create(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strUserId := p.ByName("user-id")
	userId, err := strconv.Atoi(strUserId)
	if err != nil {
		return mw.ErrInvalidUserIdParam
	}

	file, handler, err := r.FormFile("bytes")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			return mw.ErrMissingFile
		} else {
			return mw.ErrParseForm
		}
	}
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, file)
	if err != nil {
		return mw.ErrFileCopy
	}

	image := models.Image{
		ID:    uuid.NewString() + filepath.Ext(handler.Filename),
		Bytes: buf.Bytes(),
	}

	params := _pins.CreateParams{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		MediaSource: image,
		Author:      userId,
	}
	pin, err := del.serv.Create(&params)
	if err != nil {
		return mw.ErrService
	}

	response := createResponse{
		Id:          pin.Id,
		Title:       pin.Title,
		Description: pin.Description,
		MediaSource: pin.MediaSource,
		Author:      pin.Author,
	}
	data, err := json.Marshal(response)
	if err != nil {
		return mw.ErrCreateResponse
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	return err
}

func (del delivery) get(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		return mw.ErrInvalidPinIdParam
	}

	pin, err := del.serv.Get(id)
	if err != nil {
		if errors.Is(err, _pins.ErrPinNotFound) {
			return mw.ErrPinNotFound
		} else {
			return mw.ErrService
		}
	}

	response := getResponse{
		Id:          pin.Id,
		Title:       pin.Title,
		Description: pin.Description,
		MediaSource: pin.MediaSource,
		Author:      pin.Author,
	}
	data, err := json.Marshal(response)
	if err != nil {
		return mw.ErrCreateResponse
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	return err
}

func (del delivery) listByUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strUserId := p.ByName("id")
	userId, err := strconv.Atoi(strUserId)
	if err != nil {
		return mw.ErrInvalidUserIdParam
	}

	queryValues := r.URL.Query()
	page := 1
	strPage := queryValues.Get("page")
	if strPage != "" {
		page, err = strconv.Atoi(strPage)
		if err != nil {
			return mw.ErrInvalidPageParam
		} else if page < 1 {
			return mw.ErrInvalidPageParam
		}
	}

	limit := 30
	strLimit := queryValues.Get("limit")
	if strLimit != "" {
		limit, err = strconv.Atoi(strLimit)
		if err != nil {
			return mw.ErrInvalidLimitParam
		} else if limit < 0 {
			return mw.ErrInvalidLimitParam
		}
	}

	pins, err := del.serv.ListByUser(userId, page, limit)
	if err != nil {
		return mw.ErrService
	}

	response := listResponse{Pins: pins}
	data, err := json.Marshal(response)
	if err != nil {
		return mw.ErrCreateResponse
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	return err
}

func (del delivery) list(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	var err error
	queryValues := r.URL.Query()
	page := 1
	strPage := queryValues.Get("page")
	if strPage != "" {
		page, err = strconv.Atoi(strPage)
		if err != nil {
			return mw.ErrInvalidPageParam
		} else if page < 1 {
			return mw.ErrInvalidPageParam
		}
	}

	limit := 30
	strLimit := queryValues.Get("limit")
	if strLimit != "" {
		limit, err = strconv.Atoi(strLimit)
		if err != nil {
			return mw.ErrInvalidLimitParam
		} else if limit < 0 {
			return mw.ErrInvalidLimitParam
		}
	}

	pins, err := del.serv.List(page, limit)
	if err != nil {
		return mw.ErrService
	}

	response := listResponse{Pins: pins}
	data, err := json.Marshal(response)
	if err != nil {
		return mw.ErrCreateResponse
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	return err
}

func (del delivery) fullUpdate(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		return mw.ErrInvalidPinIdParam
	}

	file, handler, err := r.FormFile("bytes")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			return mw.ErrMissingFile
		} else {
			return mw.ErrParseForm
		}
	}
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, file)
	if err != nil {
		return mw.ErrFileCopy
	}

	image := models.Image{
		ID:    uuid.NewString() + filepath.Ext(handler.Filename),
		Bytes: buf.Bytes(),
	}

	params := _pins.FullUpdateParams{
		Id:          id,
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		MediaSource: image,
	}
	pin, err := del.serv.FullUpdate(&params)
	if err != nil {
		return mw.ErrService
	}

	response := fullUpdateResponse{
		Id:          pin.Id,
		Title:       pin.Title,
		Description: pin.Description,
		MediaSource: pin.MediaSource,
		Author:      pin.Author,
	}
	data, err := json.Marshal(response)
	if err != nil {
		return mw.ErrCreateResponse
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	return err
}

func (del delivery) delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		return mw.ErrInvalidPinIdParam
	}

	err = del.serv.Delete(id)
	if err != nil {
		if errors.Is(err, _pins.ErrPinNotFound) {
			return mw.ErrPinNotFound
		} else {
			return mw.ErrService
		}
	}
	return err
}
