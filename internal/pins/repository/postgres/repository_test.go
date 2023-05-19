package postgres

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"log"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/images/client/mocks"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	_pins "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pins"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
)

var err error
var logger *zap.Logger

func init() {
	logger, err = zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
}

func TestCreate(t *testing.T) {
	type fields struct {
		mock   sqlmock.Sqlmock
		s3mock *mocks.MockImageClient
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
				f.s3mock.EXPECT().UploadImage(context.Background(), &models.Image{}).Return("ms_url", nil)

				rows := sqlmock.NewRows([]string{"id", "title", "media_source", "media_source_color", "description",
					"author_id"})
				rows = rows.AddRow(1, "t1", "ms_url", "rgb(39, 102, 120)", "d1", 12)
				f.mock.
					ExpectQuery(regexp.QuoteMeta(createCmd)).
					WithArgs("t1", "ms_url", "rgb(39, 102, 120)", "d1", 12).
					WillReturnRows(rows)
			},
			params: _pins.CreateParams{Title: "t1", MediaSource: models.Image{}, Description: "d1", Author: 12},
			pin: models.Pin{Id: 1, Title: "t1", MediaSource: "ms_url", MediaSourceColor: "rgb(39, 102, 120)",
				Description: "d1", Author: 12},
			err: nil,
		},
		"query error": {
			prepare: func(f *fields) {
				f.s3mock.EXPECT().UploadImage(context.Background(), &models.Image{}).Return("ms_url", nil)

				f.mock.
					ExpectQuery(regexp.QuoteMeta(createCmd)).
					WithArgs("t1", "ms_url", "rgb(39, 102, 120)", "d1", 12).
					WillReturnError(fmt.Errorf("sql error"))
			},
			params: _pins.CreateParams{
				Title:       "t1",
				MediaSource: models.Image{},
				Description: "d1",
				Author:      12,
			},
			pin: models.Pin{},
			err: pkgErrors.ErrDb,
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

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s3Serv := mocks.NewMockImageClient(ctrl)

			repo := NewRepository(db, s3Serv, logger)

			f := fields{mock: sqlMock, s3mock: s3Serv}
			if test.prepare != nil {
				test.prepare(&f)
			}

			pin, err := repo.Create(&test.params)
			if !errors.Is(err, test.err) {
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
				rows := sqlmock.NewRows([]string{"id", "title", "description", "media_source", "media_source_color",
					"n_likes", "author_id"})
				rows = rows.AddRow(1, "t1", "d1", "ms_url1", "rgb(39, 102, 120)", 0, 12)
				rows = rows.AddRow(2, "t2", "d2", "ms_url2", "rgb(39, 102, 120)", 2, 3)
				rows = rows.AddRow(3, "t3", "d3", "ms_url3", "rgb(39, 102, 120)", 3, 10)
				f.mock.
					ExpectQuery(regexp.QuoteMeta(listCmd)).
					WithArgs(30, 0).
					WillReturnRows(rows)
			},
			page:  1,
			limit: 30,
			pins: []models.Pin{
				{Id: 1, Title: "t1", MediaSource: "ms_url1", MediaSourceColor: "rgb(39, 102, 120)", Description: "d1",
					NumLikes: 0, Author: 12},
				{Id: 2, Title: "t2", MediaSource: "ms_url2", MediaSourceColor: "rgb(39, 102, 120)", Description: "d2",
					NumLikes: 2, Author: 3},
				{Id: 3, Title: "t3", MediaSource: "ms_url3", MediaSourceColor: "rgb(39, 102, 120)", Description: "d3",
					NumLikes: 3, Author: 10},
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
			err:   pkgErrors.ErrDb,
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
			err:   pkgErrors.ErrDb,
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

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s3Serv := mocks.NewMockImageClient(ctrl)

			repo := NewRepository(db, s3Serv, logger)

			f := fields{mock: sqlMock}
			if test.prepare != nil {
				test.prepare(&f)
			}

			pins, err := repo.List(test.page, test.limit)
			if !errors.Is(err, test.err) {
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
				rows := sqlmock.NewRows([]string{"id", "title", "description", "media_source", "media_source_color",
					"n_likes", "author_id"})
				rows = rows.AddRow(1, "t1", "d1", "ms_url1", "rgb(39, 102, 120)", 0, 12)
				rows = rows.AddRow(2, "t2", "d2", "ms_url2", "rgb(39, 102, 120)", 2, 12)
				rows = rows.AddRow(3, "t3", "d3", "ms_url3", "rgb(39, 102, 120)", 3, 12)
				f.mock.
					ExpectQuery(regexp.QuoteMeta(listByUserCmd)).
					WithArgs(12, 30, 0).
					WillReturnRows(rows)
			},
			userId: 12,
			page:   1,
			limit:  30,
			pins: []models.Pin{
				{Id: 1, Title: "t1", MediaSource: "ms_url1", MediaSourceColor: "rgb(39, 102, 120)", Description: "d1",
					NumLikes: 0, Author: 12},
				{Id: 2, Title: "t2", MediaSource: "ms_url2", MediaSourceColor: "rgb(39, 102, 120)", Description: "d2",
					NumLikes: 2, Author: 12},
				{Id: 3, Title: "t3", MediaSource: "ms_url3", MediaSourceColor: "rgb(39, 102, 120)", Description: "d3",
					NumLikes: 3, Author: 12},
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
			err:    pkgErrors.ErrDb,
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
			err:    pkgErrors.ErrDb,
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

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s3Serv := mocks.NewMockImageClient(ctrl)

			repo := NewRepository(db, s3Serv, logger)

			f := fields{mock: sqlMock}
			if test.prepare != nil {
				test.prepare(&f)
			}

			pins, err := repo.ListByAuthor(test.userId, test.page, test.limit)
			if !errors.Is(err, test.err) {
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
				rows := sqlmock.NewRows([]string{"id", "title", "description", "media_source", "media_source_color",
					"n_likes", "author_id"})
				rows = rows.AddRow(3, "t1", "d1", "ms_url1", "rgb(39, 102, 120)", 3, 12)
				f.mock.
					ExpectQuery(regexp.QuoteMeta(getCmd)).
					WithArgs(3).
					WillReturnRows(rows)
			},
			id: 3,
			pin: models.Pin{Id: 3, Title: "t1", MediaSource: "ms_url1", MediaSourceColor: "rgb(39, 102, 120)",
				Description: "d1", NumLikes: 3, Author: 12},
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
			err: pkgErrors.ErrDb,
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
			err: pkgErrors.ErrDb,
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

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s3Serv := mocks.NewMockImageClient(ctrl)

			repo := NewRepository(db, s3Serv, logger)

			f := fields{mock: sqlMock}
			if test.prepare != nil {
				test.prepare(&f)
			}

			pin, err := repo.Get(test.id)
			if !errors.Is(err, test.err) {
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
