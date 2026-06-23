// Package org contains the organization tenant aggregate. An organization is the
// unit of multi-tenancy: every tenant-scoped resource (subscriptions, files,
// conversations, notifications) belongs to exactly one organization, and a member
// of one organization can never see another organization's rows.
package org

import (
	"regexp"
	"strings"
	"time"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
)

// Role is a member's role within an organization.
type Role string

const (
	// RoleOwner is the organization's owner (the user it was created for).
	RoleOwner Role = "owner"
	// RoleAdmin can manage the organization and its members.
	RoleAdmin Role = "admin"
	// RoleMember is a regular member.
	RoleMember Role = "member"
)

// slugPattern allows lowercase letters, digits and single hyphens — safe for URLs.
var slugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

// Organization is the tenant aggregate root. Fields stay exported because the
// postgres repository scans rows straight into them and response DTOs read them;
// that is the sanctioned repository/serialization boundary. Validation lives in
// the constructor so an invalid organization can never be persisted.
type Organization struct {
	ID        string
	Name      string
	Slug      string
	OwnerID   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// New creates a validated organization owned by ownerID. The name is required and
// the slug is derived from the name when not supplied explicitly.
func New(name, slug, ownerID string) (*Organization, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, domain.ErrValidation
	}
	if ownerID == "" {
		return nil, domain.ErrValidation
	}

	if slug == "" {
		slug = Slugify(name)
	}
	if !slugPattern.MatchString(slug) {
		return nil, domain.ErrValidation
	}

	now := time.Now()
	return &Organization{
		Name:      name,
		Slug:      slug,
		OwnerID:   ownerID,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Slugify converts an arbitrary name into a URL-safe slug. It lowercases, replaces
// any run of non-alphanumeric characters with a single hyphen, and trims hyphens
// from the ends. A name with no usable characters yields "org".
func Slugify(name string) string {
	var b strings.Builder
	lastHyphen := true // suppress a leading hyphen
	for _, r := range strings.ToLower(name) {
		switch {
		case (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9'):
			b.WriteRune(r)
			lastHyphen = false
		default:
			if !lastHyphen {
				b.WriteByte('-')
				lastHyphen = true
			}
		}
	}
	slug := strings.Trim(b.String(), "-")
	if slug == "" {
		return "org"
	}
	return slug
}

// Member is a user's membership in an organization, with their role.
type Member struct {
	OrgID     string
	UserID    string
	Role      Role
	CreatedAt time.Time
}

// NewMember creates a validated membership.
func NewMember(orgID, userID string, role Role) (*Member, error) {
	if orgID == "" || userID == "" {
		return nil, domain.ErrValidation
	}
	if role == "" {
		role = RoleMember
	}
	return &Member{
		OrgID:     orgID,
		UserID:    userID,
		Role:      role,
		CreatedAt: time.Now(),
	}, nil
}
