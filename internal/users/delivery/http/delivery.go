package http

import (
	"github.com/pkg/errors"
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

	mux.GET("/users/:id", mw.HandleLogger(mw.ErrorHandler(m.MetricsMiddleware(authorizer(mw.Cors(csrf(del.Get))), logger), logger), logger))
}

type delivery struct {
	serv users.Service
	log  *zap.Logger
}

func (del *delivery) Get(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrInvalidUserIdParam, err.Error())
	}

	user, err := del.serv.Get(id)
	if err != nil {
		return err
	}

	response := newGetResponse(&user)
	data, err := response.MarshalJSON()
	if err != nil {
		return errors.Wrap(pkgErrors.ErrCreateResponse, err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrCreateResponse, err.Error())
	}
	return nil
}
