package blog

import (
	"strings"
	"time"
)

type Status string

const (
	StatusDraft     Status = "draft"
	StatusPublished Status = "published"
)

type Post struct {
	ID              string
	AuthorID        string
	Title           string
	Slug            string
	Excerpt         string
	Content         string
	CoverImageURL   string
	MetaTitle       string
	MetaDescription string
	Tags            []string
	Status          Status
	PublishedAt     *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func GenerateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			return r
		}
		if r == ' ' || r == '-' {
			return '-'
		}
		return -1
	}, slug)
	// Remove consecutive dashes
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}
	return strings.Trim(slug, "-")
}
