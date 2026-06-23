package storage

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
	domainstorage "github.com/hassad/boilerplateSaaS/backend/internal/domain/storage"
	"github.com/hassad/boilerplateSaaS/backend/internal/port/repository"
	"github.com/hassad/boilerplateSaaS/backend/internal/port/service"
)

type Service struct {
	storage service.StorageService
	files   repository.FileScope
}

// NewService wires the storage use cases. files is an org-scoped unit-of-work: each
// metadata operation runs inside a transaction bound to the caller's active
// organization, so PostgreSQL row-level security enforces tenant isolation on every
// file query (in addition to the repository's own org_id filter).
func NewService(storage service.StorageService, files repository.FileScope) *Service {
	return &Service{storage: storage, files: files}
}

func (s *Service) Upload(ctx context.Context, userID, fileName, contentType string, size int64, reader io.Reader) (*domainstorage.File, error) {
	if !domainstorage.IsAllowedType(contentType) {
		return nil, fmt.Errorf("%w: file type %s not allowed", domain.ErrValidation, contentType)
	}
	if size > domainstorage.MaxFileSize {
		return nil, fmt.Errorf("%w: file exceeds maximum size of 50MB", domain.ErrValidation)
	}
	if size <= 0 {
		return nil, fmt.Errorf("%w: file is empty", domain.ErrValidation)
	}

	// Sanitize filename to prevent path traversal:
	// 1. Extract only the base filename (strip any directory components)
	// 2. Reject hidden files (starting with .)
	// 3. Generate a unique key with UUID to prevent collisions and name-based attacks
	baseName := filepath.Base(fileName)
	if baseName == "." || baseName == ".." || strings.HasPrefix(baseName, ".") {
		return nil, fmt.Errorf("%w: invalid file name", domain.ErrValidation)
	}

	ext := filepath.Ext(baseName)
	key := fmt.Sprintf("%s/%s%s", userID, uuid.New().String(), ext)

	url, err := s.storage.Upload(ctx, key, reader, contentType, size)
	if err != nil {
		return nil, fmt.Errorf("uploading file: %w", err)
	}

	file := &domainstorage.File{
		UserID:      userID,
		Name:        fileName,
		Key:         key,
		SizeBytes:   size,
		ContentType: contentType,
		URL:         url,
	}

	err = s.files.WithOrgFiles(ctx, func(files repository.FileRepository) error {
		return files.Create(ctx, file)
	})
	if err != nil {
		// Best-effort cleanup of uploaded object
		_ = s.storage.Delete(ctx, key)
		return nil, fmt.Errorf("saving file metadata: %w", err)
	}

	return file, nil
}

func (s *Service) Delete(ctx context.Context, userID, fileID string) error {
	return s.files.WithOrgFiles(ctx, func(files repository.FileRepository) error {
		file, err := files.FindByID(ctx, fileID)
		if err != nil {
			return err
		}
		if file.UserID != userID {
			return domain.ErrForbidden
		}
		if err := s.storage.Delete(ctx, file.Key); err != nil {
			return fmt.Errorf("deleting file from storage: %w", err)
		}
		return files.Delete(ctx, fileID)
	})
}

func (s *Service) List(ctx context.Context, userID string, offset, limit int) ([]*domainstorage.File, int, error) {
	var files []*domainstorage.File
	var total int
	err := s.files.WithOrgFiles(ctx, func(repo repository.FileRepository) error {
		var e error
		files, total, e = repo.ListByUserID(ctx, userID, offset, limit)
		return e
	})
	return files, total, err
}

func (s *Service) GetByID(ctx context.Context, userID, fileID string) (*domainstorage.File, error) {
	var file *domainstorage.File
	err := s.files.WithOrgFiles(ctx, func(files repository.FileRepository) error {
		f, err := files.FindByID(ctx, fileID)
		if err != nil {
			return err
		}
		if f.UserID != userID {
			return domain.ErrForbidden
		}
		file = f
		return nil
	})
	if err != nil {
		return nil, err
	}
	return file, nil
}
