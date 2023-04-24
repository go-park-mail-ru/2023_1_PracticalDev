package postgres

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pkg/errors"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
)

func TestCreate(t *testing.T) {
	type fields struct {
		mock     sqlmock.Sqlmock
		pinId    int
		authorId int
	}

	type testCase struct {
		prepare  func(f *fields)
		pinId    int
		authorId int
		err      error
	}

	tests := map[string]*testCase{
		"usual": {
			prepare: func(f *fields) {
				f.mock.
					ExpectExec(regexp.QuoteMeta(createCmd)).
					WithArgs(f.pinId, f.authorId).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			pinId:    3,
			authorId: 2,
			err:      nil,
		},
		"exec error": {
			prepare: func(f *fields) {
				f.mock.
					ExpectExec(regexp.QuoteMeta(createCmd)).
					WithArgs(f.pinId, f.authorId).
					WillReturnError(fmt.Errorf("db error"))
			},
			pinId:    3,
			authorId: 2,
			err:      pkgErrors.ErrDb,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("can't create mock: %s", err)
			}
			defer db.Close()

			f := fields{mock: mock, pinId: test.pinId, authorId: test.authorId}
			if test.prepare != nil {
				test.prepare(&f)
			}

			repo := NewRepository(db, log.New())
			err = repo.Create(test.pinId, test.authorId)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("\nThere were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestListByAuthor(t *testing.T) {
	type fields struct {
		mock     sqlmock.Sqlmock
		authorId int
		likes    []models.Like
	}

	type testCase struct {
		prepare  func(f *fields)
		authorId int
		likes    []models.Like
		err      error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"pin_id", "author_id", "created_at"})
				for _, like := range f.likes {
					rows = rows.AddRow(like.PinId, like.AuthorId, like.CreatedAt)
				}
				f.mock.
					ExpectQuery(regexp.QuoteMeta(listByAuthorCmd)).
					WithArgs(f.authorId).
					WillReturnRows(rows)
			},
			authorId: 12,
			likes: []models.Like{
				{PinId: 2, AuthorId: 12, CreatedAt: time.Unix(1681163314, 0)},
				{PinId: 3, AuthorId: 12, CreatedAt: time.Unix(1681163000, 0)},
				{PinId: 4, AuthorId: 12, CreatedAt: time.Unix(1681164555, 0)},
			},
			err: nil,
		},
		"query error": {
			prepare: func(f *fields) {
				f.mock.
					ExpectQuery(regexp.QuoteMeta(listByAuthorCmd)).
					WithArgs(f.authorId).
					WillReturnError(fmt.Errorf("db error"))
			},
			authorId: 12,
			likes:    nil,
			err:      pkgErrors.ErrDb,
		},
		"row scan error": {
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"pin_id"}).AddRow(1)
				f.mock.
					ExpectQuery(regexp.QuoteMeta(listByAuthorCmd)).
					WithArgs(f.authorId).
					WillReturnRows(rows)
			},
			authorId: 12,
			likes:    nil,
			err:      pkgErrors.ErrDb,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("can't create mock: %s", err)
			}
			defer db.Close()

			f := fields{mock: mock, authorId: test.authorId, likes: test.likes}
			if test.prepare != nil {
				test.prepare(&f)
			}

			repo := NewRepository(db, log.New())
			likes, err := repo.ListByAuthor(test.authorId)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if !reflect.DeepEqual(likes, test.likes) {
				t.Errorf("\nExpected: %v\nGot: %v", test.likes, likes)
			}
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("\nThere were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestListByPin(t *testing.T) {
	type fields struct {
		mock  sqlmock.Sqlmock
		pinId int
		likes []models.Like
	}

	type testCase struct {
		prepare func(f *fields)
		pinId   int
		likes   []models.Like
		err     error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"pin_id", "author_id", "created_at"})
				for _, like := range f.likes {
					rows = rows.AddRow(like.PinId, like.AuthorId, like.CreatedAt)
				}
				f.mock.
					ExpectQuery(regexp.QuoteMeta(listByPinCmd)).
					WithArgs(f.pinId).
					WillReturnRows(rows)
			},
			pinId: 12,
			likes: []models.Like{
				{PinId: 12, AuthorId: 4, CreatedAt: time.Unix(1681163314, 0)},
				{PinId: 12, AuthorId: 23, CreatedAt: time.Unix(1681163000, 0)},
				{PinId: 12, AuthorId: 2, CreatedAt: time.Unix(1681164555, 0)},
			},
			err: nil,
		},
		"query error": {
			prepare: func(f *fields) {
				f.mock.
					ExpectQuery(regexp.QuoteMeta(listByPinCmd)).
					WithArgs(f.pinId).
					WillReturnError(fmt.Errorf("db error"))
			},
			pinId: 12,
			likes: nil,
			err:   pkgErrors.ErrDb,
		},
		"row scan error": {
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"pin_id"}).AddRow(1)
				f.mock.
					ExpectQuery(regexp.QuoteMeta(listByPinCmd)).
					WithArgs(f.pinId).
					WillReturnRows(rows)
			},
			pinId: 12,
			likes: nil,
			err:   pkgErrors.ErrDb,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("can't create mock: %s", err)
			}
			defer db.Close()

			f := fields{mock: mock, pinId: test.pinId, likes: test.likes}
			if test.prepare != nil {
				test.prepare(&f)
			}

			repo := NewRepository(db, log.New())
			likes, err := repo.ListByPin(test.pinId)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if !reflect.DeepEqual(likes, test.likes) {
				t.Errorf("\nExpected: %v\nGot: %v", test.likes, likes)
			}
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("\nThere were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	type fields struct {
		mock     sqlmock.Sqlmock
		pinId    int
		authorId int
	}

	type testCase struct {
		prepare  func(f *fields)
		pinId    int
		authorId int
		err      error
	}

	tests := map[string]testCase{
		"good query": {
			prepare: func(f *fields) {
				f.mock.
					ExpectExec(regexp.QuoteMeta(deleteCmd)).
					WithArgs(f.pinId, f.authorId).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			pinId:    3,
			authorId: 2,
			err:      nil,
		},
		"exec error": {
			prepare: func(f *fields) {
				f.mock.
					ExpectExec(regexp.QuoteMeta(deleteCmd)).
					WithArgs(f.pinId, f.authorId).
					WillReturnError(fmt.Errorf("db error"))
			},
			pinId:    3,
			authorId: 2,
			err:      pkgErrors.ErrDb,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("can't create mock: %s", err)
			}
			defer db.Close()

			f := fields{mock: mock, pinId: test.pinId, authorId: test.authorId}
			if test.prepare != nil {
				test.prepare(&f)
			}

			repo := NewRepository(db, log.New())
			err = repo.Delete(test.pinId, test.authorId)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("\nThere were unfulfilled expectations: %s", err)
			}
		})
	}
}
