package http

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/notifications"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/constants"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
	ws "github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"strconv"

	mw "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
)

const (
	wsNotificationsUrl = "/notifications"
)

type delivery struct {
	serv     notifications.Service
	log      *zap.Logger
	upgrader ws.Upgrader
}

func RegisterHandlers(mux *httprouter.Router, logger *zap.Logger, authorizer mw.Authorizer, serv notifications.Service,
	m *mw.HttpMetricsMiddleware) {
	del := delivery{
		serv: serv,
		log:  logger,
		upgrader: ws.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}}

	// connect to websocket
	mux.GET(wsNotificationsUrl, mw.HandleLogger(mw.ErrorHandler(m.MetricsMiddleware(authorizer(
		del.notificationsHandler), logger), logger), logger))
}

func (del *delivery) notificationsHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strUserID := p.ByName("user-id")
	userID, err := strconv.Atoi(strUserID)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrInvalidUserIdParam, err.Error())
	}

	conn, err := del.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrUpgradeToWebSocket, err.Error())
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			del.log.Error(constants.FailedCloseConnection, zap.Error(err))
		}
		del.log.Debug("Websocket closed", zap.Int("user_id", userID),
			zap.String("remote_addr", conn.RemoteAddr().String()))
	}()
	del.log.Debug("Websocket connected", zap.Int("user_id", userID),
		zap.String("remote_addr", conn.RemoteAddr().String()))

	return del.serv.HandleConnection(userID, conn)
}
