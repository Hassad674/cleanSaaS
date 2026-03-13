package team

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
)

// Role represents a team member's role.
type Role string

const (
	RoleOwner  Role = "owner"
	RoleAdmin  Role = "admin"
	RoleMember Role = "member"
)

// ValidateRole checks whether the given string is a valid team role.
func ValidateRole(s string) bool {
	switch Role(s) {
	case RoleOwner, RoleAdmin, RoleMember:
		return true
	}
	return false
}

// InviteStatus represents the status of a team invitation.
type InviteStatus string

const (
	InvitePending  InviteStatus = "pending"
	InviteAccepted InviteStatus = "accepted"
	InviteDeclined InviteStatus = "declined"
)

// ValidateInviteStatus checks whether the given string is a valid invite status.
func ValidateInviteStatus(s string) bool {
	switch InviteStatus(s) {
	case InvitePending, InviteAccepted, InviteDeclined:
		return true
	}
	return false
}

const (
	nameMinLen = 2
	nameMaxLen = 50
)

// slugRegex matches only lowercase alphanumeric characters and hyphens.
var slugRegex = regexp.MustCompile(`[^a-z0-9-]`)

// Team represents an organization or team.
type Team struct {
	ID         string
	Name       string
	Slug       string
	OwnerID    string
	AvatarURL  string
	Plan       string
	MaxMembers int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// TeamMember represents a user's membership in a team.
type TeamMember struct {
	ID           string
	TeamID       string
	UserID       string
	Role         Role
	InvitedEmail string
	InviteToken  string
	InviteStatus InviteStatus
	JoinedAt     *time.Time
	CreatedAt    time.Time
}

// NewTeam creates a new Team with validated name and a generated slug.
func NewTeam(name, ownerID string) (*Team, error) {
	if err := validateName(name); err != nil {
		return nil, err
	}
	if ownerID == "" {
		return nil, fmt.Errorf("owner ID is required: %w", domain.ErrValidation)
	}

	now := time.Now()
	return &Team{
		Name:       name,
		Slug:       GenerateSlug(name),
		OwnerID:    ownerID,
		Plan:       "free",
		MaxMembers: 5,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

// NewInvite creates a pending team member invitation with a crypto-random token.
func NewInvite(teamID, email string, role Role) (*TeamMember, error) {
	if teamID == "" {
		return nil, fmt.Errorf("team ID is required: %w", domain.ErrValidation)
	}
	if email == "" {
		return nil, fmt.Errorf("invited email is required: %w", domain.ErrValidation)
	}
	if !ValidateRole(string(role)) {
		return nil, fmt.Errorf("invalid role %q: %w", role, domain.ErrValidation)
	}
	if role == RoleOwner {
		return nil, fmt.Errorf("cannot invite as owner: %w", domain.ErrValidation)
	}

	token, err := generateInviteToken()
	if err != nil {
		return nil, fmt.Errorf("generating invite token: %w", err)
	}

	return &TeamMember{
		TeamID:       teamID,
		Role:         role,
		InvitedEmail: email,
		InviteToken:  token,
		InviteStatus: InvitePending,
		CreatedAt:    time.Now(),
	}, nil
}

// IsOwner checks whether the given user ID matches the team's owner.
func (t *Team) IsOwner(userID string) bool {
	return t.OwnerID == userID
}

// CanManageMembers returns true if the role is allowed to manage members (owner or admin).
func CanManageMembers(role Role) bool {
	return role == RoleOwner || role == RoleAdmin
}

// CanDeleteTeam returns true if the role is allowed to delete the team (owner only).
func CanDeleteTeam(role Role) bool {
	return role == RoleOwner
}

// GenerateSlug converts a team name into a URL-safe slug.
// Spaces become dashes, everything is lowercased, and special characters are stripped.
func GenerateSlug(name string) string {
	slug := strings.ToLower(strings.TrimSpace(name))
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = slugRegex.ReplaceAllString(slug, "")
	// Collapse multiple consecutive hyphens
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}
	slug = strings.Trim(slug, "-")
	return slug
}

// validateName checks the team name length constraints.
func validateName(name string) error {
	trimmed := strings.TrimSpace(name)
	if len(trimmed) < nameMinLen {
		return fmt.Errorf("team name must be at least %d characters: %w", nameMinLen, domain.ErrValidation)
	}
	if len(trimmed) > nameMaxLen {
		return fmt.Errorf("team name must be at most %d characters: %w", nameMaxLen, domain.ErrValidation)
	}
	return nil
}

// generateInviteToken creates a cryptographically random 32-byte hex token.
func generateInviteToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("reading random bytes: %w", err)
	}
	return hex.EncodeToString(b), nil
}
