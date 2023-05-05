package http

import (
	"encoding/json"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	mw "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
	serv "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/search"
)

func RegisterHandlers(mux *httprouter.Router, logger *zap.Logger, authorizer mw.Authorizer, serv serv.Service) {
	del := delivery{serv, logger}
	mux.GET("/search/:query", mw.HandleLogger(mw.ErrorHandler(authorizer(mw.Cors(del.get)), logger), logger))
}

type delivery struct {
	serv.Service
	log *zap.Logger
}

func (del delivery) get(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strUserId := p.ByName("user-id")
	userId, err := strconv.Atoi(strUserId)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrInvalidUserIdParam, err.Error())
	}

	query := p.ByName("query")
	res, err := del.Service.Get(userId, query)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(res)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrCreateResponse, err.Error())
	}
	return nil
}
