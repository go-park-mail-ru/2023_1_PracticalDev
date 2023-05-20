package service

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/followings"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/notifications"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/constants"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
)

type service struct {
	rep               followings.Repository
	notificationsServ notifications.Service
}

func NewService(rep followings.Repository, notificationsServ notifications.Service) followings.Service {
	return &service{rep: rep, notificationsServ: notificationsServ}
}

func (serv *service) Follow(followerID, followeeID int) error {
	if followerID == followeeID {
		return pkgErrors.ErrSameUserId
	}

	exists, err := serv.rep.UserExists(followerID)
	if err != nil {
		return err
	}
	if !exists {
		return pkgErrors.ErrUserNotFound
	}

	exists, err = serv.rep.UserExists(followeeID)
	if err != nil {
		return err
	}
	if !exists {
		return pkgErrors.ErrUserNotFound
	}

	exists, err = serv.rep.FollowingExists(followerID, followeeID)
	if err != nil {
		return err
	}
	if exists {
		return pkgErrors.ErrFollowingAlreadyExists
	}

	err = serv.rep.Create(followerID, followeeID)
	if err != nil {
		return err
	}

	go func(followeeID int) {
		_ = serv.notificationsServ.Create(followeeID, constants.NewFollower, models.NewFollowerNotification{
			FollowerID: followerID,
		})
	}(followeeID)

	return nil
}

func (serv *service) Unfollow(followerId, followeeId int) error {
	if followerId == followeeId {
		return pkgErrors.ErrSameUserId
	}

	exists, err := serv.rep.UserExists(followerId)
	if err != nil {
		return err
	}
	if !exists {
		return pkgErrors.ErrUserNotFound
	}

	exists, err = serv.rep.UserExists(followeeId)
	if err != nil {
		return err
	}
	if !exists {
		return pkgErrors.ErrUserNotFound
	}

	exists, err = serv.rep.FollowingExists(followerId, followeeId)
	if err != nil {
		return err
	}
	if !exists {
		return pkgErrors.ErrFollowingNotFound
	}

	return serv.rep.Delete(followerId, followeeId)
}

func (serv *service) GetFollowers(userId int) ([]followings.Follower, error) {
	exists, err := serv.rep.UserExists(userId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, pkgErrors.ErrUserNotFound
	}

	return serv.rep.GetFollowers(userId)
}

func (serv *service) GetFollowees(userId int) ([]followings.Followee, error) {
	exists, err := serv.rep.UserExists(userId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, pkgErrors.ErrUserNotFound
	}

	return serv.rep.GetFollowees(userId)
}
