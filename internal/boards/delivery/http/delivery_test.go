package http

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/boards"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/boards/mocks"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
)

func TestCreate(t *testing.T) {
	type fields struct {
		serv *mocks.MockService
	}

	type testCase struct {
		prepare    func(f *fields)
		params     httprouter.Params
		request    string
		response   string
		statusCode int
		err        error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				f.serv.EXPECT().Create(&boards.CreateParams{
					Name:        "b1",
					Description: "d1",
					Privacy:     "secret",
					UserId:      3,
				}).Return(models.Board{
					Id:          1,
					Name:        "b1",
					Description: "d1",
					Privacy:     "secret",
					UserId:      3}, nil)
			},
			params:     []httprouter.Param{{Key: "user-id", Value: "3"}},
			request:    `{"name":"b1","description":"d1","privacy":"secret"}`,
			response:   `{"id":1,"name":"b1","description":"d1","privacy":"secret","user_id":3}`,
			statusCode: http.StatusOK,
			err:        nil,
		},
		"invalid user id param": {
			prepare:    func(f *fields) {},
			params:     []httprouter.Param{{Key: "user-id", Value: "a"}},
			response:   ``,
			statusCode: http.StatusBadRequest,
			err:        ErrInvalidUserId,
		},
		"missing user id param": {
			prepare:    func(f *fields) {},
			params:     []httprouter.Param{},
			response:   ``,
			statusCode: http.StatusBadRequest,
			err:        ErrInvalidUserId,
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

			logger := log.New()
			del := delivery{
				serv: f.serv,
				log:  logger,
			}

			req := httptest.NewRequest(http.MethodPost, "/boards", strings.NewReader(test.request))
			rec := httptest.NewRecorder()
			err := del.create(rec, req, test.params)
			if err != test.err {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if rec.Code != test.statusCode {
				t.Errorf("\nExpected: %d\nGot: %d", test.statusCode, rec.Code)
			}
			body, _ := io.ReadAll(rec.Body)
			if strings.Trim(string(body), "\n") != test.response {
				t.Errorf("\nExpected: %s\nGot: %s", test.response, string(body))
			}
		})
	}
}

func TestList(t *testing.T) {
	type fields struct {
		serv *mocks.MockService
	}

	type testCase struct {
		prepare    func(f *fields)
		params     httprouter.Params
		response   string
		statusCode int
		err        error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				f.serv.EXPECT().List(3).Return([]models.Board{
					{Id: 1, Name: "b1", Description: "d1", Privacy: "secret", UserId: 3},
					{Id: 2, Name: "b2", Description: "d2", Privacy: "secret", UserId: 3},
					{Id: 5, Name: "b5", Description: "d5", Privacy: "public", UserId: 3},
				}, nil)
			},
			params:     []httprouter.Param{{Key: "user-id", Value: "3"}},
			response:   `{"boards":[{"id":1,"name":"b1","description":"d1","privacy":"secret","user_id":3},{"id":2,"name":"b2","description":"d2","privacy":"secret","user_id":3},{"id":5,"name":"b5","description":"d5","privacy":"public","user_id":3}]}`,
			statusCode: http.StatusOK,
			err:        nil,
		},
		"no boards": {
			prepare: func(f *fields) {
				f.serv.EXPECT().List(3).Return([]models.Board{}, nil)
			},
			params:     []httprouter.Param{{Key: "user-id", Value: "3"}},
			response:   `{"boards":[]}`,
			statusCode: http.StatusOK,
			err:        nil,
		},
		"invalid user id param": {
			prepare:    func(f *fields) {},
			params:     []httprouter.Param{{Key: "user-id", Value: "a"}},
			response:   ``,
			statusCode: http.StatusBadRequest,
			err:        ErrInvalidUserId,
		},
		"missing user id param": {
			prepare:    func(f *fields) {},
			params:     []httprouter.Param{},
			response:   ``,
			statusCode: http.StatusBadRequest,
			err:        ErrInvalidUserId,
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

			logger := log.New()
			del := delivery{
				serv: f.serv,
				log:  logger,
			}

			req := httptest.NewRequest(http.MethodGet, "/boards", nil)
			rec := httptest.NewRecorder()
			err := del.list(rec, req, test.params)
			if err != test.err {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if rec.Code != test.statusCode {
				t.Errorf("\nExpected: %d\nGot: %d", test.statusCode, rec.Code)
			}
			body, _ := io.ReadAll(rec.Body)
			if strings.Trim(string(body), "\n") != test.response {
				t.Errorf("\nExpected: %s\nGot: %s", test.response, string(body))
			}
		})
	}
}

