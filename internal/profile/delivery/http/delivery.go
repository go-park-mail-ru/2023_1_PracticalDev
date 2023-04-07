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
	ErrMissingFile   = errors.New("missing file")
	ErrParseForm     = errors.New("parse form error")
	ErrInvalidUserId = errors.New("invalid user id")
)

func RegisterHandlers(mux *httprouter.Router, logger log.Logger, authorizer middleware.Authorizer, serv profile.Service) {
	del := delivery{serv, logger}

	mux.PUT("/profile", middleware.HandleLogger(middleware.ErrorHandler(middleware.CorsChecker(authorizer(del.fullUpdate)), logger), logger))
	mux.PATCH("/profile", middleware.HandleLogger(middleware.ErrorHandler(middleware.CorsChecker(authorizer(del.partialUpdate)), logger), logger))
}

type delivery struct {
	serv profile.Service
	log  log.Logger
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
	_, _ = io.Copy(buf, file)

	image := models.Image{
		ID:    uuid.NewString() + filepath.Ext(handler.Filename),
		Bytes: buf.Bytes(),
	}

	params := profile.FullUpdateParams{
		Id:         id,
		Username:   r.FormValue("username"),
		Name:       r.FormValue("name"),
		WebsiteUrl: r.FormValue("website_url"),
	}
	prof, err := del.serv.FullUpdate(&params, &image)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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
		return err
	}

	var request partialUpdateRequest
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err = decoder.Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	params := profile.PartialUpdateParams{Id: id}
	if request.Username != nil {
		params.UpdateUsername = true
		params.Username = *request.Username
	}
	if request.Name != nil {
		params.UpdateName = true
		params.Name = *request.Name
	}
	if request.ProfileImage != nil {
		params.UpdateProfileImage = true
		params.ProfileImage = *request.ProfileImage
	}
	if request.WebsiteUrl != nil {
		params.UpdateWebsiteUrl = true
		params.WebsiteUrl = *request.WebsiteUrl
	}

	prof, err := del.serv.PartialUpdate(&params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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
