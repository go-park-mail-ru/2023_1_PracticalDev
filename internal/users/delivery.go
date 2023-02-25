package users

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func chainLogger(handler httprouter.Handle, logger *log.Logger) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		logger.Println("New request from:", r.RemoteAddr)
		handler(w, r, p)
	}
}

func RegisterHandlers(mux *httprouter.Router, logger *log.Logger) {
	mux.GET("/users/:id", chainLogger(getUser, logger))
}

func getUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	_, err := w.Write([]byte("{user : {id:" + id + "}}"))
	if err != nil {
		return
	}
}
