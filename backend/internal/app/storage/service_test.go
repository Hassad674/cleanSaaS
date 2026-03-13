package storage

import (
	"bytes"
	"context"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
	domainstorage "github.com/hassad/boilerplateSaaS/backend/internal/domain/storage"
)

// Mocks

type mockStorageSvc struct {
	uploadFn func(ctx context.Context, key string, contentType string, size int64) (string, error)
	deleteFn func(ctx context.Context, key string) error
}

func (m *mockStorageSvc) Upload(ctx context.Context, key string, _ io.Reader, contentType string, size int64) (string, error) {
	if m.uploadFn != nil {
		return m.uploadFn(ctx, key, contentType, size)
	}
	return "https://r2.example.com/" + key, nil
}

func (m *mockStorageSvc) Delete(ctx context.Context, key string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, key)
	}
	return nil
}

func (m *mockStorageSvc) GetURL(_ context.Context, key string) (string, error) {
	return "https://r2.example.com/" + key, nil
}

type mockFileRepo struct {
	createFn       func(ctx context.Context, f *domainstorage.File) error
	findByIDFn     func(ctx context.Context, id string) (*domainstorage.File, error)
	listByUserIDFn func(ctx context.Context, userID string, offset, limit int) ([]*domainstorage.File, int, error)
	deleteFn       func(ctx context.Context, id string) error
}

func (m *mockFileRepo) Create(ctx context.Context, f *domainstorage.File) error {
	if m.createFn != nil {
		return m.createFn(ctx, f)
	}
	f.ID = "file-1"
	f.CreatedAt = time.Now()
	f.UpdatedAt = time.Now()
	return nil
}

func (m *mockFileRepo) FindByID(ctx context.Context, id string) (*domainstorage.File, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, id)
	}
	return nil, domain.ErrNotFound
}

func (m *mockFileRepo) ListByUserID(ctx context.Context, userID string, offset, limit int) ([]*domainstorage.File, int, error) {
	if m.listByUserIDFn != nil {
		return m.listByUserIDFn(ctx, userID, offset, limit)
	}
	return nil, 0, nil
}

func (m *mockFileRepo) Delete(ctx context.Context, id string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

// Tests

func TestStorageService_Upload_Success(t *testing.T) {
	var uploadedKey string
	storageMock := &mockStorageSvc{
		uploadFn: func(_ context.Context, key string, _ string, _ int64) (string, error) {
			uploadedKey = key
			return "https://r2.example.com/" + key, nil
		},
	}
	fileRepo := &mockFileRepo{}

	svc := NewService(storageMock, fileRepo)
	reader := bytes.NewReader([]byte("file content"))
	file, err := svc.Upload(context.Background(), "user-1", "photo.jpg", "image/jpeg", 1024, reader)

	assert.NoError(t, err)
	assert.Equal(t, "user-1/photo.jpg", uploadedKey)
	assert.Equal(t, "user-1", file.UserID)
	assert.Equal(t, "photo.jpg", file.Name)
	assert.Equal(t, int64(1024), file.SizeBytes)
	assert.Equal(t, "image/jpeg", file.ContentType)
}

func TestStorageService_Upload_ForbiddenType(t *testing.T) {
	svc := NewService(&mockStorageSvc{}, &mockFileRepo{})
	reader := bytes.NewReader([]byte("data"))
	_, err := svc.Upload(context.Background(), "user-1", "script.exe", "application/x-executable", 1024, reader)

	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrValidation)
}

func TestStorageService_Upload_TooLarge(t *testing.T) {
	svc := NewService(&mockStorageSvc{}, &mockFileRepo{})
	reader := bytes.NewReader([]byte("data"))
	_, err := svc.Upload(context.Background(), "user-1", "huge.pdf", "application/pdf", 60*1024*1024, reader)

	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrValidation)
}

func TestStorageService_Delete_OwnFile(t *testing.T) {
	var deletedKey string
	var deletedID string
	storageMock := &mockStorageSvc{
		deleteFn: func(_ context.Context, key string) error {
			deletedKey = key
			return nil
		},
	}
	fileRepo := &mockFileRepo{
		findByIDFn: func(_ context.Context, id string) (*domainstorage.File, error) {
			return &domainstorage.File{
				ID:     "file-1",
				UserID: "user-1",
				Key:    "user-1/photo.jpg",
			}, nil
		},
		deleteFn: func(_ context.Context, id string) error {
			deletedID = id
			return nil
		},
	}

	svc := NewService(storageMock, fileRepo)
	err := svc.Delete(context.Background(), "user-1", "file-1")

	assert.NoError(t, err)
	assert.Equal(t, "user-1/photo.jpg", deletedKey)
	assert.Equal(t, "file-1", deletedID)
}

func TestStorageService_Delete_OtherUserFile(t *testing.T) {
	fileRepo := &mockFileRepo{
		findByIDFn: func(_ context.Context, _ string) (*domainstorage.File, error) {
			return &domainstorage.File{
				ID:     "file-1",
				UserID: "user-2",
				Key:    "user-2/secret.pdf",
			}, nil
		},
	}

	svc := NewService(&mockStorageSvc{}, fileRepo)
	err := svc.Delete(context.Background(), "user-1", "file-1")

	assert.ErrorIs(t, err, domain.ErrForbidden)
}
