package service

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/followings"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/followings/mocks"
	notificationsMocks "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/notifications/mocks"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"reflect"
	"testing"
)

func TestService_Follow(t *testing.T) {
	type fields struct {
		repo              *mocks.MockRepository
		notificationsServ *notificationsMocks.MockService
		followerID        int
		followeeID        int
	}

	type testCase struct {
		prepare    func(f *fields)
		followerID int
		followeeID int
		err        error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				gomock.InOrder(
					f.repo.EXPECT().UserExists(f.followerID).Return(true, nil),
					f.repo.EXPECT().UserExists(f.followeeID).Return(true, nil),
					f.repo.EXPECT().FollowingExists(f.followerID, f.followeeID).Return(false, nil),
					f.repo.EXPECT().Create(f.followerID, f.followeeID).Return(nil),
					f.notificationsServ.EXPECT().Create(f.followeeID, gomock.Any(), gomock.Any()).Return(nil).
						MinTimes(0).MaxTimes(1),
				)
			},
			followerID: 3,
			followeeID: 2,
			err:        nil,
		},
		"follower not found": {
			prepare: func(f *fields) {
				gomock.InOrder(
					f.repo.EXPECT().UserExists(f.followerID).Return(false, nil),
				)
			},
			followerID: 3,
			followeeID: 2,
			err:        pkgErrors.ErrUserNotFound,
		},
		"followee not found": {
			prepare: func(f *fields) {
				gomock.InOrder(
					f.repo.EXPECT().UserExists(f.followerID).Return(true, nil),
					f.repo.EXPECT().UserExists(f.followeeID).Return(false, nil),
				)
			},
			followerID: 3,
			followeeID: 2,
			err:        pkgErrors.ErrUserNotFound,
		},
		"following already exists": {
			prepare: func(f *fields) {
				gomock.InOrder(
					f.repo.EXPECT().UserExists(f.followerID).Return(true, nil),
					f.repo.EXPECT().UserExists(f.followeeID).Return(true, nil),
					f.repo.EXPECT().FollowingExists(f.followerID, f.followeeID).Return(true, nil),
				)
			},
			followerID: 3,
			followeeID: 2,
			err:        pkgErrors.ErrFollowingAlreadyExists,
		},
		"db error in Create": {
			prepare: func(f *fields) {
				gomock.InOrder(
					f.repo.EXPECT().UserExists(f.followerID).Return(true, nil),
					f.repo.EXPECT().UserExists(f.followeeID).Return(true, nil),
					f.repo.EXPECT().FollowingExists(f.followerID, f.followeeID).Return(false, nil),
					f.repo.EXPECT().Create(f.followerID, f.followeeID).Return(pkgErrors.ErrDb),
				)
			},
			followerID: 3,
			followeeID: 2,
			err:        pkgErrors.ErrDb,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				repo:              mocks.NewMockRepository(ctrl),
				notificationsServ: notificationsMocks.NewMockService(ctrl),
				followerID:        test.followerID,
				followeeID:        test.followeeID,
			}
			if test.prepare != nil {
				test.prepare(&f)
			}

			serv := NewService(f.repo, f.notificationsServ)
			err := serv.Follow(test.followerID, test.followeeID)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
		})
	}
}

func TestService_Unfollow(t *testing.T) {
	type fields struct {
		repo              *mocks.MockRepository
		notificationsServ *notificationsMocks.MockService
		followerID        int
		followeeID        int
	}

	type testCase struct {
		prepare    func(f *fields)
		followerID int
		followeeID int
		err        error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				gomock.InOrder(
					f.repo.EXPECT().UserExists(f.followerID).Return(true, nil),
					f.repo.EXPECT().UserExists(f.followeeID).Return(true, nil),
					f.repo.EXPECT().FollowingExists(f.followerID, f.followeeID).Return(true, nil),
					f.repo.EXPECT().Delete(f.followerID, f.followeeID).Return(nil),
				)
			},
			followerID: 3,
			followeeID: 2,
			err:        nil,
		},
		"follower not found": {
			prepare: func(f *fields) {
				gomock.InOrder(
					f.repo.EXPECT().UserExists(f.followerID).Return(false, nil),
				)
			},
			followerID: 3,
			followeeID: 2,
			err:        pkgErrors.ErrUserNotFound,
		},
		"followee not found": {
			prepare: func(f *fields) {
				gomock.InOrder(
					f.repo.EXPECT().UserExists(f.followerID).Return(true, nil),
					f.repo.EXPECT().UserExists(f.followeeID).Return(false, nil),
				)
			},
			followerID: 3,
			followeeID: 2,
			err:        pkgErrors.ErrUserNotFound,
		},
		"following not found": {
			prepare: func(f *fields) {
				gomock.InOrder(
					f.repo.EXPECT().UserExists(f.followerID).Return(true, nil),
					f.repo.EXPECT().UserExists(f.followeeID).Return(true, nil),
					f.repo.EXPECT().FollowingExists(f.followerID, f.followeeID).Return(false, nil),
				)
			},
			followerID: 3,
			followeeID: 2,
			err:        pkgErrors.ErrFollowingNotFound,
		},
		"db error in Delete": {
			prepare: func(f *fields) {
				gomock.InOrder(
					f.repo.EXPECT().UserExists(f.followerID).Return(true, nil),
					f.repo.EXPECT().UserExists(f.followeeID).Return(true, nil),
					f.repo.EXPECT().FollowingExists(f.followerID, f.followeeID).Return(true, nil),
					f.repo.EXPECT().Delete(f.followerID, f.followeeID).Return(pkgErrors.ErrDb),
				)
			},
			followerID: 3,
			followeeID: 2,
			err:        pkgErrors.ErrDb,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				repo:              mocks.NewMockRepository(ctrl),
				notificationsServ: notificationsMocks.NewMockService(ctrl),
				followerID:        test.followerID,
				followeeID:        test.followeeID,
			}
			if test.prepare != nil {
				test.prepare(&f)
			}

			serv := NewService(f.repo, f.notificationsServ)
			err := serv.Unfollow(test.followerID, test.followeeID)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
		})
	}
}

