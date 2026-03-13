package referral

import (
	"context"
	"fmt"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
	domainreferral "github.com/hassad/boilerplateSaaS/backend/internal/domain/referral"
	"github.com/hassad/boilerplateSaaS/backend/internal/port/repository"
)

// ReferralStats holds aggregated referral statistics for a user.
type ReferralStats struct {
	TotalReferrals     int `json:"total_referrals"`
	CompletedReferrals int `json:"completed_referrals"`
	TotalRewards       int `json:"total_rewards"`
}

type Service struct {
	referrals repository.ReferralRepository
}

func NewService(referrals repository.ReferralRepository) *Service {
	return &Service{referrals: referrals}
}

// GetOrCreateCode returns the user's existing referral code, or creates one if none exists.
func (s *Service) GetOrCreateCode(ctx context.Context, userID string) (string, error) {
	existing, err := s.referrals.FindByReferrerID(ctx, userID)
	if err == nil && existing != nil {
		return existing.Code, nil
	}

	// Generate a unique code
	code, err := domainreferral.GenerateCode()
	if err != nil {
		return "", fmt.Errorf("generating referral code: %w", err)
	}

	ref := &domainreferral.Referral{
		ReferrerID: userID,
		Code:       code,
		Status:     domainreferral.StatusPending,
		RewardType: domainreferral.RewardCredit,
	}

	if err := ref.Validate(); err != nil {
		return "", err
	}

	if err := s.referrals.Create(ctx, ref); err != nil {
		return "", fmt.Errorf("creating referral code: %w", err)
	}

	return ref.Code, nil
}

// ApplyReferral links a new user to their referrer via a referral code.
func (s *Service) ApplyReferral(ctx context.Context, referredUserID, code string) error {
	ref, err := s.referrals.FindByCode(ctx, code)
	if err != nil {
		return fmt.Errorf("finding referral code: %w", err)
	}

	if err := domainreferral.CannotReferSelf(ref.ReferrerID, referredUserID); err != nil {
		return err
	}

	// Check if this user was already referred
	existing, err := s.referrals.FindByReferredID(ctx, referredUserID)
	if err == nil && existing != nil {
		return fmt.Errorf("user already has a referral: %w", domain.ErrAlreadyExists)
	}

	// Create a new referral record linking the referred user
	newRef := &domainreferral.Referral{
		ReferrerID: ref.ReferrerID,
		ReferredID: referredUserID,
		Code:       code,
		Status:     domainreferral.StatusPending,
		RewardType: domainreferral.RewardCredit,
	}

	if err := newRef.Validate(); err != nil {
		return err
	}

	if err := s.referrals.Create(ctx, newRef); err != nil {
		return fmt.Errorf("applying referral: %w", err)
	}

	return nil
}

// CompleteReferral marks a referral as completed (e.g., after the referred user's first payment).
func (s *Service) CompleteReferral(ctx context.Context, referralID string) error {
	ref, err := s.referrals.FindByCode(ctx, referralID)
	if err != nil {
		return fmt.Errorf("finding referral: %w", err)
	}

	if err := ref.Complete(); err != nil {
		return err
	}

	if err := s.referrals.Update(ctx, ref); err != nil {
		return fmt.Errorf("completing referral: %w", err)
	}

	return nil
}

// GetStats returns aggregated referral statistics for a user.
func (s *Service) GetStats(ctx context.Context, userID string) (*ReferralStats, error) {
	total, completed, err := s.referrals.CountByReferrer(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("getting referral stats: %w", err)
	}

	return &ReferralStats{
		TotalReferrals:     total,
		CompletedReferrals: completed,
		TotalRewards:       completed, // Each completed referral = 1 reward
	}, nil
}

// ListReferrals returns a paginated list of referrals for a user.
func (s *Service) ListReferrals(ctx context.Context, userID string, offset, limit int) ([]*domainreferral.Referral, int, error) {
	return s.referrals.ListByReferrer(ctx, userID, offset, limit)
}
