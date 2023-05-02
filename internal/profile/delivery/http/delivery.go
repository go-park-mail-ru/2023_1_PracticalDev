package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"

	mw "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/profile"
)

func RegisterHandlers(mux *httprouter.Router, logger log.Logger, authorizer mw.Authorizer, serv profile.Service, m *mw.HttpMetricsMiddleware) {
	del := delivery{serv, logger}

	mux.GET("/users/:id/profile", mw.HandleLogger(mw.ErrorHandler(m.MetricsMiddleware(mw.Cors(authorizer(del.getProfileByUser)), logger), logger), logger))
	mux.PUT("/profile", mw.HandleLogger(mw.ErrorHandler(m.MetricsMiddleware(mw.Cors(authorizer(del.fullUpdate)), logger), logger), logger))
	mux.PATCH("/profile", mw.HandleLogger(mw.ErrorHandler(m.MetricsMiddleware(mw.Cors(authorizer(del.partialUpdate)), logger), logger), logger))
}

type delivery struct {
	serv profile.Service
	log  log.Logger
}

func (del *delivery) getProfileByUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strUserId := p.ByName("id")
	userId, err := strconv.Atoi(strUserId)
	if err != nil {
		return pkgErrors.ErrInvalidUserIdParam
	}

	prof, err := del.serv.GetProfileByUser(userId)
	if err != nil {
		return err
	}

	response := newGetResponse(&prof)
	data, err := json.Marshal(response)
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

func (del *delivery) fullUpdate(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("user-id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		return pkgErrors.ErrInvalidUserIdParam
	}

	file, handler, err := r.FormFile("bytes")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			return pkgErrors.ErrMissingFile
		} else {
			return pkgErrors.ErrParseForm
		}
	}
	defer func() {
		err = file.Close()
		if err != nil {
			del.log.Error(err)
		}
	}()

	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, file)
	if err != nil {
		return pkgErrors.ErrFileCopy
	}

	image := models.Image{
		ID:    uuid.NewString() + filepath.Ext(handler.Filename),
		Bytes: buf.Bytes(),
	}

	params := profile.FullUpdateParams{
		Id:           id,
		Username:     r.FormValue("username"),
		Name:         r.FormValue("name"),
		ProfileImage: image,
		WebsiteUrl:   r.FormValue("website_url"),
	}
	prof, err := del.serv.FullUpdate(&params)
	if err != nil {
		return err
	}

	response := newFullUpdateResponse(&prof)
	data, err := json.Marshal(response)
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

func (del *delivery) partialUpdate(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("user-id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		return pkgErrors.ErrInvalidUserIdParam
	}

	params := profile.PartialUpdateParams{Id: id}

	file, handler, err := r.FormFile("bytes")
	if err != nil {
		if err != http.ErrMissingFile {
			return pkgErrors.ErrParseForm
		}
	} else {
		buf := bytes.NewBuffer(nil)
		_, err = io.Copy(buf, file)
		if err != nil {
			return pkgErrors.ErrFileCopy
		}

		image := models.Image{
			ID:    uuid.NewString() + filepath.Ext(handler.Filename),
			Bytes: buf.Bytes(),
		}

		params.UpdateProfileImage = true
		params.ProfileImage = image
	}

	params.UpdateUsername = r.Form.Has("username")
	if params.UpdateUsername {
		params.Username = r.Form.Get("username")
	}
	params.UpdateName = r.Form.Has("name")
	if params.UpdateName {
		params.Name = r.Form.Get("name")
	}
	params.UpdateWebsiteUrl = r.Form.Has("website_url")
	if params.UpdateWebsiteUrl {
		params.WebsiteUrl = r.Form.Get("website_url")
	}

	prof, err := del.serv.PartialUpdate(&params)
	if err != nil {
		return err
	}

	response := newPartialUpdateResponse(&prof)
	data, err := json.Marshal(response)
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
