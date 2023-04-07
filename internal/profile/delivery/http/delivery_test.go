package http

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/profile"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/profile/mocks"
)

type file struct {
	Name  string
	Bytes []byte
}

func createMultipartFormBody(values map[string]string,
	files map[string]file) (body *bytes.Buffer, contentType string, err error) {
	body = new(bytes.Buffer)
	mw := multipart.NewWriter(body)
	defer mw.Close()

	for key, value := range files {
		w, err := mw.CreateFormFile(key, value.Name)
		if err != nil {
			return nil, "", err
		}
		if _, err = w.Write(value.Bytes); err != nil {
			return nil, "", err
		}
	}

	for key, value := range values {
		err = mw.WriteField(key, value)
		if err != nil {
			return nil, "", err
		}
	}

	return body, mw.FormDataContentType(), nil
}

func TestGetProfileByUser(t *testing.T) {
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
				f.serv.EXPECT().GetProfileByUser(3).Return(profile.Profile{
					Username:     "un1",
					Name:         "n1",
					ProfileImage: "pi1",
					WebsiteUrl:   "wu1",
				}, nil)
			},
			params:     []httprouter.Param{{Key: "id", Value: "3"}},
			response:   `{"username":"un1","name":"n1","profile_image":"pi1","website_url":"wu1"}`,
			statusCode: http.StatusOK,
			err:        nil,
		},
		"invalid user id param": {
			prepare:    func(f *fields) {},
			params:     []httprouter.Param{{Key: "id", Value: "a"}},
			response:   ``,
			statusCode: http.StatusBadRequest,
			err:        ErrInvalidUserId,
		},
		"profile not found": {
			prepare: func(f *fields) {
				f.serv.EXPECT().GetProfileByUser(3).Return(profile.Profile{}, profile.ErrProfileNotFound)
			},
			params:     []httprouter.Param{{Key: "id", Value: "3"}},
			response:   ``,
			statusCode: http.StatusNotFound,
			err:        ErrProfileNotFound,
		},
		"service error": {
			prepare: func(f *fields) {
				f.serv.EXPECT().GetProfileByUser(3).Return(profile.Profile{}, profile.ErrDb)
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

			req := httptest.NewRequest(http.MethodGet, "/users/3/profile", nil)
			rec := httptest.NewRecorder()
			err := del.getProfileByUser(rec, req, test.params)
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
		formFiles  map[string]file
		response   string
		statusCode int
		err        error
	}

	tests := map[string]testCase{
		"missing file": {
			prepare: func(f *fields) {},
			params:  []httprouter.Param{{Key: "user-id", Value: "3"}},
			formValues: map[string]string{
				"username":    "username1",
				"name":        "n1",
				"website_url": "wu1",
			},
			formFiles:  map[string]file{},
			response:   ``,
			statusCode: http.StatusBadRequest,
			err:        ErrMissingFile,
		},
		"invalid user id param": {
			prepare: func(f *fields) {},
			params:  []httprouter.Param{{Key: "user-id", Value: "a"}},
			formValues: map[string]string{
				"username":    "username1",
				"name":        "n1",
				"website_url": "wu1",
			},
			formFiles: map[string]file{
				"bytes": {
					Name:  "avatar1.jpg",
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

			reqBody, contentType, err := createMultipartFormBody(test.formValues, test.formFiles)
			if err != nil {
				t.Fatal(err)
			}
			req := httptest.NewRequest(http.MethodPut, "/profile", reqBody)
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

func TestPartialUpdate(t *testing.T) {
	type fields struct {
		serv *mocks.MockService
	}

	type testCase struct {
		prepare    func(f *fields)
		params     httprouter.Params
		formValues map[string]string
		formFiles  map[string]file
		response   string
		statusCode int
		err        error
	}

	tests := map[string]testCase{
		"invalid user id param": {
			prepare:    func(f *fields) {},
			params:     []httprouter.Param{{Key: "user-id", Value: "a"}},
			formValues: map[string]string{"username": "username1"},
			formFiles:  map[string]file{},
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

			reqBody, contentType, err := createMultipartFormBody(test.formValues, test.formFiles)
			if err != nil {
				t.Fatal(err)
			}
			req := httptest.NewRequest(http.MethodPatch, "/profile", reqBody)
			req.Header.Add("Content-Type", contentType)
			rec := httptest.NewRecorder()
			err = del.partialUpdate(rec, req, test.params)
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
