package s3

import (
	"bytes"
	"context"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	images "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/images/server"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/config"
)

type s3Repository struct {
	log        *zap.Logger
	uploader   *manager.Uploader
	bucketName string
}

func NewS3Repository(log *zap.Logger) (images.Repository, error) {
	log.Info("Connecting to s3 image bucket...")

	conf, err := awsConfig.LoadDefaultConfig(context.Background(), awsConfig.WithDefaultRegion(viper.GetString(config.S3BucketConfig.DefaultRegion)))

	if err != nil {
		log.Error("Failed to create S3 connection", zap.Error(err))
		return &s3Repository{}, err
	}
	client := s3.NewFromConfig(conf, func(o *s3.Options) {
		o.EndpointResolver = s3.EndpointResolverFromURL(viper.GetString(config.S3BucketConfig.Endpoint))
	})

	uploader := manager.NewUploader(client)

	log.Info("S3 bucket connection created successfully")

	return &s3Repository{
		log:        log,
		uploader:   uploader,
		bucketName: viper.GetString(config.S3BucketConfig.BucketName),
	}, nil
}

func (rep *s3Repository) UploadImage(image *models.Image) (string, error) {
	rep.log.Debug("Initiating s3 transaction...")
	output, err := rep.uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: &rep.bucketName,
		Key:    &image.ID,
		Body:   bytes.NewReader(image.Bytes),
	})
	if err != nil {
		rep.log.Error("Failed to upload image", zap.Error(err))
		return "", err
	}
	rep.log.Debug("Successfully uploaded image", zap.String("location", output.Location))
	return output.Location, nil
}
