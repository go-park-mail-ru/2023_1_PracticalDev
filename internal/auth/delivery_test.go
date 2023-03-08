package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/db"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models/api"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/redis"
	goRedis "github.com/redis/go-redis/v9"
)

var (
	database *sql.DB
	ctx      context.Context
	rdb      *goRedis.Client
	logger   log.Logger
	del      delivery
	err      error
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

type AuthenticateTestCase struct {
	req          api.LoginParams
	expectedCode int
}

type RegisterTestCase struct {
	req          api.RegisterParams
	expectedCode int
}

type LogoutTestCase struct {
	cookie       *http.Cookie
	expectedCode int
}

func TestMain(m *testing.M) {
	logger = log.New()
	ctx = context.Background()

	if database, err = db.New(logger); err != nil {
		os.Exit(1)
	}

	if rdb, err = redis.NewRedisClient(logger, ctx); err != nil {
		os.Exit(1)
	}

	del = delivery{NewService(NewRepository(database, rdb, ctx, logger)), logger}

	os.Exit(m.Run())
}

func TestAuthenticate(t *testing.T) {
	tests := []AuthenticateTestCase{
		{
			req:          existingUsers[0],
			expectedCode: http.StatusOK,
		},
		{
			req:          existingUsers[1],
			expectedCode: http.StatusOK,
		},
		{
			req:          existingUsers[2],
			expectedCode: http.StatusOK,
		},
		{
			req:          existingUsers[3],
			expectedCode: http.StatusOK,
		},
		{
			req: api.LoginParams{
				Email:    "123@vk.com",
				Password: "12345678",
			},
			expectedCode: http.StatusNotFound,
		},
		{
			req: api.LoginParams{
				Email:    "iu7@vk.com",
				Password: "12345678910",
			},
			expectedCode: http.StatusNotFound,
		},
	}

	for testNum, test := range tests {
		url := "http://127.0.0.1/api/auth/login"
		tmp, _ := json.Marshal(test.req)
		body := strings.NewReader(string(tmp))

		req := httptest.NewRequest("POST", url, body)
		w := httptest.NewRecorder()

		err := del.Authenticate(w, req, nil)

		if w.Result().StatusCode != test.expectedCode {
			t.Errorf("[%d] wrong StatusCode: \ngot %d, \nexpected %d, \nerr %d",
				testNum, w.Code, test.expectedCode, err)
		}
	}
}

func TestRegister(t *testing.T) {
	tests := []RegisterTestCase{
		{
			req: api.RegisterParams{
				Username: "test1",
				Email:    "test1@test.ru",
				Name:     "test",
				Password: "12345",
			},
			expectedCode: http.StatusNoContent,
		},
		{
			req: api.RegisterParams{
				Username: "test2",
				Email:    "test2@test.ru",
				Name:     "test",
				Password: "12345",
			},
			expectedCode: http.StatusNoContent,
		},
		{
			req: api.RegisterParams{
				Username: "test3",
				Email:    "test1@test.ru",
				Name:     "test",
				Password: "12345",
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for testNum, test := range tests {
		url := "http://127.0.0.1/api/auth/signup"
		tmp, _ := json.Marshal(test.req)
		body := strings.NewReader(string(tmp))

		req := httptest.NewRequest("POST", url, body)
		w := httptest.NewRecorder()

		err := del.Register(w, req, nil)
		if w.Result().StatusCode != test.expectedCode {
			t.Errorf("[%d] wrong StatusCode: \ngot %d, \nexpected %d, \nerr %d",
				testNum, w.Code, test.expectedCode, err)
		}
	}
}

func TestLogout(t *testing.T) {
	tests := []LogoutTestCase{
		{
			cookie: &http.Cookie{
				Name:  "JSESSIONID",
				Value: "123456789",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			cookie: &http.Cookie{
				Name:  "JSESSIONID",
				Value: "1$23456789",
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			cookie:       &http.Cookie{},
			expectedCode: http.StatusNotFound,
		},
	}

	for _, user := range existingUsers {

		url := "http://127.0.0.1/api/auth/login"
		tmp, _ := json.Marshal(user)
		body := strings.NewReader(string(tmp))

		req := httptest.NewRequest("POST", url, body)
		w := httptest.NewRecorder()

		err := del.Authenticate(w, req, nil)
		if err != nil {
			t.Errorf("Unexpected error: %d", err)
			break
		}
		cookie := w.Result().Cookies()[0]

		tests = append(tests, LogoutTestCase{cookie, http.StatusNoContent})

	}

	for testNum, test := range tests {
		url := "http://127.0.0.1/api/auth/signup"

		req := httptest.NewRequest("DELETE", url, nil)
		req.AddCookie(test.cookie)
		w := httptest.NewRecorder()

		err := del.Logout(w, req, nil)
		if w.Result().StatusCode != test.expectedCode {
			t.Errorf("[%d] wrong StatusCode: \ngot %d, \nexpected %d, \nerr %d",
				testNum, w.Code, test.expectedCode, err)
		}
	}
}
