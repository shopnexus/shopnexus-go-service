package s3

import (
	"context"
	"io"
	"time"

	"shopnexus-go-service/config"
	s3client "shopnexus-go-service/internal/client/s3"

	s3base "github.com/aws/aws-sdk-go-v2/service/s3"
)

type Service interface {
	BaseClient() *s3base.Client
	Upload(ctx context.Context, key string, reader io.Reader, private bool) (string, error)
	Delete(ctx context.Context, key string) error
	ListObjects(ctx context.Context, prefix string) ([]string, error)
	GetPresignedURL(ctx context.Context, key string, expireIn time.Duration) (string, error)
}

type ServiceImpl struct {
	client s3client.Client
}

func NewService() (Service, error) {
	cfg := config.GetConfig().S3
	client, err := s3client.NewClient(s3client.S3Config{
		AccessKeyID:     cfg.AccessKeyID,
		SecretAccessKey: cfg.SecretAccessKey,
		Region:          cfg.Region,
		Bucket:          cfg.Bucket,
		CloudfrontURL:   cfg.CloudfrontURL,
	})
	if err != nil {
		return nil, err
	}

	return &ServiceImpl{
		client: client,
	}, nil
}

func (s *ServiceImpl) BaseClient() *s3base.Client {
	return s.client.Client()
}

func (s *ServiceImpl) Upload(ctx context.Context, key string, reader io.Reader, private bool) (string, error) {
	return s.client.Upload(ctx, key, reader, private)
}

func (s *ServiceImpl) Delete(ctx context.Context, key string) error {
	return s.client.Delete(ctx, key)
}

func (s *ServiceImpl) ListObjects(ctx context.Context, prefix string) ([]string, error) {
	return s.client.ListObjects(ctx, prefix)
}

func (s *ServiceImpl) GetPresignedURL(ctx context.Context, key string, expireIn time.Duration) (string, error) {
	return s.client.GetPresignedURL(ctx, key, expireIn)
}
