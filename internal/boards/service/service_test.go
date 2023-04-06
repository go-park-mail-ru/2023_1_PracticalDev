package service

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	_boards "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/boards"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/boards/mocks"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
)

func TestCreate(t *testing.T) {
	type fields struct {
		repo *mocks.MockRepository
	}

	type testCase struct {
		prepare func(f *fields)
		params  _boards.CreateParams
		board   models.Board
		err     error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Create(&_boards.CreateParams{
					Name:        "n1",
					Description: "d1",
					Privacy:     "secret",
					UserId:      12,
				}).Return(models.Board{
					Id:          1,
					Name:        "n1",
					Description: "d1",
					Privacy:     "secret",
					UserId:      12,
				}, nil)
			},
			params: _boards.CreateParams{Name: "n1", Description: "d1", Privacy: "secret", UserId: 12},
			board:  models.Board{Id: 1, Name: "n1", Description: "d1", Privacy: "secret", UserId: 12},
			err:    nil,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl)}
			if test.prepare != nil {
				test.prepare(&f)
			}

			serv := NewBoardsService(f.repo)

			board, err := serv.Create(&test.params)
			if err != test.err {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if board != test.board {
				t.Errorf("\nExpected: %v\nGot: %v", test.board, board)
			}
		})
	}
}

func TestList(t *testing.T) {
	type fields struct {
		repo *mocks.MockRepository
	}

	type testCase struct {
		prepare func(f *fields)
		userId  int
		boards  []models.Board
		err     error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				f.repo.EXPECT().List(12).Return([]models.Board{
					{Id: 1, Name: "b1", Description: "d1", Privacy: "secret", UserId: 12},
					{Id: 2, Name: "b2", Description: "d2", Privacy: "secret", UserId: 12},
					{Id: 5, Name: "b5", Description: "d5", Privacy: "public", UserId: 12},
				}, nil)
			},
			userId: 12,
			boards: []models.Board{
				{Id: 1, Name: "b1", Description: "d1", Privacy: "secret", UserId: 12},
				{Id: 2, Name: "b2", Description: "d2", Privacy: "secret", UserId: 12},
				{Id: 5, Name: "b5", Description: "d5", Privacy: "public", UserId: 12},
			},
			err: nil,
		},
		"no boards": {
			prepare: func(f *fields) {
				f.repo.EXPECT().List(3).Return([]models.Board{}, nil)
			},
			userId: 3,
			boards: []models.Board{},
			err:    nil,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl)}
			if test.prepare != nil {
				test.prepare(&f)
			}

			serv := NewBoardsService(f.repo)

			boards, err := serv.List(test.userId)
			if err != test.err {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if !reflect.DeepEqual(boards, test.boards) {
				t.Errorf("\nExpected: %v\nGot: %v", test.boards, boards)
			}
		})
	}
}

func TestGet(t *testing.T) {
	type fields struct {
		repo *mocks.MockRepository
	}

	type testCase struct {
		prepare func(f *fields)
		id      int
		board   models.Board
		err     error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Get(3).Return(models.Board{
					Id:          3,
					Name:        "n3",
					Description: "d3",
					Privacy:     "secret",
					UserId:      1,
				}, nil)
			},
			id: 3,
			board: models.Board{
				Id:          3,
				Name:        "n3",
				Description: "d3",
				Privacy:     "secret",
				UserId:      1,
			},
			err: nil,
		},
		"board not found": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Get(3).Return(models.Board{}, _boards.ErrBoardNotFound)
			},
			id:    3,
			board: models.Board{},
			err:   _boards.ErrBoardNotFound,
		},
		"negative board id param": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Get(-1).Return(models.Board{}, _boards.ErrBoardNotFound)
			},
			id:    -1,
			board: models.Board{},
			err:   _boards.ErrBoardNotFound,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl)}
			if test.prepare != nil {
				test.prepare(&f)
			}

			serv := NewBoardsService(f.repo)

			board, err := serv.Get(test.id)
			if err != test.err {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if board != test.board {
				t.Errorf("\nExpected: %v\nGot: %v", test.board, board)
			}
		})
	}
}

func TestFullUpdate(t *testing.T) {
	type fields struct {
		repo *mocks.MockRepository
	}

	type testCase struct {
		prepare func(f *fields)
		params  _boards.FullUpdateParams
		board   models.Board
		err     error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				f.repo.EXPECT().FullUpdate(&_boards.FullUpdateParams{
					Id:          3,
					Name:        "n1",
					Description: "d1",
					Privacy:     "secret",
				}).Return(models.Board{
					Id:          3,
					Name:        "n1",
					Description: "d1",
					Privacy:     "secret",
					UserId:      12,
				}, nil)
			},
			params: _boards.FullUpdateParams{Id: 3, Name: "n1", Description: "d1", Privacy: "secret"},
			board:  models.Board{Id: 3, Name: "n1", Description: "d1", Privacy: "secret", UserId: 12},
			err:    nil,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl)}
			if test.prepare != nil {
				test.prepare(&f)
			}

			serv := NewBoardsService(f.repo)

			board, err := serv.FullUpdate(&test.params)
			if err != test.err {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if board != test.board {
				t.Errorf("\nExpected: %v\nGot: %v", test.board, board)
			}
		})
	}
}

