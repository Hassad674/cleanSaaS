package r2

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/hassad/boilerplateSaaS/backend/pkg/ctxutil"
)

// defaultCallTimeout is the fallback per-call ceiling used when a StorageService
// is constructed without an explicit timeout.
const defaultCallTimeout = 15 * time.Second

type StorageService struct {
	client      *s3.Client
	bucketName  string
	publicURL   string
	callTimeout time.Duration
}

func NewStorageService(client *s3.Client, bucketName, publicURL string) *StorageService {
	return NewStorageServiceWithTimeout(client, bucketName, publicURL, defaultCallTimeout)
}

// NewStorageServiceWithTimeout builds a StorageService that bounds every R2/S3
// call to callTimeout (a ceiling; a nearer caller deadline still wins).
func NewStorageServiceWithTimeout(client *s3.Client, bucketName, publicURL string, callTimeout time.Duration) *StorageService {
	return &StorageService{
		client:      client,
		bucketName:  bucketName,
		publicURL:   publicURL,
		callTimeout: callTimeout,
	}
}

func (s *StorageService) Upload(ctx context.Context, key string, reader io.Reader, contentType string, size int64) (string, error) {
	ctx, cancel := ctxutil.WithTimeout(ctx, s.callTimeout)
	defer cancel()

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(s.bucketName),
		Key:           aws.String(key),
		Body:          reader,
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(size),
	})
	if err != nil {
		return "", fmt.Errorf("uploading to R2: %w", err)
	}

	url := fmt.Sprintf("%s/%s", s.publicURL, key)
	return url, nil
}

func (s *StorageService) Delete(ctx context.Context, key string) error {
	ctx, cancel := ctxutil.WithTimeout(ctx, s.callTimeout)
	defer cancel()

	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("deleting from R2: %w", err)
	}
	return nil
}

func (s *StorageService) GetURL(_ context.Context, key string) (string, error) {
	return fmt.Sprintf("%s/%s", s.publicURL, key), nil
}
