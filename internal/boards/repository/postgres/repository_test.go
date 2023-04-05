package postgres

import (
	"fmt"
	_boards "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/boards"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

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
			board: models.Board{
				Id:          1,
				Name:        "n1",
				Description: "d1",
				Privacy:     "secret",
				UserId:      12,
			},
			err: nil,
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
	const getCmd = `SELECT *
					FROM boards
					WHERE id = $1;`

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("can't create mock: %s", err)
	}
	defer db.Close()

	logger := log.New()
	repo := NewPostgresRepository(db, logger)

	const name = "Test Name"
	const description = "Test Description"
	const privacy = "secret"
	const userId = 12
	rows := sqlmock.NewRows([]string{"id", "name", "description", "privacy", "user_id"})
	expect := []models.Board{{1, name, description, privacy, userId}}
	for _, board := range expect {
		rows = rows.AddRow(board.Id, board.Name, board.Description, board.Privacy, board.UserId)
	}

	// ok query
	mock.
		ExpectQuery(regexp.QuoteMeta(getCmd)).
		WithArgs(1).
		WillReturnRows(rows)

	board, err := repo.Get(1)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	if !reflect.DeepEqual(board, expect[0]) {
		t.Errorf("results not match, expected %v, \ngot %v", expect[0], board)
		return
	}

	// query error
	mock.
		ExpectQuery(regexp.QuoteMeta(getCmd)).
		WithArgs(1).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = repo.Get(1)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// row scan error
	rows = sqlmock.NewRows([]string{"id", "name"}).AddRow(1, name)
	mock.
		ExpectQuery(regexp.QuoteMeta(getCmd)).
		WithArgs(1).
		WillReturnRows(rows)

	_, err = repo.Get(1)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestFullUpdate(t *testing.T) {
	const fullUpdateCmd = `UPDATE boards
						   SET name = $1::VARCHAR,
						   description = $2::TEXT,
						   privacy = $3::privacy
						   WHERE id = $4
						   RETURNING *;`

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("can't create mock: %s", err)
	}
	defer db.Close()

	logger := log.New()
	repo := NewPostgresRepository(db, logger)

	const (
		id          = 3
		name        = "Test Name"
		description = "Test Description"
		privacy     = "secret"
		userId      = 12
	)
	rows := sqlmock.NewRows([]string{"id", "name", "description", "privacy", "user_id"})
	rows = rows.AddRow(id, name, description, privacy, userId)

	// ok query
	mock.
		ExpectQuery(regexp.QuoteMeta(fullUpdateCmd)).
		WithArgs(name, description, privacy, id).
		WillReturnRows(rows)

	testParams := _boards.FullUpdateParams{
		Id:          id,
		Name:        name,
		Description: description,
		Privacy:     privacy,
	}
	board, err := repo.FullUpdate(&testParams)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if board.Id != id {
		t.Errorf("bad id: want %v, have %v", board.Id, 1)
		return
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// query error
	mock.
		ExpectQuery(regexp.QuoteMeta(fullUpdateCmd)).
		WithArgs(name, description, privacy, id).
		WillReturnError(fmt.Errorf("bad query"))

	_, err = repo.FullUpdate(&testParams)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPartialUpdate(t *testing.T) {
	const partialUpdateCmd = `UPDATE boards
							  SET name = CASE WHEN $1::boolean THEN $2::VARCHAR ELSE name END,
    						  description = CASE WHEN $3::boolean THEN $4::TEXT ELSE description END,
    						  privacy = CASE WHEN $5::boolean THEN $6::privacy ELSE privacy END
						      WHERE id = $7
							  RETURNING *;`

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("can't create mock: %s", err)
	}
	defer db.Close()

	logger := log.New()
	repo := NewPostgresRepository(db, logger)

	expect := models.Board{
		Id:          3,
		Name:        "Test Name",
		Description: "Test Description",
		Privacy:     "secret",
		UserId:      12,
	}
	rows := sqlmock.NewRows([]string{"id", "name", "description", "privacy", "user_id"})
	rows = rows.AddRow(expect.Id, expect.Name, expect.Description, expect.Privacy,
		expect.UserId)

	params := _boards.PartialUpdateParams{
		Id:                expect.Id,
		Name:              expect.Name,
		UpdateName:        true,
		Description:       expect.Description,
		UpdateDescription: true,
		Privacy:           expect.Privacy,
		UpdatePrivacy:     true,
	}

	// ok query
	mock.
		ExpectQuery(regexp.QuoteMeta(partialUpdateCmd)).
		WithArgs(
			params.UpdateName,
			params.Name,
			params.UpdateDescription,
			params.Description,
			params.UpdatePrivacy,
			params.Privacy,
			params.Id,
		).
		WillReturnRows(rows)

	board, err := repo.PartialUpdate(&params)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	if !reflect.DeepEqual(board, expect) {
		t.Errorf("results not match, expected %#v, \ngot %#v", expect, board)
		return
	}

	// query error
	mock.
		ExpectQuery(regexp.QuoteMeta(partialUpdateCmd)).
		WithArgs(
			params.UpdateName,
			params.Name,
			params.UpdateDescription,
			params.Description,
			params.UpdatePrivacy,
			params.Privacy,
			params.Id,
		).
		WillReturnError(fmt.Errorf("bad query"))

	_, err = repo.PartialUpdate(&params)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDelete(t *testing.T) {
	const deleteCmd = `DELETE FROM boards 
					   WHERE id = $1;`

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("can't create mock: %s", err)
	}
	defer db.Close()

	logger := log.New()
	repo := NewPostgresRepository(db, logger)

	const delId = 3

	rows := sqlmock.NewRows([]string{"id", "name", "description", "privacy", "user_id"})
	rows = rows.AddRow(delId, "Test Name", "Test Description", "secret", 12)

	// ok query
	mock.
		ExpectExec(regexp.QuoteMeta(deleteCmd)).
		WithArgs(delId).
		WillReturnResult(sqlmock.NewResult(delId, 1))

	err = repo.Delete(delId)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// query error
	mock.
		ExpectExec(regexp.QuoteMeta(deleteCmd)).
		WithArgs(delId).
		WillReturnError(fmt.Errorf("db_error"))

	err = repo.Delete(delId)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
