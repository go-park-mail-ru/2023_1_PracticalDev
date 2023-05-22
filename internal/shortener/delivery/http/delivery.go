package http

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/utils"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"strconv"

	mw "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/shortener"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

func RegisterGetHandler(mux *httprouter.Router, logger *zap.Logger, serv shortener.ShortenerService, m *mw.HttpMetricsMiddleware) {
	del := delivery{serv, logger}

	mux.GET("/:link", mw.HandleLogger(mw.ErrorHandler(m.MetricsMiddleware(del.get, logger), logger), logger))
}

func RegisterPostHandler(mux *httprouter.Router, logger *zap.Logger, authorizer mw.Authorizer, csrf mw.CSRFMiddleware, serv shortener.ShortenerService, m *mw.HttpMetricsMiddleware) {
	del := delivery{serv, logger}

	mux.POST("/share/pin/:id", mw.HandleLogger(mw.ErrorHandler(m.MetricsMiddleware(mw.Cors(authorizer(csrf(del.createPin))), logger), logger), logger))
}

type delivery struct {
	serv shortener.ShortenerService
	log  *zap.Logger
}

var shortHost = os.Getenv("SHORT_HOST")

func (del *delivery) create(w http.ResponseWriter, r *http.Request, p httprouter.Params) error { //nolint
	body, err := utils.ReadBody(r, del.log)
	if err != nil {
		return err
	}

	var request url
	err = request.UnmarshalJSON(body)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrParseJson, err.Error())
	}

	hash, err := del.serv.Create(request.URL)
	if err != nil {
		return err
	}

	response := url{URL: shortHost + "/" + hash}
	dt, err := response.MarshalJSON()
	if err != nil {
		return errors.Wrap(pkgErrors.ErrCreateResponse, err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(dt)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrCreateResponse, err.Error())
	}
	return nil
}

func (del *delivery) createPin(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	idStr := p.ByName("id")

	id, err := strconv.Atoi(idStr)
	if idStr == "" || err != nil {
		return pkgErrors.ErrInvalidLinkIDParam
	}

	hash, err := del.serv.CreatePinLink(id)
	if err != nil {
		return err
	}

	response := url{URL: shortHost + "/" + hash}
	dt, err := response.MarshalJSON()
	if err != nil {
		return errors.Wrap(pkgErrors.ErrCreateResponse, err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(dt)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrCreateResponse, err.Error())
	}
	return nil
}

func (del *delivery) get(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	link := p.ByName("link")
	if link == "" {
		return pkgErrors.ErrInvalidLinkIDParam
	}

	res, err := del.serv.Get(link)
	if err == pkgErrors.ErrLinkNotFound {
		return err
	}

	http.Redirect(w, r, res, http.StatusFound)

	return nil
}
