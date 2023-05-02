package postgres

import (
	"database/sql"

	"github.com/lib/pq"
	"github.com/pkg/errors"

	pkgChats "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/chats"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
)

type repository struct {
	db  *sql.DB
	log log.Logger
}

func NewRepository(db *sql.DB, log log.Logger) pkgChats.Repository {
	return &repository{db, log}
}

const createCmd = `INSERT INTO chats (user1_id, user2_id) 
				   VALUES ($1, $2)
				   RETURNING *;`

func (rep *repository) Create(params *pkgChats.CreateParams) (models.Chat, error) {
	row := rep.db.QueryRow(createCmd, params.User1ID, params.User2ID)

	chat := models.Chat{}
	err := row.Scan(&chat.ID, &chat.User1ID, &chat.User2ID, &chat.CreatedAt, &chat.UpdatedAt)
	if err != nil {
		pgErr, ok := err.(*pq.Error)
		if !ok {
			return models.Chat{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
		}
		if pgErr.Code == "23505" && pgErr.Constraint == "chats_user_pair" {
			return models.Chat{}, errors.Wrap(pkgErrors.ErrChatAlreadyExists, err.Error())
		}

		return models.Chat{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	return chat, nil
}

const listByUserCmd = `SELECT * 
						  FROM chats
						  WHERE user1_id = $1 OR user2_id = $1;`

func (rep *repository) ListByUser(userID int) ([]models.Chat, error) {
	rows, err := rep.db.Query(listByUserCmd, userID)
	if err != nil {
		return nil, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	chats := []models.Chat{}
	chat := models.Chat{}
	for rows.Next() {
		err = rows.Scan(&chat.ID, &chat.User1ID, &chat.User2ID, &chat.CreatedAt, &chat.UpdatedAt)
		if err != nil {
			return nil, errors.Wrap(pkgErrors.ErrDb, err.Error())
		}

		chats = append(chats, chat)
	}
	return chats, nil
}

const getCmd = `SELECT *
				 FROM chats
				 WHERE id = $1;`

func (rep *repository) Get(id int) (models.Chat, error) {
	row := rep.db.QueryRow(getCmd, id)

	chat := models.Chat{}
	err := row.Scan(&chat.ID, &chat.User1ID, &chat.User2ID, &chat.CreatedAt, &chat.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Chat{}, errors.Wrap(pkgErrors.ErrChatNotFound, err.Error())
		} else {
			return models.Chat{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
		}
	}

	return chat, nil
}

const sendMessageCmd = `
		INSERT INTO messages (author_id, chat_id, text)
		VALUES ($1, $2, $3)
		RETURNING *;
	`

func (rep *repository) SendMessage(params *pkgChats.SendMessageParams) (*models.Message, error) {

	row := rep.db.QueryRow(sendMessageCmd, params.AuthorID, params.ChatID, params.Text)

	var msg models.Message
	err := row.Scan(&msg.ID, &msg.AuthorID, &msg.ChatID, &msg.Text, &msg.CreatedAt)
	if err != nil {
		return nil, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	return &msg, nil
}

const chatExistsCmd = `
		SELECT EXISTS(SELECT id
					  FROM chats
					  WHERE user1_id = $1 AND user2_id = $2
						 OR user1_id = $2 AND user2_id = $1) AS exists;`

func (rep *repository) ChatExists(user1ID, user2ID int) (bool, error) {
	row := rep.db.QueryRow(chatExistsCmd, user1ID, user2ID)

	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}
	return exists, nil
}

const getByUsersCmd = `
		SELECT *
		FROM chats
		WHERE user1_id = $1 AND user2_id = $2
		   OR user1_id = $2 AND user2_id = $1;`

func (rep *repository) GetByUsers(user1ID, user2ID int) (models.Chat, error) {
	row := rep.db.QueryRow(getByUsersCmd, user1ID, user2ID)

	chat := models.Chat{}
	err := row.Scan(&chat.ID, &chat.User1ID, &chat.User2ID, &chat.CreatedAt, &chat.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Chat{}, errors.Wrap(pkgErrors.ErrChatNotFound, err.Error())
		} else {
			return models.Chat{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
		}
	}

	return chat, nil
}
