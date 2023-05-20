package http

import (
	"encoding/json"
	pkgNotifications "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/notifications"
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
	notificationsUrl   = "/notifications"
	wsNotificationsUrl = "/ws/notifications"
)

type delivery struct {
	serv     pkgNotifications.Service
	log      *zap.Logger
	upgrader ws.Upgrader
}

func RegisterHandlers(mux *httprouter.Router, logger *zap.Logger, authorizer mw.Authorizer, csrf mw.CSRFMiddleware,
	serv pkgNotifications.Service, m *mw.HttpMetricsMiddleware) {
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

	// notifications
	mux.GET(notificationsUrl, mw.HandleLogger(mw.ErrorHandler(mw.Cors(authorizer(csrf(del.ListUnreadByUser))), logger),
		logger))

	// connect to websocket
	mux.GET(wsNotificationsUrl, mw.HandleLogger(mw.ErrorHandler(m.MetricsMiddleware(authorizer(
		del.notificationsHandler), logger), logger), logger))
}

func (del *delivery) ListUnreadByUser(w http.ResponseWriter, _ *http.Request, p httprouter.Params) error {
	strUserID := p.ByName("user-id")
	userID, err := strconv.Atoi(strUserID)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrInvalidUserIdParam, err.Error())
	}

	notifications, err := del.serv.ListUnreadByUser(userID)
	if err != nil {
		return err
	}

	response := newListResponse(notifications)
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
