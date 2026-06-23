//go:build integration

package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
	"github.com/hassad/boilerplateSaaS/backend/internal/domain/ai"
	"github.com/hassad/boilerplateSaaS/backend/internal/domain/billing"
	"github.com/hassad/boilerplateSaaS/backend/internal/domain/notification"
	"github.com/hassad/boilerplateSaaS/backend/internal/domain/storage"
	"github.com/hassad/boilerplateSaaS/backend/internal/port/repository"
)

// The tenant-scoped repos (files, conversations, notifications, subscriptions)
// run under OrgScope: a transaction that switches to the restricted app_user role
// and sets app.current_org_id (setOrgScope), exactly as the request path does.
// Every query inside the callback is enforced by both the repository's
// WHERE org_id filter (layer 2) and PostgreSQL RLS (layer 3). The short scope
// transaction commits on a nil return, so writes persist for the next assertion.
//
// Setup (users + orgs) is done on the privileged pool — the system path — because
// app_user cannot create orgs/users, mirroring how signup runs before any request
// is scoped. Cleanup deletes the owning user, which cascades to org + tenant rows.

// TestFileScope_RoundTrip exercises the org-scoped file repository through
// OrgScope: create stamps the active org_id, find/list are org-filtered + RLS
// enforced, pagination works, and delete is org-scoped.
func TestFileScope_RoundTrip(t *testing.T) {
	db := openTestDB(t)
	scope := NewOrgScope(db)

	base := context.Background()
	owner := newUser(base, t, db)
	o := newOrg(base, t, db, owner.ID)
	ctx := orgScopedCtx(o.ID)

	var fileID string
	require.NoError(t, scope.WithOrgFiles(ctx, func(files repository.FileRepository) error {
		f := &storage.File{UserID: owner.ID, Name: "doc.txt", Key: "itest/" + uniqueTag(), SizeBytes: 12, ContentType: "text/plain", URL: "https://x.test/1"}
		if err := files.Create(ctx, f); err != nil {
			return err
		}
		fileID = f.ID
		assert.NotEmpty(t, f.ID, "Create stamps id + org_id via RETURNING under RLS")
		return nil
	}))

	// A separate scope transaction reads it back (org-filtered + RLS).
	require.NoError(t, scope.WithOrgFiles(ctx, func(files repository.FileRepository) error {
		got, err := files.FindByID(ctx, fileID)
		require.NoError(t, err)
		assert.Equal(t, owner.ID, got.UserID)

		// A few more files to verify COUNT + LIMIT/OFFSET pagination.
		for i := 0; i < 2; i++ {
			require.NoError(t, files.Create(ctx, &storage.File{UserID: owner.ID, Name: "n", Key: "itest/" + uniqueTag(), SizeBytes: 1, ContentType: "text/plain"}))
		}
		page, total, err := files.ListByUserID(ctx, owner.ID, 0, 2)
		require.NoError(t, err)
		assert.Equal(t, 3, total, "org-scoped COUNT sees exactly this org's 3 files")
		assert.Len(t, page, 2, "LIMIT 2 yields a 2-row page")
		return nil
	}))

	// Delete is org-scoped; a second delete reports not-found.
	require.NoError(t, scope.WithOrgFiles(ctx, func(files repository.FileRepository) error {
		require.NoError(t, files.Delete(ctx, fileID))
		assert.ErrorIs(t, files.Delete(ctx, fileID), domain.ErrNotFound)
		return nil
	}))
}

// TestFileScope_CrossOrgInvisible proves a file created under org A is invisible
// when the same repository runs scoped to a different org B — the org_id filter
// and RLS together deny cross-tenant reads.
func TestFileScope_CrossOrgInvisible(t *testing.T) {
	db := openTestDB(t)
	scope := NewOrgScope(db)

	base := context.Background()
	owner := newUser(base, t, db)
	orgA := newOrg(base, t, db, owner.ID)
	orgB := newOrg(base, t, db, owner.ID)

	var fileID string
	ctxA := orgScopedCtx(orgA.ID)
	require.NoError(t, scope.WithOrgFiles(ctxA, func(files repository.FileRepository) error {
		f := &storage.File{UserID: owner.ID, Name: "secret.txt", Key: "itest/" + uniqueTag(), SizeBytes: 1, ContentType: "text/plain"}
		if err := files.Create(ctxA, f); err != nil {
			return err
		}
		fileID = f.ID
		return nil
	}))

	// Under org B's scope, org A's file is not found.
	ctxB := orgScopedCtx(orgB.ID)
	require.NoError(t, scope.WithOrgFiles(ctxB, func(files repository.FileRepository) error {
		_, err := files.FindByID(ctxB, fileID)
		assert.ErrorIs(t, err, domain.ErrNotFound, "org B cannot see org A's file")
		return nil
	}))
}

