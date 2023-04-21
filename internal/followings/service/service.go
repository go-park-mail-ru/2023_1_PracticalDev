package service

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/followings"
)

type service struct {
	rep followings.Repository
}

func NewService(rep followings.Repository) followings.Service {
	return &service{rep}
}

func (serv *service) Follow(followerId, followeeId int) error {
	exists, err := serv.rep.UserExists(followerId)
	if err != nil {
		return err
	}
	if !exists {
		return followings.ErrUserNotFound
	}

	exists, err = serv.rep.UserExists(followeeId)
	if err != nil {
		return err
	}
	if !exists {
		return followings.ErrUserNotFound
	}

	exists, err = serv.rep.FollowingExists(followerId, followeeId)
	if err != nil {
		return err
	}
	if exists {
		return followings.ErrFollowingAlreadyExists
	}

	return serv.rep.Create(followerId, followeeId)
}

func (serv *service) Unfollow(followerId, followeeId int) error {
	exists, err := serv.rep.UserExists(followerId)
	if err != nil {
		return err
	}
	if !exists {
		return followings.ErrUserNotFound
	}

	exists, err = serv.rep.UserExists(followeeId)
	if err != nil {
		return err
	}
	if !exists {
		return followings.ErrUserNotFound
	}

	exists, err = serv.rep.FollowingExists(followerId, followeeId)
	if err != nil {
		return err
	}
	if !exists {
		return followings.ErrFollowingNotFound
	}

	return serv.rep.Delete(followerId, followeeId)
}

func (serv *service) GetFollowers(userId int) ([]followings.Follower, error) {
	println("dfsdgdf")
	return serv.rep.GetFollowers(userId)
}

func (serv *service) GetFollowees(userId int) ([]followings.Followee, error) {
	return serv.rep.GetFollowees(userId)
}