func TestPartialUpdate(t *testing.T) {
	type fields struct {
		repo *mocks.MockRepository
	}

	type testCase struct {
		prepare func(f *fields)
		params  _boards.PartialUpdateParams
		board   models.Board
		err     error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				f.repo.EXPECT().PartialUpdate(&_boards.PartialUpdateParams{
					Id:                3,
					Name:              "n1",
					UpdateName:        true,
					Description:       "d1",
					UpdateDescription: true,
					Privacy:           "secret",
					UpdatePrivacy:     true,
				}).Return(models.Board{
					Id:          3,
					Name:        "n1",
					Description: "d1",
					Privacy:     "secret",
					UserId:      12,
				}, nil)
			},
			params: _boards.PartialUpdateParams{
				Id:                3,
				Name:              "n1",
				UpdateName:        true,
				Description:       "d1",
				UpdateDescription: true,
				Privacy:           "secret",
				UpdatePrivacy:     true,
			},
			board: models.Board{Id: 3, Name: "n1", Description: "d1", Privacy: "secret", UserId: 12},
			err:   nil,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl)}
			if test.prepare != nil {
				test.prepare(&f)
			}

			serv := NewBoardsService(f.repo)

			board, err := serv.PartialUpdate(&test.params)
			if err != test.err {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if board != test.board {
				t.Errorf("\nExpected: %v\nGot: %v", test.board, board)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	type fields struct {
		repo *mocks.MockRepository
	}

	type testCase struct {
		prepare func(f *fields)
		id      int
		err     error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Delete(3).Return(nil)
			},
			id:  3,
			err: nil,
		},
		"board not found": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Delete(3).Return(_boards.ErrBoardNotFound)
			},
			id:  3,
			err: _boards.ErrBoardNotFound,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl)}
			if test.prepare != nil {
				test.prepare(&f)
			}

			serv := NewBoardsService(f.repo)

			err := serv.Delete(test.id)
			if err != test.err {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
		})
	}
}

func TestCheckWriteAccess(t *testing.T) {
	type fields struct {
		repo *mocks.MockRepository
	}

	type testCase struct {
		prepare func(f *fields)
		userId  string
		boardId string
		access  bool
		err     error
	}

	tests := map[string]testCase{
		"access is allowed": {
			prepare: func(f *fields) {
				f.repo.EXPECT().CheckWriteAccess("2", "3").Return(true, nil)
			},
			userId:  "2",
			boardId: "3",
			access:  true,
			err:     nil,
		},
		"access is denied": {
			prepare: func(f *fields) {
				f.repo.EXPECT().CheckWriteAccess("2", "3").Return(false, nil)
			},
			userId:  "2",
			boardId: "3",
			access:  false,
			err:     nil,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl)}
			if test.prepare != nil {
				test.prepare(&f)
			}

			serv := NewBoardsService(f.repo)

			access, err := serv.CheckWriteAccess(test.userId, test.boardId)
			if err != test.err {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if access != test.access {
				t.Errorf("\nExpected: %v\nGot: %v", test.access, access)
			}
		})
	}
}

func TestCheckReadAccess(t *testing.T) {
	type fields struct {
		repo *mocks.MockRepository
	}

	type testCase struct {
		prepare func(f *fields)
		userId  string
		boardId string
		access  bool
		err     error
	}

	tests := map[string]testCase{
		"access is allowed": {
			prepare: func(f *fields) {
				f.repo.EXPECT().CheckReadAccess("2", "3").Return(true, nil)
			},
			userId:  "2",
			boardId: "3",
			access:  true,
			err:     nil,
		},
		"access is denied": {
			prepare: func(f *fields) {
				f.repo.EXPECT().CheckReadAccess("2", "3").Return(false, nil)
			},
			userId:  "2",
			boardId: "3",
			access:  false,
			err:     nil,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl)}
			if test.prepare != nil {
				test.prepare(&f)
			}

			serv := NewBoardsService(f.repo)

			access, err := serv.CheckReadAccess(test.userId, test.boardId)
			if err != test.err {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if access != test.access {
				t.Errorf("\nExpected: %v\nGot: %v", test.access, access)
			}
		})
	}
}
