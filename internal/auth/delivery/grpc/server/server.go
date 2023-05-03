package server

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth"
	protomodels "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth/delivery/grpc/models"
	proto "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth/delivery/grpc/proto"
)

type server struct {
	proto.UnimplementedAuthenficatorServer

	rep auth.Repository
}

func NewAuthServer(rep auth.Repository) proto.AuthenficatorServer {
	return &server{
		rep: rep,
	}
}

func (serv *server) Authenticate(ctx context.Context, loginParams *proto.LoginParams) (*proto.User, error) {
	user, err := serv.rep.Authenticate(loginParams.GetEmail(), loginParams.GetPassword())

	return protomodels.NewProtoUser(&user), err
}

func (serv *server) Register(ctx context.Context, user *proto.User) (*proto.LoginParams, error) {
	err := serv.rep.Register(protomodels.NewUser(user))

	return &proto.LoginParams{Email: user.GetEmail(), Password: user.GetHashedPassword()}, err
}

func (serv *server) SetSession(ctx context.Context, params *proto.SessionSetParams) (*proto.Nothing, error) {
	err := serv.rep.SetSession(params.GetToken(), protomodels.NewSession(params.GetSession()), time.Duration(params.Experation))

	return &proto.Nothing{}, err
}

func (serv *server) CheckAuth(ctx context.Context, params *proto.SessionCheckParams) (*proto.User, error) {
	user, err := serv.rep.CheckAuth(params.GetUserId(), params.GetSessionId())

	return protomodels.NewProtoUser(&user), err
}
func (serv *server) DeleteSession(ctx context.Context, params *proto.SessionCheckParams) (*proto.Nothing, error) {
	err := serv.rep.DeleteSession(params.GetUserId(), params.GetSessionId())

	return &proto.Nothing{}, err
}
