package service

import (
	"context"
	"io"
)

type StorageService interface {
	Upload(ctx context.Context, key string, reader io.Reader, contentType string, size int64) (url string, err error)
	Delete(ctx context.Context, key string) error
	GetURL(ctx context.Context, key string) (string, error)
}
