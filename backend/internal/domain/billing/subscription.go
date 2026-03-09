package billing

import "time"

type Status string

const (
	StatusActive   Status = "active"
	StatusCanceled Status = "canceled"
	StatusPastDue  Status = "past_due"
	StatusTrialing Status = "trialing"
	StatusInactive Status = "inactive"
)

type Subscription struct {
	ID                 string
	UserID             string
	PlanID             string
	StripeSubscription string
	Status             Status
	CurrentPeriodStart time.Time
	CurrentPeriodEnd   time.Time
	CancelAtPeriodEnd  bool
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (s *Subscription) IsActive() bool {
	return s.Status == StatusActive || s.Status == StatusTrialing
}

func (s *Subscription) CanCancel() bool {
	return s.IsActive() && !s.CancelAtPeriodEnd
}

func (s *Subscription) Cancel() {
	s.CancelAtPeriodEnd = true
	s.UpdatedAt = time.Now()
}
