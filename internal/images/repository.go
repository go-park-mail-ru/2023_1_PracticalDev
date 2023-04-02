package images

import (
	"bytes"
	"context"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/config"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
)

type Repository interface {
	UploadImage(image *models.Image) (string, error)
}

type repository struct {
	log        log.Logger
	uploader   *manager.Uploader
	bucketName string
}

func NewRepository(log log.Logger) (Repository, error) {
	log.Info("Connecting to s3 image bucket...")

	conf, err := awsConfig.LoadDefaultConfig(context.Background(), awsConfig.WithDefaultRegion("ru-msk"))

	if err != nil {
		log.Error("Failed to create S3 connection, reason: ", err.Error())
		return &repository{}, err
	}
	client := s3.NewFromConfig(conf, func(o *s3.Options) {
		o.EndpointResolver = s3.EndpointResolverFromURL("https://hb.bizmrg.com")
	})

	uploader := manager.NewUploader(client)

	log.Info("S3 bucket connection created successfully")

	return &repository{
		log:        log,
		uploader:   uploader,
		bucketName: config.Get("BUCKET_NAME"),
	}, nil
}

func (rep *repository) UploadImage(image *models.Image) (string, error) {
	rep.log.Debug("Initiating s3 transaction...")
	output, err := rep.uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: &rep.bucketName,
		Key:    &image.ID,
		Body:   bytes.NewReader(image.Bytes),
	})
	if err != nil {
		rep.log.Error("Failed to upload image, reason: ", err.Error())
		return "", err
	}
	rep.log.Debug("Successfully uploaded image, location: ", output.Location)
	return output.Location, nil
}