func TestService_GetFollowers(t *testing.T) {
	type fields struct {
		repo              *mocks.MockRepository
		notificationsServ *notificationsMocks.MockService
		userID            int
		followers         []followings.Follower
	}

	type testCase struct {
		prepare   func(f *fields)
		userID    int
		followers []followings.Follower
		err       error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				gomock.InOrder(
					f.repo.EXPECT().UserExists(f.userID).Return(true, nil),
					f.repo.EXPECT().GetFollowers(f.userID).Return(f.followers, nil))
			},
			userID: 12,
			followers: []followings.Follower{
				{Id: 2, Username: "vasua", Name: "Vasya", ProfileImage: "vasya.jpg", WebsiteUrl: "vasya.com"},
				{Id: 3, Username: "kolya", Name: "Kolya", ProfileImage: "kolya.jpg", WebsiteUrl: "kolya.com"},
				{Id: 4, Username: "sasha", Name: "Sasha", ProfileImage: "sasha.jpg", WebsiteUrl: "sasha.com"},
			},
			err: nil,
		},
		"no followers": {
			prepare: func(f *fields) {
				gomock.InOrder(
					f.repo.EXPECT().UserExists(f.userID).Return(true, nil),
					f.repo.EXPECT().GetFollowers(f.userID).Return(f.followers, nil))
			},
			userID:    12,
			followers: []followings.Follower{},
			err:       nil,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				repo:              mocks.NewMockRepository(ctrl),
				notificationsServ: notificationsMocks.NewMockService(ctrl),
				userID:            test.userID,
				followers:         test.followers,
			}
			if test.prepare != nil {
				test.prepare(&f)
			}

			serv := NewService(f.repo, f.notificationsServ)
			followers, err := serv.GetFollowers(test.userID)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if !reflect.DeepEqual(followers, test.followers) {
				t.Errorf("\nExpected: %v\nGot: %v", test.followers, followers)
			}
		})
	}
}

func TestService_GetFollowees(t *testing.T) {
	type fields struct {
		repo              *mocks.MockRepository
		notificationsServ *notificationsMocks.MockService
		userID            int
		followees         []followings.Followee
	}

	type testCase struct {
		prepare   func(f *fields)
		userID    int
		followees []followings.Followee
		err       error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				gomock.InOrder(
					f.repo.EXPECT().UserExists(f.userID).Return(true, nil),
					f.repo.EXPECT().GetFollowees(f.userID).Return(f.followees, nil))
			},
			userID: 12,
			followees: []followings.Followee{
				{Id: 2, Username: "vasua", Name: "Vasya", ProfileImage: "vasya.jpg", WebsiteUrl: "vasya.com"},
				{Id: 3, Username: "kolya", Name: "Kolya", ProfileImage: "kolya.jpg", WebsiteUrl: "kolya.com"},
				{Id: 4, Username: "sasha", Name: "Sasha", ProfileImage: "sasha.jpg", WebsiteUrl: "sasha.com"},
			},
			err: nil,
		},
		"no followees": {
			prepare: func(f *fields) {
				gomock.InOrder(
					f.repo.EXPECT().UserExists(f.userID).Return(true, nil),
					f.repo.EXPECT().GetFollowees(f.userID).Return(f.followees, nil))
			},
			userID:    12,
			followees: []followings.Followee{},
			err:       nil,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				repo:              mocks.NewMockRepository(ctrl),
				notificationsServ: notificationsMocks.NewMockService(ctrl),
				userID:            test.userID,
				followees:         test.followees,
			}
			if test.prepare != nil {
				test.prepare(&f)
			}

			serv := NewService(f.repo, f.notificationsServ)
			followees, err := serv.GetFollowees(test.userID)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if !reflect.DeepEqual(followees, test.followees) {
				t.Errorf("\nExpected: %v\nGot: %v", test.followees, followees)
			}
		})
	}
}
