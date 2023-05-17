package service

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/notifications"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/connectionservice"
	ws "github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type service struct {
	rep         notifications.Repository
	connService *connectionservice.Service
	log         *zap.Logger
}

func NewService(rep notifications.Repository, logger *zap.Logger) notifications.Service {
	return &service{rep: rep, connService: connectionservice.NewService(logger), log: logger}
}

func (serv *service) HandleConnection(userID int, conn *ws.Conn) error {
	serv.connService.AddConnection(userID, conn)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			serv.log.Debug("Read from connection error", zap.Error(err), zap.Int("user_id", userID),
				zap.String("remote_addr", conn.RemoteAddr().String()))
			serv.connService.RemoveConnection(userID, conn)
			return nil
		}
	}
}

func (serv *service) Create(userID int, notificationType string, data interface{}) error {
	notificationID, err := serv.rep.Create(userID, notificationType, data)
	if err != nil {
		return err
	}

	notification, err := serv.rep.Get(notificationID)
	if err != nil {
		return err
	}

	err = serv.connService.Broadcast(notification, userID)
	if err != nil {
		return err
	}

	return nil
}
