package http

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
	"strconv"

	_boards "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/boards"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	mw "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	"github.com/julienschmidt/httprouter"
)

func RegisterHandlers(mux *httprouter.Router, logger log.Logger, authorizer mw.Authorizer, access mw.AccessChecker, serv _boards.Service) {
	del := delivery{serv, logger}

	mux.POST("/boards", mw.HandleLogger(mw.ErrorHandler(authorizer(mw.CorsChecker(del.create)), logger), logger))
	mux.GET("/boards", mw.HandleLogger(mw.ErrorHandler(authorizer(mw.CorsChecker(del.list)), logger), logger))
	mux.GET("/boards/:board_id", mw.HandleLogger(mw.ErrorHandler(mw.CorsChecker(authorizer(access.ReadChecker(del.get))), logger), logger))
	mux.PUT("/boards/:board_id", mw.HandleLogger(mw.ErrorHandler(mw.CorsChecker(authorizer(access.WriteChecker(del.fullUpdate))), logger), logger))
	mux.PATCH("/boards/:board_id", mw.HandleLogger(mw.ErrorHandler(mw.CorsChecker(authorizer(access.WriteChecker(del.partialUpdate))), logger), logger))
	mux.DELETE("/boards/:board_id", mw.HandleLogger(mw.ErrorHandler(mw.CorsChecker(authorizer(access.WriteChecker(del.delete))), logger), logger))
}

type delivery struct {
	serv _boards.Service
	log  log.Logger
}

func (del *delivery) create(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strUserId := p.ByName("user-id")
	userId, err := strconv.Atoi(strUserId)
	if err != nil {
		return mw.ErrInvalidUserIdParam
	}

	var request createRequest
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err = decoder.Decode(&request); err != nil {
		return mw.ErrParseJson
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
		if errors.Is(err, _boards.ErrInvalidPrivacy) {
			return mw.ErrBadParams
		} else {
			return mw.ErrService
		}
	}

	data, err := json.Marshal(createdBoard)
	if err != nil {
		return mw.ErrCreateResponse
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	return err
}

func (del *delivery) list(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strUserId := p.ByName("user-id")
	userId, err := strconv.Atoi(strUserId)
	if err != nil {
		return mw.ErrInvalidUserIdParam
	}

	boards, err := del.serv.List(userId)
	if err != nil {
		return mw.ErrService
	}

	response := listResponse{Boards: boards}
	data, err := json.Marshal(response)
	if err != nil {
		return mw.ErrCreateResponse
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	return err
}

func (del *delivery) get(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("board_id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		return mw.ErrInvalidBoardIdParam
	}

	board, err := del.serv.Get(id)
	if err != nil {
		return mw.ErrBoardNotFound
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
		return mw.ErrCreateResponse
	}

	_, err = w.Write(data)
	return err
}

func (del *delivery) fullUpdate(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("board_id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		return mw.ErrInvalidBoardIdParam
	}

	var request fullUpdateRequest
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err = decoder.Decode(&request); err != nil {
		return mw.ErrParseJson
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
		if errors.Is(err, _boards.ErrBoardNotFound) {
			return mw.ErrBoardNotFound
		} else {
			return mw.ErrService
		}
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
		return mw.ErrCreateResponse
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	return err
}

func (del *delivery) partialUpdate(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("board_id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		return mw.ErrInvalidBoardIdParam
	}

	var request partialUpdateRequest
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err = decoder.Decode(&request); err != nil {
		return mw.ErrParseJson
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
		if errors.Is(err, _boards.ErrBoardNotFound) {
			return mw.ErrBoardNotFound
		} else {
			return mw.ErrService
		}
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
		return mw.ErrCreateResponse
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	return err
}

func (del *delivery) delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("board_id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		return mw.ErrInvalidBoardIdParam
	}

	err = del.serv.Delete(id)
	if err != nil {
		if errors.Is(err, _boards.ErrBoardNotFound) {
			return mw.ErrBoardNotFound
		} else {
			return mw.ErrService
		}
	}
	return nil
}
