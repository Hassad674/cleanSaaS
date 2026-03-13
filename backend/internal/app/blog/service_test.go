package blog

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
	domainblog "github.com/hassad/boilerplateSaaS/backend/internal/domain/blog"
)

// Mocks

type mockBlogRepo struct {
	createFn     func(ctx context.Context, p *domainblog.Post) error
	findByIDFn   func(ctx context.Context, id string) (*domainblog.Post, error)
	findBySlugFn func(ctx context.Context, slug string) (*domainblog.Post, error)
	updateFn     func(ctx context.Context, p *domainblog.Post) error
	deleteFn     func(ctx context.Context, id string) error
	listFn       func(ctx context.Context, status, tag string, offset, limit int) ([]*domainblog.Post, int, error)
	listTagsFn   func(ctx context.Context) (map[string]int, error)
	countFn      func(ctx context.Context) (int, error)
}

func (m *mockBlogRepo) Create(ctx context.Context, p *domainblog.Post) error {
	if m.createFn != nil {
		return m.createFn(ctx, p)
	}
	p.ID = "post-1"
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	return nil
}

func (m *mockBlogRepo) FindByID(ctx context.Context, id string) (*domainblog.Post, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, id)
	}
	return nil, domain.ErrNotFound
}

func (m *mockBlogRepo) FindBySlug(ctx context.Context, slug string) (*domainblog.Post, error) {
	if m.findBySlugFn != nil {
		return m.findBySlugFn(ctx, slug)
	}
	return nil, domain.ErrNotFound
}

func (m *mockBlogRepo) Update(ctx context.Context, p *domainblog.Post) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, p)
	}
	return nil
}

func (m *mockBlogRepo) Delete(ctx context.Context, id string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

func (m *mockBlogRepo) List(ctx context.Context, status, tag string, offset, limit int) ([]*domainblog.Post, int, error) {
	if m.listFn != nil {
		return m.listFn(ctx, status, tag, offset, limit)
	}
	return nil, 0, nil
}

func (m *mockBlogRepo) ListTags(ctx context.Context) (map[string]int, error) {
	if m.listTagsFn != nil {
		return m.listTagsFn(ctx)
	}
	return map[string]int{}, nil
}

func (m *mockBlogRepo) Count(ctx context.Context) (int, error) {
	if m.countFn != nil {
		return m.countFn(ctx)
	}
	return 0, nil
}

// Tests

func TestBlogService_Create_Success(t *testing.T) {
	repo := &mockBlogRepo{}
	svc := NewService(repo)

	post, err := svc.Create(context.Background(), "user-1", "My Post", "", "excerpt", "content",
		"", "", "", []string{"go"}, false)

	assert.NoError(t, err)
	assert.Equal(t, "My Post", post.Title)
	assert.Equal(t, "my-post", post.Slug)
	assert.Equal(t, domainblog.StatusDraft, post.Status)
	assert.Nil(t, post.PublishedAt)
}

func TestBlogService_Create_WithPublish(t *testing.T) {
	repo := &mockBlogRepo{}
	svc := NewService(repo)

	post, err := svc.Create(context.Background(), "user-1", "Published Post", "pub-post", "", "content",
		"", "", "", nil, true)

	assert.NoError(t, err)
	assert.Equal(t, domainblog.StatusPublished, post.Status)
	assert.NotNil(t, post.PublishedAt)
}

func TestBlogService_Create_EmptyTitle(t *testing.T) {
	repo := &mockBlogRepo{}
	svc := NewService(repo)

	_, err := svc.Create(context.Background(), "user-1", "", "", "", "", "", "", "", nil, false)

	assert.ErrorIs(t, err, domain.ErrValidation)
}

func TestBlogService_Update_Success(t *testing.T) {
	repo := &mockBlogRepo{
		findByIDFn: func(_ context.Context, _ string) (*domainblog.Post, error) {
			return &domainblog.Post{
				ID:       "post-1",
				AuthorID: "user-1",
				Title:    "Old Title",
				Slug:     "old-title",
				Status:   domainblog.StatusDraft,
			}, nil
		},
	}
	svc := NewService(repo)

	post, err := svc.Update(context.Background(), "post-1", "New Title", "new-slug", "", "new content",
		"", "", "", []string{"updated"}, "published")

	assert.NoError(t, err)
	assert.Equal(t, "New Title", post.Title)
	assert.Equal(t, "new-slug", post.Slug)
	assert.Equal(t, domainblog.StatusPublished, post.Status)
	assert.NotNil(t, post.PublishedAt)
}

func TestBlogService_Update_NotFound(t *testing.T) {
	repo := &mockBlogRepo{}
	svc := NewService(repo)

	_, err := svc.Update(context.Background(), "nonexistent", "Title", "", "", "", "", "", "", nil, "")

	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestBlogService_ListPublished_FiltersStatus(t *testing.T) {
	var capturedStatus string
	repo := &mockBlogRepo{
		listFn: func(_ context.Context, status, _ string, _, _ int) ([]*domainblog.Post, int, error) {
			capturedStatus = status
			return []*domainblog.Post{}, 0, nil
		},
	}
	svc := NewService(repo)

	_, _, err := svc.ListPublished(context.Background(), "", 0, 10)

	assert.NoError(t, err)
	assert.Equal(t, "published", capturedStatus)
}

func TestBlogService_Delete_Success(t *testing.T) {
	deleted := false
	repo := &mockBlogRepo{
		deleteFn: func(_ context.Context, id string) error {
			assert.Equal(t, "post-1", id)
			deleted = true
			return nil
		},
	}
	svc := NewService(repo)

	err := svc.Delete(context.Background(), "post-1")

	assert.NoError(t, err)
	assert.True(t, deleted)
}
