package ping

import (
	"net/http"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	"github.com/julienschmidt/httprouter"
)

func RegisterHandlers(mux *httprouter.Router, logger log.Logger) {
	del := delivery{logger}

	mux.GET("/ping", middleware.HandleLogger(middleware.ErrorHandler(middleware.CorsChecker(del.Ping), logger), logger))
	mux.POST("/ping", middleware.HandleLogger(middleware.ErrorHandler(middleware.CorsChecker(del.Ping), logger), logger))
	mux.DELETE("/ping", middleware.HandleLogger(middleware.ErrorHandler(middleware.CorsChecker(del.Ping), logger), logger))
}

type delivery struct {
	log log.Logger
}

func (del *delivery) Ping(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	w.WriteHeader(http.StatusOK)

	w.Header().Set("Content-Type", "text/plain")
	_, err := w.Write([]byte("Pong"))
	return err
}
