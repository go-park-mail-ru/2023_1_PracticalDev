package client

import (
	"context"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	pkgPins "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pins"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/search"
	grpcModels "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/search/delivery/grpc/models"
	proto "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/search/delivery/grpc/proto"
	"google.golang.org/grpc"
)

type client struct {
	searchClient proto.SearchEngineClient

	pinServ pkgPins.Service
}

func NewSearchClient(con *grpc.ClientConn, pinServ pkgPins.Service) search.Service {
	return &client{
		searchClient: proto.NewSearchEngineClient(con),
		pinServ:      pinServ,
	}
}

func (c *client) Get(userId int, query string) (models.SearchRes, error) {
	q := grpcModels.NewProtoQuery(query)

	resp, err := c.searchClient.Get(context.TODO(), q)
	res := *grpcModels.NewQueryResult(resp)
	if err != nil {
		return res, err
	}

	for i := range res.Pins {
		err := c.pinServ.SetLikedField(&res.Pins[i], userId)
		if err != nil {
			return res, err
		}
	}

	return res, nil
}
