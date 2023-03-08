package users

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/db"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type testCase struct {
	Id         string
	Response   string
	StatusCode int
	isError    bool
}

func TestGetUser(t *testing.T) {
	cases := []testCase{
		{
			Id:         "1",
			Response:   `{"id":1,"username":"geogreck","email":"geogreck@vk.com","name":"George","profile_image":"","website_url":"","account_type":"personal"}`,
			StatusCode: http.StatusOK,
			isError:    false,
		},
		{
			Id:         "2",
			Response:   `{"id":2,"username":"kirill","email":"figma@vk.com","name":"Kirill","profile_image":"","website_url":"","account_type":"personal"}`,
			StatusCode: http.StatusOK,
			isError:    false,
		},
		{
			Id:         "0",
			Response:   "",
			StatusCode: http.StatusNotFound,
			isError:    true,
		},
		{
			Id:         "100",
			Response:   "",
			StatusCode: http.StatusNotFound,
			isError:    true,
		},
		{
			Id:         "f",
			Response:   "",
			StatusCode: http.StatusBadRequest,
			isError:    true,
		},
	}

	logger := log.New()
	db, err := db.New(logger)
	if err != nil {
		os.Exit(1)
	}
	del := delivery{NewService(NewRepository(db, logger)), logger}

	for _, item := range cases {
		url := "http://127.0.0.1/api/users/" + item.Id
		req := httptest.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()
		params := httprouter.Params{httprouter.Param{Key: "id", Value: item.Id}}

		err = del.getUser(w, req, params)

		isError := err != nil
		assert.Equal(t, isError, item.isError)
		assert.Equal(t, w.Code, item.StatusCode)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		bodyStr := string(body)

		assert.Equal(t, bodyStr, item.Response)
	}
}
