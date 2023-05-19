package client

import (
	"context"
	"fmt"
	"os"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
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
		return "", errors.RestoreHTTPError(errors.GRPCUnwrapper(err))
	}

	res := models.NewStringMessage(resp)
	return res, nil
}

func (c *client) Create(url string) (string, error) {
	resp, err := c.shortenerClient.Create(context.Background(), models.NewProtoStringMessage(url))
	if err != nil {
		return "", errors.RestoreHTTPError(errors.GRPCUnwrapper(err))
	}

	res := models.NewStringMessage(resp)
	return res, nil
}

func (c *client) CreatePinLink(id int) (string, error) {
	if os.Getenv("SHORT_HOST") == "localhost:8091" {
		return c.Create(fmt.Sprintf("http://localhost/pins/%d", id))
	}
	return c.Create(fmt.Sprintf("https://pickpin.ru/pins/%d", id))
}
