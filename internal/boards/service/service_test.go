package service

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/boards"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/boards/mocks"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	"github.com/golang/mock/gomock"
	"testing"
)

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
				f.repo.EXPECT().Get(3).Return(models.Board{}, boards.ErrBoardNotFound)
			},
			id:    3,
			board: models.Board{},
			err:   boards.ErrBoardNotFound,
		},
		"negative board id param": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Get(-1).Return(models.Board{}, boards.ErrBoardNotFound)
			},
			id:    -1,
			board: models.Board{},
			err:   boards.ErrBoardNotFound,
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
