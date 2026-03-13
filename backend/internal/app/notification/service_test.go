package notification

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
	domainnotif "github.com/hassad/boilerplateSaaS/backend/internal/domain/notification"
)

// Mocks

type mockNotifRepo struct {
	createFn       func(ctx context.Context, n *domainnotif.Notification) error
	findByIDFn     func(ctx context.Context, id string) (*domainnotif.Notification, error)
	listByUserIDFn func(ctx context.Context, userID string, unreadOnly bool, offset, limit int) ([]*domainnotif.Notification, int, error)
	markReadFn     func(ctx context.Context, id string) error
	markAllReadFn  func(ctx context.Context, userID string) error
	unreadCountFn  func(ctx context.Context, userID string) (int, error)
}

func (m *mockNotifRepo) Create(ctx context.Context, n *domainnotif.Notification) error {
	if m.createFn != nil {
		return m.createFn(ctx, n)
	}
	n.ID = "notif-1"
	n.CreatedAt = time.Now()
	return nil
}

func (m *mockNotifRepo) FindByID(ctx context.Context, id string) (*domainnotif.Notification, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, id)
	}
	return nil, domain.ErrNotFound
}

func (m *mockNotifRepo) ListByUserID(ctx context.Context, userID string, unreadOnly bool, offset, limit int) ([]*domainnotif.Notification, int, error) {
	if m.listByUserIDFn != nil {
		return m.listByUserIDFn(ctx, userID, unreadOnly, offset, limit)
	}
	return nil, 0, nil
}

func (m *mockNotifRepo) MarkRead(ctx context.Context, id string) error {
	if m.markReadFn != nil {
		return m.markReadFn(ctx, id)
	}
	return nil
}

func (m *mockNotifRepo) MarkAllRead(ctx context.Context, userID string) error {
	if m.markAllReadFn != nil {
		return m.markAllReadFn(ctx, userID)
	}
	return nil
}

func (m *mockNotifRepo) UnreadCount(ctx context.Context, userID string) (int, error) {
	if m.unreadCountFn != nil {
		return m.unreadCountFn(ctx, userID)
	}
	return 0, nil
}

// Tests

func TestNotificationService_Send_Success(t *testing.T) {
	repo := &mockNotifRepo{}
	svc := NewService(repo)

	n, err := svc.Send(context.Background(), "user-1", "info", "Welcome!", "Thanks for signing up")
	assert.NoError(t, err)
	assert.Equal(t, "user-1", n.UserID)
	assert.Equal(t, "Welcome!", n.Title)
	assert.Equal(t, "info", n.Type)
}

func TestNotificationService_MarkAsRead_Success(t *testing.T) {
	var markedID string
	repo := &mockNotifRepo{
		findByIDFn: func(_ context.Context, id string) (*domainnotif.Notification, error) {
			return &domainnotif.Notification{ID: id, UserID: "user-1"}, nil
		},
		markReadFn: func(_ context.Context, id string) error {
			markedID = id
			return nil
		},
	}
	svc := NewService(repo)

	err := svc.MarkAsRead(context.Background(), "user-1", "notif-1")
	assert.NoError(t, err)
	assert.Equal(t, "notif-1", markedID)
}

func TestNotificationService_MarkAsRead_Forbidden(t *testing.T) {
	repo := &mockNotifRepo{
		findByIDFn: func(_ context.Context, _ string) (*domainnotif.Notification, error) {
			return &domainnotif.Notification{ID: "notif-1", UserID: "user-2"}, nil
		},
	}
	svc := NewService(repo)

	err := svc.MarkAsRead(context.Background(), "user-1", "notif-1")
	assert.ErrorIs(t, err, domain.ErrForbidden)
}

func TestNotificationService_UnreadCount(t *testing.T) {
	repo := &mockNotifRepo{
		unreadCountFn: func(_ context.Context, _ string) (int, error) {
			return 5, nil
		},
	}
	svc := NewService(repo)

	count, err := svc.UnreadCount(context.Background(), "user-1")
	assert.NoError(t, err)
	assert.Equal(t, 5, count)
}
