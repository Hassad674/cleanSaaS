package user

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
)

func TestNewEmail(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		want    string
		wantErr bool
	}{
		{"simple", "user@example.com", "user@example.com", false},
		{"uppercase normalized", "User@Example.COM", "user@example.com", false},
		{"surrounding whitespace trimmed", "  user@example.com  ", "user@example.com", false},
		{"mixed case and whitespace", "  Admin@CleanSaaS.Dev ", "admin@cleansaas.dev", false},
		{"empty", "", "", true},
		{"only whitespace", "   ", "", true},
		{"no at sign", "userexample.com", "", true},
		{"no domain", "user@", "", true},
		{"no local part", "@example.com", "", true},
		{"display name form rejected", "Foo <foo@bar.com>", "", true},
		{"spaces inside rejected", "us er@example.com", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEmail(tt.raw)
			if tt.wantErr {
				assert.ErrorIs(t, err, domain.ErrValidation)
				assert.True(t, got.IsZero())
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got.String())
			assert.False(t, got.IsZero())
		})
	}
}

func TestEmail_ZeroValueIsInvalid(t *testing.T) {
	var e Email
	assert.True(t, e.IsZero())
	assert.Equal(t, "", e.String())
}