func TestGet(t *testing.T) {
	type fields struct {
		serv *mocks.MockService
	}

	type testCase struct {
		prepare    func(f *fields)
		params     httprouter.Params
		response   string
		statusCode int
		err        error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				f.serv.EXPECT().Get(3).Return(models.Board{
					Id:          3,
					Name:        "n3",
					Description: "d3",
					Privacy:     "secret",
					UserId:      1,
				}, nil)
			},
			params:     []httprouter.Param{{Key: "board_id", Value: "3"}},
			response:   `{"id":3,"name":"n3","description":"d3","privacy":"secret","user_id":1}`,
			statusCode: http.StatusOK,
			err:        nil,
		},
		"invalid board id param": {
			prepare:    func(f *fields) {},
			params:     []httprouter.Param{{Key: "board_id", Value: "a"}},
			response:   ``,
			statusCode: http.StatusBadRequest,
			err:        ErrInvalidBoardIdParam,
		},
		"missing board id param": {
			prepare:    func(f *fields) {},
			params:     []httprouter.Param{},
			response:   ``,
			statusCode: http.StatusBadRequest,
			err:        ErrInvalidBoardIdParam,
		},
		"board not found": {
			prepare: func(f *fields) {
				f.serv.EXPECT().Get(3).Return(models.Board{}, boards.ErrBoardNotFound)
			},
			params:     []httprouter.Param{{Key: "board_id", Value: "3"}},
			response:   ``,
			statusCode: http.StatusNotFound,
			err:        boards.ErrBoardNotFound,
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

			logger := log.New()
			del := delivery{
				serv: f.serv,
				log:  logger,
			}

			req := httptest.NewRequest(http.MethodGet, "/boards/3", nil)
			rec := httptest.NewRecorder()
			err := del.get(rec, req, test.params)
			if err != test.err {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}

			if rec.Code != test.statusCode {
				t.Errorf("\nExpected: %d\nGot: %d", test.statusCode, rec.Code)
			}

			body, _ := io.ReadAll(rec.Body)
			if strings.Trim(string(body), "\n") != test.response {
				t.Errorf("\nExpected: %s\nGot: %s", test.response, string(body))
			}
		})
	}
}

func TestFullUpdate(t *testing.T) {
	type fields struct {
		serv *mocks.MockService
	}

	type testCase struct {
		prepare    func(f *fields)
		params     httprouter.Params
		request    string
		response   string
		statusCode int
		err        error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				f.serv.EXPECT().FullUpdate(&boards.FullUpdateParams{
					Id:          1,
					Name:        "b1",
					Description: "d1",
					Privacy:     "secret",
				}).Return(models.Board{
					Id:          1,
					Name:        "b1",
					Description: "d1",
					Privacy:     "secret",
					UserId:      3}, nil)
			},
			params:     []httprouter.Param{{Key: "board_id", Value: "1"}},
			request:    `{"name":"b1","description":"d1","privacy":"secret"}`,
			response:   `{"id":1,"name":"b1","description":"d1","privacy":"secret","user_id":3}`,
			statusCode: http.StatusOK,
			err:        nil,
		},
		"invalid board id param": {
			prepare:    func(f *fields) {},
			params:     []httprouter.Param{{Key: "board_id", Value: "a"}},
			response:   ``,
			statusCode: http.StatusBadRequest,
			err:        ErrInvalidBoardIdParam,
		},
		"missing board id param": {
			prepare:    func(f *fields) {},
			params:     []httprouter.Param{},
			response:   ``,
			statusCode: http.StatusBadRequest,
			err:        ErrInvalidBoardIdParam,
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

			logger := log.New()
			del := delivery{
				serv: f.serv,
				log:  logger,
			}

			req := httptest.NewRequest(http.MethodPut, "/boards/1", strings.NewReader(test.request))
			rec := httptest.NewRecorder()
			err := del.fullUpdate(rec, req, test.params)
			if err != test.err {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if rec.Code != test.statusCode {
				t.Errorf("\nExpected: %d\nGot: %d", test.statusCode, rec.Code)
			}
			body, _ := io.ReadAll(rec.Body)
			if strings.Trim(string(body), "\n") != test.response {
				t.Errorf("\nExpected: %s\nGot: %s", test.response, string(body))
			}
		})
	}
}

