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

func TestListByUser(t *testing.T) {
	type fields struct {
		repo *mocks.MockRepository
	}

	type testCase struct {
		prepare func(f *fields)
		page    int
		limit   int
		userId  int
		pins    []models.Pin
		err     error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				f.repo.EXPECT().ListByUser(12, 1, 30).Return([]models.Pin{
					{Id: 1, Title: "t1", MediaSource: "ms_url1", Description: "d1", Author: 12},
					{Id: 2, Title: "t2", MediaSource: "ms_url2", Description: "d2", Author: 12},
					{Id: 3, Title: "t3", MediaSource: "ms_url3", Description: "d3", Author: 12},
				}, nil)
			},
			page:   1,
			limit:  30,
			userId: 12,
			pins: []models.Pin{
				{Id: 1, Title: "t1", MediaSource: "ms_url1", Description: "d1", Author: 12},
				{Id: 2, Title: "t2", MediaSource: "ms_url2", Description: "d2", Author: 12},
				{Id: 3, Title: "t3", MediaSource: "ms_url3", Description: "d3", Author: 12},
			},
			err: nil,
		},
		"no boards": {
			prepare: func(f *fields) {
				f.repo.EXPECT().ListByUser(12, 1, 30).Return([]models.Pin{}, nil)
			},
			userId: 12,
			page:   1,
			limit:  30,
			pins:   []models.Pin{},
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

			pins, err := serv.ListByUser(test.userId, test.page, test.limit)
			if err != test.err {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if !reflect.DeepEqual(pins, test.pins) {
				t.Errorf("\nExpected: %v\nGot: %v", test.pins, pins)
			}
		})
	}
}

func TestGet(t *testing.T) {
	type fields struct {
		repo *mocks.MockRepository
	}

	type testCase struct {
		prepare func(f *fields)
		id      int
		pin     models.Pin
		err     error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Get(3).Return(models.Pin{
					Id:          3,
					Title:       "t1",
					MediaSource: "ms_url1",
					Description: "d1",
					Author:      12,
				}, nil)
			},
			id:  3,
			pin: models.Pin{Id: 3, Title: "t1", MediaSource: "ms_url1", Description: "d1", Author: 12},
			err: nil,
		},
		"pin not found": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Get(3).Return(models.Pin{}, _pins.ErrPinNotFound)
			},
			id:  3,
			pin: models.Pin{},
			err: _pins.ErrPinNotFound,
		},
		"negative pin id param": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Get(-1).Return(models.Pin{}, _pins.ErrPinNotFound)
			},
			id:  -1,
			pin: models.Pin{},
			err: _pins.ErrPinNotFound,
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

			pin, err := serv.Get(test.id)
			if err != test.err {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if pin != test.pin {
				t.Errorf("\nExpected: %v\nGot: %v", test.pin, pin)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	type fields struct {
		repo *mocks.MockRepository
	}

	type testCase struct {
		prepare func(f *fields)
		id      int
		err     error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Delete(3).Return(nil)
			},
			id:  3,
			err: nil,
		},
		"pin not found": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Delete(3).Return(_pins.ErrPinNotFound)
			},
			id:  3,
			err: _pins.ErrPinNotFound,
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

			err := serv.Delete(test.id)
			if err != test.err {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
		})
	}
}

func TestCheckWriteAccess(t *testing.T) {
	type fields struct {
		repo *mocks.MockRepository
	}

	type testCase struct {
		prepare func(f *fields)
		userId  string
		pinId   string
		access  bool
		err     error
	}

	tests := map[string]testCase{
		"access is allowed": {
			prepare: func(f *fields) {
				f.repo.EXPECT().CheckWriteAccess("2", "3").Return(true, nil)
			},
			userId: "2",
			pinId:  "3",
			access: true,
			err:    nil,
		},
		"access is denied": {
			prepare: func(f *fields) {
				f.repo.EXPECT().CheckWriteAccess("2", "3").Return(false, nil)
			},
			userId: "2",
			pinId:  "3",
			access: false,
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

			access, err := serv.CheckWriteAccess(test.userId, test.pinId)
			if err != test.err {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if access != test.access {
				t.Errorf("\nExpected: %v\nGot: %v", test.access, access)
			}
		})
	}
}

func TestCheckReadAccess(t *testing.T) {
	type fields struct {
		repo *mocks.MockRepository
	}

	type testCase struct {
		prepare func(f *fields)
		userId  string
		pinId   string
		access  bool
		err     error
	}

	tests := map[string]testCase{
		"access is allowed": {
			prepare: func(f *fields) {
				f.repo.EXPECT().CheckReadAccess("2", "3").Return(true, nil)
			},
			userId: "2",
			pinId:  "3",
			access: true,
			err:    nil,
		},
		"access is denied": {
			prepare: func(f *fields) {
				f.repo.EXPECT().CheckReadAccess("2", "3").Return(false, nil)
			},
			userId: "2",
			pinId:  "3",
			access: false,
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

			access, err := serv.CheckReadAccess(test.userId, test.pinId)
			if err != test.err {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if access != test.access {
				t.Errorf("\nExpected: %v\nGot: %v", test.access, access)
			}
		})
	}
}
