package billing

import (
	"fmt"
	"time"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
)

type Status string

const (
	StatusActive   Status = "active"
	StatusCanceled Status = "canceled"
	StatusPastDue  Status = "past_due"
	StatusTrialing Status = "trialing"
	StatusInactive Status = "inactive"
)

// DefaultPeriod is the billing-period length used when an external source (e.g.
// a Stripe webhook) does not give us an explicit period end. It is a domain
// rule, not webhook glue: a freshly activated subscription runs for 30 days.
const DefaultPeriod = 30 * 24 * time.Hour

// Subscription is the billing aggregate root. It owns the subscription state
// machine: callers never mutate Status / period fields directly — they invoke
// the transition methods below, which enforce the invariants.
//
// Fields stay exported because the postgres repository scans rows straight into
// them and the response DTO reads them; that is the sanctioned
// repository/serialization boundary. Behavior that has invariants (status
// transitions, period validity) lives in the methods, not in the app service.
type Subscription struct {
	ID                 string
	UserID             string
	OrgID              string
	PlanID             string
	StripeSubscription string
	Status             Status
	CurrentPeriodStart time.Time
	CurrentPeriodEnd   time.Time
	CancelAtPeriodEnd  bool
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// NewSubscription creates an active subscription for a user on a plan. The
// period runs from now for the given duration; passing 0 applies DefaultPeriod.
// UserID and PlanID are required.
func NewSubscription(userID, planID, stripeSubscriptionID string, period time.Duration) (*Subscription, error) {
	if userID == "" {
		return nil, fmt.Errorf("subscription requires a user: %w", domain.ErrValidation)
	}
	if planID == "" {
		return nil, fmt.Errorf("subscription requires a plan: %w", domain.ErrValidation)
	}
	if period <= 0 {
		period = DefaultPeriod
	}

	now := time.Now()
	return &Subscription{
		UserID:             userID,
		PlanID:             planID,
		StripeSubscription: stripeSubscriptionID,
		Status:             StatusActive,
		CurrentPeriodStart: now,
		CurrentPeriodEnd:   now.Add(period),
	}, nil
}

func (s *Subscription) IsActive() bool {
	return s.Status == StatusActive || s.Status == StatusTrialing
}

func (s *Subscription) CanCancel() bool {
	return s.IsActive() && !s.CancelAtPeriodEnd
}

// Activate (re)activates the subscription on a plan with a new period end and
// clears any pending cancellation. The period end must be in the future and the
// plan is required. A subscription that has been fully canceled (terminal)
// cannot be reactivated. This is the rule behind the "subscription.updated"
// webhook.
func (s *Subscription) Activate(planID string, periodEnd time.Time) error {
	if s.Status == StatusCanceled {
		return fmt.Errorf("cannot activate a canceled subscription: %w", domain.ErrValidation)
	}
	if planID == "" {
		return fmt.Errorf("activate requires a plan: %w", domain.ErrValidation)
	}
	if !periodEnd.After(time.Now()) {
		return fmt.Errorf("period end must be in the future: %w", domain.ErrValidation)
	}

	s.PlanID = planID
	s.Status = StatusActive
	s.CurrentPeriodEnd = periodEnd
	s.CancelAtPeriodEnd = false
	s.touch()
	return nil
}

// Renew extends an active subscription to a new period end (e.g. on a paid
// invoice). The new end must be in the future, and only a live (active or
// trialing) subscription can be renewed.
func (s *Subscription) Renew(periodEnd time.Time) error {
	if !s.IsActive() {
		return fmt.Errorf("cannot renew a %s subscription: %w", s.Status, domain.ErrValidation)
	}
	if !periodEnd.After(time.Now()) {
		return fmt.Errorf("period end must be in the future: %w", domain.ErrValidation)
	}

	s.CurrentPeriodStart = time.Now()
	s.CurrentPeriodEnd = periodEnd
	s.touch()
	return nil
}

// ChangePlan moves a live subscription to a different plan. The subscription
// must be active/trialing and the new plan must differ and be non-empty.
func (s *Subscription) ChangePlan(planID string) error {
	if !s.IsActive() {
		return fmt.Errorf("cannot change plan on a %s subscription: %w", s.Status, domain.ErrValidation)
	}
	if planID == "" {
		return fmt.Errorf("change plan requires a plan: %w", domain.ErrValidation)
	}
	if planID == s.PlanID {
		return fmt.Errorf("subscription is already on plan %q: %w", planID, domain.ErrValidation)
	}

	s.PlanID = planID
	s.touch()
	return nil
}

// Cancel schedules the subscription to end at the close of the current period
// (the user keeps access until then). It is only valid on a live subscription
// that is not already scheduled for cancellation.
func (s *Subscription) Cancel() error {
	if !s.CanCancel() {
		return fmt.Errorf("subscription cannot be canceled in state %s (cancel-scheduled=%t): %w",
			s.Status, s.CancelAtPeriodEnd, domain.ErrValidation)
	}
	s.CancelAtPeriodEnd = true
	s.touch()
	return nil
}

// MarkCanceled immediately moves the subscription to the terminal canceled
// state (e.g. on a "subscription.deleted" webhook). It is idempotent: marking an
// already-canceled subscription is a no-op rather than an error, so a Stripe
// retry is safe.
func (s *Subscription) MarkCanceled() {
	if s.Status == StatusCanceled {
		return
	}
	s.Status = StatusCanceled
	s.touch()
}

func (s *Subscription) touch() {
	s.UpdatedAt = time.Now()
}
