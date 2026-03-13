package blog

import (
	"context"
	"fmt"
	"time"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
	domainblog "github.com/hassad/boilerplateSaaS/backend/internal/domain/blog"
	"github.com/hassad/boilerplateSaaS/backend/internal/port/repository"
)

type Service struct {
	posts repository.BlogRepository
}

func NewService(posts repository.BlogRepository) *Service {
	return &Service{posts: posts}
}

func (s *Service) Create(ctx context.Context, authorID, title, slug, excerpt, content, coverImageURL, metaTitle, metaDescription string, tags []string, publish bool) (*domainblog.Post, error) {
	if title == "" {
		return nil, fmt.Errorf("%w: title is required", domain.ErrValidation)
	}
	if slug == "" {
		slug = domainblog.GenerateSlug(title)
	}

	p := &domainblog.Post{
		AuthorID:        authorID,
		Title:           title,
		Slug:            slug,
		Excerpt:         excerpt,
		Content:         content,
		CoverImageURL:   coverImageURL,
		MetaTitle:       metaTitle,
		MetaDescription: metaDescription,
		Tags:            tags,
		Status:          domainblog.StatusDraft,
	}

	if publish {
		now := time.Now()
		p.Status = domainblog.StatusPublished
		p.PublishedAt = &now
	}

	if err := s.posts.Create(ctx, p); err != nil {
		return nil, fmt.Errorf("creating blog post: %w", err)
	}
	return p, nil
}

func (s *Service) Update(ctx context.Context, id, title, slug, excerpt, content, coverImageURL, metaTitle, metaDescription string, tags []string, status string) (*domainblog.Post, error) {
	p, err := s.posts.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if title != "" {
		p.Title = title
	}
	if slug != "" {
		p.Slug = slug
	}
	p.Excerpt = excerpt
	p.Content = content
	p.CoverImageURL = coverImageURL
	p.MetaTitle = metaTitle
	p.MetaDescription = metaDescription
	p.Tags = tags

	if status != "" {
		p.Status = domainblog.Status(status)
		if p.Status == domainblog.StatusPublished && p.PublishedAt == nil {
			now := time.Now()
			p.PublishedAt = &now
		}
	}

	if err := s.posts.Update(ctx, p); err != nil {
		return nil, fmt.Errorf("updating blog post: %w", err)
	}
	return p, nil
}

func (s *Service) GetBySlug(ctx context.Context, slug string) (*domainblog.Post, error) {
	return s.posts.FindBySlug(ctx, slug)
}

func (s *Service) GetByID(ctx context.Context, id string) (*domainblog.Post, error) {
	return s.posts.FindByID(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.posts.Delete(ctx, id)
}

func (s *Service) ListPublished(ctx context.Context, tag string, offset, limit int) ([]*domainblog.Post, int, error) {
	return s.posts.List(ctx, string(domainblog.StatusPublished), tag, offset, limit)
}

func (s *Service) ListAll(ctx context.Context, status, tag string, offset, limit int) ([]*domainblog.Post, int, error) {
	return s.posts.List(ctx, status, tag, offset, limit)
}

func (s *Service) ListTags(ctx context.Context) (map[string]int, error) {
	return s.posts.ListTags(ctx)
}

func (s *Service) CountPosts(ctx context.Context) (int, error) {
	return s.posts.Count(ctx)
}
