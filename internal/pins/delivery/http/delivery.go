package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	_pins "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pins"
)

func RegisterHandlers(mux *httprouter.Router, logger log.Logger, authorizer middleware.Authorizer, access middleware.AccessChecker, serv _pins.Service) {
	del := delivery{serv, logger}

	mux.POST("/pins", middleware.HandleLogger(middleware.ErrorHandler(authorizer(middleware.CorsChecker(del.createPin)), logger), logger))
	mux.GET("/pins", middleware.HandleLogger(middleware.ErrorHandler(authorizer(middleware.CorsChecker(del.getPins)), logger), logger))
	mux.GET("/pins/:id", middleware.HandleLogger(middleware.ErrorHandler(authorizer(middleware.CorsChecker(del.getPin)), logger), logger))
	mux.GET("/users/:id/pins", middleware.HandleLogger(middleware.ErrorHandler(authorizer(middleware.CorsChecker(del.getPinsByUser)), logger), logger))
	mux.GET("/boards/:id/pins", middleware.HandleLogger(middleware.ErrorHandler(authorizer(middleware.CorsChecker(del.getPinsByBoard)), logger), logger))
	mux.PUT("/pins/:id", middleware.HandleLogger(middleware.ErrorHandler(authorizer(middleware.CorsChecker(access.WriteChecker(del.updatePin))), logger), logger))
	mux.DELETE("/pins/:id", middleware.HandleLogger(middleware.ErrorHandler(authorizer(middleware.CorsChecker(access.WriteChecker(del.deletePin))), logger), logger))

	mux.POST("/boards/:id/pins", middleware.HandleLogger(middleware.ErrorHandler(authorizer(middleware.CorsChecker(access.WriteChecker(del.addPinToBoard))), logger), logger))
	mux.DELETE("/boards/:id/pins", middleware.HandleLogger(middleware.ErrorHandler(authorizer(middleware.CorsChecker(access.WriteChecker(del.removePinFromBoard))), logger), logger))
}

type delivery struct {
	serv _pins.Service
	log  log.Logger
}

func (del delivery) createPin(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strUserId := p.ByName("user-id")
	userId, err := strconv.Atoi(strUserId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
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

	createdPin, err := del.serv.Create(&params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	data, err := json.Marshal(createdPin)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	return err
}

func (del delivery) getPin(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("id")

	id, err := strconv.Atoi(strId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	pin, err := del.serv.Get(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	data, err := json.Marshal(pin)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	return err
}

func (del delivery) getPinsByUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strUserId := p.ByName("id")
	userId, err := strconv.Atoi(strUserId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
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

	data, err := json.Marshal(pins)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	return err
}

func (del delivery) getPinsByBoard(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strBoardId := p.ByName("id")
	boardId, err := strconv.Atoi(strBoardId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
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

	data, err := json.Marshal(pins)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	return err
}

func (del delivery) getPins(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
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

	data, err := json.Marshal(pins)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	return err
}

func (del delivery) updatePin(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("id")

	id, err := strconv.Atoi(strId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	pin := models.Pin{}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&pin); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	pin.Id = id

	pin, err = del.serv.Update(&pin)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	data, err := json.Marshal(pin)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	return err
}

func (del delivery) deletePin(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("id")

	id, err := strconv.Atoi(strId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}
	err = del.serv.Delete(id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	return nil
}

func (del delivery) addPinToBoard(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strBoardId := p.ByName("id")

	boardId, err := strconv.Atoi(strBoardId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	var id struct {
		Id int `json:"id"`
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	err = del.serv.AddToBoard(boardId, id.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	return nil
}

func (del delivery) removePinFromBoard(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strBoardId := p.ByName("id")

	boardId, err := strconv.Atoi(strBoardId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	var id struct {
		Id int `json:"id"`
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	err = del.serv.RemoveFromBoard(boardId, id.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	return nil
}
