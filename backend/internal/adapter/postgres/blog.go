package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/lib/pq"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
	"github.com/hassad/boilerplateSaaS/backend/internal/domain/blog"
)

type BlogRepository struct {
	db *sql.DB
}

func NewBlogRepository(db *sql.DB) *BlogRepository {
	return &BlogRepository{db: db}
}

func (r *BlogRepository) Create(ctx context.Context, p *blog.Post) error {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO blog_posts (author_id, title, slug, excerpt, content, cover_image_url, meta_title, meta_description, tags, status, published_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		 RETURNING id, created_at, updated_at`,
		p.AuthorID, p.Title, p.Slug, p.Excerpt, p.Content, p.CoverImageURL,
		p.MetaTitle, p.MetaDescription, pq.Array(p.Tags), string(p.Status), p.PublishedAt,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return fmt.Errorf("creating blog post: %w", err)
	}
	return nil
}

func (r *BlogRepository) FindByID(ctx context.Context, id string) (*blog.Post, error) {
	return r.scanPost(r.db.QueryRowContext(ctx,
		`SELECT id, author_id, title, slug, excerpt, content, cover_image_url, meta_title, meta_description, tags, status, published_at, created_at, updated_at
		 FROM blog_posts WHERE id = $1`, id,
	))
}

func (r *BlogRepository) FindBySlug(ctx context.Context, slug string) (*blog.Post, error) {
	return r.scanPost(r.db.QueryRowContext(ctx,
		`SELECT id, author_id, title, slug, excerpt, content, cover_image_url, meta_title, meta_description, tags, status, published_at, created_at, updated_at
		 FROM blog_posts WHERE slug = $1`, slug,
	))
}

func (r *BlogRepository) Update(ctx context.Context, p *blog.Post) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE blog_posts SET title = $1, slug = $2, excerpt = $3, content = $4, cover_image_url = $5,
		 meta_title = $6, meta_description = $7, tags = $8, status = $9, published_at = $10, updated_at = NOW()
		 WHERE id = $11`,
		p.Title, p.Slug, p.Excerpt, p.Content, p.CoverImageURL,
		p.MetaTitle, p.MetaDescription, pq.Array(p.Tags), string(p.Status), p.PublishedAt, p.ID,
	)
	if err != nil {
		return fmt.Errorf("updating blog post: %w", err)
	}
	return nil
}

func (r *BlogRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM blog_posts WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("deleting blog post: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *BlogRepository) List(ctx context.Context, status string, tag string, offset, limit int) ([]*blog.Post, int, error) {
	var conditions []string
	var args []interface{}
	argIdx := 1

	if status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, status)
		argIdx++
	}
	if tag != "" {
		conditions = append(conditions, fmt.Sprintf("$%d = ANY(tags)", argIdx))
		args = append(args, tag)
		argIdx++
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM blog_posts %s", where)
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("counting blog posts: %w", err)
	}

	listQuery := fmt.Sprintf(
		`SELECT id, author_id, title, slug, excerpt, content, cover_image_url, meta_title, meta_description, tags, status, published_at, created_at, updated_at
		 FROM blog_posts %s ORDER BY published_at DESC NULLS LAST, created_at DESC LIMIT $%d OFFSET $%d`,
		where, argIdx, argIdx+1,
	)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("listing blog posts: %w", err)
	}
	defer rows.Close()

	var posts []*blog.Post
	for rows.Next() {
		p, err := r.scanPostRow(rows)
		if err != nil {
			return nil, 0, err
		}
		posts = append(posts, p)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterating blog posts: %w", err)
	}

	return posts, total, nil
}

func (r *BlogRepository) ListTags(ctx context.Context) (map[string]int, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT unnest(tags) as tag, COUNT(*) as cnt
		 FROM blog_posts WHERE status = 'published'
		 GROUP BY tag ORDER BY cnt DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("listing tags: %w", err)
	}
	defer rows.Close()

	tags := make(map[string]int)
	for rows.Next() {
		var tag string
		var count int
		if err := rows.Scan(&tag, &count); err != nil {
			return nil, fmt.Errorf("scanning tag: %w", err)
		}
		tags[tag] = count
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating tags: %w", err)
	}
	return tags, nil
}

func (r *BlogRepository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM blog_posts`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("counting blog posts: %w", err)
	}
	return count, nil
}

func (r *BlogRepository) scanPost(row *sql.Row) (*blog.Post, error) {
	p := &blog.Post{}
	var tags pq.StringArray
	err := row.Scan(
		&p.ID, &p.AuthorID, &p.Title, &p.Slug, &p.Excerpt, &p.Content,
		&p.CoverImageURL, &p.MetaTitle, &p.MetaDescription, &tags,
		&p.Status, &p.PublishedAt, &p.CreatedAt, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scanning blog post: %w", err)
	}
	p.Tags = tags
	return p, nil
}

func (r *BlogRepository) scanPostRow(rows *sql.Rows) (*blog.Post, error) {
	p := &blog.Post{}
	var tags pq.StringArray
	err := rows.Scan(
		&p.ID, &p.AuthorID, &p.Title, &p.Slug, &p.Excerpt, &p.Content,
		&p.CoverImageURL, &p.MetaTitle, &p.MetaDescription, &tags,
		&p.Status, &p.PublishedAt, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scanning blog post row: %w", err)
	}
	p.Tags = tags
	return p, nil
}