func TestPartialUpdate(t *testing.T) {
	type fields struct {
		serv *mocks.MockService
	}

	type testCase struct {
		prepare    func(f *fields)
		params     httprouter.Params
		request    string
		response   string
		statusCode int
		err        error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				f.serv.EXPECT().PartialUpdate(&boards.PartialUpdateParams{
					Id:                1,
					Name:              "b1",
					UpdateName:        true,
					Description:       "d1",
					UpdateDescription: true,
					Privacy:           "secret",
					UpdatePrivacy:     true,
				}).Return(models.Board{
					Id:          1,
					Name:        "b1",
					Description: "d1",
					Privacy:     "secret",
					UserId:      3}, nil)
			},
			params:     []httprouter.Param{{Key: "board_id", Value: "1"}},
			request:    `{"name":"b1","description":"d1","privacy":"secret"}`,
			response:   `{"id":1,"name":"b1","description":"d1","privacy":"secret","user_id":3}`,
			statusCode: http.StatusOK,
			err:        nil,
		},
		"invalid board id param": {
			prepare:    func(f *fields) {},
			params:     []httprouter.Param{{Key: "board_id", Value: "a"}},
			response:   ``,
			statusCode: http.StatusBadRequest,
			err:        ErrInvalidBoardIdParam,
		},
		"missing board id param": {
			prepare:    func(f *fields) {},
			params:     []httprouter.Param{},
			response:   ``,
			statusCode: http.StatusBadRequest,
			err:        ErrInvalidBoardIdParam,
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

			logger := log.New()
			del := delivery{
				serv: f.serv,
				log:  logger,
			}

			req := httptest.NewRequest(http.MethodPatch, "/boards/1", strings.NewReader(test.request))
			rec := httptest.NewRecorder()
			err := del.partialUpdate(rec, req, test.params)
			if err != test.err {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if rec.Code != test.statusCode {
				t.Errorf("\nExpected: %d\nGot: %d", test.statusCode, rec.Code)
			}
			body, _ := io.ReadAll(rec.Body)
			if strings.Trim(string(body), "\n") != test.response {
				t.Errorf("\nExpected: %s\nGot: %s", test.response, string(body))
			}
		})
	}
}

func TestDelete(t *testing.T) {
	type fields struct {
		serv *mocks.MockService
	}

	type testCase struct {
		prepare    func(f *fields)
		params     httprouter.Params
		statusCode int
		err        error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				f.serv.EXPECT().Delete(3).Return(nil)
			},
			params:     []httprouter.Param{{Key: "board_id", Value: "3"}},
			statusCode: http.StatusOK,
			err:        nil,
		},
		"invalid board id param": {
			prepare:    func(f *fields) {},
			params:     []httprouter.Param{{Key: "board_id", Value: "a"}},
			statusCode: http.StatusBadRequest,
			err:        ErrInvalidBoardIdParam,
		},
		"missing board id param": {
			prepare:    func(f *fields) {},
			params:     []httprouter.Param{},
			statusCode: http.StatusBadRequest,
			err:        ErrInvalidBoardIdParam,
		},
		"board not found": {
			prepare: func(f *fields) {
				f.serv.EXPECT().Delete(3).Return(boards.ErrBoardNotFound)
			},
			params:     []httprouter.Param{{Key: "board_id", Value: "3"}},
			statusCode: http.StatusNotFound,
			err:        boards.ErrBoardNotFound,
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

			logger := log.New()
			del := delivery{
				serv: f.serv,
				log:  logger,
			}

			req := httptest.NewRequest(http.MethodDelete, "/boards/3", nil)
			rec := httptest.NewRecorder()
			err := del.delete(rec, req, test.params)
			if err != test.err {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if rec.Code != test.statusCode {
				t.Errorf("\nExpected: %d\nGot: %d", test.statusCode, rec.Code)
			}
		})
	}
}
