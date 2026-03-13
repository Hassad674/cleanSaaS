package storage

import "time"

var AllowedTypes = map[string]bool{
	// Images
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
	// Videos
	"video/mp4":  true,
	"video/webm": true,
	// Documents
	"application/pdf":    true,
	"text/plain":         true,
	"application/msword": true,
}

var MaxFileSize int64 = 50 * 1024 * 1024 // 50MB

type File struct {
	ID          string
	UserID      string
	Name        string
	Key         string
	SizeBytes   int64
	ContentType string
	URL         string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (f *File) IsImage() bool {
	switch f.ContentType {
	case "image/jpeg", "image/png", "image/gif", "image/webp":
		return true
	}
	return false
}

func IsAllowedType(contentType string) bool {
	return AllowedTypes[contentType]
}
