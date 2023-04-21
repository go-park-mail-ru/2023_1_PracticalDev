package followings

import (
	"github.com/pkg/errors"
)

type Follower struct {
	Id           int    `json:"id"`
	Username     string `json:"username"`
	Name         string `json:"name"`
	ProfileImage string `json:"profile_image"`
	WebsiteUrl   string `json:"website_url"`
}

type Followee struct {
	Id           int    `json:"id"`
	Username     string `json:"username"`
	Name         string `json:"name"`
	ProfileImage string `json:"profile_image"`
	WebsiteUrl   string `json:"website_url"`
}

var (
	ErrDb                     = errors.New("db error")
	ErrUserNotFound           = errors.New("no such user")
	ErrFollowingNotFound      = errors.New("no such following")
	ErrFollowingAlreadyExists = errors.New("following already exists")
)

type Repository interface {
	Create(followerId, followeeId int) error
	Delete(followerId, followeeId int) error

	GetFollowees(userId int) ([]Followee, error)
	GetFollowers(userId int) ([]Follower, error)

	FollowingExists(followerId, followeeId int) (bool, error)
	UserExists(userId int) (bool, error)
}
