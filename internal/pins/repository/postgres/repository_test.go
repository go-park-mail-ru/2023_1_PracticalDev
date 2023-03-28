package postgres

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/images/mocks"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	_pins "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pins"
)

func TestCreate(t *testing.T) {
	type fields struct {
		mock   sqlmock.Sqlmock
		s3mock *mocks.MockService
	}

	type testCase struct {
		prepare func(f *fields)
		params  _pins.CreateParams
		pin     models.Pin
		err     error
	}

	tests := map[string]testCase{
		"good query": {
			prepare: func(f *fields) {
				f.s3mock.EXPECT().UploadImage(&models.Image{}).Return("ms_url", nil)

				rows := sqlmock.NewRows([]string{"id", "title", "media_source", "description", "author_id"})
				rows = rows.AddRow(1, "t1", "ms_url", "d1", 12)
				f.mock.
					ExpectQuery(regexp.QuoteMeta(createCmd)).
					WithArgs("t1", "ms_url", "d1", 12).
					WillReturnRows(rows)
			},
			params: _pins.CreateParams{Title: "t1", MediaSource: models.Image{}, Description: "d1", Author: 12},
			pin:    models.Pin{Id: 1, Title: "t1", MediaSource: "ms_url", Description: "d1", Author: 12},
			err:    nil,
		},
		"query error": {
			prepare: func(f *fields) {
				f.s3mock.EXPECT().UploadImage(&models.Image{}).Return("ms_url", nil)

				f.mock.
					ExpectQuery(regexp.QuoteMeta(createCmd)).
					WithArgs("t1", "ms_url", "d1", 12).
					WillReturnError(fmt.Errorf("sql error"))
			},
			params: _pins.CreateParams{
				Title:       "t1",
				MediaSource: models.Image{},
				Description: "d1",
				Author:      12,
			},
			pin: models.Pin{},
			err: _pins.ErrDb,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db, sqlMock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("can't create mock: %s", err)
			}
			defer db.Close()

			logger := log.New()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s3Serv := mocks.NewMockService(ctrl)

			repo := NewRepository(db, s3Serv, logger)

			f := fields{mock: sqlMock, s3mock: s3Serv}
			if test.prepare != nil {
				test.prepare(&f)
			}

			pin, err := repo.Create(&test.params)
			if err != test.err {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if pin != test.pin {
				t.Errorf("\nExpected: %v\nGot: %v", test.pin, pin)
			}
			if err = sqlMock.ExpectationsWereMet(); err != nil {
				t.Errorf("\nThere were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestList(t *testing.T) {
	type fields struct {
		mock sqlmock.Sqlmock
	}

	type testCase struct {
		prepare func(f *fields)
		page    int
		limit   int
		pins    []models.Pin
		err     error
	}

	tests := map[string]testCase{
		"good query": {
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "media_source", "author_id"})
				rows = rows.AddRow(1, "t1", "d1", "ms_url1", 12)
				rows = rows.AddRow(2, "t2", "d2", "ms_url2", 3)
				rows = rows.AddRow(3, "t3", "d3", "ms_url3", 10)
				f.mock.
					ExpectQuery(regexp.QuoteMeta(listCmd)).
					WithArgs(30, 0).
					WillReturnRows(rows)
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
		"query error": {
			prepare: func(f *fields) {
				f.mock.
					ExpectQuery(regexp.QuoteMeta(listCmd)).
					WithArgs(30, 0).
					WillReturnError(fmt.Errorf("sql error"))
			},
			page:  1,
			limit: 30,
			pins:  nil,
			err:   _pins.ErrDb,
		},
		"row scan error": {
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"id", "title"}).AddRow(1, "t1")
				f.mock.
					ExpectQuery(regexp.QuoteMeta(listCmd)).
					WithArgs(30, 0).
					WillReturnRows(rows)
			},
			page:  1,
			limit: 30,
			pins:  nil,
			err:   _pins.ErrDb,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db, sqlMock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("can't create mock: %s", err)
			}
			defer db.Close()

			logger := log.New()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s3Serv := mocks.NewMockService(ctrl)

			repo := NewRepository(db, s3Serv, logger)

			f := fields{mock: sqlMock}
			if test.prepare != nil {
				test.prepare(&f)
			}

			pins, err := repo.List(test.page, test.limit)
			if err != test.err {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if !reflect.DeepEqual(pins, test.pins) {
				t.Errorf("\nExpected: %v\nGot: %v", test.pins, pins)
			}
			if err = sqlMock.ExpectationsWereMet(); err != nil {
				t.Errorf("\nThere were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestListByUser(t *testing.T) {
	type fields struct {
		mock sqlmock.Sqlmock
	}

	type testCase struct {
		prepare func(f *fields)
		userId  int
		page    int
		limit   int
		pins    []models.Pin
		err     error
	}

	tests := map[string]testCase{
		"good query": {
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "media_source", "author_id"})
				rows = rows.AddRow(1, "t1", "d1", "ms_url1", 12)
				rows = rows.AddRow(2, "t2", "d2", "ms_url2", 12)
				rows = rows.AddRow(3, "t3", "d3", "ms_url3", 12)
				f.mock.
					ExpectQuery(regexp.QuoteMeta(listByUserCmd)).
					WithArgs(12, 30, 0).
					WillReturnRows(rows)
			},
			userId: 12,
			page:   1,
			limit:  30,
			pins: []models.Pin{
				{Id: 1, Title: "t1", MediaSource: "ms_url1", Description: "d1", Author: 12},
				{Id: 2, Title: "t2", MediaSource: "ms_url2", Description: "d2", Author: 12},
				{Id: 3, Title: "t3", MediaSource: "ms_url3", Description: "d3", Author: 12},
			},
			err: nil,
		},
		"query error": {
			prepare: func(f *fields) {
				f.mock.
					ExpectQuery(regexp.QuoteMeta(listByUserCmd)).
					WithArgs(12, 30, 0).
					WillReturnError(fmt.Errorf("sql error"))
			},
			userId: 12,
			page:   1,
			limit:  30,
			pins:   nil,
			err:    _pins.ErrDb,
		},
		"row scan error": {
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"id", "title"}).AddRow(1, "t1")
				f.mock.
					ExpectQuery(regexp.QuoteMeta(listByUserCmd)).
					WithArgs(12, 30, 0).
					WillReturnRows(rows)
			},
			userId: 12,
			page:   1,
			limit:  30,
			pins:   nil,
			err:    _pins.ErrDb,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db, sqlMock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("can't create mock: %s", err)
			}
			defer db.Close()

			logger := log.New()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s3Serv := mocks.NewMockService(ctrl)

			repo := NewRepository(db, s3Serv, logger)

			f := fields{mock: sqlMock}
			if test.prepare != nil {
				test.prepare(&f)
			}

			pins, err := repo.ListByUser(test.userId, test.page, test.limit)
			if err != test.err {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if !reflect.DeepEqual(pins, test.pins) {
				t.Errorf("\nExpected: %v\nGot: %v", test.pins, pins)
			}
			if err = sqlMock.ExpectationsWereMet(); err != nil {
				t.Errorf("\nThere were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestListByBoard(t *testing.T) {
	type fields struct {
		mock sqlmock.Sqlmock
	}

	type testCase struct {
		prepare func(f *fields)
		boardId int
		page    int
		limit   int
		pins    []models.Pin
		err     error
	}

	tests := map[string]testCase{
		"good query": {
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "media_source", "author_id"})
				rows = rows.AddRow(1, "t1", "d1", "ms_url1", 12)
				rows = rows.AddRow(2, "t2", "d2", "ms_url2", 12)
				rows = rows.AddRow(3, "t3", "d3", "ms_url3", 12)
				f.mock.
					ExpectQuery(regexp.QuoteMeta(listByBoardCmd)).
					WithArgs(3, 30, 0).
					WillReturnRows(rows)
			},
			boardId: 3,
			page:    1,
			limit:   30,
			pins: []models.Pin{
				{Id: 1, Title: "t1", MediaSource: "ms_url1", Description: "d1", Author: 12},
				{Id: 2, Title: "t2", MediaSource: "ms_url2", Description: "d2", Author: 12},
				{Id: 3, Title: "t3", MediaSource: "ms_url3", Description: "d3", Author: 12},
			},
			err: nil,
		},
		"query error": {
			prepare: func(f *fields) {
				f.mock.
					ExpectQuery(regexp.QuoteMeta(listByBoardCmd)).
					WithArgs(3, 30, 0).
					WillReturnError(fmt.Errorf("sql error"))
			},
			boardId: 3,
			page:    1,
			limit:   30,
			pins:    nil,
			err:     _pins.ErrDb,
		},
		"row scan error": {
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"id", "title"}).AddRow(1, "t1")
				f.mock.
					ExpectQuery(regexp.QuoteMeta(listByBoardCmd)).
					WithArgs(3, 30, 0).
					WillReturnRows(rows)
			},
			boardId: 3,
			page:    1,
			limit:   30,
			pins:    nil,
			err:     _pins.ErrDb,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db, sqlMock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("can't create mock: %s", err)
			}
			defer db.Close()

			logger := log.New()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s3Serv := mocks.NewMockService(ctrl)

			repo := NewRepository(db, s3Serv, logger)

			f := fields{mock: sqlMock}
			if test.prepare != nil {
				test.prepare(&f)
			}

			pins, err := repo.ListByBoard(test.boardId, test.page, test.limit)
			if err != test.err {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if !reflect.DeepEqual(pins, test.pins) {
				t.Errorf("\nExpected: %v\nGot: %v", test.pins, pins)
			}
			if err = sqlMock.ExpectationsWereMet(); err != nil {
				t.Errorf("\nThere were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGet(t *testing.T) {
	type fields struct {
		mock sqlmock.Sqlmock
	}

	type testCase struct {
		prepare func(f *fields)
		id      int
		pin     models.Pin
		err     error
	}

	tests := map[string]testCase{
		"good query": {
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "media_source", "author_id"})
				rows = rows.AddRow(3, "t1", "d1", "ms_url1", 12)
				f.mock.
					ExpectQuery(regexp.QuoteMeta(getCmd)).
					WithArgs(3).
					WillReturnRows(rows)
			},
			id:  3,
			pin: models.Pin{Id: 3, Title: "t1", MediaSource: "ms_url1", Description: "d1", Author: 12},
			err: nil,
		},
		"query error": {
			prepare: func(f *fields) {
				f.mock.
					ExpectQuery(regexp.QuoteMeta(getCmd)).
					WithArgs(3).
					WillReturnError(fmt.Errorf("sql error"))
			},
			id:  3,
			pin: models.Pin{},
			err: _pins.ErrDb,
		},
		"row scan error": {
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"id", "title"}).AddRow(1, "t1")
				f.mock.
					ExpectQuery(regexp.QuoteMeta(getCmd)).
					WithArgs(3).
					WillReturnRows(rows)
			},
			id:  3,
			pin: models.Pin{},
			err: _pins.ErrDb,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db, sqlMock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("can't create mock: %s", err)
			}
			defer db.Close()

			logger := log.New()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s3Serv := mocks.NewMockService(ctrl)

			repo := NewRepository(db, s3Serv, logger)

			f := fields{mock: sqlMock}
			if test.prepare != nil {
				test.prepare(&f)
			}

			pin, err := repo.Get(test.id)
			if err != test.err {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if !reflect.DeepEqual(pin, test.pin) {
				t.Errorf("\nExpected: %v\nGot: %v", test.pin, pin)
			}
			if err = sqlMock.ExpectationsWereMet(); err != nil {
				t.Errorf("\nThere were unfulfilled expectations: %s", err)
			}
		})
	}
}