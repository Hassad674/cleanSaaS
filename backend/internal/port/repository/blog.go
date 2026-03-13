package repository

import (
	"context"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain/blog"
)

type BlogRepository interface {
	Create(ctx context.Context, p *blog.Post) error
	FindByID(ctx context.Context, id string) (*blog.Post, error)
	FindBySlug(ctx context.Context, slug string) (*blog.Post, error)
	Update(ctx context.Context, p *blog.Post) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, status string, tag string, offset, limit int) ([]*blog.Post, int, error)
	ListTags(ctx context.Context) (map[string]int, error)
}
