package posts

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	"github.com/julienschmidt/httprouter"
)

func RegisterHandlers(mux *httprouter.Router, logger log.Logger, authorizer auth.Authorizer, serv Service) {
	del := delivery{serv, logger}

	mux.GET("/posts", middleware.HandleLogger(middleware.ErrorHandler(authorizer(middleware.CorsChecker(del.getPosts)), logger), logger))
}

type delivery struct {
	serv Service
	log  log.Logger
}

func (del *delivery) getPosts(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	queryValues := r.URL.Query()
	page := 1
	var err error
	strPage := queryValues.Get("page")
	if strPage != "" {
		page, err = strconv.Atoi(strPage)
		if err != nil || page < 1 {
			w.WriteHeader(http.StatusBadRequest)
			return err
		}
	}

	limit := 30
	strLimit := queryValues.Get("limit")
	if strLimit != "" {
		limit, err = strconv.Atoi(strLimit)
		if err != nil || limit < 0 {
			w.WriteHeader(http.StatusBadRequest)
			return err
		}
	}

	pins, err := del.serv.GetPosts(page, limit)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return err
	}

	data, err := json.Marshal(pins)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	return err
}