// TestConversationScope_RoundTrip exercises the org-scoped conversation repo:
// create, add messages (transitively RLS-scoped), find-by-id loads messages in
// order, list pagination, update, and delete.
func TestConversationScope_RoundTrip(t *testing.T) {
	db := openTestDB(t)
	scope := NewOrgScope(db)

	base := context.Background()
	owner := newUser(base, t, db)
	o := newOrg(base, t, db, owner.ID)
	ctx := orgScopedCtx(o.ID)

	var convID string
	require.NoError(t, scope.WithOrgConversations(ctx, func(convs repository.ConversationRepository) error {
		c := &ai.Conversation{UserID: owner.ID, Title: "Chat " + uniqueTag()}
		if err := convs.Create(ctx, c); err != nil {
			return err
		}
		convID = c.ID

		// AddMessage verifies the conversation belongs to the active org first,
		// then inserts the message (RLS scopes messages transitively).
		if err := convs.AddMessage(ctx, c.ID, ai.Message{Role: ai.RoleUser, Content: "hello"}); err != nil {
			return err
		}
		return convs.AddMessage(ctx, c.ID, ai.Message{Role: ai.RoleAssistant, Content: "hi there"})
	}))

	require.NoError(t, scope.WithOrgConversations(ctx, func(convs repository.ConversationRepository) error {
		got, err := convs.FindByID(ctx, convID)
		require.NoError(t, err)
		require.Len(t, got.Messages, 2, "both messages loaded")
		assert.Equal(t, ai.RoleUser, got.Messages[0].Role, "messages ordered by created_at ASC")
		assert.Equal(t, "hi there", got.Messages[1].Content)

		// Adding a message to a non-existent conversation reports not-found.
		assert.ErrorIs(t, convs.AddMessage(ctx, uniqueTag(), ai.Message{Role: ai.RoleUser, Content: "x"}), domain.ErrNotFound)

		// Pagination over conversations for this user.
		page, total, err := convs.ListByUserID(ctx, owner.ID, 0, 10)
		require.NoError(t, err)
		assert.Equal(t, 1, total)
		assert.Len(t, page, 1)

		// Update title.
		got.Title = "Renamed"
		require.NoError(t, convs.Update(ctx, got))
		return nil
	}))

	require.NoError(t, scope.WithOrgConversations(ctx, func(convs repository.ConversationRepository) error {
		got, err := convs.FindByID(ctx, convID)
		require.NoError(t, err)
		assert.Equal(t, "Renamed", got.Title)

		require.NoError(t, convs.Delete(ctx, convID))
		assert.ErrorIs(t, convs.Delete(ctx, convID), domain.ErrNotFound)
		return nil
	}))
}

