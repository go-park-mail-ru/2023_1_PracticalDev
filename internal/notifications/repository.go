package notifications

import "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"

type Repository interface {
	Create(userID int, notificationType string, data interface{}) (int, error)
	Get(notificationID int) (*models.Notification, error)
	ListUnreadByUser(userID int) ([]models.Notification, error)
	MarkAsRead(notificationID int) error
}
