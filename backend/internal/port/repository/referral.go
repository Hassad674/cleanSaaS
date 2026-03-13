package repository

import (
	"context"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain/referral"
)

type ReferralRepository interface {
	Create(ctx context.Context, r *referral.Referral) error
	FindByCode(ctx context.Context, code string) (*referral.Referral, error)
	FindByReferrerID(ctx context.Context, referrerID string) (*referral.Referral, error)
	FindByReferredID(ctx context.Context, referredID string) (*referral.Referral, error)
	Update(ctx context.Context, r *referral.Referral) error
	CountByReferrer(ctx context.Context, referrerID string) (total int, completed int, err error)
	ListByReferrer(ctx context.Context, referrerID string, offset, limit int) ([]*referral.Referral, int, error)
}
