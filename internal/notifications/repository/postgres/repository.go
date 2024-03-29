package postgres

import (
	"database/sql"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	pkgNotifications "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/notifications"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/constants"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type repository struct {
	db  *sql.DB
	log *zap.Logger
}

func NewRepository(db *sql.DB, log *zap.Logger) pkgNotifications.Repository {
	return &repository{db: db, log: log}
}

const createNotificationCmd = `
		INSERT INTO notifications (user_id, type)
		VALUES ($1, $2)
		RETURNING id;`

const createNewPinNotificationCmd = `
		INSERT INTO new_pin_notifications (notification_id, pin_id)
		VALUES ($1, $2);`

const createNewLikeNotificationCmd = `
		INSERT INTO new_like_notifications (notification_id, pin_id, author_id)
		VALUES ($1, $2, $3);`

const createNewCommentNotificationCmd = `
		INSERT INTO new_comment_notifications (notification_id, pin_id, author_id, text)
		VALUES ($1, $2, $3, $4);`

const createNewFollowerNotificationCmd = `
		INSERT INTO new_follower_notifications (notification_id, follower_id)
		VALUES ($1, $2);`

func (rep *repository) Create(userID int, notificationType string, data interface{}) (int, error) {
	tx, err := rep.db.Begin()
	if err != nil {
		rep.log.Error("DB begin transaction failed", zap.Error(err))
		return 0, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var notificationID int
	err = tx.QueryRow(createNotificationCmd, userID, notificationType).Scan(&notificationID)
	if err != nil {
		rep.log.Error(constants.DBScanError, zap.Error(err), zap.String("sql_query", createNotificationCmd),
			zap.Int("user_id", userID), zap.String("notification_type", notificationType))

		return 0, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	switch notificationType {
	case constants.NewPin:
		np := data.(models.NewPinNotification)
		_, err = tx.Exec(createNewPinNotificationCmd, notificationID, np.PinID)
	case constants.NewLike:
		nl := data.(models.NewLikeNotification)
		_, err = tx.Exec(createNewLikeNotificationCmd, notificationID, nl.PinID, nl.AuthorID)
	case constants.NewComment:
		nc := data.(models.NewCommentNotification)
		_, err = tx.Exec(createNewCommentNotificationCmd, notificationID, nc.PinID, nc.AuthorID, nc.Text)
	case constants.NewFollower:
		nf := data.(models.NewFollowerNotification)
		_, err = tx.Exec(createNewFollowerNotificationCmd, notificationID, nf.FollowerID)
	}
	if err != nil {
		rep.log.Error(constants.DBQueryError, zap.Error(err), zap.Int("notification_id", notificationID))
		return 0, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	err = tx.Commit()
	if err != nil {
		rep.log.Error("DB commit failed", zap.Error(err))
		return 0, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}
	return notificationID, nil
}

const GetNotificationCmd = `
		SELECT n.id, n.user_id, n.created_at, n.is_read, n.type,
			np.pin_id,
			nl.pin_id, nl.author_id,
			nc.pin_id, nc.author_id, nc.text,
			nf.follower_id
		FROM notifications n
		LEFT JOIN new_pin_notifications np ON n.id = np.notification_id
		LEFT JOIN new_like_notifications nl ON n.id = nl.notification_id
		LEFT JOIN new_comment_notifications nc ON n.id = nc.notification_id
		LEFT JOIN new_follower_notifications nf ON n.id = nf.notification_id
		WHERE n.id = $1;`

func (rep *repository) Get(notificationID int) (*models.Notification, error) {
	row := rep.db.QueryRow(GetNotificationCmd, notificationID)

	var npPinID, nlPinID, nlAuthorID, ncPinID, ncAuthorID, nfFollowerID sql.NullInt32
	var ncText sql.NullString
	notification := &models.Notification{}
	err := row.Scan(&notification.ID, &notification.UserID, &notification.CreatedAt, &notification.IsRead,
		&notification.Type, &npPinID, &nlPinID, &nlAuthorID, &ncPinID, &ncAuthorID, &ncText, &nfFollowerID)
	if err != nil {
		rep.log.Error(constants.DBScanError, zap.Error(err), zap.String("sql_query", GetNotificationCmd),
			zap.Int("notification_id", notificationID))

		if err == sql.ErrNoRows {
			return nil, errors.Wrap(pkgErrors.ErrNotificationNotFound, err.Error())
		}
		return nil, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	switch notification.Type {
	case constants.NewPin:
		notification.Data = models.NewPinNotification{PinID: int(npPinID.Int32)}
	case constants.NewLike:
		notification.Data = models.NewLikeNotification{PinID: int(nlPinID.Int32), AuthorID: int(nlAuthorID.Int32)}
	case constants.NewComment:
		notification.Data = models.NewCommentNotification{PinID: int(ncPinID.Int32), AuthorID: int(ncAuthorID.Int32),
			Text: ncText.String}
	case constants.NewFollower:
		notification.Data = models.NewFollowerNotification{FollowerID: int(nfFollowerID.Int32)}
	}

	return notification, nil
}

const listUnreadByUserCmd = `
		SELECT n.id, n.user_id, n.created_at, n.is_read, n.type,
			np.pin_id,
			nl.pin_id, nl.author_id,
			nc.pin_id, nc.author_id, nc.text,
			nf.follower_id
		FROM notifications n
		LEFT JOIN new_pin_notifications np ON n.id = np.notification_id
		LEFT JOIN new_like_notifications nl ON n.id = nl.notification_id
		LEFT JOIN new_comment_notifications nc ON n.id = nc.notification_id
		LEFT JOIN new_follower_notifications nf ON n.id = nf.notification_id
		WHERE n.user_id = $1 AND n.is_read = false;`

func (rep *repository) ListUnreadByUser(userID int) ([]models.Notification, error) {
	rows, err := rep.db.Query(listUnreadByUserCmd, userID)
	if err != nil {
		return nil, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	var npPinID, nlPinID, nlAuthorID, ncPinID, ncAuthorID, nfFollowerID sql.NullInt32
	var ncText sql.NullString
	notifications := []models.Notification{}
	notification := models.Notification{}
	for rows.Next() {
		err = rows.Scan(&notification.ID, &notification.UserID, &notification.CreatedAt, &notification.IsRead,
			&notification.Type, &npPinID, &nlPinID, &nlAuthorID, &ncPinID, &ncAuthorID, &ncText, &nfFollowerID)
		if err != nil {
			rep.log.Error(constants.DBScanError, zap.Error(err), zap.String("sql_query", listUnreadByUserCmd),
				zap.Int("user_id", userID))

			return nil, errors.Wrap(pkgErrors.ErrDb, err.Error())
		}

		switch notification.Type {
		case constants.NewPin:
			notification.Data = models.NewPinNotification{PinID: int(npPinID.Int32)}
		case constants.NewLike:
			notification.Data = models.NewLikeNotification{PinID: int(nlPinID.Int32), AuthorID: int(nlAuthorID.Int32)}
		case constants.NewComment:
			notification.Data = models.NewCommentNotification{PinID: int(ncPinID.Int32), AuthorID: int(ncAuthorID.Int32),
				Text: ncText.String}
		case constants.NewFollower:
			notification.Data = models.NewFollowerNotification{FollowerID: int(nfFollowerID.Int32)}
		}

		notifications = append(notifications, notification)
	}

	return notifications, nil
}

const MarkAsReadCmd = `
		UPDATE notifications 
		SET is_read = true 
		WHERE id = $1;`

func (rep *repository) MarkAsRead(notificationID int) error {
	_, err := rep.db.Exec(MarkAsReadCmd, notificationID)
	if err != nil {
		rep.log.Error(constants.DBQueryError, zap.Error(err), zap.String("sql_query", MarkAsReadCmd),
			zap.Int("notification_id", notificationID))

		return errors.Wrap(pkgErrors.ErrDb, err.Error())
	}
	return nil
}
