package posts

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/db"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type TestCase struct {
	Parameters string
	Response   string
	StatusCode int
}

func TestGetPosts(t *testing.T) {
	cases := []TestCase{
		{
			Parameters: ``,
			Response:   `[{"id":1,"link":"","title":"Road","description":"","media_source":"","board_id":0},{"id":2,"link":"","title":"Ice","description":"","media_source":"","board_id":0},{"id":3,"link":"","title":"Future","description":"","media_source":"","board_id":0},{"id":4,"link":"","title":"Color","description":"","media_source":"","board_id":0},{"id":5,"link":"","title":"Question","description":"","media_source":"","board_id":0},{"id":6,"link":"","title":"Shops","description":"","media_source":"","board_id":0},{"id":7,"link":"","title":"School","description":"","media_source":"","board_id":0}]`,
			StatusCode: http.StatusOK,
		},
		{
			Parameters: `?page=2&limit=3`,
			Response:   `[{"id":4,"link":"","title":"Color","description":"","media_source":"","board_id":0},{"id":5,"link":"","title":"Question","description":"","media_source":"","board_id":0},{"id":6,"link":"","title":"Shops","description":"","media_source":"","board_id":0}]`,
			StatusCode: http.StatusOK,
		},
		{
			Parameters: `?page=1&limit=0`,
			Response:   `[]`,
			StatusCode: http.StatusOK,
		},
		{
			Parameters: `?page=2`,
			Response:   `[]`,
			StatusCode: http.StatusOK,
		},
		{
			Parameters: `?page=0&limit=3`,
			Response:   ``,
			StatusCode: http.StatusBadRequest,
		},
		{
			Parameters: `?page=1&limit=-1`,
			Response:   ``,
			StatusCode: http.StatusBadRequest,
		},
		{
			Parameters: `?page=1&limit=g`,
			Response:   ``,
			StatusCode: http.StatusBadRequest,
		},
	}

	logger := log.New()
	db, err := db.New(logger)
	if err != nil {
		os.Exit(1)
	}
	del := delivery{NewService(NewRepository(db, logger)), logger}

	for caseNum, item := range cases {
		url := "http://127.0.0.1/api/posts" + item.Parameters
		req := httptest.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		err = del.getPosts(w, req, nil)

		if w.Code != item.StatusCode {
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
				caseNum, w.Code, item.StatusCode)
		}

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		bodyStr := string(body)
		if bodyStr != item.Response {
			t.Errorf("[%d] wrong Response: got %+v, expected %+v",
				caseNum, bodyStr, item.Response)
		}
	}
}
