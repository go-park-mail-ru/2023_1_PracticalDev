package postgres

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	_boards "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/boards"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
)

func TestCreate(t *testing.T) {
	type fields struct {
		mock sqlmock.Sqlmock
	}

	type testCase struct {
		prepare func(f *fields)
		params  _boards.CreateParams
		board   models.Board
		err     error
	}

	const createCmd = `INSERT INTO boards (name, description, privacy, user_id) 
				       VALUES ($1, $2, $3, $4)
					   RETURNING *;`

	tests := map[string]testCase{
		"good query": {
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"id", "name", "description", "privacy", "user_id"})
				rows = rows.AddRow(1, "n1", "d1", "secret", 12)
				f.mock.
					ExpectQuery(regexp.QuoteMeta(createCmd)).
					WithArgs("n1", "d1", "secret", 12).
					WillReturnRows(rows)
			},
			params: _boards.CreateParams{
				Name:        "n1",
				Description: "d1",
				Privacy:     "secret",
				UserId:      12,
			},
			board: models.Board{Id: 1, Name: "n1", Description: "d1", Privacy: "secret", UserId: 12},
			err:   nil,
		},
		"query error": {
			prepare: func(f *fields) {
				f.mock.
					ExpectQuery(regexp.QuoteMeta(createCmd)).
					WithArgs("n1", "d1", "secret", 12).
					WillReturnError(_boards.ErrDb)
			},
			params: _boards.CreateParams{
				Name:        "n1",
				Description: "d1",
				Privacy:     "secret",
				UserId:      12,
			},
			board: models.Board{},
			err:   _boards.ErrDb,
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

			logger := log.New()
			repo := NewPostgresRepository(db, logger)

			f := fields{mock: mock}
			if test.prepare != nil {
				test.prepare(&f)
			}

			board, err := repo.Create(&test.params)
			if err != test.err {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if board != test.board {
				t.Errorf("\nExpected: %v\nGot: %v", test.board, board)
			}
			if err = mock.ExpectationsWereMet(); err != nil {
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
		userId  int
		boards  []models.Board
		err     error
	}

	const listCmd = `SELECT *
					 FROM boards
					 WHERE user_id = $1;`

	tests := map[string]testCase{
		"good query": {
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"id", "name", "description", "privacy", "user_id"})
				expect := []models.Board{
					{Id: 1, Name: "b1", Description: "d1", Privacy: "secret", UserId: 12},
					{Id: 2, Name: "b2", Description: "d2", Privacy: "secret", UserId: 12},
					{Id: 5, Name: "b5", Description: "d5", Privacy: "public", UserId: 12},
				}
				for _, board := range expect {
					rows = rows.AddRow(board.Id, board.Name, board.Description, board.Privacy, board.UserId)
				}
				f.mock.
					ExpectQuery(regexp.QuoteMeta(listCmd)).
					WithArgs(12).
					WillReturnRows(rows)
			},
			userId: 12,
			boards: []models.Board{
				{Id: 1, Name: "b1", Description: "d1", Privacy: "secret", UserId: 12},
				{Id: 2, Name: "b2", Description: "d2", Privacy: "secret", UserId: 12},
				{Id: 5, Name: "b5", Description: "d5", Privacy: "public", UserId: 12},
			},
			err: nil,
		},
		"query error": {
			prepare: func(f *fields) {
				f.mock.
					ExpectQuery(regexp.QuoteMeta(listCmd)).
					WithArgs(12).
					WillReturnError(fmt.Errorf("db error"))
			},
			userId: 12,
			boards: nil,
			err:    _boards.ErrDb,
		},
		"row scan error": {
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "b1")
				f.mock.
					ExpectQuery(regexp.QuoteMeta(listCmd)).
					WithArgs(12).
					WillReturnRows(rows)
			},
			userId: 12,
			boards: nil,
			err:    _boards.ErrDb,
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

			logger := log.New()
			repo := NewPostgresRepository(db, logger)

			f := fields{mock: mock}
			if test.prepare != nil {
				test.prepare(&f)
			}

			boards, err := repo.List(test.userId)
			if err != test.err {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if !reflect.DeepEqual(boards, test.boards) {
				t.Errorf("\nExpected: %v\nGot: %v", test.boards, boards)
			}
			if err = mock.ExpectationsWereMet(); err != nil {
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
		board   models.Board
		err     error
	}

	const getCmd = `SELECT *
					FROM boards
					WHERE id = $1;`

	tests := map[string]testCase{
		"good query": {
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"id", "name", "description", "privacy", "user_id"})
				rows = rows.AddRow(3, "n1", "d1", "secret", 12)
				f.mock.
					ExpectQuery(regexp.QuoteMeta(getCmd)).
					WithArgs(3).
					WillReturnRows(rows)
			},
			id:    3,
			board: models.Board{Id: 3, Name: "n1", Description: "d1", Privacy: "secret", UserId: 12},
			err:   nil,
		},
		"query error": {
			prepare: func(f *fields) {
				f.mock.
					ExpectQuery(regexp.QuoteMeta(getCmd)).
					WithArgs(3).
					WillReturnError(fmt.Errorf("db error"))
			},
			id:    3,
			board: models.Board{},
			err:   _boards.ErrDb,
		},
		"row scan error": {
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "b1")
				f.mock.
					ExpectQuery(regexp.QuoteMeta(getCmd)).
					WithArgs(1).
					WillReturnRows(rows)
			},
			id:    1,
			board: models.Board{},
			err:   _boards.ErrDb,
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

			logger := log.New()
			repo := NewPostgresRepository(db, logger)

			f := fields{mock: mock}
			if test.prepare != nil {
				test.prepare(&f)
			}

			board, err := repo.Get(test.id)
			if err != test.err {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if board != test.board {
				t.Errorf("\nExpected: %v\nGot: %v", test.board, board)
			}
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("\nThere were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestFullUpdate(t *testing.T) {
	type fields struct {
		mock sqlmock.Sqlmock
	}

	type testCase struct {
		prepare func(f *fields)
		params  _boards.FullUpdateParams
		board   models.Board
		err     error
	}

	const fullUpdateCmd = `UPDATE boards
						   SET name = $1::VARCHAR,
						   description = $2::TEXT,
						   privacy = $3::privacy
						   WHERE id = $4
						   RETURNING *;`

	tests := map[string]testCase{
		"good query": {
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"id", "name", "description", "privacy", "user_id"})
				rows = rows.AddRow(3, "upd_n1", "upd_d1", "secret", 12)
				f.mock.
					ExpectQuery(regexp.QuoteMeta(fullUpdateCmd)).
					WithArgs("upd_n1", "upd_d1", "secret", 3).
					WillReturnRows(rows)
			},
			params: _boards.FullUpdateParams{
				Id:          3,
				Name:        "upd_n1",
				Description: "upd_d1",
				Privacy:     "secret",
			},
			board: models.Board{
				Id:          3,
				Name:        "upd_n1",
				Description: "upd_d1",
				Privacy:     "secret",
				UserId:      12,
			},
			err: nil,
		},
		"query error": {
			prepare: func(f *fields) {
				f.mock.
					ExpectQuery(regexp.QuoteMeta(fullUpdateCmd)).
					WithArgs("upd_n1", "upd_d1", "secret", 3).
					WillReturnError(fmt.Errorf("db error"))
			},
			params: _boards.FullUpdateParams{
				Id:          3,
				Name:        "upd_n1",
				Description: "upd_d1",
				Privacy:     "secret",
			},
			board: models.Board{},
			err:   _boards.ErrDb,
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

			logger := log.New()
			repo := NewPostgresRepository(db, logger)

			f := fields{mock: mock}
			if test.prepare != nil {
				test.prepare(&f)
			}

			board, err := repo.FullUpdate(&test.params)
			if err != test.err {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if board != test.board {
				t.Errorf("\nExpected: %v\nGot: %v", test.board, board)
			}
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("\nThere were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestPartialUpdate(t *testing.T) {
	type fields struct {
		mock sqlmock.Sqlmock
	}

	type testCase struct {
		prepare func(f *fields)
		params  _boards.PartialUpdateParams
		board   models.Board
		err     error
	}

	const partialUpdateCmd = `UPDATE boards
							  SET name = CASE WHEN $1::boolean THEN $2::VARCHAR ELSE name END,
    						  description = CASE WHEN $3::boolean THEN $4::TEXT ELSE description END,
    						  privacy = CASE WHEN $5::boolean THEN $6::privacy ELSE privacy END
						      WHERE id = $7
							  RETURNING *;`

	tests := map[string]testCase{
		"good query": {
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"id", "name", "description", "privacy", "user_id"})
				rows = rows.AddRow(3, "upd_n1", "upd_d1", "secret", 12)
				f.mock.
					ExpectQuery(regexp.QuoteMeta(partialUpdateCmd)).
					WithArgs(true, "upd_n1", true, "upd_d1", true, "secret", 3).
					WillReturnRows(rows)
			},
			params: _boards.PartialUpdateParams{
				Id:                3,
				Name:              "upd_n1",
				UpdateName:        true,
				Description:       "upd_d1",
				UpdateDescription: true,
				Privacy:           "secret",
				UpdatePrivacy:     true,
			},
			board: models.Board{
				Id:          3,
				Name:        "upd_n1",
				Description: "upd_d1",
				Privacy:     "secret",
				UserId:      12,
			},
			err: nil,
		},
		"query error": {
			prepare: func(f *fields) {
				f.mock.
					ExpectQuery(regexp.QuoteMeta(partialUpdateCmd)).
					WithArgs(true, "upd_n1", true, "upd_d1", true, "secret", 3).
					WillReturnError(fmt.Errorf("db error"))
			},
			params: _boards.PartialUpdateParams{
				Id:                3,
				Name:              "upd_n1",
				UpdateName:        true,
				Description:       "upd_d1",
				UpdateDescription: true,
				Privacy:           "secret",
				UpdatePrivacy:     true,
			},
			board: models.Board{},
			err:   _boards.ErrDb,
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

			logger := log.New()
			repo := NewPostgresRepository(db, logger)

			f := fields{mock: mock}
			if test.prepare != nil {
				test.prepare(&f)
			}

			board, err := repo.PartialUpdate(&test.params)
			if err != test.err {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if board != test.board {
				t.Errorf("\nExpected: %v\nGot: %v", test.board, board)
			}
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("\nThere were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	type fields struct {
		mock sqlmock.Sqlmock
	}

	type testCase struct {
		prepare func(f *fields)
		id      int
		err     error
	}

	const deleteCmd = `DELETE FROM boards 
					   WHERE id = $1;`

	tests := map[string]testCase{
		"good query": {
			prepare: func(f *fields) {
				f.mock.
					ExpectExec(regexp.QuoteMeta(deleteCmd)).
					WithArgs(3).
					WillReturnResult(sqlmock.NewResult(3, 1))
			},
			id:  3,
			err: nil,
		},
		"query error": {
			prepare: func(f *fields) {
				f.mock.
					ExpectExec(regexp.QuoteMeta(deleteCmd)).
					WithArgs(3).
					WillReturnError(fmt.Errorf("db error"))
			},
			id:  3,
			err: _boards.ErrDb,
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

			logger := log.New()
			repo := NewPostgresRepository(db, logger)

			f := fields{mock: mock}
			if test.prepare != nil {
				test.prepare(&f)
			}

			err = repo.Delete(test.id)
			if err != test.err {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("\nThere were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestCheckWriteAccess(t *testing.T) {
	type fields struct {
		mock sqlmock.Sqlmock
	}

	type testCase struct {
		prepare func(f *fields)
		userId  string
		boardId string
		access  bool
		err     error
	}

	const checkCommand = `SELECT EXISTS(SELECT id
     			          				FROM boards
              			  				WHERE id = $1 AND user_id = $2) AS access;`

	tests := map[string]testCase{
		"good query": {
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"access"}).AddRow(true)
				f.mock.
					ExpectQuery(regexp.QuoteMeta(checkCommand)).
					WithArgs("3", "2").
					WillReturnRows(rows)
			},
			userId:  "2",
			boardId: "3",
			access:  true,
			err:     nil,
		},
		"query error": {
			prepare: func(f *fields) {
				f.mock.
					ExpectQuery(regexp.QuoteMeta(checkCommand)).
					WithArgs("3", "2").
					WillReturnError(fmt.Errorf("db error"))
			},
			userId:  "2",
			boardId: "3",
			access:  false,
			err:     _boards.ErrDb,
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

			logger := log.New()
			repo := NewPostgresRepository(db, logger)

			f := fields{mock: mock}
			if test.prepare != nil {
				test.prepare(&f)
			}

			access, err := repo.CheckWriteAccess(test.userId, test.boardId)
			if err != test.err {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if access != test.access {
				t.Errorf("\nExpected: %v\nGot: %v", test.access, access)
			}
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("\nThere were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestCheckReadAccess(t *testing.T) {
	type fields struct {
		mock sqlmock.Sqlmock
	}

	type testCase struct {
		prepare func(f *fields)
		userId  string
		boardId string
		access  bool
		err     error
	}

	const checkCommand = `SELECT EXISTS(SELECT
              							FROM boards
              							WHERE id = $1 AND (privacy = 'public' OR user_id = $2)) AS access;`

	tests := map[string]testCase{
		"good query": {
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"access"}).AddRow(true)
				f.mock.
					ExpectQuery(regexp.QuoteMeta(checkCommand)).
					WithArgs("3", "2").
					WillReturnRows(rows)
			},
			userId:  "2",
			boardId: "3",
			access:  true,
			err:     nil,
		},
		"query error": {
			prepare: func(f *fields) {
				f.mock.
					ExpectQuery(regexp.QuoteMeta(checkCommand)).
					WithArgs("3", "2").
					WillReturnError(fmt.Errorf("db error"))
			},
			userId:  "2",
			boardId: "3",
			access:  false,
			err:     _boards.ErrDb,
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

			logger := log.New()
			repo := NewPostgresRepository(db, logger)

			f := fields{mock: mock}
			if test.prepare != nil {
				test.prepare(&f)
			}

			access, err := repo.CheckReadAccess(test.userId, test.boardId)
			if err != test.err {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if access != test.access {
				t.Errorf("\nExpected: %v\nGot: %v", test.access, access)
			}
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("\nThere were unfulfilled expectations: %s", err)
			}
		})
	}
}
