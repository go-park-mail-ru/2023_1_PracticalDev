package http

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/followings"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/xss"
)

//go:generate easyjson -all -snake_case api_models.go

// API responses
type followersResponse struct {
	Followers []followings.Follower `json:"followers"`
}

func newFollowersResponse(followers []followings.Follower) *followersResponse {
	for i := range followers {
		followers[i].Username = xss.Sanitize(followers[i].Username)
		followers[i].Name = xss.Sanitize(followers[i].Name)
	}

	return &followersResponse{
		Followers: followers,
	}
}

type followeesResponse struct {
	Followees []followings.Followee `json:"followees"`
}

func newFolloweesResponse(followees []followings.Followee) *followeesResponse {
	for i := range followees {
		followees[i].Username = xss.Sanitize(followees[i].Username)
		followees[i].Name = xss.Sanitize(followees[i].Name)
	}

	return &followeesResponse{
		Followees: followees,
	}
}
