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

type testCase struct {
	Parameters string
	Response   string
	StatusCode int
}

func TestDeliveryGetPosts(t *testing.T) {
	cases := []testCase{
		{
			Parameters: ``,
			Response:   `[{"id":1,"link":"","title":"Road","description":"","media_source":"https://wg.grechkogv.ru/assets/pet7.webp","board_id":1},{"id":2,"link":"","title":"Ice","description":"","media_source":"https://wg.grechkogv.ru/assets/armorChest4.webp","board_id":1},{"id":3,"link":"","title":"Future","description":"","media_source":"https://wg.grechkogv.ru/assets/pet6.webp","board_id":1},{"id":4,"link":"","title":"Color","description":"","media_source":"https://wg.grechkogv.ru/assets/pet8.webp","board_id":2},{"id":5,"link":"","title":"Shops","description":"","media_source":"https://i.pinimg.com/564x/2f/93/56/2f9356b9346e82c14bf286c6a107bc7a.jpg","board_id":3},{"id":6,"link":"","title":"Shops","description":"","media_source":"https://i.pinimg.com/564x/32/ff/71/32ff717c3cd3bd3d1886c775b59f0769.jpg","board_id":3},{"id":7,"link":"","title":"Shops","description":"","media_source":"https://i.pinimg.com/564x/ce/e3/01/cee3011f3e19de4377dbf98f397c027b.jpg","board_id":3},{"id":8,"link":"","title":"Shops","description":"","media_source":"https://i.pinimg.com/564x/a6/ba/55/a6ba553df2a0c0f3894ef328a86fb373.jpg","board_id":3},{"id":9,"link":"","title":"Shops","description":"","media_source":"https://i.pinimg.com/564x/43/2d/3b/432d3b28d1661439245422e9218ffcce.jpg","board_id":3},{"id":10,"link":"","title":"School","description":"","media_source":"https://i.pinimg.com/564x/98/9d/3f/989d3f5c158dcac7ca4d115bff866d84.jpg","board_id":4}]`,
			StatusCode: http.StatusOK,
		},
		{
			Parameters: `?page=2&limit=3`,
			Response:   `[{"id":4,"link":"","title":"Color","description":"","media_source":"https://wg.grechkogv.ru/assets/pet8.webp","board_id":2},{"id":5,"link":"","title":"Shops","description":"","media_source":"https://i.pinimg.com/564x/2f/93/56/2f9356b9346e82c14bf286c6a107bc7a.jpg","board_id":3},{"id":6,"link":"","title":"Shops","description":"","media_source":"https://i.pinimg.com/564x/32/ff/71/32ff717c3cd3bd3d1886c775b59f0769.jpg","board_id":3}]`,
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
