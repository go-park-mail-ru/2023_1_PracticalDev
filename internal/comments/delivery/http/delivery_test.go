package http

import (
	pkgComments "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/comments"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/comments/mocks"
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

func TestDelivery_Create(t *testing.T) {
	type fields struct {
		serv *mocks.MockService
	}

	type testCase struct {
		prepare  func(f *fields)
		params   httprouter.Params
		request  string
		response string
		err      error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				f.serv.EXPECT().Create(&pkgComments.CreateParams{
					AuthorID: 27,
					PinID:    21,
					Text:     "Good pin!",
				}).Return(models.Comment{
					ID:       2,
					AuthorID: 27,
					PinID:    21,
					Text:     "Good pin!",
				}, nil)
			},
			params: []httprouter.Param{
				{Key: "user-id", Value: "27"},
				{Key: "id", Value: "21"},
			},
			request:  `{"text":"Good pin!"}`,
			response: `{"id":2,"author_id":27,"pin_id":21,"text":"Good pin!","created_at":"0001-01-01T00:00:00Z"}`,
			err:      nil,
		},
		"invalid user id param": {
			prepare: func(f *fields) {},
			params: []httprouter.Param{
				{Key: "user-id", Value: "a"},
				{Key: "id", Value: "21"},
			},
			request:  `{"text":"Good pin!"}`,
			response: ``,
			err:      pkgErrors.ErrInvalidUserIdParam,
		},
		"missing user id param": {
			prepare:  func(f *fields) {},
			params:   []httprouter.Param{{Key: "id", Value: "21"}},
			request:  `{"text":"Good pin!"}`,
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

			del := delivery{serv: f.serv, log: logger}

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.request))
			rec := httptest.NewRecorder()
			err = del.Create(rec, req, test.params)
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

func TestDelivery_List(t *testing.T) {
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
				f.serv.EXPECT().List(21).Return([]models.Comment{
					{ID: 2, AuthorID: 27, PinID: 21, Text: "Good pin!"},
					{ID: 3, AuthorID: 28, PinID: 21, Text: "Yeah!"},
					{ID: 4, AuthorID: 27, PinID: 21, Text: "Fantastic!"},
				}, nil)
			},
			params:   []httprouter.Param{{Key: "id", Value: "21"}},
			response: `{"items":[{"id":2,"author_id":27,"pin_id":21,"text":"Good pin!","created_at":"0001-01-01T00:00:00Z"},{"id":3,"author_id":28,"pin_id":21,"text":"Yeah!","created_at":"0001-01-01T00:00:00Z"},{"id":4,"author_id":27,"pin_id":21,"text":"Fantastic!","created_at":"0001-01-01T00:00:00Z"}]}`,
			err:      nil,
		},
		"empty result": {
			prepare: func(f *fields) {
				f.serv.EXPECT().List(21).Return([]models.Comment{}, nil)
			},
			params:   []httprouter.Param{{Key: "id", Value: "21"}},
			response: `{"items":[]}`,
			err:      nil,
		},
		"invalid pin id param": {
			prepare:  nil,
			params:   []httprouter.Param{{Key: "id", Value: "a"}},
			response: ``,
			err:      pkgErrors.ErrInvalidPinIdParam,
		},
		"missing pin id param": {
			prepare:  nil,
			params:   []httprouter.Param{},
			response: ``,
			err:      pkgErrors.ErrInvalidPinIdParam,
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
			err = del.List(rec, req, test.params)
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
