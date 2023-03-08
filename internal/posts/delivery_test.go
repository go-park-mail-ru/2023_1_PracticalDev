package posts

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/db"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
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
			Response:   `[{"id":1,"link":"","title":"Road","description":"","media_source":"https://wg.grechkogv.ru/assets/pet7.webp","board_id":1},{"id":2,"link":"","title":"Ice","description":"","media_source":"https://wg.grechkogv.ru/assets/armorChest4.webp","board_id":1},{"id":3,"link":"","title":"Future","description":"","media_source":"https://wg.grechkogv.ru/assets/pet6.webp","board_id":1},{"id":4,"link":"","title":"Color","description":"","media_source":"https://wg.grechkogv.ru/assets/pet8.webp","board_id":2},{"id":5,"link":"","title":"Question","description":"","media_source":"https://wg.grechkogv.ru/assets/weapon5.webp","board_id":2},{"id":6,"link":"","title":"Shops","description":"","media_source":"https://wg.grechkogv.ru/assets/weapon1.webp","board_id":3},{"id":7,"link":"","title":"School","description":"","media_source":"https://wg.grechkogv.ru/assets/armorBeing3.webp","board_id":4}]`,
			StatusCode: http.StatusOK,
		},
		{
			Parameters: `?page=2&limit=3`,
			Response:   `[{"id":4,"link":"","title":"Color","description":"","media_source":"https://wg.grechkogv.ru/assets/pet8.webp","board_id":2},{"id":5,"link":"","title":"Question","description":"","media_source":"https://wg.grechkogv.ru/assets/weapon5.webp","board_id":2},{"id":6,"link":"","title":"Shops","description":"","media_source":"https://wg.grechkogv.ru/assets/weapon1.webp","board_id":3}]`,
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
			t.Errorf("[%d] wrong StatusCode: \ngot %d, \nexpected %d, \nerr %d",
				caseNum, w.Code, item.StatusCode, err)
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
