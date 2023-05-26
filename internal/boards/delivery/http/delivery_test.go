package http

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	_boards "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/boards"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/boards/mocks"
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

func TestCreate(t *testing.T) {
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
				f.serv.EXPECT().Create(&_boards.CreateParams{
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
			params:   []httprouter.Param{{Key: "user-id", Value: "3"}},
			request:  `{"name":"b1","description":"d1","privacy":"secret"}`,
			response: `{"id":1,"name":"b1","description":"d1","privacy":"secret","user_id":3}`,
			err:      nil,
		},
		"invalid user id param": {
			prepare:  func(f *fields) {},
			params:   []httprouter.Param{{Key: "user-id", Value: "a"}},
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

			req := httptest.NewRequest(http.MethodPost, "/boards", strings.NewReader(test.request))
			rec := httptest.NewRecorder()
			err := del.create(rec, req, test.params)
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

func TestList(t *testing.T) {
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
				f.serv.EXPECT().List(3).Return([]models.Board{
					{Id: 1, Name: "b1", Description: "d1", Privacy: "secret", UserId: 3},
					{Id: 2, Name: "b2", Description: "d2", Privacy: "secret", UserId: 3},
					{Id: 5, Name: "b5", Description: "d5", Privacy: "public", UserId: 3},
				}, nil)
			},
			params:   []httprouter.Param{{Key: "user-id", Value: "3"}},
			response: `{"boards":[{"id":1,"name":"b1","description":"d1","privacy":"secret","user_id":3},{"id":2,"name":"b2","description":"d2","privacy":"secret","user_id":3},{"id":5,"name":"b5","description":"d5","privacy":"public","user_id":3}]}`,
			err:      nil,
		},
		"no boards": {
			prepare: func(f *fields) {
				f.serv.EXPECT().List(3).Return([]models.Board{}, nil)
			},
			params:   []httprouter.Param{{Key: "user-id", Value: "3"}},
			response: `{"boards":[]}`,
			err:      nil,
		},
		"invalid user id param": {
			prepare:  func(f *fields) {},
			params:   []httprouter.Param{{Key: "user-id", Value: "a"}},
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

			del := delivery{
				serv: f.serv,
				log:  logger,
			}

			req := httptest.NewRequest(http.MethodGet, "/boards", nil)
			rec := httptest.NewRecorder()
			err := del.list(rec, req, test.params)
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

func TestGet(t *testing.T) {
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
				f.serv.EXPECT().Get(3).Return(models.Board{
					Id:          3,
					Name:        "n3",
					Description: "d3",
					Privacy:     "secret",
					UserId:      1,
				}, nil)
			},
			params:   []httprouter.Param{{Key: "id", Value: "3"}},
			response: `{"id":3,"name":"n3","description":"d3","privacy":"secret","user_id":1}`,
			err:      nil,
		},
		"invalid board id param": {
			prepare:  func(f *fields) {},
			params:   []httprouter.Param{{Key: "id", Value: "a"}},
			response: ``,
			err:      pkgErrors.ErrInvalidBoardIdParam,
		},
		"missing board id param": {
			prepare:  func(f *fields) {},
			params:   []httprouter.Param{},
			response: ``,
			err:      pkgErrors.ErrInvalidBoardIdParam,
		},
		"board not found": {
			prepare: func(f *fields) {
				f.serv.EXPECT().Get(3).Return(models.Board{}, pkgErrors.ErrBoardNotFound)
			},
			params:   []httprouter.Param{{Key: "id", Value: "3"}},
			response: ``,
			err:      pkgErrors.ErrBoardNotFound,
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

			del := delivery{
				serv: f.serv,
				log:  logger,
			}

			req := httptest.NewRequest(http.MethodGet, "/boards/3", nil)
			rec := httptest.NewRecorder()
			err := del.get(rec, req, test.params)
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

func TestFullUpdate(t *testing.T) {
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
				f.serv.EXPECT().FullUpdate(&_boards.FullUpdateParams{
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
			params:   []httprouter.Param{{Key: "id", Value: "1"}},
			request:  `{"name":"b1","description":"d1","privacy":"secret"}`,
			response: `{"id":1,"name":"b1","description":"d1","privacy":"secret","user_id":3}`,
			err:      nil,
		},
		"invalid board id param": {
			prepare:  func(f *fields) {},
			params:   []httprouter.Param{{Key: "id", Value: "a"}},
			response: ``,
			err:      pkgErrors.ErrInvalidBoardIdParam,
		},
		"missing board id param": {
			prepare:  func(f *fields) {},
			params:   []httprouter.Param{},
			response: ``,
			err:      pkgErrors.ErrInvalidBoardIdParam,
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

			del := delivery{
				serv: f.serv,
				log:  logger,
			}

			req := httptest.NewRequest(http.MethodPut, "/boards/1", strings.NewReader(test.request))
			rec := httptest.NewRecorder()
			err := del.fullUpdate(rec, req, test.params)
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

func TestPartialUpdate(t *testing.T) {
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
				f.serv.EXPECT().PartialUpdate(&_boards.PartialUpdateParams{
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
			params:   []httprouter.Param{{Key: "id", Value: "1"}},
			request:  `{"name":"b1","description":"d1","privacy":"secret"}`,
			response: `{"id":1,"name":"b1","description":"d1","privacy":"secret","user_id":3}`,
			err:      nil,
		},
		"invalid board id param": {
			prepare:  func(f *fields) {},
			params:   []httprouter.Param{{Key: "id", Value: "a"}},
			response: ``,
			err:      pkgErrors.ErrInvalidBoardIdParam,
		},
		"missing board id param": {
			prepare:  func(f *fields) {},
			params:   []httprouter.Param{},
			response: ``,
			err:      pkgErrors.ErrInvalidBoardIdParam,
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

			del := delivery{
				serv: f.serv,
				log:  logger,
			}

			req := httptest.NewRequest(http.MethodPatch, "/boards/1", strings.NewReader(test.request))
			rec := httptest.NewRecorder()
			err := del.partialUpdate(rec, req, test.params)
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

func TestDelete(t *testing.T) {
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
				f.serv.EXPECT().Delete(3).Return(nil)
			},
			params: []httprouter.Param{{Key: "id", Value: "3"}},
			err:    pkgErrors.ErrNoContent,
		},
		"invalid board id param": {
			prepare: func(f *fields) {},
			params:  []httprouter.Param{{Key: "id", Value: "a"}},
			err:     pkgErrors.ErrInvalidBoardIdParam,
		},
		"missing board id param": {
			prepare: func(f *fields) {},
			params:  []httprouter.Param{},
			err:     pkgErrors.ErrInvalidBoardIdParam,
		},
		"board not found": {
			prepare: func(f *fields) {
				f.serv.EXPECT().Delete(3).Return(pkgErrors.ErrBoardNotFound)
			},
			params: []httprouter.Param{{Key: "id", Value: "3"}},
			err:    pkgErrors.ErrBoardNotFound,
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

			del := delivery{
				serv: f.serv,
				log:  logger,
			}

			req := httptest.NewRequest(http.MethodDelete, "/boards/3", nil)
			rec := httptest.NewRecorder()
			err := del.delete(rec, req, test.params)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
		})
	}
}

func TestAddPin(t *testing.T) {
	type fields struct {
		serv *mocks.MockService
	}

	type testCase struct {
		prepare func(f *fields)
		params  httprouter.Params
		request string
		err     error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				f.serv.EXPECT().AddPin(3, 2).Return(nil)
			},
			params: []httprouter.Param{
				{Key: "id", Value: "3"},
				{Key: "pin_id", Value: "2"},
			},
			err: pkgErrors.ErrNoContent,
		},
		"invalid board id param": {
			prepare: func(f *fields) {},
			params: []httprouter.Param{
				{Key: "id", Value: "a"},
				{Key: "pin_id", Value: "3"},
			},
			err: pkgErrors.ErrInvalidBoardIdParam,
		},
		"invalid pin id": {
			prepare: func(f *fields) {},
			params: []httprouter.Param{
				{Key: "id", Value: "3"},
				{Key: "pin_id", Value: "a"},
			},
			err: pkgErrors.ErrInvalidPinIdParam,
		},
		"db error": {
			prepare: func(f *fields) {
				f.serv.EXPECT().AddPin(3, 2).Return(pkgErrors.ErrDb)
			},
			params: []httprouter.Param{
				{Key: "id", Value: "3"},
				{Key: "pin_id", Value: "2"},
			},
			err: pkgErrors.ErrDb,
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

			del := delivery{
				serv: f.serv,
				log:  logger,
			}

			req := httptest.NewRequest(http.MethodPost, "/boards/3/pins/2", strings.NewReader(test.request))
			rec := httptest.NewRecorder()
			err := del.addPin(rec, req, test.params)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
		})
	}
}

func TestPinsList(t *testing.T) {
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
				f.serv.EXPECT().PinsList(3, 12, 1, 30).Return([]models.Pin{
					{Id: 1, Title: "t1", MediaSource: "ms_url1", MediaSourceColor: "rgb(39, 102, 120)",
						Description: "d1", Author: models.Profile{
							Id:           12,
							Username:     "un1",
							Name:         "n1",
							ProfileImage: "pi1",
							WebsiteUrl:   "wu1",
						}},
					{Id: 2, Title: "t2", MediaSource: "ms_url2", MediaSourceColor: "rgb(39, 102, 120)",
						Description: "d2", Author: models.Profile{
							Id:           13,
							Username:     "un2",
							Name:         "n2",
							ProfileImage: "pi2",
							WebsiteUrl:   "wu2",
						}},
					{Id: 3, Title: "t3", MediaSource: "ms_url3", MediaSourceColor: "rgb(39, 102, 120)",
						Description: "d3", Author: models.Profile{
							Id:           14,
							Username:     "un3",
							Name:         "n3",
							ProfileImage: "pi3",
							WebsiteUrl:   "wu3",
						}},
				}, nil)
			},
			params: []httprouter.Param{
				{Key: "page", Value: "1"},
				{Key: "limit", Value: "30"},
				{Key: "id", Value: "12"},
				{Key: "user-id", Value: "3"},
			},
			response: `{"pins":[{"id":1,"title":"t1","description":"d1","media_source":"ms_url1","media_source_color":"rgb(39, 102, 120)","n_likes":0,"liked":false,"author":{"id":12,"username":"un1","name":"n1","profile_image":"pi1","website_url":"wu1"}},{"id":2,"title":"t2","description":"d2","media_source":"ms_url2","media_source_color":"rgb(39, 102, 120)","n_likes":0,"liked":false,"author":{"id":13,"username":"un2","name":"n2","profile_image":"pi2","website_url":"wu2"}},{"id":3,"title":"t3","description":"d3","media_source":"ms_url3","media_source_color":"rgb(39, 102, 120)","n_likes":0,"liked":false,"author":{"id":14,"username":"un3","name":"n3","profile_image":"pi3","website_url":"wu3"}}]}`,
			err:      nil,
		},
		"no pins": {
			prepare: func(f *fields) {
				f.serv.EXPECT().PinsList(3, 12, 1, 30).Return([]models.Pin{}, nil)
			},
			params: []httprouter.Param{
				{Key: "page", Value: "1"},
				{Key: "limit", Value: "30"},
				{Key: "id", Value: "12"},
				{Key: "user-id", Value: "3"},
			},
			response: `{"pins":[]}`,
			err:      nil,
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

			del := delivery{
				serv: f.serv,
				log:  logger,
			}

			req := httptest.NewRequest(http.MethodGet, "/boards/12/pins", nil)
			rec := httptest.NewRecorder()
			err := del.pinsList(rec, req, test.params)
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
