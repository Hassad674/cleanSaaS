package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
	"github.com/hassad/boilerplateSaaS/backend/internal/domain/storage"
	"github.com/hassad/boilerplateSaaS/backend/pkg/orgctx"
)

// FileRepository implements repository.FileRepository. It holds a DBTX so the same
// code runs on the pool (system path) or an org-scoped transaction (request path).
//
// Tenant isolation is defense-in-depth: every query also filters by the active
// org_id read from the context (layer 2), on top of the PostgreSQL row-level
// security the org-scoped transaction enforces (layer 3). On insert the active
// org_id is stamped onto the row.
type FileRepository struct {
	db DBTX
}

func NewFileRepository(db *sql.DB) *FileRepository {
	return &FileRepository{db: db}
}

// newFileRepositoryTx binds the repository to an open transaction (the org scope).
func newFileRepositoryTx(tx DBTX) *FileRepository {
	return &FileRepository{db: tx}
}

func (r *FileRepository) Create(ctx context.Context, f *storage.File) error {
	orgID, ok := orgctx.OrgID(ctx)
	if !ok {
		return fmt.Errorf("creating file record: %w", domain.ErrForbidden)
	}
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO files (user_id, org_id, name, key, size_bytes, content_type, url)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, created_at, updated_at`,
		f.UserID, orgID, f.Name, f.Key, f.SizeBytes, f.ContentType, f.URL,
	).Scan(&f.ID, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return fmt.Errorf("creating file record: %w", err)
	}
	return nil
}

func (r *FileRepository) FindByID(ctx context.Context, id string) (*storage.File, error) {
	f := &storage.File{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, name, key, size_bytes, content_type, url, created_at, updated_at
		 FROM files WHERE id = $1 AND org_id = $2`, id, orgFilter(ctx),
	).Scan(&f.ID, &f.UserID, &f.Name, &f.Key, &f.SizeBytes, &f.ContentType, &f.URL, &f.CreatedAt, &f.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("finding file by id: %w", err)
	}
	return f, nil
}

func (r *FileRepository) ListByUserID(ctx context.Context, userID string, offset, limit int) ([]*storage.File, int, error) {
	org := orgFilter(ctx)
	var total int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM files WHERE user_id = $1 AND org_id = $2`, userID, org,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("counting files: %w", err)
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, name, key, size_bytes, content_type, url, created_at, updated_at
		 FROM files WHERE user_id = $1 AND org_id = $2
		 ORDER BY created_at DESC
		 LIMIT $3 OFFSET $4`, userID, org, limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("listing files: %w", err)
	}
	defer rows.Close()

	var files []*storage.File
	for rows.Next() {
		f := &storage.File{}
		if err := rows.Scan(&f.ID, &f.UserID, &f.Name, &f.Key, &f.SizeBytes, &f.ContentType, &f.URL, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("scanning file row: %w", err)
		}
		files = append(files, f)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterating files: %w", err)
	}

	return files, total, nil
}

func (r *FileRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM files WHERE id = $1 AND org_id = $2`, id, orgFilter(ctx))
	if err != nil {
		return fmt.Errorf("deleting file: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}
