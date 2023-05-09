package postgres

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"reflect"
	"regexp"
	"testing"
)

func TestRepository_Get(t *testing.T) {
	type fields struct {
		mock sqlmock.Sqlmock
		id   int
		user *models.User
	}

	type testCase struct {
		prepare func(f *fields)
		id      int
		user    models.User
		err     error
	}

	tests := map[string]testCase{
		"good query": {
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"id", "username", "email", "hashed_password", "name", "profile_image",
					"website_url", "account_type"})
				rows = rows.AddRow(f.user.Id, f.user.Username, f.user.Email, f.user.HashedPassword, f.user.Name,
					f.user.ProfileImage, f.user.WebsiteUrl, f.user.AccountType)
				f.mock.
					ExpectQuery(regexp.QuoteMeta(getCmd)).
					WithArgs(f.id).
					WillReturnRows(rows)
			},
			id: 3,
			user: models.User{Id: 3, Username: "petya", Email: "petya@vk.com", Name: "Petya",
				ProfileImage: "petya.jpg", WebsiteUrl: "petya.com", AccountType: "personal"},
			err: nil,
		},
		"query error": {
			prepare: func(f *fields) {
				f.mock.
					ExpectQuery(regexp.QuoteMeta(getCmd)).
					WithArgs(f.id).
					WillReturnError(fmt.Errorf("db error"))
			},
			id:   3,
			user: models.User{},
			err:  pkgErrors.ErrDb,
		},
		"row scan error": {
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"username", "name"}).AddRow("petya", "Petya")
				f.mock.
					ExpectQuery(regexp.QuoteMeta(getCmd)).
					WithArgs(f.id).
					WillReturnRows(rows)
			},
			id:   3,
			user: models.User{},
			err:  pkgErrors.ErrDb,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			logger, err := zap.NewDevelopment()
			if err != nil {
				t.Fatalf("can't create logger: %s", err)
			}

			db, sqlMock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("can't create mock: %s", err)
			}
			defer db.Close()

			f := fields{mock: sqlMock, id: test.id, user: &test.user}
			if test.prepare != nil {
				test.prepare(&f)
			}

			repo := NewRepository(db, logger)
			user, err := repo.Get(test.id)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if !reflect.DeepEqual(user, test.user) {
				t.Errorf("\nExpected: %v\nGot: %v", test.user, user)
			}
			if err = sqlMock.ExpectationsWereMet(); err != nil {
				t.Errorf("\nThere were unfulfilled expectations: %s", err)
			}
		})
	}
}
