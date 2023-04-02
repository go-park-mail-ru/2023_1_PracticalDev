package posts

import (
	"os"
	"reflect"
	"testing"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/db"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
)

type serviceTestCase struct {
	Page    int
	Limit   int
	Pins    []models.Pin
	IsError bool
}

func TestServiceGetPosts(t *testing.T) {
	cases := []serviceTestCase{
		{
			Page:  1,
			Limit: 30,
			Pins: []models.Pin{
				{
					Id:          1,
					Title:       "Road",
					Description: "",
					MediaSource: "https://wg.grechkogv.ru/assets/pet7.webp",
				},
				{
					Id:          2,
					Title:       "Ice",
					Description: "",
					MediaSource: "https://wg.grechkogv.ru/assets/armorChest4.webp",
				},
				{
					Id:          3,
					Title:       "Future",
					Description: "",
					MediaSource: "https://wg.grechkogv.ru/assets/pet6.webp",
				},
				{
					Id:          4,
					Title:       "Color",
					Description: "",
					MediaSource: "https://wg.grechkogv.ru/assets/pet8.webp",
				},
				{
					Id:          5,
					Title:       "Shops",
					Description: "",
					MediaSource: "https://i.pinimg.com/564x/2f/93/56/2f9356b9346e82c14bf286c6a107bc7a.jpg",
				},
				{
					Id:          6,
					Title:       "Shops",
					Description: "",
					MediaSource: "https://i.pinimg.com/564x/32/ff/71/32ff717c3cd3bd3d1886c775b59f0769.jpg",
				},
				{
					Id:          7,
					Title:       "Shops",
					Description: "",
					MediaSource: "https://i.pinimg.com/564x/ce/e3/01/cee3011f3e19de4377dbf98f397c027b.jpg",
				},
				{
					Id:          8,
					Title:       "Shops",
					Description: "",
					MediaSource: "https://i.pinimg.com/564x/a6/ba/55/a6ba553df2a0c0f3894ef328a86fb373.jpg",
				},
				{
					Id:          9,
					Title:       "Shops",
					Description: "",
					MediaSource: "https://i.pinimg.com/564x/43/2d/3b/432d3b28d1661439245422e9218ffcce.jpg",
				},
				{
					Id:          10,
					Title:       "School",
					Description: "",
					MediaSource: "https://i.pinimg.com/564x/98/9d/3f/989d3f5c158dcac7ca4d115bff866d84.jpg",
				},
			},
			IsError: false,
		},
		{
			Page:  2,
			Limit: 3,
			Pins: []models.Pin{
				{
					Id:          4,
					Title:       "Color",
					Description: "",
					MediaSource: "https://wg.grechkogv.ru/assets/pet8.webp",
				},
				{
					Id:          5,
					Title:       "Shops",
					Description: "",
					MediaSource: "https://i.pinimg.com/564x/2f/93/56/2f9356b9346e82c14bf286c6a107bc7a.jpg",
				},
				{
					Id:          6,
					Title:       "Shops",
					Description: "",
					MediaSource: "https://i.pinimg.com/564x/32/ff/71/32ff717c3cd3bd3d1886c775b59f0769.jpg",
				},
			},
			IsError: false,
		},
		{
			Page:    1,
			Limit:   0,
			Pins:    []models.Pin{},
			IsError: false,
		},
		{
			Page:    2,
			Limit:   30,
			Pins:    []models.Pin{},
			IsError: false,
		},
		{
			Page:    -1,
			Limit:   30,
			Pins:    []models.Pin{},
			IsError: true,
		},
		{
			Page:    1,
			Limit:   -1,
			Pins:    []models.Pin{},
			IsError: true,
		},
	}

	logger := log.New()
	db, err := db.New(logger)
	if err != nil {
		os.Exit(1)
	}
	serv := NewService(NewRepository(db, logger))

	var pins []models.Pin
	for caseNum, item := range cases {
		pins, err = serv.GetPosts(item.Page, item.Limit)

		if err != nil && !item.IsError {
			t.Errorf("[%d] unexpected error: %#v", caseNum, err)
		}
		if err == nil && item.IsError {
			t.Errorf("[%d] expected error, got nil", caseNum)
		}
		if !reflect.DeepEqual(item.Pins, pins) {
			t.Errorf("[%d] wrong result, \nexpected %#v, \ngot %#v", caseNum, item.Pins, pins)
		}
	}
}
