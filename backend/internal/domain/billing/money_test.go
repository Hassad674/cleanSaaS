package billing

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
)

func TestNewMoney(t *testing.T) {
	tests := []struct {
		name         string
		cents        int64
		currency     string
		wantErr      bool
		wantCents    int64
		wantCurrency string
	}{
		{"usd amount", 1900, "usd", false, 1900, "usd"},
		{"uppercase currency normalized", 1900, "USD", false, 1900, "usd"},
		{"currency whitespace trimmed", 1900, " eur ", false, 1900, "eur"},
		{"empty currency defaults to usd", 4900, "", false, 4900, "usd"},
		{"zero amount allowed (free plan)", 0, "usd", false, 0, "usd"},
		{"negative amount rejected", -1, "usd", true, 0, ""},
		{"too-short currency rejected", 100, "us", true, 0, ""},
		{"too-long currency rejected", 100, "usdd", true, 0, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := NewMoney(tt.cents, tt.currency)
			if tt.wantErr {
				assert.ErrorIs(t, err, domain.ErrValidation)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantCents, m.Cents())
			assert.Equal(t, tt.wantCurrency, m.Currency())
		})
	}
}

func TestMoney_IsZero(t *testing.T) {
	zero, err := NewMoney(0, "usd")
	assert.NoError(t, err)
	assert.True(t, zero.IsZero())

	nonZero, err := NewMoney(1, "usd")
	assert.NoError(t, err)
	assert.False(t, nonZero.IsZero())
}

func TestMoney_String(t *testing.T) {
	tests := []struct {
		cents    int64
		currency string
		want     string
	}{
		{1900, "usd", "19.00 USD"},
		{49900, "usd", "499.00 USD"},
		{5, "eur", "0.05 EUR"},
		{0, "usd", "0.00 USD"},
		{100, "gbp", "1.00 GBP"},
	}
	for _, tt := range tests {
		m, err := NewMoney(tt.cents, tt.currency)
		assert.NoError(t, err)
		assert.Equal(t, tt.want, m.String())
	}
}
