package followings

type Service interface {
	Follow(followerId, followeeId int) error
	Unfollow(followerId, followeeId int) error

	GetFollowees(userId int) ([]Followee, error)
	GetFollowers(userId int) ([]Follower, error)
}
