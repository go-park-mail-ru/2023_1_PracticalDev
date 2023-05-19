package service

import (
	"context"

	proto "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/images/proto"
	images "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/images/server"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
)

type service struct {
	proto.UnimplementedImageUploaderServer

	rep images.Repository
}

func NewS3Service(rep images.Repository) proto.ImageUploaderServer {
	return &service{
		rep: rep,
	}
}

func (serv *service) UploadImage(ctx context.Context, image *proto.Image) (*proto.Url, error) {
	img := &models.Image{
		ID:    image.GetID(),
		Bytes: image.GetBytes(),
	}
	url, err := serv.rep.UploadImage(img)
	if err != nil {
		return &proto.Url{}, errors.GRPCWrapper(err)
	}

	return &proto.Url{URL: url}, nil
}
