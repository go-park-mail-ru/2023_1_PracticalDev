package http

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pins"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pins/mocks"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/utils"
)

func TestCreate(t *testing.T) {
	type fields struct {
		serv *mocks.MockService
	}

	type testCase struct {
		prepare    func(f *fields)
		params     httprouter.Params
		formValues map[string]string
		formFiles  map[string]utils.File
		response   string
		statusCode int
		err        error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				f.serv.EXPECT().Create(gomock.Any()).Return(models.Pin{
					Id:          1,
					Title:       "t1",
					Description: "d1",
					MediaSource: "ms_url",
					Author:      3,
				}, nil)
			},
			params: []httprouter.Param{{Key: "user-id", Value: "3"}},
			formValues: map[string]string{
				"title":       "t1",
				"description": "d1",
			},
			formFiles: map[string]utils.File{
				"bytes": {
					Name:  "test.jpg",
					Bytes: make([]byte, 3),
				},
			},
			response:   `{"id":1,"title":"t1","description":"d1","media_source":"ms_url","author_id":3}`,
			statusCode: http.StatusOK,
			err:        nil,
		},
		"invalid user id param": {
			prepare: func(f *fields) {},
			params:  []httprouter.Param{{Key: "user-id", Value: "a"}},
			formValues: map[string]string{
				"title":       "t1",
				"description": "d1",
			},
			formFiles: map[string]utils.File{
				"bytes": {
					Name:  "test.jpg",
					Bytes: make([]byte, 3),
				},
			},
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

			reqBody, contentType, err := utils.CreateMultipartFormBody(test.formValues, test.formFiles)
			if err != nil {
				t.Fatal(err)
			}
			req := httptest.NewRequest(http.MethodPost, "/pins", reqBody)
			req.Header.Add("Content-Type", contentType)
			rec := httptest.NewRecorder()
			err = del.create(rec, req, test.params)
			if err != test.err {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if rec.Code != test.statusCode {
				t.Errorf("\nExpected: %d\nGot: %d", test.statusCode, rec.Code)
			}
			respBody, _ := io.ReadAll(rec.Body)
			if strings.Trim(string(respBody), "\n") != test.response {
				t.Errorf("\nExpected: %s\nGot: %s", test.response, string(respBody))
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
				f.serv.EXPECT().List(1, 30).Return([]models.Pin{
					{Id: 1, Title: "t1", MediaSource: "ms_url1", Description: "d1", Author: 12},
					{Id: 2, Title: "t2", MediaSource: "ms_url2", Description: "d2", Author: 3},
					{Id: 3, Title: "t3", MediaSource: "ms_url3", Description: "d3", Author: 10},
				}, nil)
			},
			params: []httprouter.Param{
				{Key: "page", Value: "1"},
				{Key: "limit", Value: "30"},
			},
			response:   `{"pins":[{"id":1,"title":"t1","description":"d1","media_source":"ms_url1","author_id":12},{"id":2,"title":"t2","description":"d2","media_source":"ms_url2","author_id":3},{"id":3,"title":"t3","description":"d3","media_source":"ms_url3","author_id":10}]}`,
			statusCode: http.StatusOK,
			err:        nil,
		},
		"no pins": {
			prepare: func(f *fields) {
				f.serv.EXPECT().List(1, 30).Return([]models.Pin{}, nil)
			},
			params:     []httprouter.Param{{Key: "user-id", Value: "3"}},
			response:   `{"pins":[]}`,
			statusCode: http.StatusOK,
			err:        nil,
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

			req := httptest.NewRequest(http.MethodGet, "/pins", nil)
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

func TestListByUser(t *testing.T) {
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
				f.serv.EXPECT().ListByUser(12, 1, 30).Return([]models.Pin{
					{Id: 1, Title: "t1", MediaSource: "ms_url1", Description: "d1", Author: 12},
					{Id: 2, Title: "t2", MediaSource: "ms_url2", Description: "d2", Author: 12},
					{Id: 3, Title: "t3", MediaSource: "ms_url3", Description: "d3", Author: 12},
				}, nil)
			},
			params: []httprouter.Param{
				{Key: "page", Value: "1"},
				{Key: "limit", Value: "30"},
				{Key: "id", Value: "12"},
			},
			response:   `{"pins":[{"id":1,"title":"t1","description":"d1","media_source":"ms_url1","author_id":12},{"id":2,"title":"t2","description":"d2","media_source":"ms_url2","author_id":12},{"id":3,"title":"t3","description":"d3","media_source":"ms_url3","author_id":12}]}`,
			statusCode: http.StatusOK,
			err:        nil,
		},
		"no pins": {
			prepare: func(f *fields) {
				f.serv.EXPECT().ListByUser(12, 1, 30).Return([]models.Pin{}, nil)
			},
			params: []httprouter.Param{
				{Key: "page", Value: "1"},
				{Key: "limit", Value: "30"},
				{Key: "id", Value: "12"},
			},
			response:   `{"pins":[]}`,
			statusCode: http.StatusOK,
			err:        nil,
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

			req := httptest.NewRequest(http.MethodGet, "/users/12/pins", nil)
			rec := httptest.NewRecorder()
			err := del.listByUser(rec, req, test.params)
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

func TestListByBoard(t *testing.T) {
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
				f.serv.EXPECT().ListByBoard(12, 1, 30).Return([]models.Pin{
					{Id: 1, Title: "t1", MediaSource: "ms_url1", Description: "d1", Author: 12},
					{Id: 2, Title: "t2", MediaSource: "ms_url2", Description: "d2", Author: 10},
					{Id: 3, Title: "t3", MediaSource: "ms_url3", Description: "d3", Author: 3},
				}, nil)
			},
			params: []httprouter.Param{
				{Key: "page", Value: "1"},
				{Key: "limit", Value: "30"},
				{Key: "board_id", Value: "12"},
			},
			response:   `{"pins":[{"id":1,"title":"t1","description":"d1","media_source":"ms_url1","author_id":12},{"id":2,"title":"t2","description":"d2","media_source":"ms_url2","author_id":10},{"id":3,"title":"t3","description":"d3","media_source":"ms_url3","author_id":3}]}`,
			statusCode: http.StatusOK,
			err:        nil,
		},
		"no pins": {
			prepare: func(f *fields) {
				f.serv.EXPECT().ListByBoard(12, 1, 30).Return([]models.Pin{}, nil)
			},
			params: []httprouter.Param{
				{Key: "page", Value: "1"},
				{Key: "limit", Value: "30"},
				{Key: "board_id", Value: "12"},
			},
			response:   `{"pins":[]}`,
			statusCode: http.StatusOK,
			err:        nil,
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

			req := httptest.NewRequest(http.MethodGet, "/boards/12/pins", nil)
			rec := httptest.NewRecorder()
			err := del.listByBoard(rec, req, test.params)
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
				f.serv.EXPECT().Get(3).Return(models.Pin{
					Id:          3,
					Title:       "t1",
					MediaSource: "ms_url1",
					Description: "d1",
					Author:      12,
				}, nil)
			},
			params:     []httprouter.Param{{Key: "id", Value: "3"}},
			response:   `{"id":3,"title":"t1","description":"d1","media_source":"ms_url1","author_id":12}`,
			statusCode: http.StatusOK,
			err:        nil,
		},
		"invalid pin id param": {
			prepare:    func(f *fields) {},
			params:     []httprouter.Param{{Key: "id", Value: "a"}},
			response:   ``,
			statusCode: http.StatusBadRequest,
			err:        ErrInvalidPinId,
		},
		"pin not found": {
			prepare: func(f *fields) {
				f.serv.EXPECT().Get(3).Return(models.Pin{}, pins.ErrPinNotFound)
			},
			params:     []httprouter.Param{{Key: "id", Value: "3"}},
			response:   ``,
			statusCode: http.StatusNotFound,
			err:        ErrPinNotFound,
		},
		"service error": {
			prepare: func(f *fields) {
				f.serv.EXPECT().Get(3).Return(models.Pin{}, pins.ErrDb)
			},
			params:     []httprouter.Param{{Key: "id", Value: "3"}},
			response:   ``,
			statusCode: http.StatusInternalServerError,
			err:        ErrService,
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

			req := httptest.NewRequest(http.MethodGet, "/pins/3", nil)
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
		formValues map[string]string
		formFiles  map[string]utils.File
		response   string
		statusCode int
		err        error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				f.serv.EXPECT().FullUpdate(gomock.Any()).Return(models.Pin{
					Id:          3,
					Title:       "t1",
					Description: "d1",
					MediaSource: "ms_url",
					Author:      12,
				}, nil)
			},
			params: []httprouter.Param{{Key: "id", Value: "3"}},
			formValues: map[string]string{
				"title":       "t1",
				"description": "d1",
			},
			formFiles: map[string]utils.File{
				"bytes": {
					Name:  "test.jpg",
					Bytes: make([]byte, 3),
				},
			},
			response:   `{"id":3,"title":"t1","description":"d1","media_source":"ms_url","author_id":12}`,
			statusCode: http.StatusOK,
			err:        nil,
		},
		"missing file": {
			prepare: func(f *fields) {},
			params:  []httprouter.Param{{Key: "id", Value: "3"}},
			formValues: map[string]string{
				"title":       "t1",
				"description": "d1",
			},
			formFiles:  map[string]utils.File{},
			response:   ``,
			statusCode: http.StatusBadRequest,
			err:        ErrMissingFile,
		},
		"invalid user id param": {
			prepare: func(f *fields) {},
			params:  []httprouter.Param{{Key: "id", Value: "a"}},
			formValues: map[string]string{
				"title":       "t1",
				"description": "d1",
			},
			formFiles: map[string]utils.File{
				"bytes": {
					Name:  "test.jpg",
					Bytes: make([]byte, 3),
				},
			},
			response:   ``,
			statusCode: http.StatusBadRequest,
			err:        ErrInvalidPinId,
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

			reqBody, contentType, err := utils.CreateMultipartFormBody(test.formValues, test.formFiles)
			if err != nil {
				t.Fatal(err)
			}
			req := httptest.NewRequest(http.MethodPut, "/pins/3", reqBody)
			req.Header.Add("Content-Type", contentType)
			rec := httptest.NewRecorder()
			err = del.fullUpdate(rec, req, test.params)
			if err != test.err {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if rec.Code != test.statusCode {
				t.Errorf("\nExpected: %d\nGot: %d", test.statusCode, rec.Code)
			}
			respBody, _ := io.ReadAll(rec.Body)
			if strings.Trim(string(respBody), "\n") != test.response {
				t.Errorf("\nExpected: %s\nGot: %s", test.response, string(respBody))
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
			params:     []httprouter.Param{{Key: "id", Value: "3"}},
			statusCode: http.StatusOK,
			err:        nil,
		},
		"invalid pin id param": {
			prepare:    func(f *fields) {},
			params:     []httprouter.Param{{Key: "id", Value: "a"}},
			statusCode: http.StatusBadRequest,
			err:        ErrInvalidPinId,
		},
		"missing pin id param": {
			prepare:    func(f *fields) {},
			params:     []httprouter.Param{},
			statusCode: http.StatusBadRequest,
			err:        ErrInvalidPinId,
		},
		"pin not found": {
			prepare: func(f *fields) {
				f.serv.EXPECT().Delete(3).Return(pins.ErrPinNotFound)
			},
			params:     []httprouter.Param{{Key: "id", Value: "3"}},
			statusCode: http.StatusNotFound,
			err:        ErrPinNotFound,
		},
		"service error": {
			prepare: func(f *fields) {
				f.serv.EXPECT().Delete(3).Return(pins.ErrDb)
			},
			params:     []httprouter.Param{{Key: "id", Value: "3"}},
			statusCode: http.StatusInternalServerError,
			err:        ErrService,
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

			req := httptest.NewRequest(http.MethodDelete, "/pins/3", nil)
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

func TestAddToBoard(t *testing.T) {
	type fields struct {
		serv *mocks.MockService
	}

	type testCase struct {
		prepare    func(f *fields)
		params     httprouter.Params
		request    string
		statusCode int
		err        error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				f.serv.EXPECT().AddToBoard(3, 2).Return(nil)
			},
			params: []httprouter.Param{
				{Key: "board_id", Value: "3"},
				{Key: "id", Value: "2"},
			},
			statusCode: http.StatusOK,
			err:        nil,
		},
		"invalid board id param": {
			prepare: func(f *fields) {},
			params: []httprouter.Param{
				{Key: "board_id", Value: "a"},
				{Key: "id", Value: "3"},
			},
			statusCode: http.StatusBadRequest,
			err:        ErrInvalidBoardId,
		},
		"invalid pin id": {
			prepare: func(f *fields) {},
			params: []httprouter.Param{
				{Key: "board_id", Value: "3"},
				{Key: "id", Value: "a"},
			},
			statusCode: http.StatusBadRequest,
			err:        ErrInvalidPinId,
		},
		"service error": {
			prepare: func(f *fields) {
				f.serv.EXPECT().AddToBoard(3, 2).Return(pins.ErrDb)
			},
			params: []httprouter.Param{
				{Key: "board_id", Value: "3"},
				{Key: "id", Value: "2"},
			},
			statusCode: http.StatusInternalServerError,
			err:        ErrService,
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

			req := httptest.NewRequest(http.MethodPost, "/boards/3/pins/2", strings.NewReader(test.request))
			rec := httptest.NewRecorder()
			err := del.addToBoard(rec, req, test.params)
			if err != test.err {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if rec.Code != test.statusCode {
				t.Errorf("\nExpected: %d\nGot: %d", test.statusCode, rec.Code)
			}
		})
	}
}