// TestNotificationScope_RoundTrip exercises the org-scoped notification repo:
// create (JSONB data), find, unread filtering + pagination, mark read, mark all
// read, and unread counts.
func TestNotificationScope_RoundTrip(t *testing.T) {
	db := openTestDB(t)
	scope := NewOrgScope(db)

	base := context.Background()
	owner := newUser(base, t, db)
	o := newOrg(base, t, db, owner.ID)
	ctx := orgScopedCtx(o.ID)

	var firstID string
	require.NoError(t, scope.WithOrgNotifications(ctx, func(notifs repository.NotificationRepository) error {
		n := &notification.Notification{UserID: owner.ID, Type: "info", Title: "First", Message: "m1", Data: map[string]interface{}{"k": "v"}}
		if err := notifs.Create(ctx, n); err != nil {
			return err
		}
		firstID = n.ID
		// A second notification, no data (exercises the {} default path).
		return notifs.Create(ctx, &notification.Notification{UserID: owner.ID, Type: "info", Title: "Second", Message: "m2"})
	}))

	require.NoError(t, scope.WithOrgNotifications(ctx, func(notifs repository.NotificationRepository) error {
		got, err := notifs.FindByID(ctx, firstID)
		require.NoError(t, err)
		assert.Equal(t, "First", got.Title)
		assert.Equal(t, "v", got.Data["k"], "JSONB data round-trips")

		count, err := notifs.UnreadCount(ctx, owner.ID)
		require.NoError(t, err)
		assert.Equal(t, 2, count, "both notifications start unread")

		// unreadOnly pagination.
		page, total, err := notifs.ListByUserID(ctx, owner.ID, true, 0, 1)
		require.NoError(t, err)
		assert.Equal(t, 2, total)
		assert.Len(t, page, 1)

		// Mark one read; unread count drops to 1.
		require.NoError(t, notifs.MarkRead(ctx, firstID))
		count, err = notifs.UnreadCount(ctx, owner.ID)
		require.NoError(t, err)
		assert.Equal(t, 1, count)

		// Marking a non-existent notification read reports not-found.
		assert.ErrorIs(t, notifs.MarkRead(ctx, uniqueTag()), domain.ErrNotFound)

		// Mark all read; unread count is then 0.
		require.NoError(t, notifs.MarkAllRead(ctx, owner.ID))
		count, err = notifs.UnreadCount(ctx, owner.ID)
		require.NoError(t, err)
		assert.Equal(t, 0, count)
		return nil
	}))
}

// TestSubscriptionScope_OrgScopedRead proves the authenticated READ path:
// FindByUserID is org-filtered, so it returns the subscription created for the
// active org and NOT a subscription belonging to a different org of the same
// user. The subscription is written on the system path (privileged pool, RLS
// bypassed) as the webhook does, then read back under the org scope.
func TestSubscriptionScope_OrgScopedRead(t *testing.T) {
	db := openTestDB(t)
	scope := NewOrgScope(db)
	sysRepo := NewSubscriptionRepository(db)
	planRepo := NewPlanRepository(db)

	base := context.Background()
	owner := newUser(base, t, db)
	orgA := newOrg(base, t, db, owner.ID)
	orgB := newOrg(base, t, db, owner.ID)
	plan := firstSeededPlan(base, t, planRepo)

	// Two subscriptions for the SAME user, one per org (system writes).
	now := time.Now()
	subA := &billing.Subscription{UserID: owner.ID, OrgID: orgA.ID, PlanID: plan.ID, StripeSubscription: "sub_" + uniqueTag(), Status: billing.StatusActive, CurrentPeriodStart: now, CurrentPeriodEnd: now.Add(30 * 24 * time.Hour)}
	require.NoError(t, sysRepo.Create(base, subA))
	subB := &billing.Subscription{UserID: owner.ID, OrgID: orgB.ID, PlanID: plan.ID, StripeSubscription: "sub_" + uniqueTag(), Status: billing.StatusTrialing, CurrentPeriodStart: now, CurrentPeriodEnd: now.Add(30 * 24 * time.Hour)}
	require.NoError(t, sysRepo.Create(base, subB))

	// Under org A's scope, FindByUserID returns ONLY org A's subscription.
	ctxA := orgScopedCtx(orgA.ID)
	require.NoError(t, scope.WithOrgSubscriptions(ctxA, func(subs repository.SubscriptionRepository) error {
		got, err := subs.FindByUserID(ctxA, owner.ID)
		require.NoError(t, err)
		assert.Equal(t, subA.ID, got.ID, "org scope returns org A's subscription")
		assert.Equal(t, billing.StatusActive, got.Status)
		return nil
	}))

	// Under org B's scope, the same query returns org B's subscription.
	ctxB := orgScopedCtx(orgB.ID)
	require.NoError(t, scope.WithOrgSubscriptions(ctxB, func(subs repository.SubscriptionRepository) error {
		got, err := subs.FindByUserID(ctxB, owner.ID)
		require.NoError(t, err)
		assert.Equal(t, subB.ID, got.ID, "org scope returns org B's subscription")
		assert.Equal(t, billing.StatusTrialing, got.Status)
		return nil
	}))
}
