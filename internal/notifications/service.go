package notifications

import ws "github.com/gorilla/websocket"

type Service interface {
	HandleConnection(userID int, conn *ws.Conn) error
	Create(userID int, notificationType string, data interface{}) error
}
