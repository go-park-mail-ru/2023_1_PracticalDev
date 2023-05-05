package http

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/likes/mocks"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
)

var err error
var logger *zap.Logger

func init() {
	logger, err = zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
}

func TestLike(t *testing.T) {
	type fields struct {
		serv *mocks.MockService
	}

	type testCase struct {
		prepare func(f *fields)
		params  httprouter.Params
		err     error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				f.serv.EXPECT().Like(3, 2).Return(nil)
			},
			params: []httprouter.Param{
				{Key: "id", Value: "3"},
				{Key: "user-id", Value: "2"},
			},
			err: pkgErrors.ErrNoContent,
		},
		"invalid user id param": {
			prepare: func(f *fields) {},
			params: []httprouter.Param{
				{Key: "id", Value: "3"},
				{Key: "user-id", Value: "a"},
			},
			err: pkgErrors.ErrInvalidUserIdParam,
		},
		"missing user id param": {
			prepare: func(f *fields) {},
			params:  []httprouter.Param{{Key: "id", Value: "3"}},
			err:     pkgErrors.ErrInvalidUserIdParam,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{serv: mocks.NewMockService(ctrl)}
			if test.prepare != nil {
				test.prepare(&f)
			}

			del := delivery{serv: f.serv, log: logger}
			req := httptest.NewRequest(http.MethodPost, "/pins/3/like", nil)
			rec := httptest.NewRecorder()
			err := del.like(rec, req, test.params)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
		})
	}
}

func TestListByAuthor(t *testing.T) {
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
				f.serv.EXPECT().ListByAuthor(12).Return([]models.Like{
					{PinId: 2, AuthorId: 12, CreatedAt: time.Unix(1681163314, 0)},
					{PinId: 3, AuthorId: 12, CreatedAt: time.Unix(1681163000, 0)},
					{PinId: 4, AuthorId: 12, CreatedAt: time.Unix(1681164555, 0)},
				}, nil)
			},
			params:   []httprouter.Param{{Key: "id", Value: "12"}},
			response: `{"likes":[{"pin_id":2,"author_id":12,"created_at":"2023-04-10T21:48:34Z"},{"pin_id":3,"author_id":12,"created_at":"2023-04-10T21:43:20Z"},{"pin_id":4,"author_id":12,"created_at":"2023-04-10T22:09:15Z"}]}`,
			err:      nil,
		},
		"no likes": {
			prepare: func(f *fields) {
				f.serv.EXPECT().ListByAuthor(12).Return([]models.Like{}, nil)
			},
			params:   []httprouter.Param{{Key: "id", Value: "12"}},
			response: `{"likes":[]}`,
			err:      nil,
		},
		"invalid user id param": {
			prepare:  func(f *fields) {},
			params:   []httprouter.Param{{Key: "id", Value: "a"}},
			response: ``,
			err:      pkgErrors.ErrInvalidUserIdParam,
		},
		"missing user id param": {
			prepare:  func(f *fields) {},
			params:   []httprouter.Param{},
			response: ``,
			err:      pkgErrors.ErrInvalidUserIdParam,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{serv: mocks.NewMockService(ctrl)}
			if test.prepare != nil {
				test.prepare(&f)
			}

			del := delivery{serv: f.serv, log: logger}
			req := httptest.NewRequest(http.MethodGet, "/users/3/likes", nil)
			rec := httptest.NewRecorder()
			err := del.listByAuthor(rec, req, test.params)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			body, _ := io.ReadAll(rec.Body)
			if strings.Trim(string(body), "\n") != test.response {
				t.Errorf("\nExpected:\n%s\nGot:\n%s", test.response, string(body))
			}
		})
	}
}

func TestListByPin(t *testing.T) {
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
				f.serv.EXPECT().ListByPin(12).Return([]models.Like{
					{PinId: 12, AuthorId: 4, CreatedAt: time.Unix(1681163314, 0)},
					{PinId: 12, AuthorId: 23, CreatedAt: time.Unix(1681163000, 0)},
					{PinId: 12, AuthorId: 2, CreatedAt: time.Unix(1681164555, 0)},
				}, nil)
			},
			params:   []httprouter.Param{{Key: "id", Value: "12"}},
			response: `{"likes":[{"pin_id":12,"author_id":4,"created_at":"2023-04-10T21:48:34Z"},{"pin_id":12,"author_id":23,"created_at":"2023-04-10T21:43:20Z"},{"pin_id":12,"author_id":2,"created_at":"2023-04-10T22:09:15Z"}]}`,
			err:      nil,
		},
		"no likes": {
			prepare: func(f *fields) {
				f.serv.EXPECT().ListByPin(12).Return([]models.Like{}, nil)
			},
			params:   []httprouter.Param{{Key: "id", Value: "12"}},
			response: `{"likes":[]}`,
			err:      nil,
		},
		"invalid user id param": {
			prepare:  func(f *fields) {},
			params:   []httprouter.Param{{Key: "id", Value: "a"}},
			response: ``,
			err:      pkgErrors.ErrInvalidPinIdParam,
		},
		"missing user id param": {
			prepare:  func(f *fields) {},
			params:   []httprouter.Param{},
			response: ``,
			err:      pkgErrors.ErrInvalidPinIdParam,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{serv: mocks.NewMockService(ctrl)}
			if test.prepare != nil {
				test.prepare(&f)
			}

			del := delivery{serv: f.serv, log: logger}
			req := httptest.NewRequest(http.MethodGet, "/pins/3/likes", nil)
			rec := httptest.NewRecorder()
			err := del.listByPin(rec, req, test.params)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			body, _ := io.ReadAll(rec.Body)
			if strings.Trim(string(body), "\n") != test.response {
				t.Errorf("\nExpected:\n%s\nGot:\n%s", test.response, string(body))
			}
		})
	}
}
func TestUnlike(t *testing.T) {
	type fields struct {
		serv *mocks.MockService
	}

	type testCase struct {
		prepare func(f *fields)
		params  httprouter.Params
		err     error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				f.serv.EXPECT().Unlike(3, 2).Return(nil)
			},
			params: []httprouter.Param{
				{Key: "id", Value: "3"},
				{Key: "user-id", Value: "2"},
			},
			err: pkgErrors.ErrNoContent,
		},
		"invalid user id param": {
			prepare: func(f *fields) {},
			params: []httprouter.Param{
				{Key: "id", Value: "3"},
				{Key: "user-id", Value: "a"},
			},
			err: pkgErrors.ErrInvalidUserIdParam,
		},
		"missing user id param": {
			prepare: func(f *fields) {},
			params:  []httprouter.Param{{Key: "id", Value: "3"}},
			err:     pkgErrors.ErrInvalidUserIdParam,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{serv: mocks.NewMockService(ctrl)}
			if test.prepare != nil {
				test.prepare(&f)
			}

			del := delivery{serv: f.serv, log: logger}
			req := httptest.NewRequest(http.MethodDelete, "/pins/3/like", nil)
			rec := httptest.NewRecorder()
			err := del.unlike(rec, req, test.params)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected:\n%s\nGot:\n%s", test.err, err)
			}
		})
	}
}
