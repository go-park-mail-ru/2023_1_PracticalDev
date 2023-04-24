package service

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/followings"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
)

type service struct {
	rep followings.Repository
}

func NewService(rep followings.Repository) followings.Service {
	return &service{rep}
}

func (serv *service) Follow(followerId, followeeId int) error {
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
	if exists {
		return pkgErrors.ErrFollowingAlreadyExists
	}

	return serv.rep.Create(followerId, followeeId)
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
