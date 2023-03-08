package posts

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/db"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	"os"
	"reflect"
	"testing"
)

type ServiceTestCase struct {
	Page    int
	Limit   int
	Pins    []models.Pin
	IsError bool
}

func TestServiceGetPosts(t *testing.T) {
	cases := []ServiceTestCase{
		{
			Page:  1,
			Limit: 30,
			Pins: []models.Pin{
				{
					Id:          1,
					Link:        "",
					Title:       "Road",
					Description: "",
					MediaSource: "https://wg.grechkogv.ru/assets/pet7.webp",
					BoardId:     1,
				},
				{
					Id:          2,
					Link:        "",
					Title:       "Ice",
					Description: "",
					MediaSource: "https://wg.grechkogv.ru/assets/armorChest4.webp",
					BoardId:     1,
				},
				{
					Id:          3,
					Link:        "",
					Title:       "Future",
					Description: "",
					MediaSource: "https://wg.grechkogv.ru/assets/pet6.webp",
					BoardId:     1,
				},
				{
					Id:          4,
					Link:        "",
					Title:       "Color",
					Description: "",
					MediaSource: "https://wg.grechkogv.ru/assets/pet8.webp",
					BoardId:     2,
				},
				{
					Id:          5,
					Link:        "",
					Title:       "Shops",
					Description: "",
					MediaSource: "https://i.pinimg.com/564x/2f/93/56/2f9356b9346e82c14bf286c6a107bc7a.jpg",
					BoardId:     3,
				},
				{
					Id:          6,
					Link:        "",
					Title:       "Shops",
					Description: "",
					MediaSource: "https://i.pinimg.com/564x/32/ff/71/32ff717c3cd3bd3d1886c775b59f0769.jpg",
					BoardId:     3,
				},
				{
					Id:          7,
					Link:        "",
					Title:       "Shops",
					Description: "",
					MediaSource: "https://i.pinimg.com/564x/ce/e3/01/cee3011f3e19de4377dbf98f397c027b.jpg",
					BoardId:     3,
				},
				{
					Id:          8,
					Link:        "",
					Title:       "Shops",
					Description: "",
					MediaSource: "https://i.pinimg.com/564x/a6/ba/55/a6ba553df2a0c0f3894ef328a86fb373.jpg",
					BoardId:     3,
				},
				{
					Id:          9,
					Link:        "",
					Title:       "Shops",
					Description: "",
					MediaSource: "https://i.pinimg.com/564x/43/2d/3b/432d3b28d1661439245422e9218ffcce.jpg",
					BoardId:     3,
				},
				{
					Id:          10,
					Link:        "",
					Title:       "School",
					Description: "",
					MediaSource: "https://i.pinimg.com/564x/98/9d/3f/989d3f5c158dcac7ca4d115bff866d84.jpg",
					BoardId:     4,
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
					Link:        "",
					Title:       "Color",
					Description: "",
					MediaSource: "https://wg.grechkogv.ru/assets/pet8.webp",
					BoardId:     2,
				},
				{
					Id:          5,
					Link:        "",
					Title:       "Shops",
					Description: "",
					MediaSource: "https://i.pinimg.com/564x/2f/93/56/2f9356b9346e82c14bf286c6a107bc7a.jpg",
					BoardId:     3,
				},
				{
					Id:          6,
					Link:        "",
					Title:       "Shops",
					Description: "",
					MediaSource: "https://i.pinimg.com/564x/32/ff/71/32ff717c3cd3bd3d1886c775b59f0769.jpg",
					BoardId:     3,
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
