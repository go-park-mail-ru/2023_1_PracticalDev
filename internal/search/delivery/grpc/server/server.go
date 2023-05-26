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

func (s *server) Search(ctx context.Context, q *proto.Query) (*proto.SearchResult, error) {
	req := grpcModels.NewQuery(q)

	res, err := s.rep.Search(req)
	if err != nil {
		return &proto.SearchResult{}, errors.GRPCWrapper(err)
	}

	return grpcModels.NewProtoSearchResult(&res), err
}

func (s *server) Suggestions(ctx context.Context, q *proto.Query) (*proto.SuggestionsResult, error) {
	req := grpcModels.NewQuery(q)

	res, err := s.rep.Suggestions(req)
	if err != nil {
		return &proto.SuggestionsResult{}, errors.GRPCWrapper(err)
	}

	return grpcModels.NewProtoSuggestionsResult(res), err
}
