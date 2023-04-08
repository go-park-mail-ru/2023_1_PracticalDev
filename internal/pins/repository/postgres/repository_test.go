package postgres

import (
	"fmt"
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
			params: _pins.CreateParams{
				Title:       "t1",
				MediaSource: models.Image{},
				Description: "d1",
				Author:      12,
			},
			pin: models.Pin{Id: 1, Title: "t1", MediaSource: "ms_url", Description: "d1", Author: 12},
			err: nil,
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

			pin, err := repo.CreatePin(&test.params)
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
