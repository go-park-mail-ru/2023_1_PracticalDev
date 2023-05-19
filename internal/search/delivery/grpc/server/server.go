package server

import (
	"context"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/search"
	grpcModels "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/search/delivery/grpc/models"
	proto "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/search/delivery/grpc/proto"
)

type server struct {
	proto.UnimplementedSearchEngineServer

	rep search.Repository
}

func NewSearchServer(rep search.Repository) proto.SearchEngineServer {
	return &server{
		rep: rep,
	}
}

func (s *server) Get(ctx context.Context, q *proto.Query) (*proto.QueryResult, error) {
	req := grpcModels.NewQuery(q)

	res, err := s.rep.Get(req)
	if err != nil {
		return &proto.QueryResult{}, errors.GRPCWrapper(err)
	}

	return grpcModels.NewProtoQueryResult(&res), err
}
