package referral

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
)

func TestGenerateCode_Length(t *testing.T) {
	code, err := GenerateCode()
	assert.NoError(t, err)
	assert.Len(t, code, codeLength)
}

func TestGenerateCode_Uniqueness(t *testing.T) {
	codes := make(map[string]bool)
	for i := 0; i < 100; i++ {
		code, err := GenerateCode()
		assert.NoError(t, err)
		codes[code] = true
	}
	// With 36^8 possible codes, 100 should all be unique
	assert.Len(t, codes, 100)
}

func TestGenerateCode_CharacterSet(t *testing.T) {
	code, err := GenerateCode()
	assert.NoError(t, err)
	for _, c := range code {
		assert.Contains(t, codeCharset, string(c))
	}
}

func TestReferral_Validate_Success(t *testing.T) {
	code, _ := GenerateCode()
	r := &Referral{
		ReferrerID: "user-1",
		Code:       code,
	}
	assert.NoError(t, r.Validate())
}

func TestReferral_Validate_MissingReferrerID(t *testing.T) {
	code, _ := GenerateCode()
	r := &Referral{
		Code: code,
	}
	err := r.Validate()
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrValidation)
}

func TestReferral_Validate_MissingCode(t *testing.T) {
	r := &Referral{
		ReferrerID: "user-1",
	}
	err := r.Validate()
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrValidation)
}

func TestReferral_Validate_InvalidCodeLength(t *testing.T) {
	r := &Referral{
		ReferrerID: "user-1",
		Code:       "ABC",
	}
	err := r.Validate()
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrValidation)
}

func TestCannotReferSelf(t *testing.T) {
	err := CannotReferSelf("user-1", "user-1")
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrValidation)
}

func TestCannotReferSelf_DifferentUsers(t *testing.T) {
	err := CannotReferSelf("user-1", "user-2")
	assert.NoError(t, err)
}

func TestValidateStatus(t *testing.T) {
	assert.True(t, ValidateStatus("pending"))
	assert.True(t, ValidateStatus("completed"))
	assert.True(t, ValidateStatus("rewarded"))
	assert.False(t, ValidateStatus("invalid"))
	assert.False(t, ValidateStatus(""))
}

func TestValidateRewardType(t *testing.T) {
	assert.True(t, ValidateRewardType("credit"))
	assert.True(t, ValidateRewardType("trial_extension"))
	assert.False(t, ValidateRewardType("invalid"))
	assert.False(t, ValidateRewardType(""))
}

func TestReferral_Complete_Success(t *testing.T) {
	r := &Referral{Status: StatusPending}
	err := r.Complete()
	assert.NoError(t, err)
	assert.Equal(t, StatusCompleted, r.Status)
	assert.NotNil(t, r.CompletedAt)
}

func TestReferral_Complete_NotPending(t *testing.T) {
	r := &Referral{Status: StatusCompleted}
	err := r.Complete()
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrValidation)
}

func TestReferral_MarkRewarded_Success(t *testing.T) {
	r := &Referral{Status: StatusCompleted}
	err := r.MarkRewarded(RewardCredit, 500)
	assert.NoError(t, err)
	assert.Equal(t, StatusRewarded, r.Status)
	assert.Equal(t, RewardCredit, r.RewardType)
	assert.Equal(t, 500, r.RewardAmount)
	assert.NotNil(t, r.RewardedAt)
}

func TestReferral_MarkRewarded_NotCompleted(t *testing.T) {
	r := &Referral{Status: StatusPending}
	err := r.MarkRewarded(RewardCredit, 500)
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrValidation)
}

func TestStatus_Values(t *testing.T) {
	assert.Equal(t, Status("pending"), StatusPending)
	assert.Equal(t, Status("completed"), StatusCompleted)
	assert.Equal(t, Status("rewarded"), StatusRewarded)
}

func TestRewardType_Values(t *testing.T) {
	assert.Equal(t, RewardType("credit"), RewardCredit)
	assert.Equal(t, RewardType("trial_extension"), RewardTrialExtension)
}
