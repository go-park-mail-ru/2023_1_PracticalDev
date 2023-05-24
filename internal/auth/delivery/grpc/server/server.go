package server

import (
	"context"
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth"
	protomodels "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth/delivery/grpc/models"
	proto "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth/delivery/grpc/proto"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
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
	fmt.Println("here", err)
	if err != nil {
		return &proto.User{}, pkgErrors.GRPCWrapper(err)
	}
	return protomodels.NewProtoUser(&user), err
}

func (serv *server) Register(ctx context.Context, user *proto.User) (*proto.LoginParams, error) {
	err := serv.rep.Register(protomodels.NewUser(user))
	if err != nil {
		return &proto.LoginParams{}, pkgErrors.GRPCWrapper(err)
	}
	return &proto.LoginParams{Email: user.GetEmail(), Password: user.GetHashedPassword()}, err
}

func (serv *server) SetSession(ctx context.Context, params *proto.SessionSetParams) (*proto.Nothing, error) {
	err := serv.rep.SetSession(params.GetToken(), protomodels.NewSession(params.GetSession()), time.Duration(params.Experation))
	if err != nil {
		return &proto.Nothing{}, pkgErrors.GRPCWrapper(err)
	}
	return &proto.Nothing{}, err
}

func (serv *server) CheckAuth(ctx context.Context, params *proto.SessionCheckParams) (*proto.User, error) {
	user, err := serv.rep.CheckAuth(params.GetUserId(), params.GetSessionId())
	if err != nil {
		return &proto.User{}, pkgErrors.GRPCWrapper(err)
	}
	return protomodels.NewProtoUser(&user), err
}
func (serv *server) DeleteSession(ctx context.Context, params *proto.SessionCheckParams) (*proto.Nothing, error) {
	err := serv.rep.DeleteSession(params.GetUserId(), params.GetSessionId())
	if err != nil {
		return &proto.Nothing{}, pkgErrors.GRPCWrapper(err)
	}
	return &proto.Nothing{}, err
}
