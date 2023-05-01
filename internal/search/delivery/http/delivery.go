package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	mw "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/log"
	serv "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/search"
	"github.com/julienschmidt/httprouter"

	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
)

func RegisterHandlers(mux *httprouter.Router, logger log.Logger, authorizer mw.Authorizer, serv serv.Service) {
	del := delivery{serv, logger}
	mux.GET("/search/:query", mw.HandleLogger(mw.ErrorHandler(authorizer(mw.Cors(del.get)), logger), logger))
}

type delivery struct {
	serv.Service
	log log.Logger
}

func (del delivery) get(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strUserId := p.ByName("user-id")
	userId, err := strconv.Atoi(strUserId)
	if err != nil {
		return pkgErrors.ErrInvalidUserIdParam
	}

	query := p.ByName("query")
	res, err := del.Service.Get(userId, query)
	encoder := json.NewEncoder(w)
	encoder.Encode(res)
	return err
}
