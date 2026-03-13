package user

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
)

func TestNew_Success(t *testing.T) {
	u, err := New("test@example.com", "Test User", "hashedpw")
	assert.NoError(t, err)
	assert.Equal(t, "test@example.com", u.Email)
	assert.Equal(t, "Test User", u.Name)
	assert.Equal(t, RoleMember, u.Role)
	assert.Equal(t, "email", u.Provider)
	assert.False(t, u.EmailVerified)
}

func TestNew_EmptyEmail(t *testing.T) {
	_, err := New("", "Test User", "hashedpw")
	assert.ErrorIs(t, err, domain.ErrValidation)
}

func TestNew_EmptyName(t *testing.T) {
	_, err := New("test@example.com", "", "hashedpw")
	assert.ErrorIs(t, err, domain.ErrValidation)
}

func TestUser_IsAdmin(t *testing.T) {
	u, _ := New("admin@test.com", "Admin", "pw")
	assert.False(t, u.IsAdmin())

	u.Role = RoleAdmin
	assert.True(t, u.IsAdmin())
}

func TestUser_VerifyEmail(t *testing.T) {
	u, _ := New("test@test.com", "User", "pw")
	assert.False(t, u.EmailVerified)

	u.VerifyEmail()
	assert.True(t, u.EmailVerified)
}

func TestRole_IsValid(t *testing.T) {
	assert.True(t, RoleAdmin.IsValid())
	assert.True(t, RoleMember.IsValid())
	assert.False(t, Role("invalid").IsValid())
}
