package response

import (
	"github.com/hassad/boilerplateSaaS/backend/internal/domain/blog"
)

type BlogPostResponse struct {
	ID              string   `json:"id"`
	AuthorID        string   `json:"author_id"`
	Title           string   `json:"title"`
	Slug            string   `json:"slug"`
	Excerpt         string   `json:"excerpt"`
	Content         string   `json:"content"`
	CoverImageURL   string   `json:"cover_image_url"`
	MetaTitle       string   `json:"meta_title"`
	MetaDescription string   `json:"meta_description"`
	Tags            []string `json:"tags"`
	Status          string   `json:"status"`
	PublishedAt     string   `json:"published_at"`
	CreatedAt       string   `json:"created_at"`
	UpdatedAt       string   `json:"updated_at"`
}

func BlogPostFromDomain(p *blog.Post) BlogPostResponse {
	publishedAt := ""
	if p.PublishedAt != nil {
		publishedAt = p.PublishedAt.Format("2006-01-02T15:04:05Z")
	}
	tags := p.Tags
	if tags == nil {
		tags = []string{}
	}
	return BlogPostResponse{
		ID:              p.ID,
		AuthorID:        p.AuthorID,
		Title:           p.Title,
		Slug:            p.Slug,
		Excerpt:         p.Excerpt,
		Content:         p.Content,
		CoverImageURL:   p.CoverImageURL,
		MetaTitle:       p.MetaTitle,
		MetaDescription: p.MetaDescription,
		Tags:            tags,
		Status:          string(p.Status),
		PublishedAt:     publishedAt,
		CreatedAt:       p.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:       p.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
