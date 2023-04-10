package users

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	mw "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	"github.com/julienschmidt/httprouter"
)

func RegisterHandlers(mux *httprouter.Router, logger log.Logger, authorizer mw.Authorizer, serv Service) {
	del := delivery{serv, logger}

	mux.GET("/users/:id", mw.HandleLogger(mw.ErrorHandler(authorizer(mw.CorsChecker(del.getUser)), logger), logger))
}

type delivery struct {
	serv Service
	log  log.Logger
}

func (del delivery) getUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		return mw.ErrInvalidUserIdParam
	}

	user, err := del.serv.GetUser(id)
	if err != nil {
		return mw.ErrService
	}

	data, err := json.Marshal(user)
	if err != nil {
		return mw.ErrCreateResponse
	}

	_, err = w.Write(data)
	return err
}
