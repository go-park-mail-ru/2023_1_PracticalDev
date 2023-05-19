package server

import (
	"context"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/shortener"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/shortener/delivery/grpc/models"
	proto "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/shortener/delivery/grpc/proto"
)

type server struct {
	proto.UnimplementedShortenerServer

	rep shortener.ShortenerRepository
}

func NewShortenerServer(rep shortener.ShortenerRepository) proto.ShortenerServer {
	return &server{
		rep: rep,
	}
}

func (s *server) Get(ctx context.Context, url *proto.StringMessage) (*proto.StringMessage, error) {
	hash := models.NewStringMessage(url)
	res, err := s.rep.Get(hash)
	if err != nil {
		return &proto.StringMessage{}, errors.GRPCWrapper(err)
	}
	return models.NewProtoStringMessage(res), err
}

func (s *server) Create(ctx context.Context, url *proto.StringMessage) (*proto.StringMessage, error) {
	msg := models.NewStringMessage(url)
	hash, err := s.rep.Create(msg)
	if err != nil {
		return &proto.StringMessage{}, errors.GRPCWrapper(err)
	}
	return models.NewProtoStringMessage(hash), err
}
