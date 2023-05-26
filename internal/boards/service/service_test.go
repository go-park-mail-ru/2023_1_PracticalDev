package service

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"

	_boards "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/boards"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/boards/mocks"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	pinsMock "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pins/mocks"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
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
			serv := NewBoardsService(f.repo, pinsMock.NewMockService(ctrl))

			board, err := serv.Create(&test.params)
			if !errors.Is(err, test.err) {
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

			serv := NewBoardsService(f.repo, pinsMock.NewMockService(ctrl))

			boards, err := serv.List(test.userId)
			if !errors.Is(err, test.err) {
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
				f.repo.EXPECT().Get(3).Return(models.Board{}, pkgErrors.ErrBoardNotFound)
			},
			id:    3,
			board: models.Board{},
			err:   pkgErrors.ErrBoardNotFound,
		},
		"negative board id param": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Get(-1).Return(models.Board{}, pkgErrors.ErrBoardNotFound)
			},
			id:    -1,
			board: models.Board{},
			err:   pkgErrors.ErrBoardNotFound,
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

			serv := NewBoardsService(f.repo, pinsMock.NewMockService(ctrl))

			board, err := serv.Get(test.id)
			if !errors.Is(err, test.err) {
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

			serv := NewBoardsService(f.repo, pinsMock.NewMockService(ctrl))

			board, err := serv.FullUpdate(&test.params)
			if !errors.Is(err, test.err) {
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

			serv := NewBoardsService(f.repo, pinsMock.NewMockService(ctrl))

			board, err := serv.PartialUpdate(&test.params)
			if !errors.Is(err, test.err) {
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
				f.repo.EXPECT().Delete(3).Return(pkgErrors.ErrBoardNotFound)
			},
			id:  3,
			err: pkgErrors.ErrBoardNotFound,
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

			serv := NewBoardsService(f.repo, pinsMock.NewMockService(ctrl))

			err := serv.Delete(test.id)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
		})
	}
}

func TestPinsList(t *testing.T) {
	type fields struct {
		repo     *mocks.MockRepository
		pinsServ *pinsMock.MockService
	}

	type testCase struct {
		prepare func(f *fields)
		page    int
		limit   int
		boardId int
		userId  int
		pins    []models.Pin
		err     error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				gomock.InOrder(
					f.repo.EXPECT().PinsList(3, 1, 30).Return([]models.Pin{
						{Id: 1, Title: "t1", MediaSource: "ms_url1", Description: "d1",
							Author: models.Profile{
								Id:           12,
								Username:     "un1",
								Name:         "n1",
								ProfileImage: "pi1",
								WebsiteUrl:   "wu1",
							}},
						{Id: 2, Title: "t2", MediaSource: "ms_url2", Description: "d2",
							Author: models.Profile{
								Id:           13,
								Username:     "un2",
								Name:         "n2",
								ProfileImage: "pi2",
								WebsiteUrl:   "wu2",
							}},
						{Id: 3, Title: "t3", MediaSource: "ms_url3", Description: "d3",
							Author: models.Profile{
								Id:           14,
								Username:     "un3",
								Name:         "n3",
								ProfileImage: "pi3",
								WebsiteUrl:   "wu3",
							}},
					}, nil),
					f.pinsServ.EXPECT().SetLikedField(gomock.Any(), 10).Return(nil).Times(3),
				)
			},
			page:    1,
			limit:   30,
			boardId: 3,
			userId:  10,
			pins: []models.Pin{
				{Id: 1, Title: "t1", MediaSource: "ms_url1", Description: "d1",
					Author: models.Profile{
						Id:           12,
						Username:     "un1",
						Name:         "n1",
						ProfileImage: "pi1",
						WebsiteUrl:   "wu1",
					}},
				{Id: 2, Title: "t2", MediaSource: "ms_url2", Description: "d2",
					Author: models.Profile{
						Id:           13,
						Username:     "un2",
						Name:         "n2",
						ProfileImage: "pi2",
						WebsiteUrl:   "wu2",
					}},
				{Id: 3, Title: "t3", MediaSource: "ms_url3", Description: "d3",
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
		"no boards": {
			prepare: func(f *fields) {
				f.repo.EXPECT().PinsList(3, 1, 30).Return([]models.Pin{}, nil)
			},
			boardId: 3,
			userId:  10,
			page:    1,
			limit:   30,
			pins:    []models.Pin{},
			err:     nil,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl), pinsServ: pinsMock.NewMockService(ctrl)}
			if test.prepare != nil {
				test.prepare(&f)
			}

			serv := NewBoardsService(f.repo, f.pinsServ)

			pins, err := serv.PinsList(test.userId, test.boardId, test.page, test.limit)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if !reflect.DeepEqual(pins, test.pins) {
				t.Errorf("\nExpected: %v\nGot: %v", test.pins, pins)
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

			serv := NewBoardsService(f.repo, pinsMock.NewMockService(ctrl))

			access, err := serv.CheckWriteAccess(test.userId, test.boardId)
			if !errors.Is(err, test.err) {
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

			serv := NewBoardsService(f.repo, pinsMock.NewMockService(ctrl))

			access, err := serv.CheckReadAccess(test.userId, test.boardId)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if access != test.access {
				t.Errorf("\nExpected: %v\nGot: %v", test.access, access)
			}
		})
	}
}
