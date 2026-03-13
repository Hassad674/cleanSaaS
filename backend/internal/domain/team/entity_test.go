package team

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
)

// --- NewTeam ---

func TestNewTeam_Success(t *testing.T) {
	tm, err := NewTeam("My Team", "owner-1")
	assert.NoError(t, err)
	assert.Equal(t, "My Team", tm.Name)
	assert.Equal(t, "my-team", tm.Slug)
	assert.Equal(t, "owner-1", tm.OwnerID)
	assert.Equal(t, "free", tm.Plan)
	assert.Equal(t, 5, tm.MaxMembers)
	assert.NotZero(t, tm.CreatedAt)
	assert.NotZero(t, tm.UpdatedAt)
}

func TestNewTeam_EmptyName(t *testing.T) {
	_, err := NewTeam("", "owner-1")
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrValidation)
}

func TestNewTeam_NameTooShort(t *testing.T) {
	_, err := NewTeam("A", "owner-1")
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrValidation)
}

func TestNewTeam_NameTooLong(t *testing.T) {
	longName := ""
	for i := 0; i < 51; i++ {
		longName += "A"
	}
	_, err := NewTeam(longName, "owner-1")
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrValidation)
}

func TestNewTeam_NameExactlyMinLength(t *testing.T) {
	tm, err := NewTeam("AB", "owner-1")
	assert.NoError(t, err)
	assert.Equal(t, "AB", tm.Name)
}

func TestNewTeam_NameExactlyMaxLength(t *testing.T) {
	name := ""
	for i := 0; i < 50; i++ {
		name += "A"
	}
	tm, err := NewTeam(name, "owner-1")
	assert.NoError(t, err)
	assert.Equal(t, name, tm.Name)
}

func TestNewTeam_EmptyOwnerID(t *testing.T) {
	_, err := NewTeam("My Team", "")
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrValidation)
}

// --- Slug generation ---

func TestGenerateSlug_SpacesToDashes(t *testing.T) {
	assert.Equal(t, "my-cool-team", GenerateSlug("My Cool Team"))
}

func TestGenerateSlug_Lowercase(t *testing.T) {
	assert.Equal(t, "uppercase", GenerateSlug("UPPERCASE"))
}

func TestGenerateSlug_SpecialCharsStripped(t *testing.T) {
	assert.Equal(t, "hello-world", GenerateSlug("Hello! @World#"))
}

func TestGenerateSlug_MultipleSpaces(t *testing.T) {
	assert.Equal(t, "a-b", GenerateSlug("a   b"))
}

func TestGenerateSlug_LeadingTrailingSpaces(t *testing.T) {
	assert.Equal(t, "trimmed", GenerateSlug("  trimmed  "))
}

func TestGenerateSlug_AccentsAndSpecials(t *testing.T) {
	slug := GenerateSlug("Café & Co.")
	assert.Equal(t, "caf-co", slug)
}

func TestGenerateSlug_AllSpecialChars(t *testing.T) {
	slug := GenerateSlug("!@#$%^&*()")
	assert.Equal(t, "", slug)
}

// --- NewInvite ---

func TestNewInvite_Success(t *testing.T) {
	inv, err := NewInvite("team-1", "user@example.com", RoleMember)
	assert.NoError(t, err)
	assert.Equal(t, "team-1", inv.TeamID)
	assert.Equal(t, "user@example.com", inv.InvitedEmail)
	assert.Equal(t, RoleMember, inv.Role)
	assert.Equal(t, InvitePending, inv.InviteStatus)
	assert.Len(t, inv.InviteToken, 64) // 32 bytes = 64 hex chars
	assert.NotZero(t, inv.CreatedAt)
}

func TestNewInvite_AdminRole(t *testing.T) {
	inv, err := NewInvite("team-1", "admin@example.com", RoleAdmin)
	assert.NoError(t, err)
	assert.Equal(t, RoleAdmin, inv.Role)
}

func TestNewInvite_EmptyTeamID(t *testing.T) {
	_, err := NewInvite("", "user@example.com", RoleMember)
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrValidation)
}

func TestNewInvite_EmptyEmail(t *testing.T) {
	_, err := NewInvite("team-1", "", RoleMember)
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrValidation)
}

func TestNewInvite_InvalidRole(t *testing.T) {
	_, err := NewInvite("team-1", "user@example.com", Role("superadmin"))
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrValidation)
}

func TestNewInvite_CannotInviteAsOwner(t *testing.T) {
	_, err := NewInvite("team-1", "user@example.com", RoleOwner)
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrValidation)
}

func TestNewInvite_UniqueTokens(t *testing.T) {
	tokens := make(map[string]bool)
	for i := 0; i < 50; i++ {
		inv, err := NewInvite("team-1", "user@example.com", RoleMember)
		assert.NoError(t, err)
		tokens[inv.InviteToken] = true
	}
	assert.Len(t, tokens, 50)
}

// --- Role permissions ---

func TestIsOwner(t *testing.T) {
	tm := &Team{OwnerID: "user-1"}
	assert.True(t, tm.IsOwner("user-1"))
	assert.False(t, tm.IsOwner("user-2"))
}

func TestCanManageMembers(t *testing.T) {
	assert.True(t, CanManageMembers(RoleOwner))
	assert.True(t, CanManageMembers(RoleAdmin))
	assert.False(t, CanManageMembers(RoleMember))
}

func TestCanDeleteTeam(t *testing.T) {
	assert.True(t, CanDeleteTeam(RoleOwner))
	assert.False(t, CanDeleteTeam(RoleAdmin))
	assert.False(t, CanDeleteTeam(RoleMember))
}

// --- ValidateRole ---

func TestValidateRole(t *testing.T) {
	assert.True(t, ValidateRole("owner"))
	assert.True(t, ValidateRole("admin"))
	assert.True(t, ValidateRole("member"))
	assert.False(t, ValidateRole("superadmin"))
	assert.False(t, ValidateRole(""))
}

// --- ValidateInviteStatus ---

func TestValidateInviteStatus(t *testing.T) {
	assert.True(t, ValidateInviteStatus("pending"))
	assert.True(t, ValidateInviteStatus("accepted"))
	assert.True(t, ValidateInviteStatus("declined"))
	assert.False(t, ValidateInviteStatus("expired"))
	assert.False(t, ValidateInviteStatus(""))
}

// --- Role and InviteStatus constants ---

func TestRole_Values(t *testing.T) {
	assert.Equal(t, Role("owner"), RoleOwner)
	assert.Equal(t, Role("admin"), RoleAdmin)
	assert.Equal(t, Role("member"), RoleMember)
}

func TestInviteStatus_Values(t *testing.T) {
	assert.Equal(t, InviteStatus("pending"), InvitePending)
	assert.Equal(t, InviteStatus("accepted"), InviteAccepted)
	assert.Equal(t, InviteStatus("declined"), InviteDeclined)
}
