package postgres

import (
	"github.com/DATA-DOG/go-sqlmock"
	pkgChats "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/chats"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"reflect"
	"regexp"
	"testing"
)

func TestRepository_Create(t *testing.T) {
	type fields struct {
		mock   sqlmock.Sqlmock
		params *pkgChats.CreateParams
		chat   *models.Chat
	}

	type testCase struct {
		prepare func(f *fields)
		params  *pkgChats.CreateParams
		chat    models.Chat
		err     error
	}

	tests := map[string]testCase{
		"good query": {
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"id", "user1_id", "user2_id", "created_at", "updated_at"})
				rows = rows.AddRow(f.chat.ID, f.chat.User1ID, f.chat.User2ID, f.chat.CreatedAt, f.chat.UpdatedAt)
				f.mock.
					ExpectQuery(regexp.QuoteMeta(createCmd)).
					WithArgs(f.params.User1ID, f.params.User2ID).
					WillReturnRows(rows)
			},
			params: &pkgChats.CreateParams{User1ID: 2, User2ID: 3},
			chat:   models.Chat{ID: 1, User1ID: 2, User2ID: 3},
			err:    nil,
		},
		"query error": {
			prepare: func(f *fields) {
				f.mock.
					ExpectQuery(regexp.QuoteMeta(createCmd)).
					WithArgs(f.params.User1ID, f.params.User2ID).
					WillReturnError(&pq.Error{Message: "sql error"})
			},
			params: &pkgChats.CreateParams{User1ID: 2, User2ID: 3},
			chat:   models.Chat{},
			err:    pkgErrors.ErrDb,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			logger, err := zap.NewDevelopment()
			if err != nil {
				t.Fatalf("can't create logger: %s", err)
			}

			db, sqlMock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("can't create mock: %s", err)
			}
			defer db.Close()

			f := fields{mock: sqlMock, params: test.params, chat: &test.chat}
			if test.prepare != nil {
				test.prepare(&f)
			}

			repo := NewRepository(db, logger)
			chat, err := repo.Create(test.params)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if chat != test.chat {
				t.Errorf("\nExpected: %v\nGot: %v", test.chat, chat)
			}
			if err = sqlMock.ExpectationsWereMet(); err != nil {
				t.Errorf("\nThere were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestRepository_ListByUser(t *testing.T) {
	type fields struct {
		mock   sqlmock.Sqlmock
		userID int
		chats  []models.Chat
	}

	type testCase struct {
		prepare func(f *fields)
		userID  int
		chats   []models.Chat
		err     error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"id", "user1_id", "user2_id", "created_at", "updated_at"})
				for _, chat := range f.chats {
					rows = rows.AddRow(chat.ID, chat.User1ID, chat.User2ID, chat.CreatedAt, chat.UpdatedAt)
				}
				f.mock.
					ExpectQuery(regexp.QuoteMeta(listByUserCmd)).
					WithArgs(f.userID).
					WillReturnRows(rows)
			},
			userID: 2,
			chats: []models.Chat{
				{ID: 2, User1ID: 2, User2ID: 3},
				{ID: 3, User1ID: 8, User2ID: 2},
				{ID: 4, User1ID: 2, User2ID: 4},
			},
			err: nil,
		},
		"query error": {
			prepare: func(f *fields) {
				f.mock.
					ExpectQuery(regexp.QuoteMeta(listByUserCmd)).
					WithArgs(f.userID).
					WillReturnError(&pq.Error{Message: "sql error"})
			},
			userID: 2,
			chats:  nil,
			err:    pkgErrors.ErrDb,
		},
		"row scan error": {
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"id", "user1_id"}).AddRow(1, 2)
				f.mock.
					ExpectQuery(regexp.QuoteMeta(listByUserCmd)).
					WithArgs(f.userID).
					WillReturnRows(rows)
			},
			userID: 2,
			chats:  nil,
			err:    pkgErrors.ErrDb,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			logger, err := zap.NewDevelopment()
			if err != nil {
				t.Fatalf("can't create logger: %s", err)
			}

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("can't create mock: %s", err)
			}
			defer db.Close()

			f := fields{mock: mock, userID: test.userID, chats: test.chats}
			if test.prepare != nil {
				test.prepare(&f)
			}

			repo := NewRepository(db, logger)
			chats, err := repo.ListByUser(test.userID)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if !reflect.DeepEqual(chats, test.chats) {
				t.Errorf("\nExpected: %v\nGot: %v", test.chats, chats)
			}
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("\nThere were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestRepository_Get(t *testing.T) {
	type fields struct {
		mock   sqlmock.Sqlmock
		chatID int
		chat   *models.Chat
	}

	type testCase struct {
		prepare func(f *fields)
		chatID  int
		chat    models.Chat
		err     error
	}

	tests := map[string]testCase{
		"good query": {
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"id", "user1_id", "user2_id", "created_at", "updated_at"})
				rows = rows.AddRow(f.chat.ID, f.chat.User1ID, f.chat.User2ID, f.chat.CreatedAt, f.chat.UpdatedAt)
				f.mock.
					ExpectQuery(regexp.QuoteMeta(getCmd)).
					WithArgs(f.chatID).
					WillReturnRows(rows)
			},
			chatID: 2,
			chat:   models.Chat{ID: 2, User1ID: 2, User2ID: 3},
			err:    nil,
		},
		"query error": {
			prepare: func(f *fields) {
				f.mock.
					ExpectQuery(regexp.QuoteMeta(getCmd)).
					WithArgs(f.chatID).
					WillReturnError(&pq.Error{Message: "sql error"})
			},
			chatID: 2,
			chat:   models.Chat{},
			err:    pkgErrors.ErrDb,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			logger, err := zap.NewDevelopment()
			if err != nil {
				t.Fatalf("can't create logger: %s", err)
			}

			db, sqlMock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("can't create mock: %s", err)
			}
			defer db.Close()

			f := fields{mock: sqlMock, chatID: test.chatID, chat: &test.chat}
			if test.prepare != nil {
				test.prepare(&f)
			}

			repo := NewRepository(db, logger)
			chat, err := repo.Get(test.chatID)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if chat != test.chat {
				t.Errorf("\nExpected: %v\nGot: %v", test.chat, chat)
			}
			if err = sqlMock.ExpectationsWereMet(); err != nil {
				t.Errorf("\nThere were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestRepository_MessagesList(t *testing.T) {
	type fields struct {
		mock     sqlmock.Sqlmock
		chatID   int
		messages []models.Message
	}

	type testCase struct {
		prepare  func(f *fields)
		chatID   int
		messages []models.Message
		err      error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"id", "author_id", "chat_id", "text", "created_at"})
				for _, message := range f.messages {
					rows = rows.AddRow(message.ID, message.AuthorID, message.ChatID, message.Text, message.CreatedAt)
				}
				f.mock.
					ExpectQuery(regexp.QuoteMeta(messagesListCmd)).
					WithArgs(f.chatID).
					WillReturnRows(rows)
			},
			chatID: 2,
			messages: []models.Message{
				{ID: 1, AuthorID: 3, ChatID: 2, Text: "msg 1"},
				{ID: 2, AuthorID: 4, ChatID: 2, Text: "msg 2"},
				{ID: 3, AuthorID: 3, ChatID: 2, Text: "msg 3"},
			},
			err: nil,
		},
		"query error": {
			prepare: func(f *fields) {
				f.mock.
					ExpectQuery(regexp.QuoteMeta(messagesListCmd)).
					WithArgs(f.chatID).
					WillReturnError(&pq.Error{Message: "sql error"})
			},
			chatID:   2,
			messages: nil,
			err:      pkgErrors.ErrDb,
		},
		"row scan error": {
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"id", "author_id"}).AddRow(1, 2)
				f.mock.
					ExpectQuery(regexp.QuoteMeta(messagesListCmd)).
					WithArgs(f.chatID).
					WillReturnRows(rows)
			},
			chatID:   2,
			messages: nil,
			err:      pkgErrors.ErrDb,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			logger, err := zap.NewDevelopment()
			if err != nil {
				t.Fatalf("can't create logger: %s", err)
			}

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("can't create mock: %s", err)
			}
			defer db.Close()

			f := fields{mock: mock, chatID: test.chatID, messages: test.messages}
			if test.prepare != nil {
				test.prepare(&f)
			}

			repo := NewRepository(db, logger)
			messages, err := repo.MessagesList(test.chatID)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if !reflect.DeepEqual(messages, test.messages) {
				t.Errorf("\nExpected: %v\nGot: %v", test.messages, messages)
			}
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("\nThere were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestRepository_CreateMessage(t *testing.T) {
	type fields struct {
		mock    sqlmock.Sqlmock
		params  *pkgChats.CreateMessageParams
		message *models.Message
	}

	type testCase struct {
		prepare func(f *fields)
		params  *pkgChats.CreateMessageParams
		message models.Message
		err     error
	}

	tests := map[string]testCase{
		"good query": {
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"id", "author_id", "chat_id", "text", "created_at"})
				rows = rows.AddRow(f.message.ID, f.message.AuthorID, f.message.ChatID, f.message.Text,
					f.message.CreatedAt)
				f.mock.
					ExpectQuery(regexp.QuoteMeta(createMessageCmd)).
					WithArgs(f.params.AuthorID, f.params.ChatID, f.params.Text).
					WillReturnRows(rows)
			},
			params:  &pkgChats.CreateMessageParams{AuthorID: 2, ChatID: 3, Text: "hello!"},
			message: models.Message{ID: 1, AuthorID: 2, ChatID: 3, Text: "hello!"},
			err:     nil,
		},
		"query error": {
			prepare: func(f *fields) {
				f.mock.
					ExpectQuery(regexp.QuoteMeta(createMessageCmd)).
					WithArgs(f.params.AuthorID, f.params.ChatID, f.params.Text).
					WillReturnError(&pq.Error{Message: "sql error"})
			},
			params:  &pkgChats.CreateMessageParams{AuthorID: 2, ChatID: 3, Text: "hello!"},
			message: models.Message{},
			err:     pkgErrors.ErrDb,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			logger, err := zap.NewDevelopment()
			if err != nil {
				t.Fatalf("can't create logger: %s", err)
			}

			db, sqlMock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("can't create mock: %s", err)
			}
			defer db.Close()

			f := fields{mock: sqlMock, params: test.params, message: &test.message}
			if test.prepare != nil {
				test.prepare(&f)
			}

			repo := NewRepository(db, logger)
			message, err := repo.CreateMessage(test.params)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if err != nil {
				if message != nil {
					t.Errorf("\nExpected: %v\nGot: %v", nil, message)
				}
			} else if *message != test.message {
				t.Errorf("\nExpected: %v\nGot: %v", test.message, message)
			}
			if err = sqlMock.ExpectationsWereMet(); err != nil {
				t.Errorf("\nThere were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestRepository_ChatExists(t *testing.T) {
	type fields struct {
		mock    sqlmock.Sqlmock
		user1ID int
		user2ID int
		exists  bool
	}

	type testCase struct {
		prepare func(f *fields)
		user1ID int
		user2ID int
		exists  bool
		err     error
	}

	tests := map[string]testCase{
		"good query": {
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"exists"}).AddRow(f.exists)
				f.mock.
					ExpectQuery(regexp.QuoteMeta(chatExistsCmd)).
					WithArgs(f.user1ID, f.user2ID).
					WillReturnRows(rows)
			},
			user1ID: 2,
			user2ID: 3,
			exists:  true,
			err:     nil,
		},
		"query error": {
			prepare: func(f *fields) {
				f.mock.
					ExpectQuery(regexp.QuoteMeta(chatExistsCmd)).
					WithArgs(f.user1ID, f.user2ID).
					WillReturnError(&pq.Error{Message: "sql error"})
			},
			user1ID: 2,
			user2ID: 3,
			exists:  false,
			err:     pkgErrors.ErrDb,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			logger, err := zap.NewDevelopment()
			if err != nil {
				t.Fatalf("can't create logger: %s", err)
			}

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("can't create mock: %s", err)
			}
			defer db.Close()

			f := fields{mock: mock, user1ID: test.user1ID, user2ID: test.user2ID, exists: test.exists}
			if test.prepare != nil {
				test.prepare(&f)
			}

			repo := NewRepository(db, logger)
			exists, err := repo.ChatExists(test.user1ID, test.user2ID)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if exists != test.exists {
				t.Errorf("\nExpected: %v\nGot: %v", test.exists, exists)
			}
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("\nThere were unfulfilled expectations: %s", err)
			}
		})
	}
}
