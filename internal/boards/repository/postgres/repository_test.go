package postgres

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	pkgBoards "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/boards"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
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
		mock sqlmock.Sqlmock
	}

	type testCase struct {
		prepare func(f *fields)
		params  pkgBoards.CreateParams
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
			params: pkgBoards.CreateParams{
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
					WillReturnError(pkgErrors.ErrDb)
			},
			params: pkgBoards.CreateParams{
				Name:        "n1",
				Description: "d1",
				Privacy:     "secret",
				UserId:      12,
			},
			board: models.Board{},
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

			repo := NewPostgresRepository(db, logger)

			f := fields{mock: mock}
			if test.prepare != nil {
				test.prepare(&f)
			}

			board, err := repo.Create(&test.params)
			if !errors.Is(err, test.err) {
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
			err:    pkgErrors.ErrDb,
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
			err:    pkgErrors.ErrDb,
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

			repo := NewPostgresRepository(db, logger)

			f := fields{mock: mock}
			if test.prepare != nil {
				test.prepare(&f)
			}

			boards, err := repo.List(test.userId)
			if !errors.Is(err, test.err) {
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
			err:   pkgErrors.ErrDb,
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

			repo := NewPostgresRepository(db, logger)

			f := fields{mock: mock}
			if test.prepare != nil {
				test.prepare(&f)
			}

			board, err := repo.Get(test.id)
			if !errors.Is(err, test.err) {
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
		params  pkgBoards.FullUpdateParams
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
			params: pkgBoards.FullUpdateParams{
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
			params: pkgBoards.FullUpdateParams{
				Id:          3,
				Name:        "upd_n1",
				Description: "upd_d1",
				Privacy:     "secret",
			},
			board: models.Board{},
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

			repo := NewPostgresRepository(db, logger)

			f := fields{mock: mock}
			if test.prepare != nil {
				test.prepare(&f)
			}

			board, err := repo.FullUpdate(&test.params)
			if !errors.Is(err, test.err) {
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
		params  pkgBoards.PartialUpdateParams
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
			params: pkgBoards.PartialUpdateParams{
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
			params: pkgBoards.PartialUpdateParams{
				Id:                3,
				Name:              "upd_n1",
				UpdateName:        true,
				Description:       "upd_d1",
				UpdateDescription: true,
				Privacy:           "secret",
				UpdatePrivacy:     true,
			},
			board: models.Board{},
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

			repo := NewPostgresRepository(db, logger)

			f := fields{mock: mock}
			if test.prepare != nil {
				test.prepare(&f)
			}

			board, err := repo.PartialUpdate(&test.params)
			if !errors.Is(err, test.err) {
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
			err: pkgErrors.ErrDb,
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

			repo := NewPostgresRepository(db, logger)

			f := fields{mock: mock}
			if test.prepare != nil {
				test.prepare(&f)
			}

			err = repo.Delete(test.id)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("\nThere were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestPinsList(t *testing.T) {
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
				rows := sqlmock.NewRows([]string{"id", "title", "description", "media_source", "media_source_color",
					"author_id"})
				rows = rows.AddRow(1, "t1", "d1", "ms_url1", "rgb(39, 102, 120)", 12)
				rows = rows.AddRow(2, "t2", "d2", "ms_url2", "rgb(39, 102, 120)", 12)
				rows = rows.AddRow(3, "t3", "d3", "ms_url3", "rgb(39, 102, 120)", 12)
				f.mock.
					ExpectQuery(regexp.QuoteMeta(pinsListCmd)).
					WithArgs(3, 30, 0).
					WillReturnRows(rows)
			},
			boardId: 3,
			page:    1,
			limit:   30,
			pins: []models.Pin{
				{Id: 1, Title: "t1", MediaSource: "ms_url1", MediaSourceColor: "rgb(39, 102, 120)", Description: "d1",
					Author: models.Profile{
						Id:           12,
						Username:     "un1",
						Name:         "n1",
						ProfileImage: "pi1",
						WebsiteUrl:   "wu1",
					}},
				{Id: 2, Title: "t2", MediaSource: "ms_url2", MediaSourceColor: "rgb(39, 102, 120)", Description: "d2",
					Author: models.Profile{
						Id:           13,
						Username:     "un2",
						Name:         "n2",
						ProfileImage: "pi2",
						WebsiteUrl:   "wu2",
					}},
				{Id: 3, Title: "t3", MediaSource: "ms_url3", MediaSourceColor: "rgb(39, 102, 120)", Description: "d3",
					Author: models.Profile{
						Id:           14,
						Username:     "un3",
						Name:         "n3",
						ProfileImage: "pi3",
						WebsiteUrl:   "wu3",
					}},
			},
			err: nil,
		},
		"query error": {
			prepare: func(f *fields) {
				f.mock.
					ExpectQuery(regexp.QuoteMeta(pinsListCmd)).
					WithArgs(3, 30, 0).
					WillReturnError(fmt.Errorf("sql error"))
			},
			boardId: 3,
			page:    1,
			limit:   30,
			pins:    nil,
			err:     pkgErrors.ErrDb,
		},
		"row scan error": {
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"id", "title"}).AddRow(1, "t1")
				f.mock.
					ExpectQuery(regexp.QuoteMeta(pinsListCmd)).
					WithArgs(3, 30, 0).
					WillReturnRows(rows)
			},
			boardId: 3,
			page:    1,
			limit:   30,
			pins:    nil,
			err:     pkgErrors.ErrDb,
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

			repo := NewPostgresRepository(db, logger)

			f := fields{mock: sqlMock}
			if test.prepare != nil {
				test.prepare(&f)
			}

			pins, err := repo.PinsList(test.boardId, test.page, test.limit)
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
			err:     pkgErrors.ErrDb,
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

			repo := NewPostgresRepository(db, logger)

			f := fields{mock: mock}
			if test.prepare != nil {
				test.prepare(&f)
			}

			access, err := repo.CheckWriteAccess(test.userId, test.boardId)
			if !errors.Is(err, test.err) {
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
			err:     pkgErrors.ErrDb,
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

			repo := NewPostgresRepository(db, logger)

			f := fields{mock: mock}
			if test.prepare != nil {
				test.prepare(&f)
			}

			access, err := repo.CheckReadAccess(test.userId, test.boardId)
			if !errors.Is(err, test.err) {
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
