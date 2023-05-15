package http

import (
	"encoding/json"
	"net/http"
	"os"

	mw "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/constants"
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

	mux.POST("/share", mw.HandleLogger(mw.ErrorHandler(m.MetricsMiddleware(mw.Cors(authorizer(csrf(del.create))), logger), logger), logger))
}

type delivery struct {
	serv shortener.ShortenerService
	log  *zap.Logger
}

type url struct {
	URL string `json:"url"`
}

var shortHost = os.Getenv("SHORT_HOST")

func (del *delivery) create(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	decoder := json.NewDecoder(r.Body)
	defer func() {
		err := r.Body.Close()
		if err != nil {
			del.log.Error(constants.FailedCloseRequestBody, zap.Error(err))
		}
	}()

	data := url{}
	if err := decoder.Decode(&data); err != nil {
		return pkgErrors.ErrBadRequest
	}

	hash, err := del.serv.Create(data.URL)
	if err != nil {
		return err
	}

	dt, err := json.Marshal(url{
		URL: shortHost + "/" + hash,
	})
	if err != nil {
		return pkgErrors.ErrCreateResponse
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(dt)
	if err != nil {
		return pkgErrors.ErrCreateResponse
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
