package blog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateSlug(t *testing.T) {
	tests := []struct {
		title    string
		expected string
	}{
		{"Hello World", "hello-world"},
		{"My First Post!", "my-first-post"},
		{"  Spaces  Everywhere  ", "spaces-everywhere"},
		{"Already-slugified", "already-slugified"},
		{"UPPERCASE TITLE", "uppercase-title"},
		{"Special @#$ Characters", "special-characters"},
		{"Numbers 123 Here", "numbers-123-here"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			assert.Equal(t, tt.expected, GenerateSlug(tt.title))
		})
	}
}

func TestStatus_Values(t *testing.T) {
	assert.Equal(t, Status("draft"), StatusDraft)
	assert.Equal(t, Status("published"), StatusPublished)
}
