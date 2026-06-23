package billing

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
)

func future() time.Time { return time.Now().Add(DefaultPeriod) }
func past() time.Time   { return time.Now().Add(-time.Hour) }

func TestNewSubscription(t *testing.T) {
	t.Run("defaults to a 30-day active period", func(t *testing.T) {
		sub, err := NewSubscription("u1", "p1", "sub_1", 0)
		assert.NoError(t, err)
		assert.Equal(t, StatusActive, sub.Status)
		assert.Equal(t, "u1", sub.UserID)
		assert.Equal(t, "p1", sub.PlanID)
		assert.False(t, sub.CancelAtPeriodEnd)
		// Period end is ~30 days out from start.
		assert.WithinDuration(t, sub.CurrentPeriodStart.Add(DefaultPeriod), sub.CurrentPeriodEnd, time.Second)
		assert.True(t, sub.CurrentPeriodEnd.After(time.Now()))
	})

	t.Run("honors an explicit period", func(t *testing.T) {
		sub, err := NewSubscription("u1", "p1", "sub_1", 7*24*time.Hour)
		assert.NoError(t, err)
		assert.WithinDuration(t, sub.CurrentPeriodStart.Add(7*24*time.Hour), sub.CurrentPeriodEnd, time.Second)
	})

	t.Run("requires a user", func(t *testing.T) {
		_, err := NewSubscription("", "p1", "sub_1", 0)
		assert.ErrorIs(t, err, domain.ErrValidation)
	})

	t.Run("requires a plan", func(t *testing.T) {
		_, err := NewSubscription("u1", "", "sub_1", 0)
		assert.ErrorIs(t, err, domain.ErrValidation)
	})
}

func TestSubscription_Activate(t *testing.T) {
	t.Run("reactivates and clears pending cancellation", func(t *testing.T) {
		sub := &Subscription{Status: StatusPastDue, PlanID: "p_old", CancelAtPeriodEnd: true}
		err := sub.Activate("p_new", future())
		assert.NoError(t, err)
		assert.Equal(t, StatusActive, sub.Status)
		assert.Equal(t, "p_new", sub.PlanID)
		assert.False(t, sub.CancelAtPeriodEnd)
	})

	t.Run("rejects a canceled (terminal) subscription", func(t *testing.T) {
		sub := &Subscription{Status: StatusCanceled, PlanID: "p1"}
		err := sub.Activate("p2", future())
		assert.ErrorIs(t, err, domain.ErrValidation)
		assert.Equal(t, StatusCanceled, sub.Status)
	})

	t.Run("rejects an empty plan", func(t *testing.T) {
		sub := &Subscription{Status: StatusInactive}
		err := sub.Activate("", future())
		assert.ErrorIs(t, err, domain.ErrValidation)
	})

	t.Run("rejects a past period end", func(t *testing.T) {
		sub := &Subscription{Status: StatusInactive}
		err := sub.Activate("p1", past())
		assert.ErrorIs(t, err, domain.ErrValidation)
	})
}

func TestSubscription_Renew(t *testing.T) {
	t.Run("extends an active subscription", func(t *testing.T) {
		sub := &Subscription{Status: StatusActive, CurrentPeriodEnd: time.Now()}
		newEnd := future()
		err := sub.Renew(newEnd)
		assert.NoError(t, err)
		assert.Equal(t, newEnd, sub.CurrentPeriodEnd)
		assert.WithinDuration(t, time.Now(), sub.CurrentPeriodStart, time.Second)
	})

	t.Run("renews a trialing subscription", func(t *testing.T) {
		sub := &Subscription{Status: StatusTrialing}
		assert.NoError(t, sub.Renew(future()))
		assert.Equal(t, StatusTrialing, sub.Status)
	})

	t.Run("rejects renewing a canceled subscription", func(t *testing.T) {
		sub := &Subscription{Status: StatusCanceled}
		err := sub.Renew(future())
		assert.ErrorIs(t, err, domain.ErrValidation)
	})

	t.Run("rejects a past period end", func(t *testing.T) {
		sub := &Subscription{Status: StatusActive}
		err := sub.Renew(past())
		assert.ErrorIs(t, err, domain.ErrValidation)
	})
}

func TestSubscription_ChangePlan(t *testing.T) {
	t.Run("moves an active subscription to a new plan", func(t *testing.T) {
		sub := &Subscription{Status: StatusActive, PlanID: "p_old"}
		err := sub.ChangePlan("p_new")
		assert.NoError(t, err)
		assert.Equal(t, "p_new", sub.PlanID)
	})

	t.Run("rejects changing to the same plan", func(t *testing.T) {
		sub := &Subscription{Status: StatusActive, PlanID: "p1"}
		err := sub.ChangePlan("p1")
		assert.ErrorIs(t, err, domain.ErrValidation)
	})

	t.Run("rejects an empty plan", func(t *testing.T) {
		sub := &Subscription{Status: StatusActive, PlanID: "p1"}
		err := sub.ChangePlan("")
		assert.ErrorIs(t, err, domain.ErrValidation)
	})

	t.Run("rejects changing plan on a canceled subscription", func(t *testing.T) {
		sub := &Subscription{Status: StatusCanceled, PlanID: "p1"}
		err := sub.ChangePlan("p2")
		assert.ErrorIs(t, err, domain.ErrValidation)
		assert.Equal(t, "p1", sub.PlanID)
	})
}

func TestSubscription_Cancel_Transitions(t *testing.T) {
	t.Run("schedules cancellation on a live subscription", func(t *testing.T) {
		sub := &Subscription{Status: StatusActive}
		assert.NoError(t, sub.Cancel())
		assert.True(t, sub.CancelAtPeriodEnd)
		assert.Equal(t, StatusActive, sub.Status, "access kept until period end")
	})

	t.Run("rejects canceling an already-canceling subscription", func(t *testing.T) {
		sub := &Subscription{Status: StatusActive, CancelAtPeriodEnd: true}
		err := sub.Cancel()
		assert.ErrorIs(t, err, domain.ErrValidation)
	})

	t.Run("rejects canceling a terminal (canceled) subscription", func(t *testing.T) {
		sub := &Subscription{Status: StatusCanceled}
		err := sub.Cancel()
		assert.ErrorIs(t, err, domain.ErrValidation)
	})

	t.Run("rejects canceling an inactive subscription", func(t *testing.T) {
		sub := &Subscription{Status: StatusInactive}
		err := sub.Cancel()
		assert.ErrorIs(t, err, domain.ErrValidation)
	})
}

func TestSubscription_MarkCanceled(t *testing.T) {
	t.Run("moves a live subscription to terminal", func(t *testing.T) {
		sub := &Subscription{Status: StatusActive, UpdatedAt: time.Now().Add(-time.Hour)}
		before := sub.UpdatedAt
		sub.MarkCanceled()
		assert.Equal(t, StatusCanceled, sub.Status)
		assert.True(t, sub.UpdatedAt.After(before))
	})

	t.Run("is idempotent on an already-canceled subscription", func(t *testing.T) {
		sub := &Subscription{Status: StatusCanceled, UpdatedAt: time.Now().Add(-time.Hour)}
		before := sub.UpdatedAt
		sub.MarkCanceled()
		assert.Equal(t, StatusCanceled, sub.Status)
		assert.Equal(t, before, sub.UpdatedAt, "no-op must not touch UpdatedAt")
	})
}
