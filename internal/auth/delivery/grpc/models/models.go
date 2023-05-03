package models

import (
	"time"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth"
	proto "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth/delivery/grpc/proto"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
)

func NewProtoUser(user *models.User) *proto.User {
	return &proto.User{
		ID:             int64(user.Id),
		Username:       user.Username,
		Email:          user.Email,
		HashedPassword: user.HashedPassword,
		Name:           user.Name,
		ProfileImage:   user.ProfileImage,
		WebsiteUrl:     user.WebsiteUrl,
		AccountType:    user.AccountType,
	}
}

func NewUser(user *proto.User) *models.User {
	return &models.User{
		Id:             int(user.GetID()),
		Username:       user.GetUsername(),
		Email:          user.GetEmail(),
		HashedPassword: user.GetHashedPassword(),
		Name:           user.GetName(),
		ProfileImage:   user.GetProfileImage(),
		WebsiteUrl:     user.GetWebsiteUrl(),
		AccountType:    user.GetAccountType(),
	}
}

func NewProtoSessionParams(params *auth.SessionParams) *proto.SessionParams {
	return &proto.SessionParams{
		Token:      params.Token,
		LivingTime: params.LivingTime.Nanoseconds(),
	}
}

func NewSessionParams(params *proto.SessionParams) *auth.SessionParams {
	return &auth.SessionParams{
		Token:      params.GetToken(),
		LivingTime: time.Duration(params.LivingTime),
	}
}

func NewProtoSession(session *models.Session) *proto.Session {
	return &proto.Session{
		UserId:    int64(session.UserId),
		UserEmail: session.UserEmail,
	}
}

func NewSession(session *proto.Session) *models.Session {
	return &models.Session{
		UserId:    int(session.GetUserId()),
		UserEmail: session.GetUserEmail(),
	}
}
