package http

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/chats/mocks"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDelivery_ListByUser(t *testing.T) {
	type fields struct {
		serv *mocks.MockService
	}

	type testCase struct {
		prepare  func(f *fields)
		params   httprouter.Params
		response string
		err      error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				f.serv.EXPECT().ListByUser(2).Return([]models.Chat{
					{ID: 2, User1ID: 2, User2ID: 3},
					{ID: 3, User1ID: 8, User2ID: 2},
					{ID: 4, User1ID: 2, User2ID: 4},
				}, nil)
			},
			params:   []httprouter.Param{{Key: "user-id", Value: "2"}},
			response: `{"items":[{"id":2,"user1_id":2,"user2_id":3,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"},{"id":3,"user1_id":8,"user2_id":2,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"},{"id":4,"user1_id":2,"user2_id":4,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}]}`,
			err:      nil,
		},
		"no chats": {
			prepare: func(f *fields) {
				f.serv.EXPECT().ListByUser(3).Return([]models.Chat{}, nil)
			},
			params:   []httprouter.Param{{Key: "user-id", Value: "3"}},
			response: `{"items":[]}`,
			err:      nil,
		},
		"invalid user id param": {
			prepare:  nil,
			params:   []httprouter.Param{{Key: "user-id", Value: "a"}},
			response: ``,
			err:      pkgErrors.ErrInvalidUserIdParam,
		},
		"missing user id param": {
			prepare:  nil,
			params:   []httprouter.Param{},
			response: ``,
			err:      pkgErrors.ErrInvalidUserIdParam,
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

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{serv: mocks.NewMockService(ctrl)}
			if test.prepare != nil {
				test.prepare(&f)
			}

			del := delivery{
				serv: f.serv,
				log:  logger,
			}

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			err = del.ListByUser(rec, req, test.params)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			body, _ := io.ReadAll(rec.Body)
			if strings.Trim(string(body), "\n") != test.response {
				t.Errorf("\nExpected: %s\nGot: %s", test.response, string(body))
			}
		})
	}
}

func TestDelivery_Get(t *testing.T) {
	type fields struct {
		serv *mocks.MockService
	}

	type testCase struct {
		prepare  func(f *fields)
		params   httprouter.Params
		response string
		err      error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				f.serv.EXPECT().Get(2).Return(models.Chat{ID: 2, User1ID: 2, User2ID: 3}, nil)
			},
			params:   []httprouter.Param{{Key: "id", Value: "2"}},
			response: `{"id":2,"user1_id":2,"user2_id":3,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}`,
			err:      nil,
		},
		"invalid chat id param": {
			prepare:  nil,
			params:   []httprouter.Param{{Key: "id", Value: "a"}},
			response: ``,
			err:      pkgErrors.ErrInvalidChatIDParam,
		},
		"missing chat id param": {
			prepare:  nil,
			params:   []httprouter.Param{},
			response: ``,
			err:      pkgErrors.ErrInvalidChatIDParam,
		},
		"chat not found": {
			prepare: func(f *fields) {
				f.serv.EXPECT().Get(3).Return(models.Chat{}, pkgErrors.ErrChatNotFound)
			},
			params:   []httprouter.Param{{Key: "id", Value: "3"}},
			response: ``,
			err:      pkgErrors.ErrChatNotFound,
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

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{serv: mocks.NewMockService(ctrl)}
			if test.prepare != nil {
				test.prepare(&f)
			}

			del := delivery{
				serv: f.serv,
				log:  logger,
			}

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			err = del.Get(rec, req, test.params)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			body, _ := io.ReadAll(rec.Body)
			if strings.Trim(string(body), "\n") != test.response {
				t.Errorf("\nExpected: %s\nGot: %s", test.response, string(body))
			}
		})
	}
}

func TestDelivery_MessagesList(t *testing.T) {
	type fields struct {
		serv *mocks.MockService
	}

	type testCase struct {
		prepare  func(f *fields)
		params   httprouter.Params
		response string
		err      error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				f.serv.EXPECT().MessagesList(2).Return([]models.Message{
					{ID: 1, AuthorID: 3, ChatID: 2, Text: "msg 1"},
					{ID: 2, AuthorID: 4, ChatID: 2, Text: "msg 2"},
					{ID: 3, AuthorID: 3, ChatID: 2, Text: "msg 3"},
				}, nil)
			},
			params:   []httprouter.Param{{Key: "id", Value: "2"}},
			response: `{"items":[{"id":1,"author_id":3,"chat_id":2,"text":"msg 1","created_at":"0001-01-01T00:00:00Z"},{"id":2,"author_id":4,"chat_id":2,"text":"msg 2","created_at":"0001-01-01T00:00:00Z"},{"id":3,"author_id":3,"chat_id":2,"text":"msg 3","created_at":"0001-01-01T00:00:00Z"}]}`,
			err:      nil,
		},
		"no messages": {
			prepare: func(f *fields) {
				f.serv.EXPECT().MessagesList(3).Return([]models.Message{}, nil)
			},
			params:   []httprouter.Param{{Key: "id", Value: "3"}},
			response: `{"items":[]}`,
			err:      nil,
		},
		"invalid chat id param": {
			prepare:  nil,
			params:   []httprouter.Param{{Key: "id", Value: "a"}},
			response: ``,
			err:      pkgErrors.ErrInvalidChatIDParam,
		},
		"missing chat id param": {
			prepare:  nil,
			params:   []httprouter.Param{},
			response: ``,
			err:      pkgErrors.ErrInvalidChatIDParam,
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

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{serv: mocks.NewMockService(ctrl)}
			if test.prepare != nil {
				test.prepare(&f)
			}

			del := delivery{
				serv: f.serv,
				log:  logger,
			}

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			err = del.MessagesList(rec, req, test.params)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			body, _ := io.ReadAll(rec.Body)
			if strings.Trim(string(body), "\n") != test.response {
				t.Errorf("\nExpected: %s\nGot: %s", test.response, string(body))
			}
		})
	}
}
