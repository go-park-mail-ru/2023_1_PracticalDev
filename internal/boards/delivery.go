package boards

import (
	"database/sql"
	"encoding/json"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type apiBoard struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Privacy     *string `json:"privacy"`
}

type apiPartialUpdateBoard struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Privacy     *string `json:"privacy"`
}

func RegisterHandlers(mux *httprouter.Router, logger log.Logger, authorizer middleware.Authorizer, access middleware.AccessChecker, serv Service) {
	del := delivery{serv, logger}

	mux.POST("/boards", middleware.HandleLogger(middleware.ErrorHandler(authorizer(middleware.CorsChecker(del.createBoard)), logger), logger))
	mux.GET("/boards", middleware.HandleLogger(middleware.ErrorHandler(authorizer(middleware.CorsChecker(del.getBoards)), logger), logger))
	mux.GET("/boards/:id", middleware.HandleLogger(middleware.ErrorHandler(middleware.CorsChecker(authorizer(access.ReadChecker(del.getBoard))), logger), logger))
	mux.PUT("/boards/:id", middleware.HandleLogger(middleware.ErrorHandler(middleware.CorsChecker(authorizer(access.WriteChecker(del.fullUpdateBoard))), logger), logger))
	mux.PATCH("/boards/:id", middleware.HandleLogger(middleware.ErrorHandler(middleware.CorsChecker(authorizer(access.WriteChecker(del.partialUpdateBoard))), logger), logger))
	mux.DELETE("/boards/:id", middleware.HandleLogger(middleware.ErrorHandler(middleware.CorsChecker(authorizer(access.WriteChecker(del.deleteBoard))), logger), logger))
}

type delivery struct {
	serv Service
	log  log.Logger
}

func (del *delivery) createBoard(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strUserId := p.ByName("user-id")
	userId, err := strconv.Atoi(strUserId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	var request apiBoard
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err = decoder.Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	params := createBoardParams{
		Name:        request.Name,
		Description: request.Description,
		Privacy:     "secret",
		UserId:      userId,
	}
	if request.Privacy != nil {
		params.Privacy = *request.Privacy
	}

	createdBoard, err := del.serv.CreateBoard(&params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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

func (del *delivery) getBoards(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strUserId := p.ByName("user-id")
	userId, err := strconv.Atoi(strUserId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	boards, err := del.serv.GetBoards(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	response := map[string]interface{}{"items": boards}
	data, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	return err
}

func (del *delivery) getBoard(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	board, err := del.serv.GetBoard(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return err
	}

	data, err := json.Marshal(board)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	_, err = w.Write(data)
	return err
}

func (del *delivery) fullUpdateBoard(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	var request apiBoard
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err = decoder.Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	params := FullUpdateBoardParams{
		Id:          id,
		Name:        request.Name,
		Description: request.Description,
		Privacy:     "secret",
	}
	if request.Privacy != nil {
		params.Privacy = *request.Privacy
	}

	updatedBoard, err := del.serv.FullUpdateBoard(&params)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			return err
		}
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	data, err := json.Marshal(updatedBoard)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	return err
}

func (del *delivery) partialUpdateBoard(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	var request apiPartialUpdateBoard
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err = decoder.Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	params := partialUpdateBoardParams{
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

	updatedBoard, err := del.serv.PartialUpdateBoard(&params)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return err
	}

	data, err := json.Marshal(updatedBoard)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	return err
}

func (del *delivery) deleteBoard(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	err = del.serv.DeleteBoard(id)
	switch err {
	case ErrDeleteBoard:
		w.WriteHeader(http.StatusInternalServerError)
	case ErrBoardNotFound:
		w.WriteHeader(http.StatusNotFound)
	}
	return err
}
