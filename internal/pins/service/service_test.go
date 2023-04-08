package service

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	_pins "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pins"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pins/mocks"
)

func TestCreate(t *testing.T) {
	type fields struct {
		repo *mocks.MockRepository
	}

	type testCase struct {
		prepare func(f *fields)
		params  _pins.CreateParams
		pin     models.Pin
		err     error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Create(&_pins.CreateParams{
					Title:       "t1",
					MediaSource: models.Image{},
					Description: "d1",
					Author:      12,
				}).Return(models.Pin{Id: 1,
					Title:       "t1",
					MediaSource: "ms_url",
					Description: "d1",
					Author:      12,
				}, nil)
			},
			params: _pins.CreateParams{Title: "t1", MediaSource: models.Image{}, Description: "d1", Author: 12},
			pin:    models.Pin{Id: 1, Title: "t1", MediaSource: "ms_url", Description: "d1", Author: 12},
			err:    nil,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl)}
			if test.prepare != nil {
				test.prepare(&f)
			}

			serv := NewService(f.repo)

			pin, err := serv.Create(&test.params)
			if err != test.err {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if pin != test.pin {
				t.Errorf("\nExpected: %v\nGot: %v", test.pin, pin)
			}
		})
	}
}

func TestList(t *testing.T) {
	type fields struct {
		repo *mocks.MockRepository
	}

	type testCase struct {
		prepare func(f *fields)
		page    int
		limit   int
		pins    []models.Pin
		err     error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				f.repo.EXPECT().List(1, 30).Return([]models.Pin{
					{Id: 1, Title: "t1", MediaSource: "ms_url1", Description: "d1", Author: 12},
					{Id: 2, Title: "t2", MediaSource: "ms_url2", Description: "d2", Author: 3},
					{Id: 3, Title: "t3", MediaSource: "ms_url3", Description: "d3", Author: 10},
				}, nil)
			},
			page:  1,
			limit: 30,
			pins: []models.Pin{
				{Id: 1, Title: "t1", MediaSource: "ms_url1", Description: "d1", Author: 12},
				{Id: 2, Title: "t2", MediaSource: "ms_url2", Description: "d2", Author: 3},
				{Id: 3, Title: "t3", MediaSource: "ms_url3", Description: "d3", Author: 10},
			},
			err: nil,
		},
		"no boards": {
			prepare: func(f *fields) {
				f.repo.EXPECT().List(1, 30).Return([]models.Pin{}, nil)
			},
			page:  1,
			limit: 30,
			pins:  []models.Pin{},
			err:   nil,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl)}
			if test.prepare != nil {
				test.prepare(&f)
			}

			serv := NewService(f.repo)

			pins, err := serv.List(test.page, test.limit)
			if err != test.err {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if !reflect.DeepEqual(pins, test.pins) {
				t.Errorf("\nExpected: %v\nGot: %v", test.pins, pins)
			}
		})
	}
}
