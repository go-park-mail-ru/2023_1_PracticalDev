package http

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/profile"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/profile/mocks"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/utils"
)

func TestGetProfileByUser(t *testing.T) {
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
				f.serv.EXPECT().GetProfileByUser(3).Return(profile.Profile{
					Username:     "un1",
					Name:         "n1",
					ProfileImage: "pi1",
					WebsiteUrl:   "wu1",
				}, nil)
			},
			params:   []httprouter.Param{{Key: "id", Value: "3"}},
			response: `{"username":"un1","name":"n1","profile_image":"pi1","website_url":"wu1"}`,
			err:      nil,
		},
		"invalid user id param": {
			prepare:  func(f *fields) {},
			params:   []httprouter.Param{{Key: "id", Value: "a"}},
			response: ``,
			err:      pkgErrors.ErrInvalidUserIdParam,
		},
		"profile not found": {
			prepare: func(f *fields) {
				f.serv.EXPECT().GetProfileByUser(3).Return(profile.Profile{}, pkgErrors.ErrProfileNotFound)
			},
			params:   []httprouter.Param{{Key: "id", Value: "3"}},
			response: ``,
			err:      pkgErrors.ErrProfileNotFound,
		},
		"service error": {
			prepare: func(f *fields) {
				f.serv.EXPECT().GetProfileByUser(3).Return(profile.Profile{}, pkgErrors.ErrDb)
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

			req := httptest.NewRequest(http.MethodGet, "/users/3/profile", nil)
			rec := httptest.NewRecorder()
			err := del.getProfileByUser(rec, req, test.params)
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
		prepare    func(f *fields)
		params     httprouter.Params
		formValues map[string]string
		formFiles  map[string]utils.File
		response   string
		err        error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				f.serv.EXPECT().FullUpdate(gomock.Any()).Return(profile.Profile{
					Username:     "username1",
					Name:         "n1",
					ProfileImage: "pi_url",
					WebsiteUrl:   "wu1",
				}, nil)
			},
			params: []httprouter.Param{{Key: "user-id", Value: "3"}},
			formValues: map[string]string{
				"username":    "username1",
				"name":        "n1",
				"website_url": "wu1",
			},
			formFiles: map[string]utils.File{
				"bytes": {
					Name:  "test.jpg",
					Bytes: make([]byte, 3),
				},
			},
			response: `{"username":"username1","name":"n1","profile_image":"pi_url","website_url":"wu1"}`,
			err:      nil,
		},
		"missing file": {
			prepare: func(f *fields) {},
			params:  []httprouter.Param{{Key: "user-id", Value: "3"}},
			formValues: map[string]string{
				"username":    "username1",
				"name":        "n1",
				"website_url": "wu1",
			},
			formFiles: map[string]utils.File{},
			response:  ``,
			err:       pkgErrors.ErrMissingFile,
		},
		"invalid user id param": {
			prepare: func(f *fields) {},
			params:  []httprouter.Param{{Key: "user-id", Value: "a"}},
			formValues: map[string]string{
				"username":    "username1",
				"name":        "n1",
				"website_url": "wu1",
			},
			formFiles: map[string]utils.File{
				"bytes": {
					Name:  "avatar1.jpg",
					Bytes: make([]byte, 3),
				},
			},
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

			logger := log.New()
			del := delivery{
				serv: f.serv,
				log:  logger,
			}

			reqBody, contentType, err := utils.CreateMultipartFormBody(test.formValues, test.formFiles)
			if err != nil {
				t.Fatal(err)
			}
			req := httptest.NewRequest(http.MethodPut, "/profile", reqBody)
			req.Header.Add("Content-Type", contentType)
			rec := httptest.NewRecorder()
			err = del.fullUpdate(rec, req, test.params)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			respBody, _ := io.ReadAll(rec.Body)
			if strings.Trim(string(respBody), "\n") != test.response {
				t.Errorf("\nExpected: %s\nGot: %s", test.response, string(respBody))
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
		formValues map[string]string
		formFiles  map[string]utils.File
		response   string
		err        error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				f.serv.EXPECT().PartialUpdate(gomock.Any()).Return(profile.Profile{
					Username:     "username1",
					Name:         "n1",
					ProfileImage: "pi_url",
					WebsiteUrl:   "wu1",
				}, nil)
			},
			params: []httprouter.Param{{Key: "user-id", Value: "3"}},
			formValues: map[string]string{
				"username":    "username1",
				"name":        "n1",
				"website_url": "wu1",
			},
			formFiles: map[string]utils.File{
				"bytes": {
					Name:  "avatar1.jpg",
					Bytes: make([]byte, 3),
				},
			},
			response: `{"username":"username1","name":"n1","profile_image":"pi_url","website_url":"wu1"}`,
			err:      nil,
		},
		"invalid user id param": {
			prepare:    func(f *fields) {},
			params:     []httprouter.Param{{Key: "user-id", Value: "a"}},
			formValues: map[string]string{"username": "username1"},
			formFiles:  map[string]utils.File{},
			response:   ``,
			err:        pkgErrors.ErrInvalidUserIdParam,
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
			req := httptest.NewRequest(http.MethodPatch, "/profile", reqBody)
			req.Header.Add("Content-Type", contentType)
			rec := httptest.NewRecorder()
			err = del.partialUpdate(rec, req, test.params)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			respBody, _ := io.ReadAll(rec.Body)
			if strings.Trim(string(respBody), "\n") != test.response {
				t.Errorf("\nExpected: %s\nGot: %s", test.response, string(respBody))
			}
		})
	}
}
