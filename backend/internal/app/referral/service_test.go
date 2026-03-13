package referral

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
	domainreferral "github.com/hassad/boilerplateSaaS/backend/internal/domain/referral"
)

// Mocks

type mockReferralRepo struct {
	createFn         func(ctx context.Context, r *domainreferral.Referral) error
	findByCodeFn     func(ctx context.Context, code string) (*domainreferral.Referral, error)
	findByReferrerFn func(ctx context.Context, referrerID string) (*domainreferral.Referral, error)
	findByReferredFn func(ctx context.Context, referredID string) (*domainreferral.Referral, error)
	updateFn         func(ctx context.Context, r *domainreferral.Referral) error
	countByReferrerFn func(ctx context.Context, referrerID string) (int, int, error)
	listByReferrerFn func(ctx context.Context, referrerID string, offset, limit int) ([]*domainreferral.Referral, int, error)
}

func (m *mockReferralRepo) Create(ctx context.Context, r *domainreferral.Referral) error {
	if m.createFn != nil {
		return m.createFn(ctx, r)
	}
	r.ID = "ref-1"
	return nil
}

func (m *mockReferralRepo) FindByCode(ctx context.Context, code string) (*domainreferral.Referral, error) {
	if m.findByCodeFn != nil {
		return m.findByCodeFn(ctx, code)
	}
	return nil, domain.ErrNotFound
}

func (m *mockReferralRepo) FindByReferrerID(ctx context.Context, referrerID string) (*domainreferral.Referral, error) {
	if m.findByReferrerFn != nil {
		return m.findByReferrerFn(ctx, referrerID)
	}
	return nil, domain.ErrNotFound
}

func (m *mockReferralRepo) FindByReferredID(ctx context.Context, referredID string) (*domainreferral.Referral, error) {
	if m.findByReferredFn != nil {
		return m.findByReferredFn(ctx, referredID)
	}
	return nil, domain.ErrNotFound
}

func (m *mockReferralRepo) Update(ctx context.Context, r *domainreferral.Referral) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, r)
	}
	return nil
}

func (m *mockReferralRepo) CountByReferrer(ctx context.Context, referrerID string) (int, int, error) {
	if m.countByReferrerFn != nil {
		return m.countByReferrerFn(ctx, referrerID)
	}
	return 0, 0, nil
}

func (m *mockReferralRepo) ListByReferrer(ctx context.Context, referrerID string, offset, limit int) ([]*domainreferral.Referral, int, error) {
	if m.listByReferrerFn != nil {
		return m.listByReferrerFn(ctx, referrerID, offset, limit)
	}
	return nil, 0, nil
}

// Tests

func TestReferralService_GetOrCreateCode_NewCode(t *testing.T) {
	repo := &mockReferralRepo{}
	svc := NewService(repo)

	code, err := svc.GetOrCreateCode(context.Background(), "user-1")
	assert.NoError(t, err)
	assert.Len(t, code, 8)
}

func TestReferralService_GetOrCreateCode_ExistingCode(t *testing.T) {
	repo := &mockReferralRepo{
		findByReferrerFn: func(_ context.Context, _ string) (*domainreferral.Referral, error) {
			return &domainreferral.Referral{Code: "ABCD1234"}, nil
		},
	}
	svc := NewService(repo)

	code, err := svc.GetOrCreateCode(context.Background(), "user-1")
	assert.NoError(t, err)
	assert.Equal(t, "ABCD1234", code)
}

func TestReferralService_ApplyReferral_Success(t *testing.T) {
	var createdRef *domainreferral.Referral
	repo := &mockReferralRepo{
		findByCodeFn: func(_ context.Context, _ string) (*domainreferral.Referral, error) {
			return &domainreferral.Referral{
				ID:         "ref-1",
				ReferrerID: "user-1",
				Code:       "ABCD1234",
			}, nil
		},
		createFn: func(_ context.Context, r *domainreferral.Referral) error {
			createdRef = r
			r.ID = "ref-2"
			return nil
		},
	}
	svc := NewService(repo)

	err := svc.ApplyReferral(context.Background(), "user-2", "ABCD1234")
	assert.NoError(t, err)
	assert.NotNil(t, createdRef)
	assert.Equal(t, "user-1", createdRef.ReferrerID)
	assert.Equal(t, "user-2", createdRef.ReferredID)
}

func TestReferralService_ApplyReferral_CannotReferSelf(t *testing.T) {
	repo := &mockReferralRepo{
		findByCodeFn: func(_ context.Context, _ string) (*domainreferral.Referral, error) {
			return &domainreferral.Referral{
				ReferrerID: "user-1",
				Code:       "ABCD1234",
			}, nil
		},
	}
	svc := NewService(repo)

	err := svc.ApplyReferral(context.Background(), "user-1", "ABCD1234")
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrValidation)
}

func TestReferralService_ApplyReferral_AlreadyReferred(t *testing.T) {
	repo := &mockReferralRepo{
		findByCodeFn: func(_ context.Context, _ string) (*domainreferral.Referral, error) {
			return &domainreferral.Referral{
				ReferrerID: "user-1",
				Code:       "ABCD1234",
			}, nil
		},
		findByReferredFn: func(_ context.Context, _ string) (*domainreferral.Referral, error) {
			return &domainreferral.Referral{ID: "existing-ref"}, nil
		},
	}
	svc := NewService(repo)

	err := svc.ApplyReferral(context.Background(), "user-2", "ABCD1234")
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrAlreadyExists)
}

func TestReferralService_ApplyReferral_CodeNotFound(t *testing.T) {
	repo := &mockReferralRepo{}
	svc := NewService(repo)

	err := svc.ApplyReferral(context.Background(), "user-2", "INVALID1")
	assert.Error(t, err)
}

func TestReferralService_GetStats(t *testing.T) {
	repo := &mockReferralRepo{
		countByReferrerFn: func(_ context.Context, _ string) (int, int, error) {
			return 10, 5, nil
		},
	}
	svc := NewService(repo)

	stats, err := svc.GetStats(context.Background(), "user-1")
	assert.NoError(t, err)
	assert.Equal(t, 10, stats.TotalReferrals)
	assert.Equal(t, 5, stats.CompletedReferrals)
	assert.Equal(t, 5, stats.TotalRewards)
}

func TestReferralService_ListReferrals(t *testing.T) {
	refs := []*domainreferral.Referral{
		{ID: "ref-1", ReferrerID: "user-1", Code: "ABCD1234", Status: domainreferral.StatusPending},
		{ID: "ref-2", ReferrerID: "user-1", Code: "ABCD1234", Status: domainreferral.StatusCompleted},
	}
	repo := &mockReferralRepo{
		listByReferrerFn: func(_ context.Context, _ string, _ int, _ int) ([]*domainreferral.Referral, int, error) {
			return refs, 2, nil
		},
	}
	svc := NewService(repo)

	result, total, err := svc.ListReferrals(context.Background(), "user-1", 0, 20)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, 2, total)
}
