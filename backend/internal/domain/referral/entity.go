package referral

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusCompleted Status = "completed"
	StatusRewarded  Status = "rewarded"
)

type RewardType string

const (
	RewardCredit         RewardType = "credit"
	RewardTrialExtension RewardType = "trial_extension"
)

const codeLength = 8
const codeCharset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type Referral struct {
	ID           string
	ReferrerID   string
	ReferredID   string
	Code         string
	Status       Status
	RewardType   RewardType
	RewardAmount int
	CreatedAt    time.Time
	CompletedAt  *time.Time
	RewardedAt   *time.Time
}

// GenerateCode creates a cryptographically random 8-character alphanumeric code.
func GenerateCode() (string, error) {
	code := make([]byte, codeLength)
	for i := range code {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(codeCharset))))
		if err != nil {
			return "", fmt.Errorf("generating referral code: %w", err)
		}
		code[i] = codeCharset[idx.Int64()]
	}
	return string(code), nil
}

// Validate checks the referral entity for required fields and business rules.
func (r *Referral) Validate() error {
	if r.ReferrerID == "" {
		return fmt.Errorf("referrer ID is required: %w", domain.ErrValidation)
	}
	if r.Code == "" {
		return fmt.Errorf("referral code is required: %w", domain.ErrValidation)
	}
	if len(r.Code) != codeLength {
		return fmt.Errorf("referral code must be %d characters: %w", codeLength, domain.ErrValidation)
	}
	return nil
}

// ValidateStatus checks that a given status string is valid.
func ValidateStatus(s string) bool {
	switch Status(s) {
	case StatusPending, StatusCompleted, StatusRewarded:
		return true
	}
	return false
}

// ValidateRewardType checks that a given reward type string is valid.
func ValidateRewardType(s string) bool {
	switch RewardType(s) {
	case RewardCredit, RewardTrialExtension:
		return true
	}
	return false
}

// CannotReferSelf checks the business rule that a user cannot use their own referral code.
func CannotReferSelf(referrerID, referredID string) error {
	if referrerID == referredID {
		return fmt.Errorf("cannot use your own referral code: %w", domain.ErrValidation)
	}
	return nil
}

// Complete transitions the referral to completed status.
func (r *Referral) Complete() error {
	if r.Status != StatusPending {
		return fmt.Errorf("referral is not pending: %w", domain.ErrValidation)
	}
	now := time.Now()
	r.Status = StatusCompleted
	r.CompletedAt = &now
	return nil
}

// MarkRewarded transitions the referral to rewarded status.
func (r *Referral) MarkRewarded(rewardType RewardType, amount int) error {
	if r.Status != StatusCompleted {
		return fmt.Errorf("referral is not completed: %w", domain.ErrValidation)
	}
	now := time.Now()
	r.Status = StatusRewarded
	r.RewardType = rewardType
	r.RewardAmount = amount
	r.RewardedAt = &now
	return nil
}
