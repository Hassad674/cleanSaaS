package response

import (
	"github.com/hassad/boilerplateSaaS/backend/internal/domain/referral"
)

type ReferralResponse struct {
	ID           string  `json:"id"`
	ReferrerID   string  `json:"referrer_id"`
	ReferredID   string  `json:"referred_id,omitempty"`
	Code         string  `json:"code"`
	Status       string  `json:"status"`
	RewardType   string  `json:"reward_type"`
	RewardAmount int     `json:"reward_amount"`
	CreatedAt    string  `json:"created_at"`
	CompletedAt  *string `json:"completed_at,omitempty"`
	RewardedAt   *string `json:"rewarded_at,omitempty"`
}

func ReferralFromDomain(r *referral.Referral) ReferralResponse {
	resp := ReferralResponse{
		ID:           r.ID,
		ReferrerID:   r.ReferrerID,
		ReferredID:   r.ReferredID,
		Code:         r.Code,
		Status:       string(r.Status),
		RewardType:   string(r.RewardType),
		RewardAmount: r.RewardAmount,
		CreatedAt:    r.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
	if r.CompletedAt != nil {
		s := r.CompletedAt.Format("2006-01-02T15:04:05Z")
		resp.CompletedAt = &s
	}
	if r.RewardedAt != nil {
		s := r.RewardedAt.Format("2006-01-02T15:04:05Z")
		resp.RewardedAt = &s
	}
	return resp
}
