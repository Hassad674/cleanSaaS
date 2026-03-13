package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsAllowedType(t *testing.T) {
	assert.True(t, IsAllowedType("image/jpeg"))
	assert.True(t, IsAllowedType("image/png"))
	assert.True(t, IsAllowedType("application/pdf"))
	assert.True(t, IsAllowedType("video/mp4"))
	assert.False(t, IsAllowedType("application/octet-stream"))
	assert.False(t, IsAllowedType("text/html"))
	assert.False(t, IsAllowedType(""))
}

func TestFile_IsImage(t *testing.T) {
	f := &File{ContentType: "image/jpeg"}
	assert.True(t, f.IsImage())

	f.ContentType = "image/png"
	assert.True(t, f.IsImage())

	f.ContentType = "application/pdf"
	assert.False(t, f.IsImage())

	f.ContentType = "video/mp4"
	assert.False(t, f.IsImage())
}

func TestMaxFileSize(t *testing.T) {
	assert.Equal(t, int64(50*1024*1024), MaxFileSize)
}
