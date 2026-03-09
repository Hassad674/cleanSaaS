package storage

import "time"

var AllowedImageTypes = []string{"image/jpeg", "image/png", "image/webp"}
var AllowedDocTypes = []string{"application/pdf"}
var MaxFileSize int64 = 10 * 1024 * 1024 // 10MB

type File struct {
	ID          string
	UserID      string
	Name        string
	ContentType string
	Size        int64
	URL         string
	CreatedAt   time.Time
}

func (f *File) IsImage() bool {
	for _, t := range AllowedImageTypes {
		if f.ContentType == t {
			return true
		}
	}
	return false
}
