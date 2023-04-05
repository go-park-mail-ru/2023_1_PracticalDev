package http

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/boards"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/boards/mocks"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

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
			params:     []httprouter.Param{{Key: "id", Value: "3"}},
			response:   `{"id":3,"name":"n3","description":"d3","privacy":"secret","user_id":1}`,
			statusCode: http.StatusOK,
			err:        nil,
		},
		"invalid board id param": {
			prepare:    func(f *fields) {},
			params:     []httprouter.Param{{Key: "id", Value: "a"}},
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
			params:     []httprouter.Param{{Key: "id", Value: "3"}},
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
