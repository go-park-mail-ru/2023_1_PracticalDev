package posts

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/db"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	"os"
	"reflect"
	"testing"
)

type RepositoryTestCase struct {
	Limit   int
	Offset  int
	Pins    []models.Pin
	IsError bool
}

func TestRepositoryGetPosts(t *testing.T) {
	cases := []RepositoryTestCase{
		{
			Limit:  30,
			Offset: 0,
			Pins: []models.Pin{
				{
					Id:          1,
					Link:        "",
					Title:       "Road",
					Description: "",
					MediaSource: "",
					BoardId:     1,
				},
				{
					Id:          2,
					Link:        "",
					Title:       "Ice",
					Description: "",
					MediaSource: "",
					BoardId:     1,
				},
				{
					Id:          3,
					Link:        "",
					Title:       "Future",
					Description: "",
					MediaSource: "",
					BoardId:     1,
				},
				{
					Id:          4,
					Link:        "",
					Title:       "Color",
					Description: "",
					MediaSource: "",
					BoardId:     2,
				},
				{
					Id:          5,
					Link:        "",
					Title:       "Question",
					Description: "",
					MediaSource: "",
					BoardId:     2,
				},
				{
					Id:          6,
					Link:        "",
					Title:       "Shops",
					Description: "",
					MediaSource: "",
					BoardId:     3,
				},
				{
					Id:          7,
					Link:        "",
					Title:       "School",
					Description: "",
					MediaSource: "",
					BoardId:     4,
				},
			},
			IsError: false,
		},
		{
			Limit:  3,
			Offset: 3,
			Pins: []models.Pin{
				{
					Id:          4,
					Link:        "",
					Title:       "Color",
					Description: "",
					MediaSource: "",
					BoardId:     2,
				},
				{
					Id:          5,
					Link:        "",
					Title:       "Question",
					Description: "",
					MediaSource: "",
					BoardId:     2,
				},
				{
					Id:          6,
					Link:        "",
					Title:       "Shops",
					Description: "",
					MediaSource: "",
					BoardId:     3,
				},
			},
			IsError: false,
		},
		{
			Limit:   0,
			Offset:  3,
			Pins:    []models.Pin{},
			IsError: false,
		},
		{
			Limit:   30,
			Offset:  30,
			Pins:    []models.Pin{},
			IsError: false,
		},
		{
			Limit:   -1,
			Offset:  3,
			Pins:    []models.Pin{},
			IsError: true,
		},
		{
			Limit:   30,
			Offset:  -1,
			Pins:    []models.Pin{},
			IsError: true,
		},
	}

	logger := log.New()
	db, err := db.New(logger)
	if err != nil {
		os.Exit(1)
	}
	rep := NewRepository(db, logger)

	var pins []models.Pin
	for caseNum, item := range cases {
		pins, err = rep.GetPosts(item.Limit, item.Offset)

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
