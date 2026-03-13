package user

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
	domainuser "github.com/hassad/boilerplateSaaS/backend/internal/domain/user"
	"github.com/hassad/boilerplateSaaS/backend/pkg/hash"
)

// mockUserRepo is a manual mock for repository.UserRepository.
type mockUserRepo struct {
	findByIDFn func(ctx context.Context, id string) (*domainuser.User, error)
	updateFn   func(ctx context.Context, u *domainuser.User) error
	deleteFn   func(ctx context.Context, id string) error
}

func (m *mockUserRepo) Create(_ context.Context, _ *domainuser.User) error { return nil }
func (m *mockUserRepo) FindByID(ctx context.Context, id string) (*domainuser.User, error) {
	return m.findByIDFn(ctx, id)
}
func (m *mockUserRepo) FindByEmail(_ context.Context, _ string) (*domainuser.User, error) {
	return nil, domain.ErrNotFound
}
func (m *mockUserRepo) FindByProvider(_ context.Context, _, _ string) (*domainuser.User, error) {
	return nil, domain.ErrNotFound
}
func (m *mockUserRepo) Update(ctx context.Context, u *domainuser.User) error {
	return m.updateFn(ctx, u)
}
func (m *mockUserRepo) Delete(ctx context.Context, id string) error {
	return m.deleteFn(ctx, id)
}
func (m *mockUserRepo) List(_ context.Context, _, _ int) ([]*domainuser.User, int, error) {
	return nil, 0, nil
}

func TestService_ChangePassword_Success(t *testing.T) {
	hashed, _ := hash.Password("oldpassword")
	repo := &mockUserRepo{
		findByIDFn: func(_ context.Context, _ string) (*domainuser.User, error) {
			return &domainuser.User{ID: "user-1", PasswordHash: hashed}, nil
		},
		updateFn: func(_ context.Context, _ *domainuser.User) error {
			return nil
		},
	}

	svc := NewService(repo)
	err := svc.ChangePassword(context.Background(), "user-1", "oldpassword", "newpassword")
	assert.NoError(t, err)
}

func TestService_ChangePassword_WrongOldPassword(t *testing.T) {
	hashed, _ := hash.Password("correctpassword")
	repo := &mockUserRepo{
		findByIDFn: func(_ context.Context, _ string) (*domainuser.User, error) {
			return &domainuser.User{ID: "user-1", PasswordHash: hashed}, nil
		},
	}

	svc := NewService(repo)
	err := svc.ChangePassword(context.Background(), "user-1", "wrongpassword", "newpassword")
	assert.ErrorIs(t, err, domain.ErrUnauthorized)
}

func TestService_ChangePassword_SamePassword(t *testing.T) {
	hashed, _ := hash.Password("samepassword")
	repo := &mockUserRepo{
		findByIDFn: func(_ context.Context, _ string) (*domainuser.User, error) {
			return &domainuser.User{ID: "user-1", PasswordHash: hashed}, nil
		},
	}

	svc := NewService(repo)
	err := svc.ChangePassword(context.Background(), "user-1", "samepassword", "samepassword")
	assert.ErrorIs(t, err, domain.ErrValidation)
}
