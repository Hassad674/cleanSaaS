package response

import (
	"github.com/hassad/boilerplateSaaS/backend/internal/domain/storage"
)

type FileResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Key         string `json:"key"`
	SizeBytes   int64  `json:"size_bytes"`
	ContentType string `json:"content_type"`
	URL         string `json:"url"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func FileFromDomain(f *storage.File) FileResponse {
	return FileResponse{
		ID:          f.ID,
		Name:        f.Name,
		Key:         f.Key,
		SizeBytes:   f.SizeBytes,
		ContentType: f.ContentType,
		URL:         f.URL,
		CreatedAt:   f.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   f.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
