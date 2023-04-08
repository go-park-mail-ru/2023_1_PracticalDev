package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	_pins "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pins"
)

var (
	ErrFileCopy       = errors.New("file copy error")
	ErrMissingFile    = errors.New("missing file")
	ErrParseForm      = errors.New("parse form error")
	ErrInvalidPinId   = errors.New("invalid pin id")
	ErrInvalidBoardId = errors.New("invalid board id")
	ErrInvalidUserId  = errors.New("invalid user id")
	ErrPinNotFound    = errors.New("pin not found")
	ErrService        = errors.New("service error")
)

func RegisterHandlers(mux *httprouter.Router, logger log.Logger, authorizer middleware.Authorizer, access middleware.AccessChecker, serv _pins.Service) {
	del := delivery{serv, logger}

	mux.POST("/pins", middleware.HandleLogger(middleware.ErrorHandler(authorizer(middleware.CorsChecker(del.create)), logger), logger))
	mux.GET("/pins", middleware.HandleLogger(middleware.ErrorHandler(authorizer(middleware.CorsChecker(del.list)), logger), logger))
	mux.GET("/pins/:id", middleware.HandleLogger(middleware.ErrorHandler(authorizer(middleware.CorsChecker(del.get)), logger), logger))
	mux.GET("/users/:id/pins", middleware.HandleLogger(middleware.ErrorHandler(authorizer(middleware.CorsChecker(del.listByUser)), logger), logger))
	mux.GET("/boards/:id/pins", middleware.HandleLogger(middleware.ErrorHandler(authorizer(middleware.CorsChecker(del.listByBoard)), logger), logger))
	mux.PUT("/pins/:id", middleware.HandleLogger(middleware.ErrorHandler(authorizer(middleware.CorsChecker(access.WriteChecker(del.fullUpdate))), logger), logger))
	mux.DELETE("/pins/:id", middleware.HandleLogger(middleware.ErrorHandler(authorizer(middleware.CorsChecker(access.WriteChecker(del.delete))), logger), logger))

	mux.POST("/boards/:id/pins", middleware.HandleLogger(middleware.ErrorHandler(authorizer(middleware.CorsChecker(access.WriteChecker(del.addToBoard))), logger), logger))
	mux.DELETE("/boards/:board_id/pins/:id", middleware.HandleLogger(middleware.ErrorHandler(authorizer(middleware.CorsChecker(access.WriteChecker(del.removeFromBoard))), logger), logger))
}

type delivery struct {
	serv _pins.Service
	log  log.Logger
}

func (del delivery) create(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strUserId := p.ByName("user-id")
	userId, err := strconv.Atoi(strUserId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return ErrInvalidUserId
	}

	file, _, err := r.FormFile("bytes")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}
	defer file.Close()
	buf := bytes.NewBuffer(nil)
	_, _ = io.Copy(buf, file)

	image := models.Image{
		ID:    uuid.NewString() + ".jpg",
		Bytes: buf.Bytes(),
	}
	del.log.Debug(r.FormValue("title"))

	params := _pins.CreateParams{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		MediaSource: image,
		Author:      userId,
	}
	pin, err := del.serv.Create(&params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
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
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	return err
}

func (del delivery) get(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return ErrInvalidPinId
	}

	pin, err := del.serv.Get(id)
	if err != nil {
		if err == _pins.ErrPinNotFound {
			err = ErrPinNotFound
			w.WriteHeader(http.StatusNotFound)
		} else {
			err = ErrService
			w.WriteHeader(http.StatusInternalServerError)
		}
		return err
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
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	return err
}

func (del delivery) listByUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strUserId := p.ByName("id")
	userId, err := strconv.Atoi(strUserId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return ErrInvalidPinId
	}

	queryValues := r.URL.Query()
	page := 1
	strPage := queryValues.Get("page")
	if strPage != "" {
		page, err = strconv.Atoi(strPage)
		if err != nil || page < 1 {
			w.WriteHeader(http.StatusBadRequest)
			return err
		}
	}

	limit := 30
	strLimit := queryValues.Get("limit")
	if strLimit != "" {
		limit, err = strconv.Atoi(strLimit)
		if err != nil || limit < 0 {
			w.WriteHeader(http.StatusBadRequest)
			return err
		}
	}

	pins, err := del.serv.ListByUser(userId, page, limit)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return err
	}

	response := listResponse{Pins: pins}
	data, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	return err
}

func (del delivery) listByBoard(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strBoardId := p.ByName("id")
	boardId, err := strconv.Atoi(strBoardId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return ErrInvalidPinId
	}

	queryValues := r.URL.Query()
	page := 1
	strPage := queryValues.Get("page")
	if strPage != "" {
		page, err = strconv.Atoi(strPage)
		if err != nil || page < 1 {
			w.WriteHeader(http.StatusBadRequest)
			return err
		}
	}

	limit := 30
	strLimit := queryValues.Get("limit")
	if strLimit != "" {
		limit, err = strconv.Atoi(strLimit)
		if err != nil || limit < 0 {
			w.WriteHeader(http.StatusBadRequest)
			return err
		}
	}

	pins, err := del.serv.ListByBoard(boardId, page, limit)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return err
	}

	response := listResponse{Pins: pins}
	data, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
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
		if err != nil || page < 1 {
			w.WriteHeader(http.StatusBadRequest)
			return err
		}
	}

	limit := 30
	strLimit := queryValues.Get("limit")
	if strLimit != "" {
		limit, err = strconv.Atoi(strLimit)
		if err != nil || limit < 0 {
			w.WriteHeader(http.StatusBadRequest)
			return err
		}
	}

	pins, err := del.serv.List(page, limit)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return err
	}

	response := listResponse{Pins: pins}
	data, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	return err
}

func (del delivery) fullUpdate(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return ErrInvalidPinId
	}

	file, handler, err := r.FormFile("bytes")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err == http.ErrMissingFile {
			err = ErrMissingFile
		} else {
			err = ErrParseForm
		}
		return err
	}
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return ErrFileCopy
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
		w.WriteHeader(http.StatusBadRequest)
		return err
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
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	return err
}

func (del delivery) delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return ErrInvalidPinId
	}

	err = del.serv.Delete(id)
	switch err {
	case _pins.ErrDb:
		w.WriteHeader(http.StatusInternalServerError)
		err = ErrService
	case _pins.ErrPinNotFound:
		w.WriteHeader(http.StatusNotFound)
		err = ErrPinNotFound
	}
	return err
}

func (del delivery) addToBoard(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("board_id")
	boardId, err := strconv.Atoi(strId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return ErrInvalidBoardId
	}

	strId = p.ByName("id")
	pinId, err := strconv.Atoi(strId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return ErrInvalidPinId
	}

	err = del.serv.AddToBoard(boardId, pinId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return ErrService
	}
	return nil
}

func (del delivery) removeFromBoard(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("board_id")
	boardId, err := strconv.Atoi(strId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return ErrInvalidBoardId
	}

	strId = p.ByName("id")
	pinId, err := strconv.Atoi(strId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return ErrInvalidPinId
	}

	err = del.serv.RemoveFromBoard(boardId, pinId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	return nil
}
