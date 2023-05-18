package http

import (
	"encoding/json"
	pkgComments "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/comments"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/constants"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"strconv"

	mw "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
	"github.com/julienschmidt/httprouter"
)

const (
	commentsUrl = "/pins/:id/comments"
)

type delivery struct {
	serv pkgComments.Service
	log  *zap.Logger
}

func RegisterHandlers(mux *httprouter.Router, logger *zap.Logger, authorizer mw.Authorizer, csrf mw.CSRFMiddleware, serv pkgComments.Service, m *mw.HttpMetricsMiddleware) {
	del := delivery{serv, logger}

	mux.POST(commentsUrl, mw.HandleLogger(mw.ErrorHandler(m.MetricsMiddleware(authorizer(mw.Cors(csrf(del.Create))), logger), logger), logger))
	mux.GET(commentsUrl, mw.HandleLogger(mw.ErrorHandler(m.MetricsMiddleware(authorizer(mw.Cors(csrf(del.List))), logger), logger), logger))
}

func (del delivery) Create(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strUserID := p.ByName("user-id")
	userID, err := strconv.Atoi(strUserID)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrInvalidUserIdParam, err.Error())
	}

	strPinID := p.ByName("id")
	pinID, err := strconv.Atoi(strPinID)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrInvalidPinIdParam, err.Error())
	}

	var request createRequest
	decoder := json.NewDecoder(r.Body)
	defer func() {
		err = r.Body.Close()
		if err != nil {
			del.log.Error(constants.FailedCloseRequestBody, zap.Error(err))
		}
	}()
	if err = decoder.Decode(&request); err != nil {
		return errors.Wrap(pkgErrors.ErrParseJson, err.Error())
	}

	params := pkgComments.CreateParams{
		AuthorID: userID,
		PinID:    pinID,
		Text:     request.Text,
	}
	comment, err := del.serv.Create(&params)
	if err != nil {
		return err
	}

	response := newCreateResponse(&comment)
	data, err := json.Marshal(response)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrCreateResponse, err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrCreateResponse, err.Error())
	}
	return err
}

func (del delivery) List(w http.ResponseWriter, _ *http.Request, p httprouter.Params) error {
	strPinID := p.ByName("id")
	pinID, err := strconv.Atoi(strPinID)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrInvalidPinIdParam, err.Error())
	}

	comments, err := del.serv.List(pinID)
	if err != nil {
		return err
	}

	response := newListResponse(comments)
	data, err := json.Marshal(response)
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
