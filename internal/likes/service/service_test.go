package service

import (
	notificationsMocks "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/notifications/mocks"
	pinsMocks "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pins/mocks"
	"go.uber.org/zap"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/likes/mocks"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
)

func TestLike(t *testing.T) {
	type fields struct {
		repo              *mocks.MockRepository
		notificationsServ *notificationsMocks.MockService
		pinsRepo          *pinsMocks.MockRepository
		pinId             int
		authorId          int
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
				gomock.InOrder(
					f.repo.EXPECT().PinExists(f.pinId).Return(true, nil),
					f.repo.EXPECT().UserExists(f.authorId).Return(true, nil),
					f.repo.EXPECT().LikeExists(f.pinId, f.authorId).Return(false, nil),
					f.repo.EXPECT().Create(f.pinId, f.authorId).Return(nil),
					f.pinsRepo.EXPECT().Get(f.pinId).Return(models.Pin{Author: 12}, nil).
						MinTimes(0).MaxTimes(1),
					f.notificationsServ.EXPECT().Create(12, gomock.Any(), gomock.Any()).Return(nil).
						MinTimes(0).MaxTimes(1),
				)
			},
			pinId:    3,
			authorId: 2,
			err:      nil,
		},
		"like already exists": {
			prepare: func(f *fields) {
				gomock.InOrder(
					f.repo.EXPECT().PinExists(f.pinId).Return(true, nil),
					f.repo.EXPECT().UserExists(f.authorId).Return(true, nil),
					f.repo.EXPECT().LikeExists(f.pinId, f.authorId).Return(true, nil),
				)
			},
			pinId:    3,
			authorId: 2,
			err:      pkgErrors.ErrLikeAlreadyExists,
		},
		"db error in Create": {
			prepare: func(f *fields) {
				gomock.InOrder(
					f.repo.EXPECT().PinExists(f.pinId).Return(true, nil),
					f.repo.EXPECT().UserExists(f.authorId).Return(true, nil),
					f.repo.EXPECT().LikeExists(f.pinId, f.authorId).Return(false, nil),
					f.repo.EXPECT().Create(f.pinId, f.authorId).Return(pkgErrors.ErrDb),
				)
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

			logger, err := zap.NewDevelopment()
			if err != nil {
				t.Fatalf("can't create logger: %s", err)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				repo:              mocks.NewMockRepository(ctrl),
				notificationsServ: notificationsMocks.NewMockService(ctrl),
				pinsRepo:          pinsMocks.NewMockRepository(ctrl),
				pinId:             test.pinId,
				authorId:          test.authorId,
			}
			if test.prepare != nil {
				test.prepare(&f)
			}

			serv := NewService(f.repo, f.notificationsServ, f.pinsRepo, logger)
			err = serv.Like(test.pinId, test.authorId)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
		})
	}
}

func TestListByAuthor(t *testing.T) {
	type fields struct {
		repo              *mocks.MockRepository
		notificationsServ *notificationsMocks.MockService
		pinsRepo          *pinsMocks.MockRepository
		authorId          int
		likes             []models.Like
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

			logger, err := zap.NewDevelopment()
			if err != nil {
				t.Fatalf("can't create logger: %s", err)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				repo:              mocks.NewMockRepository(ctrl),
				notificationsServ: notificationsMocks.NewMockService(ctrl),
				pinsRepo:          pinsMocks.NewMockRepository(ctrl),
				authorId:          test.authorId,
				likes:             test.likes,
			}
			if test.prepare != nil {
				test.prepare(&f)
			}

			serv := NewService(f.repo, f.notificationsServ, f.pinsRepo, logger)
			likes, err := serv.ListByAuthor(test.authorId)
			if !errors.Is(err, test.err) {
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
		repo              *mocks.MockRepository
		notificationsServ *notificationsMocks.MockService
		pinsRepo          *pinsMocks.MockRepository
		pinId             int
		likes             []models.Like
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

			logger, err := zap.NewDevelopment()
			if err != nil {
				t.Fatalf("can't create logger: %s", err)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				repo:              mocks.NewMockRepository(ctrl),
				notificationsServ: notificationsMocks.NewMockService(ctrl),
				pinsRepo:          pinsMocks.NewMockRepository(ctrl),
				pinId:             test.pinId,
				likes:             test.likes,
			}
			if test.prepare != nil {
				test.prepare(&f)
			}

			serv := NewService(f.repo, f.notificationsServ, f.pinsRepo, logger)
			likes, err := serv.ListByPin(test.pinId)
			if !errors.Is(err, test.err) {
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
		repo              *mocks.MockRepository
		notificationsServ *notificationsMocks.MockService
		pinsRepo          *pinsMocks.MockRepository
		pinId             int
		authorId          int
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
				gomock.InOrder(
					f.repo.EXPECT().PinExists(f.pinId).Return(true, nil),
					f.repo.EXPECT().UserExists(f.authorId).Return(true, nil),
					f.repo.EXPECT().LikeExists(f.pinId, f.authorId).Return(true, nil),
					f.repo.EXPECT().Delete(f.pinId, f.authorId).Return(nil),
				)
			},
			pinId:    3,
			authorId: 2,
			err:      nil,
		},
		"like not found": {
			prepare: func(f *fields) {
				gomock.InOrder(
					f.repo.EXPECT().PinExists(f.pinId).Return(true, nil),
					f.repo.EXPECT().UserExists(f.authorId).Return(true, nil),
					f.repo.EXPECT().LikeExists(f.pinId, f.authorId).Return(false, nil),
				)
			},
			pinId:    3,
			authorId: 2,
			err:      pkgErrors.ErrLikeNotFound,
		},
		"db error in Delete": {
			prepare: func(f *fields) {
				gomock.InOrder(
					f.repo.EXPECT().PinExists(f.pinId).Return(true, nil),
					f.repo.EXPECT().UserExists(f.authorId).Return(true, nil),
					f.repo.EXPECT().LikeExists(f.pinId, f.authorId).Return(true, nil),
					f.repo.EXPECT().Delete(f.pinId, f.authorId).Return(pkgErrors.ErrDb),
				)
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

			logger, err := zap.NewDevelopment()
			if err != nil {
				t.Fatalf("can't create logger: %s", err)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				repo:              mocks.NewMockRepository(ctrl),
				notificationsServ: notificationsMocks.NewMockService(ctrl),
				pinsRepo:          pinsMocks.NewMockRepository(ctrl),
				pinId:             test.pinId,
				authorId:          test.authorId,
			}
			if test.prepare != nil {
				test.prepare(&f)
			}

			serv := NewService(f.repo, f.notificationsServ, f.pinsRepo, logger)
			err = serv.Unlike(test.pinId, test.authorId)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
		})
	}
}
