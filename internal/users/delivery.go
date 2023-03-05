package users

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	"github.com/julienschmidt/httprouter"
)

func RegisterHandlers(mux *httprouter.Router, logger log.Logger, authorizer auth.Authorizer, serv Service) {
	del := delivery{serv, logger}

	mux.GET("/users/:id", middleware.Logger(middleware.ErrorHandler(authorizer(middleware.CorsChecker(del.getUser)), logger), logger))
}

type delivery struct {
	serv Service
	log  log.Logger
}

func (del delivery) getUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	str_id := p.ByName("id")
	id, err := strconv.Atoi(str_id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	user, err := del.serv.GetUser(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return err
	}

	data, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	_, err = w.Write(data)
	return err
}
