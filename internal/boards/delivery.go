package boards

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func RegisterHandlers(mux *httprouter.Router, logger log.Logger, authorizer middleware.Authorizer, serv Service) {
	del := delivery{serv, logger}

	mux.GET("/boards", middleware.HandleLogger(middleware.ErrorHandler(authorizer(middleware.CorsChecker(del.getBoards)), logger), logger))
}

type delivery struct {
	serv Service
	log  log.Logger
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

	data, err := json.Marshal(boards)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	return err
}
