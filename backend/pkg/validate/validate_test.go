package validate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmail_Valid(t *testing.T) {
	tests := []struct {
		name  string
		email string
		valid bool
	}{
		{"Simple email", "test@example.com", true},
		{"With subdomain", "user@mail.example.com", true},
		{"With plus", "user+tag@example.com", true},
		{"With dots", "first.last@example.com", true},
		{"Empty string", "", false},
		{"Missing at sign", "testexample.com", false},
		{"Missing domain", "test@", false},
		{"Missing local part", "@example.com", false},
		{"Double at sign", "test@@example.com", false},
		{"Spaces", "test @example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.valid, Email(tt.email))
		})
	}
}

func TestMinLength(t *testing.T) {
	tests := []struct {
		name  string
		s     string
		min   int
		valid bool
	}{
		{"Exact length", "abc", 3, true},
		{"Over minimum", "abcdef", 3, true},
		{"Under minimum", "ab", 3, false},
		{"Empty string with min 0", "", 0, true},
		{"Empty string with min 1", "", 1, false},
		{"Zero length requirement", "anything", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.valid, MinLength(tt.s, tt.min))
		})
	}
}

func TestSlug_Valid(t *testing.T) {
	tests := []struct {
		name  string
		slug  string
		valid bool
	}{
		{"Simple slug", "hello-world", true},
		{"Single word", "hello", true},
		{"With numbers", "post-123", true},
		{"Numbers only", "123", true},
		{"Empty string", "", false},
		{"Starts with dash", "-hello", false},
		{"Ends with dash", "hello-", false},
		{"Double dash", "hello--world", false},
		{"Uppercase", "Hello-World", false},
		{"With spaces", "hello world", false},
		{"With special chars", "hello@world", false},
		{"Underscore", "hello_world", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.valid, Slug(tt.slug))
		})
	}
}
