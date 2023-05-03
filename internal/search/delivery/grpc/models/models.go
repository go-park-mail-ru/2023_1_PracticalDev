package models

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	proto "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/search/delivery/grpc/proto"
)

func NewProtoQuery(q string) *proto.Query {
	return &proto.Query{
		Query: q,
	}
}

func NewQuery(q *proto.Query) string {
	return q.Query
}

func NewProtoQueryResult(q *models.SearchRes) *proto.QueryResult {
	pins := make([]*proto.Pin, 0)
	for _, pin := range q.Pins {
		pins = append(pins, &proto.Pin{
			Id:          int64(pin.Id),
			Title:       pin.Title,
			Description: pin.Description,
			MediaSource: pin.MediaSource,
			NumLikes:    int64(pin.NumLikes),
			Liked:       pin.Liked,
			Author:      int64(pin.Author),
		})
	}
	boards := make([]*proto.Board, 0)
	for _, board := range q.Boards {
		boards = append(boards, &proto.Board{
			Id:          int64(board.Id),
			Name:        board.Name,
			Description: board.Description,
			Privacy:     board.Privacy,
			UserId:      int64(board.UserId),
		})
	}
	users := make([]*proto.Profile, 0)
	for _, user := range q.Users {
		users = append(users, &proto.Profile{
			Id:           int64(user.Id),
			Username:     user.Username,
			Name:         user.Name,
			ProfileImage: user.ProfileImage,
			WebsiteUrl:   user.WebsiteUrl,
		})
	}
	return &proto.QueryResult{
		Pins:   pins,
		Boards: boards,
		Users:  users,
	}
}

func NewQueryResult(q *proto.QueryResult) *models.SearchRes {
	pins := make([]models.Pin, 0)
	for _, pin := range q.Pins {
		pins = append(pins, models.Pin{
			Id:          int(pin.Id),
			Title:       pin.Title,
			Description: pin.Description,
			MediaSource: pin.MediaSource,
			NumLikes:    int(pin.NumLikes),
			Liked:       pin.Liked,
			Author:      int(pin.Author),
		})
	}
	boards := make([]models.Board, 0)
	for _, board := range q.Boards {
		boards = append(boards, models.Board{
			Id:          int(board.Id),
			Name:        board.Name,
			Description: board.Description,
			Privacy:     board.Privacy,
			UserId:      int(board.UserId),
		})
	}
	users := make([]models.Profile, 0)
	for _, user := range q.Users {
		users = append(users, models.Profile{
			Id:           int(user.Id),
			Username:     user.Username,
			Name:         user.Name,
			ProfileImage: user.ProfileImage,
			WebsiteUrl:   user.WebsiteUrl,
		})
	}
	return &models.SearchRes{
		Users:  users,
		Boards: boards,
		Pins:   pins,
	}
}
