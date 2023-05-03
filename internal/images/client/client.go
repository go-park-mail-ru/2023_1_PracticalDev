package client

import (
	"context"
	"fmt"

	proto "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/images/proto"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	"google.golang.org/grpc"
)

type ImageClient interface {
	UploadImage(ctx context.Context, image *models.Image) (string, error)
}

type client struct {
	imageClient proto.ImageUploaderClient
}

func NewImageUploaderClient(con *grpc.ClientConn) ImageClient {
	return &client{imageClient: proto.NewImageUploaderClient(con)}
}

func (client *client) UploadImage(ctx context.Context, image *models.Image) (string, error) {
	img := proto.Image{
		ID:    image.ID,
		Bytes: image.Bytes,
	}
	url, err := client.imageClient.UploadImage(ctx, &img)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return url.GetURL(), nil
}
