package postgres

import (
	"github.com/DATA-DOG/go-sqlmock"
	pkgComments "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/comments"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"reflect"
	"regexp"
	"testing"
)

func TestRepository_Create(t *testing.T) {
	type fields struct {
		mock    sqlmock.Sqlmock
		params  *pkgComments.CreateParams
		comment *models.Comment
	}

	type testCase struct {
		prepare func(f *fields)
		params  *pkgComments.CreateParams
		comment models.Comment
		err     error
	}

	tests := map[string]testCase{
		"good query": {
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"id", "author_id", "pin_id", "text", "created_at"})
				rows = rows.AddRow(f.comment.ID, f.comment.AuthorID, f.comment.PinID, f.comment.Text,
					f.comment.CreatedAt)
				f.mock.
					ExpectQuery(regexp.QuoteMeta(createCmd)).
					WithArgs(f.params.AuthorID, f.params.PinID, f.params.Text).
					WillReturnRows(rows)
			},
			params:  &pkgComments.CreateParams{AuthorID: 27, PinID: 21, Text: "Good pin!"},
			comment: models.Comment{ID: 2, AuthorID: 27, PinID: 21, Text: "Good pin!"},
			err:     nil,
		},
		"query error": {
			prepare: func(f *fields) {
				f.mock.
					ExpectQuery(regexp.QuoteMeta(createCmd)).
					WithArgs(f.params.AuthorID, f.params.PinID, f.params.Text).
					WillReturnError(&pq.Error{Message: "sql error"})
			},
			params:  &pkgComments.CreateParams{AuthorID: 27, PinID: 21, Text: "Good pin!"},
			comment: models.Comment{},
			err:     pkgErrors.ErrDb,
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

			f := fields{mock: sqlMock, params: test.params, comment: &test.comment}
			if test.prepare != nil {
				test.prepare(&f)
			}

			repo := NewRepository(db, logger)
			comment, err := repo.Create(test.params)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if comment != test.comment {
				t.Errorf("\nExpected: %v\nGot: %v", test.comment, comment)
			}
			if err = sqlMock.ExpectationsWereMet(); err != nil {
				t.Errorf("\nThere were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestRepository_List(t *testing.T) {
	type fields struct {
		mock     sqlmock.Sqlmock
		pinID    int
		comments []models.Comment
	}

	type testCase struct {
		prepare  func(f *fields)
		pinID    int
		comments []models.Comment
		err      error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"id", "author_id", "pin_id", "text", "created_at"})
				for _, comment := range f.comments {
					rows = rows.AddRow(comment.ID, comment.AuthorID, comment.PinID, comment.Text, comment.CreatedAt)
				}
				f.mock.
					ExpectQuery(regexp.QuoteMeta(listCmd)).
					WithArgs(f.pinID).
					WillReturnRows(rows)
			},
			pinID: 21,
			comments: []models.Comment{
				{ID: 2, AuthorID: 27, PinID: 21, Text: "Good pin!"},
				{ID: 3, AuthorID: 28, PinID: 21, Text: "Yeah!"},
				{ID: 4, AuthorID: 27, PinID: 21, Text: "Fantastic!"},
			},
			err: nil,
		},
		"query error": {
			prepare: func(f *fields) {
				f.mock.
					ExpectQuery(regexp.QuoteMeta(listCmd)).
					WithArgs(f.pinID).
					WillReturnError(&pq.Error{Message: "sql error"})
			},
			pinID:    21,
			comments: nil,
			err:      pkgErrors.ErrDb,
		},
		"row scan error": {
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"id", "author_id"}).AddRow(3, 2)
				f.mock.
					ExpectQuery(regexp.QuoteMeta(listCmd)).
					WithArgs(f.pinID).
					WillReturnRows(rows)
			},
			pinID:    21,
			comments: nil,
			err:      pkgErrors.ErrDb,
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

			f := fields{mock: mock, pinID: test.pinID, comments: test.comments}
			if test.prepare != nil {
				test.prepare(&f)
			}

			repo := NewRepository(db, logger)
			comments, err := repo.List(test.pinID)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if !reflect.DeepEqual(comments, test.comments) {
				t.Errorf("\nExpected: %v\nGot: %v", test.comments, comments)
			}
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("\nThere were unfulfilled expectations: %s", err)
			}
		})
	}
}
