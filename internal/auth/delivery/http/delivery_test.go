package http

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth/tokens"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth"
	authMocks "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth/mocks"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	mw "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models/api"
)

var (
	logger log.Logger
	err    error
)

var existingUsers []api.LoginParams = []api.LoginParams{
	{
		Email:    "geogreck@vk.com",
		Password: "12345678",
	},
	{
		Email:    "figma@vk.com",
		Password: "12345678",
	},
	{
		Email:    "iu7@vk.com",
		Password: "12345678",
	},
	{
		Email:    "test@vk.com",
		Password: "12345678",
	},
}

type fields struct {
	serv *authMocks.MockService
}

type AuthenticateTestCase struct {
	prepare func(f *fields)
	req     api.LoginParams
	err     error
}

type RegisterTestCase struct {
	prepare func(f *fields)
	req     api.RegisterParams
	err     error
}

type LogoutTestCase struct {
	prepare func(f *fields)
	cookie  *http.Cookie
	err     error
}

func TestAuthenticate(t *testing.T) {
	tests := []AuthenticateTestCase{
		{
			prepare: func(f *fields) {
				f.serv.EXPECT().Authenticate(existingUsers[0].Email, existingUsers[0].Password).
					Return(models.User{
						Id:             2,
						Username:       "vitya",
						Email:          existingUsers[0].Email,
						HashedPassword: "hashed_pswd",
						Name:           "Vitya",
						ProfileImage:   "img.png",
						WebsiteUrl:     "www.vk.ru",
						AccountType:    "personal",
					}, auth.SessionParams{}, nil)
			},
			req: existingUsers[0],
			err: nil,
		},
		{
			prepare: func(f *fields) {
				f.serv.EXPECT().Authenticate("123@vk.com", "12345678").
					Return(models.User{}, auth.SessionParams{}, auth.WrongPasswordOrLoginError)
			},
			req: api.LoginParams{
				Email:    "123@vk.com",
				Password: "12345678",
			},
			err: mw.ErrUserNotFound,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for testNum, test := range tests {
		f := fields{serv: authMocks.NewMockService(ctrl)}
		if test.prepare != nil {
			test.prepare(&f)
		}

		del := delivery{f.serv, logger, tokens.NewHMACHashToken("test_secret")}

		url := "http://127.0.0.1/api/auth/login"
		tmp, _ := json.Marshal(test.req)
		body := strings.NewReader(string(tmp))

		req := httptest.NewRequest(http.MethodPost, url, body)
		w := httptest.NewRecorder()
		err = del.Authenticate(w, req, nil)
		if err != test.err {
			t.Errorf("\n[%d] \nExpected: %s\nGot: %s", testNum, test.err, err)
		}
	}
}

func TestRegister(t *testing.T) {
	tests := []RegisterTestCase{
		{
			prepare: func(f *fields) {
				f.serv.EXPECT().Register(&api.RegisterParams{
					Username: "test1",
					Email:    "test1@test.ru",
					Name:     "test",
					Password: "12345",
				}).Return(models.User{
					Id:             2,
					Username:       "test1",
					Email:          "test1@test.ru",
					HashedPassword: "hashed_pswd",
					Name:           "test",
					ProfileImage:   "",
					WebsiteUrl:     "",
					AccountType:    "personal",
				}, auth.SessionParams{}, nil)
			},
			req: api.RegisterParams{
				Username: "test1",
				Email:    "test1@test.ru",
				Name:     "test",
				Password: "12345",
			},
			err: nil,
		},
		{
			prepare: func(f *fields) {
				f.serv.EXPECT().Register(&api.RegisterParams{
					Username: "test3",
					Email:    "test1@test.ru",
					Name:     "test",
					Password: "12345",
				}).Return(models.User{
					Id:             2,
					Username:       "test3",
					Email:          "test1@test.ru",
					HashedPassword: "hashed_pswd",
					Name:           "test",
					ProfileImage:   "",
					WebsiteUrl:     "",
					AccountType:    "personal",
				}, auth.SessionParams{}, auth.UserAlreadyExistsError)
			},
			req: api.RegisterParams{
				Username: "test3",
				Email:    "test1@test.ru",
				Name:     "test",
				Password: "12345",
			},
			err: mw.ErrUserAlreadyExists,
		},
		{
			prepare: func(f *fields) {
				f.serv.EXPECT().Register(&api.RegisterParams{
					Username: "test3",
					Email:    "test1@test.ru",
					Name:     "test",
					Password: "12345",
				}).Return(models.User{
					Id:             2,
					Username:       "test3",
					Email:          "test1@test.ru",
					HashedPassword: "hashed_pswd",
					Name:           "test",
					ProfileImage:   "",
					WebsiteUrl:     "",
					AccountType:    "personal",
				}, auth.SessionParams{}, auth.DBConnectionError)
			},
			req: api.RegisterParams{
				Username: "test3",
				Email:    "test1@test.ru",
				Name:     "test",
				Password: "12345",
			},
			err: mw.ErrService,
		},
		{
			prepare: func(f *fields) {
				f.serv.EXPECT().Register(&api.RegisterParams{
					Username: "test3",
					Email:    "test1@test.ru",
					Name:     "test",
					Password: "12345",
				}).Return(models.User{
					Id:             2,
					Username:       "test3",
					Email:          "test1@test.ru",
					HashedPassword: "hashed_pswd",
					Name:           "test",
					ProfileImage:   "",
					WebsiteUrl:     "",
					AccountType:    "personal",
				}, auth.SessionParams{}, auth.UserCreationError)
			},
			req: api.RegisterParams{
				Username: "test3",
				Email:    "test1@test.ru",
				Name:     "test",
				Password: "12345",
			},
			err: mw.ErrService,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for testNum, test := range tests {
		f := fields{serv: authMocks.NewMockService(ctrl)}
		if test.prepare != nil {
			test.prepare(&f)
		}

		del := delivery{f.serv, logger, tokens.NewHMACHashToken("test_secret")}

		url := "http://127.0.0.1/api/auth/signup"
		tmp, _ := json.Marshal(test.req)
		body := strings.NewReader(string(tmp))

		req := httptest.NewRequest(http.MethodPost, url, body)
		w := httptest.NewRecorder()
		err = del.Register(w, req, nil)
		if err != test.err {
			t.Errorf("\n[%d] \nExpected: %s\nGot: %s", testNum, test.err, err)
		}
	}
}

func TestLogout(t *testing.T) {
	tests := []LogoutTestCase{
		{
			prepare: func(f *fields) {},
			cookie: &http.Cookie{
				Name:  "JSESSIONID",
				Value: "123456789",
			},
			err: mw.ErrBadRequest,
		},
		{
			prepare: func(f *fields) {
				f.serv.EXPECT().DeleteSession("1", "1$23456789").Return(auth.WrongPasswordOrLoginError)
			},
			cookie: &http.Cookie{
				Name:  "JSESSIONID",
				Value: "1$23456789",
			},
			err: mw.ErrUnauthorized,
		},
		{
			prepare: func(f *fields) {},
			cookie:  &http.Cookie{},
			err:     mw.ErrUserNotFound,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for testNum, test := range tests {
		f := fields{serv: authMocks.NewMockService(ctrl)}
		if test.prepare != nil {
			test.prepare(&f)
		}

		del := delivery{f.serv, logger, tokens.NewHMACHashToken("test_secret")}

		const url = "http://127.0.0.1/api/auth/logout"
		req := httptest.NewRequest(http.MethodDelete, url, nil)
		req.AddCookie(test.cookie)
		w := httptest.NewRecorder()
		err = del.Logout(w, req, nil)
		if err != test.err {
			t.Errorf("\n[%d] \nExpected: %s\nGot: %s", testNum, test.err, err)
		}
	}
}
