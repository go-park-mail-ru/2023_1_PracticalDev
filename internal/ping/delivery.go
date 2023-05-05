package ping

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
)

func RegisterHandlers(mux *httprouter.Router, logger *zap.Logger) {
	del := delivery{logger}

	mux.GET("/ping", middleware.HandleLogger(middleware.ErrorHandler(middleware.Cors(del.Ping), logger), logger))
	mux.POST("/ping", middleware.HandleLogger(middleware.ErrorHandler(middleware.Cors(del.Ping), logger), logger))
	mux.DELETE("/ping", middleware.HandleLogger(middleware.ErrorHandler(middleware.Cors(del.Ping), logger), logger))
}

type delivery struct {
	log *zap.Logger
}

func (del *delivery) Ping(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	w.Header().Set("Content-Type", "text/plain")
	_, err := w.Write([]byte("Pong"))
	return err
}
