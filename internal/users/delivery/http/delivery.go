package http

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	mw "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/users"
)

func RegisterHandlers(mux *httprouter.Router, logger *zap.Logger, authorizer mw.Authorizer, csrf mw.CSRFMiddleware, serv users.Service, m *mw.HttpMetricsMiddleware) {
	del := delivery{serv, logger}

	mux.GET("/users/:id", mw.HandleLogger(mw.ErrorHandler(m.MetricsMiddleware(authorizer(mw.Cors(csrf(del.get))), logger), logger), logger))
}

type delivery struct {
	serv users.Service
	log  *zap.Logger
}

func (del *delivery) get(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		return pkgErrors.ErrInvalidUserIdParam
	}

	user, err := del.serv.Get(id)
	if err != nil {
		return err
	}

	data, err := json.Marshal(user)
	if err != nil {
		return pkgErrors.ErrCreateResponse
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return pkgErrors.ErrCreateResponse
	}
	return nil
}
