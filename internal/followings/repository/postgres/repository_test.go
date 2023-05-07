package postgres

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/followings"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"reflect"
	"regexp"
	"testing"
)

func TestRepository_Create(t *testing.T) {
	type fields struct {
		mock       sqlmock.Sqlmock
		followerID int
		followeeID int
	}

	type testCase struct {
		prepare    func(f *fields)
		followerID int
		followeeID int
		err        error
	}

	tests := map[string]*testCase{
		"usual": {
			prepare: func(f *fields) {
				f.mock.
					ExpectExec(regexp.QuoteMeta(createCmd)).
					WithArgs(f.followerID, f.followeeID).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			followerID: 3,
			followeeID: 2,
			err:        nil,
		},
		"exec error": {
			prepare: func(f *fields) {
				f.mock.
					ExpectExec(regexp.QuoteMeta(createCmd)).
					WithArgs(f.followerID, f.followeeID).
					WillReturnError(fmt.Errorf("db error"))
			},
			followerID: 3,
			followeeID: 2,
			err:        pkgErrors.ErrDb,
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

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("can't create mock: %s", err)
			}
			defer db.Close()

			f := fields{mock: mock, followerID: test.followerID, followeeID: test.followeeID}
			if test.prepare != nil {
				test.prepare(&f)
			}

			repo := NewRepository(db, logger)
			err = repo.Create(test.followerID, test.followeeID)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("\nThere were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestRepository_Delete(t *testing.T) {
	type fields struct {
		mock       sqlmock.Sqlmock
		followerID int
		followeeID int
	}

	type testCase struct {
		prepare    func(f *fields)
		followerID int
		followeeID int
		err        error
	}

	tests := map[string]testCase{
		"good query": {
			prepare: func(f *fields) {
				f.mock.
					ExpectExec(regexp.QuoteMeta(deleteCmd)).
					WithArgs(f.followerID, f.followeeID).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			followerID: 3,
			followeeID: 2,
			err:        nil,
		},
		"exec error": {
			prepare: func(f *fields) {
				f.mock.
					ExpectExec(regexp.QuoteMeta(deleteCmd)).
					WithArgs(f.followerID, f.followeeID).
					WillReturnError(fmt.Errorf("db error"))
			},
			followerID: 3,
			followeeID: 2,
			err:        pkgErrors.ErrDb,
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

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("can't create mock: %s", err)
			}
			defer db.Close()

			f := fields{mock: mock, followerID: test.followerID, followeeID: test.followeeID}
			if test.prepare != nil {
				test.prepare(&f)
			}

			repo := NewRepository(db, logger)
			err = repo.Delete(test.followerID, test.followeeID)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("\nThere were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestRepository_GetFollowees(t *testing.T) {
	type fields struct {
		mock      sqlmock.Sqlmock
		userID    int
		followees []followings.Followee
	}

	type testCase struct {
		prepare   func(f *fields)
		userID    int
		followees []followings.Followee
		err       error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"id", "username", "name", "profile_image", "website_url"})
				for _, followee := range f.followees {
					rows = rows.AddRow(followee.Id, followee.Username, followee.Name, followee.ProfileImage,
						followee.WebsiteUrl)
				}
				f.mock.
					ExpectQuery(regexp.QuoteMeta(getFolloweesCmd)).
					WithArgs(f.userID).
					WillReturnRows(rows)
			},
			userID: 12,
			followees: []followings.Followee{
				{Id: 2, Username: "vasua", Name: "Vasya", ProfileImage: "vasya.jpg", WebsiteUrl: "vasya.com"},
				{Id: 3, Username: "kolya", Name: "Kolya", ProfileImage: "kolya.jpg", WebsiteUrl: "kolya.com"},
				{Id: 4, Username: "sasha", Name: "Sasha", ProfileImage: "sasha.jpg", WebsiteUrl: "sasha.com"},
			},
			err: nil,
		},
		"query error": {
			prepare: func(f *fields) {
				f.mock.
					ExpectQuery(regexp.QuoteMeta(getFolloweesCmd)).
					WithArgs(f.userID).
					WillReturnError(fmt.Errorf("db error"))
			},
			userID:    12,
			followees: nil,
			err:       pkgErrors.ErrDb,
		},
		"row scan error": {
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"id", "username"}).AddRow(1, "gosha")
				f.mock.
					ExpectQuery(regexp.QuoteMeta(getFolloweesCmd)).
					WithArgs(f.userID).
					WillReturnRows(rows)
			},
			userID:    12,
			followees: nil,
			err:       pkgErrors.ErrDb,
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

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("can't create mock: %s", err)
			}
			defer db.Close()

			f := fields{mock: mock, userID: test.userID, followees: test.followees}
			if test.prepare != nil {
				test.prepare(&f)
			}

			repo := NewRepository(db, logger)
			followees, err := repo.GetFollowees(test.userID)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if !reflect.DeepEqual(followees, test.followees) {
				t.Errorf("\nExpected: %v\nGot: %v", test.followees, followees)
			}
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("\nThere were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestRepository_GetFollowers(t *testing.T) {
	type fields struct {
		mock      sqlmock.Sqlmock
		userID    int
		followers []followings.Follower
	}

	type testCase struct {
		prepare   func(f *fields)
		userID    int
		followers []followings.Follower
		err       error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"id", "username", "name", "profile_image", "website_url"})
				for _, followee := range f.followers {
					rows = rows.AddRow(followee.Id, followee.Username, followee.Name, followee.ProfileImage,
						followee.WebsiteUrl)
				}
				f.mock.
					ExpectQuery(regexp.QuoteMeta(getFollowersCmd)).
					WithArgs(f.userID).
					WillReturnRows(rows)
			},
			userID: 12,
			followers: []followings.Follower{
				{Id: 2, Username: "vasua", Name: "Vasya", ProfileImage: "vasya.jpg", WebsiteUrl: "vasya.com"},
				{Id: 3, Username: "kolya", Name: "Kolya", ProfileImage: "kolya.jpg", WebsiteUrl: "kolya.com"},
				{Id: 4, Username: "sasha", Name: "Sasha", ProfileImage: "sasha.jpg", WebsiteUrl: "sasha.com"},
			},
			err: nil,
		},
		"query error": {
			prepare: func(f *fields) {
				f.mock.
					ExpectQuery(regexp.QuoteMeta(getFollowersCmd)).
					WithArgs(f.userID).
					WillReturnError(fmt.Errorf("db error"))
			},
			userID:    12,
			followers: nil,
			err:       pkgErrors.ErrDb,
		},
		"row scan error": {
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"id", "username"}).AddRow(1, "gosha")
				f.mock.
					ExpectQuery(regexp.QuoteMeta(getFollowersCmd)).
					WithArgs(f.userID).
					WillReturnRows(rows)
			},
			userID:    12,
			followers: nil,
			err:       pkgErrors.ErrDb,
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

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("can't create mock: %s", err)
			}
			defer db.Close()

			f := fields{mock: mock, userID: test.userID, followers: test.followers}
			if test.prepare != nil {
				test.prepare(&f)
			}

			repo := NewRepository(db, logger)
			followers, err := repo.GetFollowers(test.userID)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if !reflect.DeepEqual(followers, test.followers) {
				t.Errorf("\nExpected: %v\nGot: %v", test.followers, followers)
			}
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("\nThere were unfulfilled expectations: %s", err)
			}
		})
	}
}
