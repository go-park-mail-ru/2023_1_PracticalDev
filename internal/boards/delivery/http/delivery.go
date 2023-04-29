package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	pkgBoards "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/boards"
	mw "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/log"
	"github.com/julienschmidt/httprouter"
)

func RegisterHandlers(mux *httprouter.Router, logger log.Logger, authorizer mw.Authorizer, access mw.AccessChecker, serv pkgBoards.Service) {
	del := delivery{serv, logger}

	mux.POST("/boards", mw.HandleLogger(mw.ErrorHandler(authorizer(mw.Cors(del.create)), logger), logger))
	mux.GET("/boards", mw.HandleLogger(mw.ErrorHandler(authorizer(mw.Cors(del.list)), logger), logger))
	mux.GET("/boards/:id", mw.HandleLogger(mw.ErrorHandler(mw.Cors(authorizer(access.ReadChecker(del.get))), logger), logger))
	mux.PUT("/boards/:id", mw.HandleLogger(mw.ErrorHandler(mw.Cors(authorizer(access.WriteChecker(del.fullUpdate))), logger), logger))
	mux.PATCH("/boards/:id", mw.HandleLogger(mw.ErrorHandler(mw.Cors(authorizer(access.WriteChecker(del.partialUpdate))), logger), logger))
	mux.DELETE("/boards/:id", mw.HandleLogger(mw.ErrorHandler(mw.Cors(authorizer(access.WriteChecker(del.delete))), logger), logger))

	mux.POST("/boards/:id/pins/:pin_id", mw.HandleLogger(mw.ErrorHandler(authorizer(mw.Cors(access.WriteChecker(del.addPin))), logger), logger))
	mux.GET("/boards/:id/pins", mw.HandleLogger(mw.ErrorHandler(authorizer(mw.Cors(access.ReadChecker(del.pinsList))), logger), logger))
	mux.DELETE("/boards/:id/pins/:pin_id", mw.HandleLogger(mw.ErrorHandler(authorizer(mw.Cors(access.WriteChecker(del.removePin))), logger), logger))
}

type delivery struct {
	serv pkgBoards.Service
	log  log.Logger
}

func (del *delivery) create(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strUserId := p.ByName("user-id")
	userId, err := strconv.Atoi(strUserId)
	if err != nil {
		return pkgErrors.ErrInvalidUserIdParam
	}

	var request createRequest
	decoder := json.NewDecoder(r.Body)
	defer func() {
		err = r.Body.Close()
		if err != nil {
			del.log.Error(err)
		}
	}()
	if err = decoder.Decode(&request); err != nil {
		return pkgErrors.ErrParseJson
	}

	params := pkgBoards.CreateParams{
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
		return err
	}

	data, err := json.Marshal(createdBoard)
	if err != nil {
		return pkgErrors.ErrCreateResponse
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	return err
}

func (del *delivery) list(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strUserId := p.ByName("user-id")
	userId, err := strconv.Atoi(strUserId)
	if err != nil {
		return pkgErrors.ErrInvalidUserIdParam
	}

	boards, err := del.serv.List(userId)
	if err != nil {
		return pkgErrors.ErrService
	}

	response := listResponse{Boards: boards}
	data, err := json.Marshal(response)
	if err != nil {
		return pkgErrors.ErrCreateResponse
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	return err
}

func (del *delivery) get(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		return pkgErrors.ErrInvalidBoardIdParam
	}

	board, err := del.serv.Get(id)
	if err != nil {
		return pkgErrors.ErrBoardNotFound
	}

	response := newGetResponse(&board)
	data, err := json.Marshal(response)
	if err != nil {
		return pkgErrors.ErrCreateResponse
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	return err
}

func (del *delivery) fullUpdate(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		return pkgErrors.ErrInvalidBoardIdParam
	}

	var request fullUpdateRequest
	decoder := json.NewDecoder(r.Body)
	defer func() {
		err = r.Body.Close()
		if err != nil {
			del.log.Error(err)
		}
	}()
	if err = decoder.Decode(&request); err != nil {
		return pkgErrors.ErrParseJson
	}

	params := pkgBoards.FullUpdateParams{
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
		return err
	}

	response := newFullUpdateResponse(&board)
	data, err := json.Marshal(response)
	if err != nil {
		return pkgErrors.ErrCreateResponse
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	return err
}

func (del *delivery) partialUpdate(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		return pkgErrors.ErrInvalidBoardIdParam
	}

	var request partialUpdateRequest
	decoder := json.NewDecoder(r.Body)
	defer func() {
		err = r.Body.Close()
		if err != nil {
			del.log.Error(err)
		}
	}()
	if err = decoder.Decode(&request); err != nil {
		return pkgErrors.ErrParseJson
	}

	params := pkgBoards.PartialUpdateParams{
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
		return err
	}

	response := newPartialUpdateResponse(&board)
	data, err := json.Marshal(response)
	if err != nil {
		return pkgErrors.ErrCreateResponse
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	return err
}

func (del *delivery) delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		return pkgErrors.ErrInvalidBoardIdParam
	}

	err = del.serv.Delete(id)
	if err != nil {
		return err
	}
	return pkgErrors.ErrNoContent
}

func (del *delivery) pinsList(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strBoardId := p.ByName("id")
	boardId, err := strconv.Atoi(strBoardId)
	if err != nil {
		return pkgErrors.ErrInvalidBoardIdParam
	}

	queryValues := r.URL.Query()
	page := 1
	strPage := queryValues.Get("page")
	if strPage != "" {
		page, err = strconv.Atoi(strPage)
		if err != nil {
			return pkgErrors.ErrInvalidPageParam
		} else if page < 1 {
			return pkgErrors.ErrInvalidPageParam
		}
	}

	limit := 30
	strLimit := queryValues.Get("limit")
	if strLimit != "" {
		limit, err = strconv.Atoi(strLimit)
		if err != nil {
			return pkgErrors.ErrInvalidLimitParam
		} else if limit < 0 {
			return pkgErrors.ErrInvalidLimitParam
		}
	}

	pins, err := del.serv.PinsList(boardId, page, limit)
	if err != nil {
		return pkgErrors.ErrService
	}

	response := pinListResponse{Pins: pins}
	data, err := json.Marshal(response)
	if err != nil {
		return pkgErrors.ErrCreateResponse
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	return err
}

func (del *delivery) addPin(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("id")
	boardId, err := strconv.Atoi(strId)
	if err != nil {
		return pkgErrors.ErrInvalidBoardIdParam
	}

	strId = p.ByName("pin_id")
	pinId, err := strconv.Atoi(strId)
	if err != nil {
		return pkgErrors.ErrInvalidPinIdParam
	}

	err = del.serv.AddPin(boardId, pinId)
	if err != nil {
		return err
	}
	return pkgErrors.ErrNoContent
}

func (del *delivery) removePin(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("id")
	boardId, err := strconv.Atoi(strId)
	if err != nil {
		return pkgErrors.ErrInvalidBoardIdParam
	}

	strId = p.ByName("pin_id")
	pinId, err := strconv.Atoi(strId)
	if err != nil {
		return pkgErrors.ErrInvalidPinIdParam
	}

	err = del.serv.RemovePin(boardId, pinId)
	if err != nil {
		return err
	}
	return pkgErrors.ErrNoContent
}
