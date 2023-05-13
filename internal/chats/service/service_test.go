package service

import (
	pkgChats "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/chats"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/chats/mocks"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"reflect"
	"testing"
)

func TestService_Create(t *testing.T) {
	type fields struct {
		repo   *mocks.MockRepository
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
		"usual": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Create(f.params).Return(*f.chat, nil)
			},
			params: &pkgChats.CreateParams{User1ID: 2, User2ID: 3},
			chat:   models.Chat{ID: 1, User1ID: 2, User2ID: 3},
			err:    nil,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl), params: test.params, chat: &test.chat}
			if test.prepare != nil {
				test.prepare(&f)
			}

			serv := NewService(f.repo)
			chat, err := serv.Create(test.params)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if chat != test.chat {
				t.Errorf("\nExpected: %v\nGot: %v", test.chat, chat)
			}
		})
	}
}

func TestService_ListByUser(t *testing.T) {
	type fields struct {
		repo   *mocks.MockRepository
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
				f.repo.EXPECT().ListByUser(f.userID).Return(f.chats, nil)
			},
			userID: 2,
			chats: []models.Chat{
				{ID: 2, User1ID: 2, User2ID: 3},
				{ID: 3, User1ID: 8, User2ID: 2},
				{ID: 4, User1ID: 2, User2ID: 4},
			},
			err: nil,
		},
		"no chats": {
			prepare: func(f *fields) {
				f.repo.EXPECT().ListByUser(f.userID).Return(f.chats, nil)
			},
			userID: 2,
			chats:  []models.Chat{},
			err:    nil,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl), userID: test.userID, chats: test.chats}
			if test.prepare != nil {
				test.prepare(&f)
			}

			serv := NewService(f.repo)
			chats, err := serv.ListByUser(test.userID)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if !reflect.DeepEqual(chats, test.chats) {
				t.Errorf("\nExpected: %v\nGot: %v", test.chats, chats)
			}
		})
	}
}

func TestService_Get(t *testing.T) {
	type fields struct {
		repo   *mocks.MockRepository
		chatID int
		chat   models.Chat
	}

	type testCase struct {
		prepare func(f *fields)
		chatID  int
		chat    models.Chat
		err     error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Get(f.chatID).Return(f.chat, nil)
			},
			chatID: 2,
			chat:   models.Chat{ID: 2, User1ID: 2, User2ID: 3},
			err:    nil,
		},
		"chat not found": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Get(f.chatID).Return(f.chat, pkgErrors.ErrChatNotFound)
			},
			chatID: 3,
			chat:   models.Chat{},
			err:    pkgErrors.ErrChatNotFound,
		},
		"negative chat id param": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Get(f.chatID).Return(f.chat, pkgErrors.ErrChatNotFound)
			},
			chatID: -1,
			chat:   models.Chat{},
			err:    pkgErrors.ErrChatNotFound,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl), chatID: test.chatID, chat: test.chat}
			if test.prepare != nil {
				test.prepare(&f)
			}

			serv := NewService(f.repo)
			board, err := serv.Get(test.chatID)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if board != test.chat {
				t.Errorf("\nExpected: %v\nGot: %v", test.chat, board)
			}
		})
	}
}

func TestService_MessagesList(t *testing.T) {
	type fields struct {
		repo     *mocks.MockRepository
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
				f.repo.EXPECT().MessagesList(f.chatID).Return(f.messages, nil)
			},
			chatID: 2,
			messages: []models.Message{
				{ID: 1, AuthorID: 3, ChatID: 2, Text: "msg 1"},
				{ID: 2, AuthorID: 4, ChatID: 2, Text: "msg 2"},
				{ID: 3, AuthorID: 3, ChatID: 2, Text: "msg 3"},
			},
			err: nil,
		},
		"no messages": {
			prepare: func(f *fields) {
				f.repo.EXPECT().MessagesList(f.chatID).Return(f.messages, nil)
			},
			chatID:   2,
			messages: []models.Message{},
			err:      nil,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl), chatID: test.chatID, messages: test.messages}
			if test.prepare != nil {
				test.prepare(&f)
			}

			serv := NewService(f.repo)
			messages, err := serv.MessagesList(test.chatID)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if !reflect.DeepEqual(messages, test.messages) {
				t.Errorf("\nExpected: %v\nGot: %v", test.messages, messages)
			}
		})
	}
}

func TestService_CreateMessage(t *testing.T) {
	type fields struct {
		repo    *mocks.MockRepository
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
		"usual": {
			prepare: func(f *fields) {
				f.repo.EXPECT().CreateMessage(f.params).Return(f.message, nil)
			},
			params:  &pkgChats.CreateMessageParams{AuthorID: 2, ChatID: 3, Text: "hello!"},
			message: models.Message{ID: 1, AuthorID: 2, ChatID: 3, Text: "hello!"},
			err:     nil,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl), params: test.params, message: &test.message}
			if test.prepare != nil {
				test.prepare(&f)
			}

			serv := NewService(f.repo)
			message, err := serv.CreateMessage(test.params)
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
		})
	}
}

func TestService_ChatExists(t *testing.T) {
	type fields struct {
		repo    *mocks.MockRepository
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
		"chat exists": {
			prepare: func(f *fields) {
				f.repo.EXPECT().ChatExists(f.user1ID, f.user2ID).Return(f.exists, nil)
			},
			user1ID: 2,
			user2ID: 3,
			exists:  true,
			err:     nil,
		},
		"chat not exists": {
			prepare: func(f *fields) {
				f.repo.EXPECT().ChatExists(f.user1ID, f.user2ID).Return(f.exists, nil)
			},
			user1ID: 2,
			user2ID: 3,
			exists:  false,
			err:     nil,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl), user1ID: test.user1ID, user2ID: test.user2ID, exists: test.exists}
			if test.prepare != nil {
				test.prepare(&f)
			}

			serv := NewService(f.repo)
			exists, err := serv.ChatExists(test.user1ID, test.user2ID)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if exists != test.exists {
				t.Errorf("\nExpected: %v\nGot: %v", test.exists, exists)
			}
		})
	}
}
