package http

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	_boards "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/boards"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	"github.com/julienschmidt/httprouter"
)

var (
	ErrInvalidBoardIdParam = errors.New("invalid board id param")
	ErrInvalidUserId       = errors.New("invalid user id")
)

func RegisterHandlers(mux *httprouter.Router, logger log.Logger, authorizer middleware.Authorizer, access middleware.AccessChecker, serv _boards.Service) {
	del := delivery{serv, logger}

	mux.POST("/boards", middleware.HandleLogger(middleware.ErrorHandler(authorizer(middleware.CorsChecker(del.create)), logger), logger))
	mux.GET("/boards", middleware.HandleLogger(middleware.ErrorHandler(authorizer(middleware.CorsChecker(del.list)), logger), logger))
	mux.GET("/boards/:id", middleware.HandleLogger(middleware.ErrorHandler(middleware.CorsChecker(authorizer(access.ReadChecker(del.get))), logger), logger))
	mux.PUT("/boards/:id", middleware.HandleLogger(middleware.ErrorHandler(middleware.CorsChecker(authorizer(access.WriteChecker(del.fullUpdate))), logger), logger))
	mux.PATCH("/boards/:id", middleware.HandleLogger(middleware.ErrorHandler(middleware.CorsChecker(authorizer(access.WriteChecker(del.partialUpdate))), logger), logger))
	mux.DELETE("/boards/:id", middleware.HandleLogger(middleware.ErrorHandler(middleware.CorsChecker(authorizer(access.WriteChecker(del.delete))), logger), logger))
}

type delivery struct {
	serv _boards.Service
	log  log.Logger
}

func (del *delivery) create(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strUserId := p.ByName("user-id")
	userId, err := strconv.Atoi(strUserId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return ErrInvalidUserId
	}

	var request createRequest
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err = decoder.Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	params := _boards.CreateParams{
		Name:        request.Name,
		Description: request.Description,
		Privacy:     "secret",
		UserId:      userId,
	}
	if request.Privacy != nil {
		params.Privacy = *request.Privacy
	}

	createdBoard, err := del.serv.Create(&params)
	if err != nil {
		switch err {
		case _boards.ErrDb:
			w.WriteHeader(http.StatusInternalServerError)
		case _boards.ErrInvalidPrivacy:
			w.WriteHeader(http.StatusBadRequest)
		}
		return err
	}

	data, err := json.Marshal(createdBoard)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	return err
}

func (del *delivery) list(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strUserId := p.ByName("user-id")
	userId, err := strconv.Atoi(strUserId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return ErrInvalidUserId
	}

	boards, err := del.serv.List(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	response := listResponse{Boards: boards}
	data, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	return err
}

func (del *delivery) get(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return ErrInvalidBoardIdParam
	}

	board, err := del.serv.Get(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return err
	}

	response := getResponse{
		Id:          board.Id,
		Name:        board.Name,
		Description: board.Description,
		Privacy:     board.Privacy,
		UserId:      board.UserId,
	}
	data, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	_, err = w.Write(data)
	return err
}

func (del *delivery) fullUpdate(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return ErrInvalidBoardIdParam
	}

	var request fullUpdateRequest
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err = decoder.Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	params := _boards.FullUpdateParams{
		Id:          id,
		Name:        request.Name,
		Description: request.Description,
		Privacy:     "secret",
	}
	if request.Privacy != nil {
		params.Privacy = *request.Privacy
	}

	board, err := del.serv.FullUpdate(&params)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			return err
		}
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	response := fullUpdateResponse{
		Id:          board.Id,
		Name:        board.Name,
		Description: board.Description,
		Privacy:     board.Privacy,
		UserId:      board.UserId,
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

func (del *delivery) partialUpdate(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return ErrInvalidBoardIdParam
	}

	var request partialUpdateRequest
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err = decoder.Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	params := _boards.PartialUpdateParams{
		Id:      id,
		Privacy: "secret",
	}
	if request.Name != nil {
		params.UpdateName = true
		params.Name = *request.Name
	}
	if request.Description != nil {
		params.UpdateDescription = true
		params.Description = *request.Description
	}
	if request.Privacy != nil {
		params.UpdatePrivacy = true
		params.Privacy = *request.Privacy
	}

	board, err := del.serv.PartialUpdate(&params)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return err
	}

	response := partialUpdateResponse{
		Id:          board.Id,
		Name:        board.Name,
		Description: board.Description,
		Privacy:     board.Privacy,
		UserId:      board.UserId,
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

func (del *delivery) delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return ErrInvalidBoardIdParam
	}

	err = del.serv.Delete(id)
	switch err {
	case _boards.ErrDb:
		w.WriteHeader(http.StatusInternalServerError)
	case _boards.ErrBoardNotFound:
		w.WriteHeader(http.StatusNotFound)
	}
	return err
}
