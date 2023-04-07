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
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/profile"
)

func TestGetProfileByUser(t *testing.T) {
	type fields struct {
		mock sqlmock.Sqlmock
	}

	type testCase struct {
		prepare func(f *fields)
		id      int
		profile profile.Profile
		err     error
	}

	const getCmd = `SELECT username, name, profile_image, website_url 
					FROM users 
					WHERE id = $1;`

	tests := map[string]testCase{
		"good query": {
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"username", "name", "profile_image", "website_url"})
				rows = rows.AddRow("un1", "n1", "pi1", "wu1")
				f.mock.
					ExpectQuery(regexp.QuoteMeta(getCmd)).
					WithArgs(3).
					WillReturnRows(rows)
			},
			id:      3,
			profile: profile.Profile{Username: "un1", Name: "n1", ProfileImage: "pi1", WebsiteUrl: "wu1"},
			err:     nil,
		},
		"query error": {
			prepare: func(f *fields) {
				f.mock.
					ExpectQuery(regexp.QuoteMeta(getCmd)).
					WithArgs(3).
					WillReturnError(fmt.Errorf("db error"))
			},
			id:      3,
			profile: profile.Profile{},
			err:     profile.ErrDb,
		},
		"row scan error": {
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"username", "name"}).AddRow("un1", "n1")
				f.mock.
					ExpectQuery(regexp.QuoteMeta(getCmd)).
					WithArgs(3).
					WillReturnRows(rows)
			},
			id:      3,
			profile: profile.Profile{},
			err:     profile.ErrDb,
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

			repo := NewPostgresRepository(db, s3Serv, logger)

			f := fields{mock: sqlMock}
			if test.prepare != nil {
				test.prepare(&f)
			}

			prof, err := repo.GetProfileByUser(test.id)
			if err != test.err {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if prof != test.profile {
				t.Errorf("\nExpected: %v\nGot: %v", test.profile, prof)
			}
			if err = sqlMock.ExpectationsWereMet(); err != nil {
				t.Errorf("\nThere were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestFullUpdate(t *testing.T) {
	type fields struct {
		mock   sqlmock.Sqlmock
		s3mock *mocks.MockService
	}

	type testCase struct {
		prepare func(f *fields)
		params  profile.FullUpdateParams
		profile profile.Profile
		err     error
	}

	const fullUpdateCmd = `UPDATE users
							SET username = $1::VARCHAR,
							name = $2::VARCHAR,
							profile_image = $3::VARCHAR,
							website_url = $4::VARCHAR
							WHERE id = $5
							RETURNING username, name, profile_image, website_url;`

	tests := map[string]testCase{
		"good query": {
			prepare: func(f *fields) {
				f.s3mock.EXPECT().UploadImage(&models.Image{}).Return("pi_url", nil)

				rows := sqlmock.NewRows([]string{"username", "name", "profile_image", "website_url"})
				rows = rows.AddRow("un1", "n1", "pi_url", "wu1")
				f.mock.
					ExpectQuery(regexp.QuoteMeta(fullUpdateCmd)).
					WithArgs("un1", "n1", "pi_url", "wu1", 3).
					WillReturnRows(rows)
			},
			params: profile.FullUpdateParams{
				Id:           3,
				Username:     "un1",
				Name:         "n1",
				ProfileImage: models.Image{},
				WebsiteUrl:   "wu1",
			},
			profile: profile.Profile{Username: "un1", Name: "n1", ProfileImage: "pi_url", WebsiteUrl: "wu1"},
			err:     nil,
		},
		"query error": {
			prepare: func(f *fields) {
				f.s3mock.EXPECT().UploadImage(&models.Image{}).Return("pi_url", nil)

				f.mock.
					ExpectQuery(regexp.QuoteMeta(fullUpdateCmd)).
					WithArgs("un1", "n1", "pi_url", "wu1", 3).
					WillReturnError(fmt.Errorf("db error"))
			},
			params: profile.FullUpdateParams{
				Id:           3,
				Username:     "un1",
				Name:         "n1",
				ProfileImage: models.Image{},
				WebsiteUrl:   "wu1",
			},
			profile: profile.Profile{},
			err:     profile.ErrDb,
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

			repo := NewPostgresRepository(db, s3Serv, logger)

			f := fields{mock: sqlMock, s3mock: s3Serv}
			if test.prepare != nil {
				test.prepare(&f)
			}

			prof, err := repo.FullUpdate(&test.params)
			if err != test.err {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if prof != test.profile {
				t.Errorf("\nExpected: %v\nGot: %v", test.profile, prof)
			}
			if err = sqlMock.ExpectationsWereMet(); err != nil {
				t.Errorf("\nThere were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestPartialUpdate(t *testing.T) {
	type fields struct {
		mock   sqlmock.Sqlmock
		s3mock *mocks.MockService
	}

	type testCase struct {
		prepare func(f *fields)
		params  profile.PartialUpdateParams
		profile profile.Profile
		err     error
	}

	const partialUpdateCmd = `UPDATE users
								SET username = CASE WHEN $1::BOOLEAN THEN $2::VARCHAR ELSE username END,
								name = CASE WHEN $3::BOOLEAN THEN $4::VARCHAR ELSE name END,
    							profile_image = CASE WHEN $5::BOOLEAN THEN $6::VARCHAR ELSE profile_image END,
    							website_url = CASE WHEN $7::BOOLEAN THEN $8::VARCHAR ELSE website_url END
								WHERE id = $9
								RETURNING username, name, profile_image, website_url;`

	tests := map[string]testCase{
		"good query": {
			prepare: func(f *fields) {
				f.s3mock.EXPECT().UploadImage(&models.Image{}).Return("pi_url", nil)

				rows := sqlmock.NewRows([]string{"username", "name", "profile_image", "website_url"})
				rows = rows.AddRow("un1", "n1", "pi_url", "wu1")
				f.mock.
					ExpectQuery(regexp.QuoteMeta(partialUpdateCmd)).
					WithArgs(true, "un1", true, "n1", true, "pi_url", true, "wu1", 3).
					WillReturnRows(rows)
			},
			params: profile.PartialUpdateParams{
				Id:                 3,
				Username:           "un1",
				UpdateUsername:     true,
				Name:               "n1",
				UpdateName:         true,
				ProfileImage:       models.Image{},
				UpdateProfileImage: true,
				WebsiteUrl:         "wu1",
				UpdateWebsiteUrl:   true,
			},
			profile: profile.Profile{Username: "un1", Name: "n1", ProfileImage: "pi_url", WebsiteUrl: "wu1"},
			err:     nil,
		},
		"query error": {
			prepare: func(f *fields) {
				f.s3mock.EXPECT().UploadImage(&models.Image{}).Return("pi_url", nil)

				f.mock.
					ExpectQuery(regexp.QuoteMeta(partialUpdateCmd)).
					WithArgs(true, "un1", true, "n1", true, "pi_url", true, "wu1", 3).
					WillReturnError(fmt.Errorf("db error"))
			},
			params: profile.PartialUpdateParams{
				Id:                 3,
				Username:           "un1",
				UpdateUsername:     true,
				Name:               "n1",
				UpdateName:         true,
				ProfileImage:       models.Image{},
				UpdateProfileImage: true,
				WebsiteUrl:         "wu1",
				UpdateWebsiteUrl:   true,
			},
			profile: profile.Profile{},
			err:     profile.ErrDb,
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

			repo := NewPostgresRepository(db, s3Serv, logger)

			f := fields{mock: sqlMock, s3mock: s3Serv}
			if test.prepare != nil {
				test.prepare(&f)
			}

			prof, err := repo.PartialUpdate(&test.params)
			if err != test.err {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if prof != test.profile {
				t.Errorf("\nExpected: %v\nGot: %v", test.profile, prof)
			}
			if err = sqlMock.ExpectationsWereMet(); err != nil {
				t.Errorf("\nThere were unfulfilled expectations: %s", err)
			}
		})
	}
}
