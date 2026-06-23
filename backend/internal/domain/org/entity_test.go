package org

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
)

func TestNew_Success(t *testing.T) {
	o, err := New("Acme Inc", "", "user-1")
	assert.NoError(t, err)
	assert.Equal(t, "Acme Inc", o.Name)
	assert.Equal(t, "acme-inc", o.Slug)
	assert.Equal(t, "user-1", o.OwnerID)
	assert.False(t, o.CreatedAt.IsZero())
}

func TestNew_ExplicitSlug(t *testing.T) {
	o, err := New("Acme Inc", "acme", "user-1")
	assert.NoError(t, err)
	assert.Equal(t, "acme", o.Slug)
}

func TestNew_EmptyName(t *testing.T) {
	_, err := New("   ", "", "user-1")
	assert.ErrorIs(t, err, domain.ErrValidation)
}

func TestNew_EmptyOwner(t *testing.T) {
	_, err := New("Acme", "", "")
	assert.ErrorIs(t, err, domain.ErrValidation)
}

func TestNew_InvalidSlug(t *testing.T) {
	_, err := New("Acme", "Invalid Slug!", "user-1")
	assert.ErrorIs(t, err, domain.ErrValidation)
}

func TestSlugify(t *testing.T) {
	cases := map[string]string{
		"Acme Inc":         "acme-inc",
		"  Hello  World  ": "hello-world",
		"Foo___Bar":        "foo-bar",
		"UPPER":            "upper",
		"!!!":              "org",
		"a-b-c":            "a-b-c",
		"Café & Co":        "caf-co",
		"123 Numbers":      "123-numbers",
	}
	for in, want := range cases {
		assert.Equal(t, want, Slugify(in), "Slugify(%q)", in)
	}
}

func TestNewMember_Success(t *testing.T) {
	m, err := NewMember("org-1", "user-1", RoleOwner)
	assert.NoError(t, err)
	assert.Equal(t, RoleOwner, m.Role)
}

func TestNewMember_DefaultsToMember(t *testing.T) {
	m, err := NewMember("org-1", "user-1", "")
	assert.NoError(t, err)
	assert.Equal(t, RoleMember, m.Role)
}

func TestNewMember_Validation(t *testing.T) {
	_, err := NewMember("", "user-1", RoleMember)
	assert.ErrorIs(t, err, domain.ErrValidation)
	_, err = NewMember("org-1", "", RoleMember)
	assert.ErrorIs(t, err, domain.ErrValidation)
}
