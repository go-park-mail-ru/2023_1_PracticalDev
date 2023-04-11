package service

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	pkgLikes "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/likes"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/likes/mocks"
)

func TestLike(t *testing.T) {
	type fields struct {
		repo     *mocks.MockRepository
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
		"usual": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Create(f.pinId, f.authorId).Return(nil)
			},
			pinId:    3,
			authorId: 2,
			err:      nil,
		},
		"like already exists": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Create(f.pinId, f.authorId).Return(pkgLikes.ErrLikeAlreadyExists)
			},
			pinId:    3,
			authorId: 2,
			err:      pkgLikes.ErrLikeAlreadyExists,
		},
		"db error": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Create(f.pinId, f.authorId).Return(pkgLikes.ErrDb)
			},
			pinId:    3,
			authorId: 2,
			err:      pkgLikes.ErrDb,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl), pinId: test.pinId, authorId: test.authorId}
			if test.prepare != nil {
				test.prepare(&f)
			}

			serv := NewService(f.repo)
			err := serv.Like(test.pinId, test.authorId)
			if err != test.err {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
		})
	}
}

func TestListByAuthor(t *testing.T) {
	type fields struct {
		repo     *mocks.MockRepository
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
				f.repo.EXPECT().ListByAuthor(f.authorId).Return(f.likes, nil)
			},
			authorId: 12,
			likes: []models.Like{
				{PinId: 2, AuthorId: 12, CreatedAt: time.Unix(1681163314, 0)},
				{PinId: 3, AuthorId: 12, CreatedAt: time.Unix(1681163000, 0)},
				{PinId: 4, AuthorId: 12, CreatedAt: time.Unix(1681164555, 0)},
			},
			err: nil,
		},
		"no likes": {
			prepare: func(f *fields) {
				f.repo.EXPECT().ListByAuthor(f.authorId).Return(f.likes, nil)
			},
			authorId: 12,
			likes:    []models.Like{},
			err:      nil,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl), authorId: test.authorId, likes: test.likes}
			if test.prepare != nil {
				test.prepare(&f)
			}

			serv := NewService(f.repo)
			likes, err := serv.ListByAuthor(test.authorId)
			if err != test.err {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if !reflect.DeepEqual(likes, test.likes) {
				t.Errorf("\nExpected: %v\nGot: %v", test.likes, likes)
			}
		})
	}
}

func TestListByPin(t *testing.T) {
	type fields struct {
		repo  *mocks.MockRepository
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
				f.repo.EXPECT().ListByPin(f.pinId).Return(f.likes, nil)
			},
			pinId: 12,
			likes: []models.Like{
				{PinId: 12, AuthorId: 4, CreatedAt: time.Unix(1681163314, 0)},
				{PinId: 12, AuthorId: 23, CreatedAt: time.Unix(1681163000, 0)},
				{PinId: 12, AuthorId: 2, CreatedAt: time.Unix(1681164555, 0)},
			},
			err: nil,
		},
		"no likes": {
			prepare: func(f *fields) {
				f.repo.EXPECT().ListByPin(f.pinId).Return(f.likes, nil)
			},
			pinId: 12,
			likes: []models.Like{},
			err:   nil,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl), pinId: test.pinId, likes: test.likes}
			if test.prepare != nil {
				test.prepare(&f)
			}

			serv := NewService(f.repo)
			likes, err := serv.ListByPin(test.pinId)
			if err != test.err {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if !reflect.DeepEqual(likes, test.likes) {
				t.Errorf("\nExpected: %v\nGot: %v", test.likes, likes)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	type fields struct {
		repo     *mocks.MockRepository
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
		"usual": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Delete(f.pinId, f.authorId).Return(nil)
			},
			pinId:    3,
			authorId: 2,
			err:      nil,
		},
		"like not found": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Delete(f.pinId, f.authorId).Return(pkgLikes.ErrLikeNotFound)
			},
			pinId:    3,
			authorId: 2,
			err:      pkgLikes.ErrLikeNotFound,
		},
		"db error": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Delete(f.pinId, f.authorId).Return(pkgLikes.ErrDb)
			},
			pinId:    3,
			authorId: 2,
			err:      pkgLikes.ErrDb,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl), pinId: test.pinId, authorId: test.authorId}
			if test.prepare != nil {
				test.prepare(&f)
			}

			serv := NewService(f.repo)
			err := serv.Unlike(test.pinId, test.authorId)
			if err != test.err {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
		})
	}
}
