package repository

import (
	"context"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain/storage"
)

type FileRepository interface {
	Create(ctx context.Context, f *storage.File) error
	FindByID(ctx context.Context, id string) (*storage.File, error)
	ListByUserID(ctx context.Context, userID string, offset, limit int) ([]*storage.File, int, error)
	Delete(ctx context.Context, id string) error
}
