package http

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/users/mocks"
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
				f.serv.EXPECT().Get(3).Return(models.User{
					Id:           3,
					Username:     "petya",
					Email:        "petya@vk.com",
					Name:         "Petya",
					ProfileImage: "petya.jpg",
					WebsiteUrl:   "petya.com",
					AccountType:  "personal"}, nil)
			},
			params:   []httprouter.Param{{Key: "id", Value: "3"}},
			response: `{"id":3,"username":"petya","email":"petya@vk.com","name":"Petya","profile_image":"petya.jpg","website_url":"petya.com","account_type":"personal"}`,
			err:      nil,
		},
		"invalid user id param": {
			prepare:  nil,
			params:   []httprouter.Param{{Key: "id", Value: "a"}},
			response: ``,
			err:      pkgErrors.ErrInvalidUserIdParam,
		},
		"profile not found": {
			prepare: func(f *fields) {
				f.serv.EXPECT().Get(3).Return(models.User{}, pkgErrors.ErrProfileNotFound)
			},
			params:   []httprouter.Param{{Key: "id", Value: "3"}},
			response: ``,
			err:      pkgErrors.ErrProfileNotFound,
		},
		"service error": {
			prepare: func(f *fields) {
				f.serv.EXPECT().Get(3).Return(models.User{}, pkgErrors.ErrDb)
			},
			params:   []httprouter.Param{{Key: "id", Value: "3"}},
			response: ``,
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

			req := httptest.NewRequest(http.MethodGet, "/users/3", nil)
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
