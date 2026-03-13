package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
	"github.com/hassad/boilerplateSaaS/backend/internal/domain/referral"
)

type ReferralRepository struct {
	db *sql.DB
}

func NewReferralRepository(db *sql.DB) *ReferralRepository {
	return &ReferralRepository{db: db}
}

func (r *ReferralRepository) Create(ctx context.Context, ref *referral.Referral) error {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO referrals (referrer_id, referred_id, code, status, reward_type, reward_amount)
		 VALUES ($1, NULLIF($2, ''), $3, $4, $5, $6)
		 RETURNING id, created_at`,
		ref.ReferrerID, ref.ReferredID, ref.Code, ref.Status, ref.RewardType, ref.RewardAmount,
	).Scan(&ref.ID, &ref.CreatedAt)
	if err != nil {
		return fmt.Errorf("creating referral: %w", err)
	}
	return nil
}

func (r *ReferralRepository) FindByCode(ctx context.Context, code string) (*referral.Referral, error) {
	ref := &referral.Referral{}
	var referredID sql.NullString
	var completedAt sql.NullTime
	var rewardedAt sql.NullTime

	err := r.db.QueryRowContext(ctx,
		`SELECT id, referrer_id, referred_id, code, status, reward_type, reward_amount, created_at, completed_at, rewarded_at
		 FROM referrals WHERE code = $1 AND referred_id IS NULL
		 ORDER BY created_at ASC LIMIT 1`, code,
	).Scan(&ref.ID, &ref.ReferrerID, &referredID, &ref.Code, &ref.Status, &ref.RewardType, &ref.RewardAmount, &ref.CreatedAt, &completedAt, &rewardedAt)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("finding referral by code: %w", err)
	}

	if referredID.Valid {
		ref.ReferredID = referredID.String
	}
	if completedAt.Valid {
		ref.CompletedAt = &completedAt.Time
	}
	if rewardedAt.Valid {
		ref.RewardedAt = &rewardedAt.Time
	}

	return ref, nil
}

func (r *ReferralRepository) FindByReferrerID(ctx context.Context, referrerID string) (*referral.Referral, error) {
	ref := &referral.Referral{}
	var referredID sql.NullString
	var completedAt sql.NullTime
	var rewardedAt sql.NullTime

	err := r.db.QueryRowContext(ctx,
		`SELECT id, referrer_id, referred_id, code, status, reward_type, reward_amount, created_at, completed_at, rewarded_at
		 FROM referrals WHERE referrer_id = $1 AND referred_id IS NULL
		 ORDER BY created_at ASC LIMIT 1`, referrerID,
	).Scan(&ref.ID, &ref.ReferrerID, &referredID, &ref.Code, &ref.Status, &ref.RewardType, &ref.RewardAmount, &ref.CreatedAt, &completedAt, &rewardedAt)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("finding referral by referrer: %w", err)
	}

	if referredID.Valid {
		ref.ReferredID = referredID.String
	}
	if completedAt.Valid {
		ref.CompletedAt = &completedAt.Time
	}
	if rewardedAt.Valid {
		ref.RewardedAt = &rewardedAt.Time
	}

	return ref, nil
}

func (r *ReferralRepository) FindByReferredID(ctx context.Context, referredID string) (*referral.Referral, error) {
	ref := &referral.Referral{}
	var refID sql.NullString
	var completedAt sql.NullTime
	var rewardedAt sql.NullTime

	err := r.db.QueryRowContext(ctx,
		`SELECT id, referrer_id, referred_id, code, status, reward_type, reward_amount, created_at, completed_at, rewarded_at
		 FROM referrals WHERE referred_id = $1
		 ORDER BY created_at ASC LIMIT 1`, referredID,
	).Scan(&ref.ID, &ref.ReferrerID, &refID, &ref.Code, &ref.Status, &ref.RewardType, &ref.RewardAmount, &ref.CreatedAt, &completedAt, &rewardedAt)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("finding referral by referred: %w", err)
	}

	if refID.Valid {
		ref.ReferredID = refID.String
	}
	if completedAt.Valid {
		ref.CompletedAt = &completedAt.Time
	}
	if rewardedAt.Valid {
		ref.RewardedAt = &rewardedAt.Time
	}

	return ref, nil
}

func (r *ReferralRepository) Update(ctx context.Context, ref *referral.Referral) error {
	result, err := r.db.ExecContext(ctx,
		`UPDATE referrals SET status = $1, reward_type = $2, reward_amount = $3, completed_at = $4, rewarded_at = $5
		 WHERE id = $6`,
		ref.Status, ref.RewardType, ref.RewardAmount, ref.CompletedAt, ref.RewardedAt, ref.ID,
	)
	if err != nil {
		return fmt.Errorf("updating referral: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *ReferralRepository) CountByReferrer(ctx context.Context, referrerID string) (int, int, error) {
	var total, completed int
	err := r.db.QueryRowContext(ctx,
		`SELECT
			COUNT(*) FILTER (WHERE referred_id IS NOT NULL),
			COUNT(*) FILTER (WHERE referred_id IS NOT NULL AND status IN ('completed', 'rewarded'))
		 FROM referrals WHERE referrer_id = $1`, referrerID,
	).Scan(&total, &completed)
	if err != nil {
		return 0, 0, fmt.Errorf("counting referrals: %w", err)
	}
	return total, completed, nil
}

func (r *ReferralRepository) ListByReferrer(ctx context.Context, referrerID string, offset, limit int) ([]*referral.Referral, int, error) {
	var total int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM referrals WHERE referrer_id = $1 AND referred_id IS NOT NULL`, referrerID,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("counting referrals: %w", err)
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, referrer_id, referred_id, code, status, reward_type, reward_amount, created_at, completed_at, rewarded_at
		 FROM referrals WHERE referrer_id = $1 AND referred_id IS NOT NULL
		 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		referrerID, limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("listing referrals: %w", err)
	}
	defer rows.Close()

	var referrals []*referral.Referral
	for rows.Next() {
		ref := &referral.Referral{}
		var referredID sql.NullString
		var completedAt sql.NullTime
		var rewardedAt sql.NullTime

		if err := rows.Scan(
			&ref.ID, &ref.ReferrerID, &referredID, &ref.Code, &ref.Status, &ref.RewardType, &ref.RewardAmount,
			&ref.CreatedAt, &completedAt, &rewardedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scanning referral: %w", err)
		}

		if referredID.Valid {
			ref.ReferredID = referredID.String
		}
		if completedAt.Valid {
			ref.CompletedAt = &completedAt.Time
		}
		if rewardedAt.Valid {
			ref.RewardedAt = &rewardedAt.Time
		}

		referrals = append(referrals, ref)
	}

	return referrals, total, nil
}
