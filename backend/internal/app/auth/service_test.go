package auth

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
	"github.com/hassad/boilerplateSaaS/backend/internal/domain/user"
	"github.com/hassad/boilerplateSaaS/backend/internal/port/repository"
	"github.com/hassad/boilerplateSaaS/backend/internal/port/service"
	"github.com/hassad/boilerplateSaaS/backend/pkg/hash"
	"github.com/hassad/boilerplateSaaS/backend/pkg/jwt"
)

// Mock implementations

type mockUserRepo struct {
	findByEmailFn func(ctx context.Context, email string) (*user.User, error)
	findByIDFn    func(ctx context.Context, id string) (*user.User, error)
	createFn      func(ctx context.Context, u *user.User) error
	updateFn      func(ctx context.Context, u *user.User) error
}

func (m *mockUserRepo) Create(ctx context.Context, u *user.User) error {
	if m.createFn != nil {
		return m.createFn(ctx, u)
	}
	return nil
}
func (m *mockUserRepo) FindByID(ctx context.Context, id string) (*user.User, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, id)
	}
	return nil, domain.ErrNotFound
}
func (m *mockUserRepo) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	if m.findByEmailFn != nil {
		return m.findByEmailFn(ctx, email)
	}
	return nil, domain.ErrNotFound
}
func (m *mockUserRepo) FindByProvider(_ context.Context, _, _ string) (*user.User, error) {
	return nil, domain.ErrNotFound
}
func (m *mockUserRepo) Update(ctx context.Context, u *user.User) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, u)
	}
	return nil
}
func (m *mockUserRepo) Delete(_ context.Context, _ string) error { return nil }
func (m *mockUserRepo) List(_ context.Context, _, _ int) ([]*user.User, int, error) {
	return nil, 0, nil
}

type mockResetRepo struct {
	createFn      func(ctx context.Context, pr *repository.PasswordReset) error
	findByTokenFn func(ctx context.Context, token string) (*repository.PasswordReset, error)
	markUsedFn    func(ctx context.Context, id string) error
}

func (m *mockResetRepo) Create(ctx context.Context, pr *repository.PasswordReset) error {
	if m.createFn != nil {
		return m.createFn(ctx, pr)
	}
	return nil
}
func (m *mockResetRepo) FindByToken(ctx context.Context, token string) (*repository.PasswordReset, error) {
	if m.findByTokenFn != nil {
		return m.findByTokenFn(ctx, token)
	}
	return nil, domain.ErrNotFound
}
func (m *mockResetRepo) MarkUsed(ctx context.Context, id string) error {
	if m.markUsedFn != nil {
		return m.markUsedFn(ctx, id)
	}
	return nil
}
func (m *mockResetRepo) DeleteExpired(_ context.Context) error { return nil }

type mockEmailSvc struct {
	sendFn         func(ctx context.Context, email service.Email) error
	sendTemplateFn func(ctx context.Context, to, template string, data map[string]string) error
}

func (m *mockEmailSvc) Send(ctx context.Context, email service.Email) error {
	if m.sendFn != nil {
		return m.sendFn(ctx, email)
	}
	return nil
}
func (m *mockEmailSvc) SendTemplate(ctx context.Context, to, template string, data map[string]string) error {
	if m.sendTemplateFn != nil {
		return m.sendTemplateFn(ctx, to, template, data)
	}
	return nil
}

func TestAuthService_ForgotPassword_UserExists(t *testing.T) {
	var sentEmail bool
	userRepo := &mockUserRepo{
		findByEmailFn: func(_ context.Context, _ string) (*user.User, error) {
			return &user.User{ID: "u1", Email: "test@test.com", Name: "Test"}, nil
		},
	}
	resetRepo := &mockResetRepo{
		createFn: func(_ context.Context, _ *repository.PasswordReset) error { return nil },
	}
	emailSvc := &mockEmailSvc{
		sendTemplateFn: func(_ context.Context, _, _ string, _ map[string]string) error {
			sentEmail = true
			return nil
		},
	}

	svc := NewService(userRepo, resetRepo, emailSvc, jwt.NewMaker("secret"), "http://localhost:3006")
	err := svc.ForgotPassword(context.Background(), "test@test.com", "")
	assert.NoError(t, err)
	assert.True(t, sentEmail, "email should have been sent")
}

func TestAuthService_ForgotPassword_UserNotFound(t *testing.T) {
	userRepo := &mockUserRepo{
		findByEmailFn: func(_ context.Context, _ string) (*user.User, error) {
			return nil, domain.ErrNotFound
		},
	}
	resetRepo := &mockResetRepo{}

	svc := NewService(userRepo, resetRepo, nil, jwt.NewMaker("secret"), "http://localhost:3006")
	err := svc.ForgotPassword(context.Background(), "noone@test.com", "")
	assert.NoError(t, err, "should not return error for non-existent user")
}

func TestAuthService_ResetPassword_ValidToken(t *testing.T) {
	var updated bool
	hashed, _ := hash.Password("oldpass")

	userRepo := &mockUserRepo{
		findByIDFn: func(_ context.Context, _ string) (*user.User, error) {
			return &user.User{ID: "u1", PasswordHash: hashed}, nil
		},
		updateFn: func(_ context.Context, _ *user.User) error {
			updated = true
			return nil
		},
	}
	resetRepo := &mockResetRepo{
		findByTokenFn: func(_ context.Context, _ string) (*repository.PasswordReset, error) {
			return &repository.PasswordReset{
				ID:        "r1",
				UserID:    "u1",
				Token:     "valid-token",
				ExpiresAt: time.Now().Add(time.Hour),
				Used:      false,
			}, nil
		},
		markUsedFn: func(_ context.Context, _ string) error { return nil },
	}

	svc := NewService(userRepo, resetRepo, nil, jwt.NewMaker("secret"), "http://localhost:3006")
	err := svc.ResetPassword(context.Background(), "valid-token", "newpassword")
	assert.NoError(t, err)
	assert.True(t, updated, "user password should have been updated")
}

func TestAuthService_ResetPassword_ExpiredToken(t *testing.T) {
	resetRepo := &mockResetRepo{
		findByTokenFn: func(_ context.Context, _ string) (*repository.PasswordReset, error) {
			return &repository.PasswordReset{
				ID:        "r1",
				UserID:    "u1",
				Token:     "expired-token",
				ExpiresAt: time.Now().Add(-time.Hour),
				Used:      false,
			}, nil
		},
	}

	svc := NewService(&mockUserRepo{}, resetRepo, nil, jwt.NewMaker("secret"), "http://localhost:3006")
	err := svc.ResetPassword(context.Background(), "expired-token", "newpassword")
	assert.ErrorIs(t, err, domain.ErrExpiredToken)
}

func TestAuthService_ResetPassword_UsedToken(t *testing.T) {
	resetRepo := &mockResetRepo{
		findByTokenFn: func(_ context.Context, _ string) (*repository.PasswordReset, error) {
			return &repository.PasswordReset{
				ID:        "r1",
				UserID:    "u1",
				Token:     "used-token",
				ExpiresAt: time.Now().Add(time.Hour),
				Used:      true,
			}, nil
		},
	}

	svc := NewService(&mockUserRepo{}, resetRepo, nil, jwt.NewMaker("secret"), "http://localhost:3006")
	err := svc.ResetPassword(context.Background(), "used-token", "newpassword")
	assert.ErrorIs(t, err, domain.ErrInvalidToken)
}
