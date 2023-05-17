package notifications

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	ws "github.com/gorilla/websocket"
)

type Service interface {
	HandleConnection(userID int, conn *ws.Conn) error
	Create(userID int, notificationType string, data interface{}) error
	ListUnreadByUser(userID int) ([]models.Notification, error)
}
