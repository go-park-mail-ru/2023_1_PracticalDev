package service

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
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
		_, message, err := conn.ReadMessage()
		if err != nil {
			serv.log.Debug("Read from connection error", zap.Error(err), zap.Int("user_id", userID),
				zap.String("remote_addr", conn.RemoteAddr().String()))
			serv.connService.RemoveConnection(userID, conn)
			return nil
		}

		var req Request
		err = req.UnmarshalJSON(message)
		if err != nil {
			serv.log.Debug("Failed to unmarshal message", zap.Error(err), zap.Int("user_id", userID),
				zap.String("remote_addr", conn.RemoteAddr().String()))

			msg := Message{Type: "response", Content: ResponseContent{Message: "json unmarshal error", Code: 40}}
			data, err := msg.MarshalJSON()
			if err != nil {
				serv.log.Error("Marshal json failed", zap.Error(err), zap.Int("user_id", userID),
					zap.String("remote_addr", conn.RemoteAddr().String()), zap.Any("message", msg))

				serv.connService.RemoveConnection(userID, conn)
				return nil
			}

			err = conn.WriteMessage(ws.TextMessage, data)
			if err != nil {
				serv.log.Debug("Write json failed", zap.Error(err), zap.Int("user_id", userID),
					zap.String("remote_addr", conn.RemoteAddr().String()), zap.Any("message", msg))

				serv.connService.RemoveConnection(userID, conn)
				return nil
			}

			continue
		}
		serv.log.Debug("Got message", zap.Int("user_id", userID),
			zap.String("remote_addr", conn.RemoteAddr().String()), zap.Any("message", req))

		var msg Message
		err = serv.rep.MarkAsRead(req.ID)
		if err != nil {
			msg = Message{Type: "response", Content: ResponseContent{Message: "mark as read error", Code: 50}}
		} else {
			msg = Message{Type: "response", Content: ResponseContent{
				Message: "notification mark as read successfully",
				Code:    20,
			}}
		}

		data, err := msg.MarshalJSON()
		if err != nil {
			serv.log.Debug("Marshal json failed", zap.Error(err), zap.Int("user_id", userID),
				zap.String("remote_addr", conn.RemoteAddr().String()), zap.Any("message", msg))

			serv.connService.RemoveConnection(userID, conn)
			return nil
		}

		err = conn.WriteMessage(ws.TextMessage, data)
		if err != nil {
			serv.log.Debug("Write json failed", zap.Error(err), zap.Int("user_id", userID),
				zap.String("remote_addr", conn.RemoteAddr().String()), zap.Any("message", msg))

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

	msg := Message{Type: "notification", Content: notification}
	return serv.connService.Broadcast(msg, userID)
}

func (serv *service) ListUnreadByUser(userID int) ([]models.Notification, error) {
	return serv.rep.ListUnreadByUser(userID)
}
