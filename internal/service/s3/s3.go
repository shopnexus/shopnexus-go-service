package s3

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"shopnexus-go-service/config"
	"shopnexus-go-service/internal/repository"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type S3Service struct {
	Client        *awsS3.Client
	bucket        string
	cloudfrontURL string
}

var _ S3ServiceInterface = (*S3Service)(nil)

type S3ServiceInterface interface {
	Upload(ctx context.Context, key string, reader io.Reader, private bool) (string, error)
	Delete(ctx context.Context, key string) error
	ListObjects(ctx context.Context, prefix string) ([]string, error)
	GetPresignedURL(ctx context.Context, key string, expireIn time.Duration) (string, error)
}

// NewS3Service creates a new S3 client with the provided credentials
func NewS3Service(repository *repository.RepositoryImpl) *S3Service {
	cfg := config.GetConfig().S3

	// Create custom credentials provider
	credProvider := credentials.NewStaticCredentialsProvider(
		cfg.AccessKeyID,
		cfg.SecretAccessKey,
		"", // Session token is optional and usually not needed for regular access keys
	)

	// Load AWS configuration with custom credentials
	awsCfg, err := awsConfig.LoadDefaultConfig(
		context.Background(),
		awsConfig.WithRegion(cfg.Region),
		awsConfig.WithCredentialsProvider(credProvider),
	)
	if err != nil {
		panic(fmt.Errorf("failed to load AWS configuration: %w", err))
	}

	return &S3Service{
		Client:        s3.NewFromConfig(awsCfg),
		bucket:        cfg.Bucket,
		cloudfrontURL: cfg.CloudfrontURL,
	}
}

func (s *S3Service) Upload(ctx context.Context, key string, body io.Reader, private bool) (string, error) {
	prefix := "public/"
	if private {
		prefix = "private/"
	}

	if !strings.HasPrefix(key, prefix) {
		key = prefix + key
	}

	_, err := s.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   body,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	// if !private {
	// 	return fmt.Sprintf("https://%s/%s", s.cloudfrontURL, key), nil
	// }

	return key, nil
}

func (s *S3Service) Delete(ctx context.Context, key string) error {
	_, err := s.Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}

	return nil
}

func (s *S3Service) ListObjects(ctx context.Context, prefix string) ([]string, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(prefix),
	}

	var keys []string
	paginator := s3.NewListObjectsV2Paginator(s.Client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list objects from S3: %w", err)
		}

		for _, obj := range page.Contents {
			keys = append(keys, *obj.Key)
		}
	}

	return keys, nil
}

func (s *S3Service) GetPresignedURL(ctx context.Context, key string, expireIn time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(s.Client)

	request, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expireIn))
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return request.URL, nil
}

// GenKey creates a structured and unique S3 file key.
func GenKey(userID int64, originalFilename string) string {
	// Extract file extension
	ext := filepath.Ext(originalFilename)

	// Generate a unique identifier
	uniqueID := uuid.New().String()

	// Format: {userID}-{uuid}{ext}
	return fmt.Sprintf(
		"%d-%s-%s%s",
		userID,
		time.Now().Format("20060102150405"),
		uniqueID,
		ext,
	)
}

type S3KeyData struct {
	UserID    int64
	UniqueID  string
	Extension string
}

// ParseKey extracts the structured data from an S3 file key.
func ParseKey(key string) (S3KeyData, error) {
	parts := strings.Split(key, "-")
	if len(parts) != 3 {
		return S3KeyData{}, fmt.Errorf("invalid key format")
	}

	userID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return S3KeyData{}, fmt.Errorf("failed to parse user ID: %w", err)
	}

	uniqueID := parts[1]

	ext := filepath.Ext(parts[2])

	return S3KeyData{
		UserID:    userID,
		UniqueID:  uniqueID,
		Extension: ext,
	}, nil
}
