package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/profile"
)

var (
	ErrFileCopy        = errors.New("file copy error")
	ErrMissingFile     = errors.New("missing file")
	ErrParseForm       = errors.New("parse form error")
	ErrInvalidUserId   = errors.New("invalid user id")
	ErrProfileNotFound = errors.New("profile not found")
	ErrService         = errors.New("service error")
)

func RegisterHandlers(mux *httprouter.Router, logger log.Logger, authorizer middleware.Authorizer, serv profile.Service) {
	del := delivery{serv, logger}

	mux.GET("/users/:id/profile", middleware.HandleLogger(middleware.ErrorHandler(middleware.CorsChecker(authorizer(del.getProfileByUser)), logger), logger))
	mux.PUT("/profile", middleware.HandleLogger(middleware.ErrorHandler(middleware.CorsChecker(authorizer(del.fullUpdate)), logger), logger))
	mux.PATCH("/profile", middleware.HandleLogger(middleware.ErrorHandler(middleware.CorsChecker(authorizer(del.partialUpdate)), logger), logger))
}

type delivery struct {
	serv profile.Service
	log  log.Logger
}

func (del *delivery) getProfileByUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strUserId := p.ByName("id")
	userId, err := strconv.Atoi(strUserId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return ErrInvalidUserId
	}

	prof, err := del.serv.GetProfileByUser(userId)
	if err != nil {
		if err == profile.ErrProfileNotFound {
			err = ErrProfileNotFound
			w.WriteHeader(http.StatusNotFound)
		} else {
			err = ErrService
			w.WriteHeader(http.StatusInternalServerError)
		}
		return err
	}

	response := getProfileResponse{
		Username:     prof.Username,
		Name:         prof.Name,
		ProfileImage: prof.ProfileImage,
		WebsiteUrl:   prof.WebsiteUrl,
	}
	data, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	return err
}

func (del *delivery) fullUpdate(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("user-id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return ErrInvalidUserId
	}

	file, handler, err := r.FormFile("bytes")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err == http.ErrMissingFile {
			err = ErrMissingFile
		} else {
			err = ErrParseForm
		}
		return err
	}
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return ErrFileCopy
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
		if err == profile.ErrUsernameAlreadyExists ||
			err == profile.ErrTooLongUsername ||
			err == profile.ErrTooShortUsername ||
			err == profile.ErrTooLongName ||
			err == profile.ErrEmptyName {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return err
	}

	response := fullUpdateResponse{
		Username:     prof.Username,
		Name:         prof.Name,
		ProfileImage: prof.ProfileImage,
		WebsiteUrl:   prof.WebsiteUrl,
	}
	data, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	return err
}

func (del *delivery) partialUpdate(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strId := p.ByName("user-id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return ErrInvalidUserId
	}

	params := profile.PartialUpdateParams{Id: id}

	file, handler, err := r.FormFile("bytes")
	if err != nil {
		if err != http.ErrMissingFile {
			w.WriteHeader(http.StatusBadRequest)
			err = ErrParseForm
			return err
		}
	} else {
		buf := bytes.NewBuffer(nil)
		_, err = io.Copy(buf, file)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return ErrFileCopy
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
		if err == profile.ErrUsernameAlreadyExists ||
			err == profile.ErrTooLongUsername ||
			err == profile.ErrTooShortUsername ||
			err == profile.ErrTooLongName ||
			err == profile.ErrEmptyName {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return err
	}

	response := partialUpdateResponse{
		Username:     prof.Username,
		Name:         prof.Name,
		ProfileImage: prof.ProfileImage,
		WebsiteUrl:   prof.WebsiteUrl,
	}
	data, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	return err
}
