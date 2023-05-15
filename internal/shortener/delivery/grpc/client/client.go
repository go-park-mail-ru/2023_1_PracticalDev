package client

import (
	"context"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/shortener"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/shortener/delivery/grpc/models"
	proto "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/shortener/delivery/grpc/proto"
	"google.golang.org/grpc"
)

type client struct {
	shortenerClient proto.ShortenerClient
}

func NewShortenerClient(con *grpc.ClientConn) shortener.ShortenerService {
	return &client{
		shortenerClient: proto.NewShortenerClient(con),
	}
}

func (c *client) Get(hash string) (string, error) {
	resp, err := c.shortenerClient.Get(context.Background(), models.NewProtoStringMessage(hash))
	if err != nil {
		return "", err
	}

	res := models.NewStringMessage(resp)
	return res, nil
}

func (c *client) Create(url string) (string, error) {
	resp, err := c.shortenerClient.Create(context.Background(), models.NewProtoStringMessage(url))
	if err != nil {
		return "", err
	}

	res := models.NewStringMessage(resp)
	return res, nil
}
