package http

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/followings"
)

// API responses
type followersResponse struct {
	Followers []followings.Follower `json:"followers"`
}

type followeesResponse struct {
	Followees []followings.Followee `json:"followees"`
}
