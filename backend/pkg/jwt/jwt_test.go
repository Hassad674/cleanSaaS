package jwt

import (
	"testing"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMaker_Generate_ValidToken(t *testing.T) {
	maker := NewMaker("test-secret-key")
	token, err := maker.Generate("user-123", "member")

	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestMaker_Validate_ValidToken(t *testing.T) {
	maker := NewMaker("test-secret-key")
	token, err := maker.Generate("user-123", "admin")
	require.NoError(t, err)

	claims, err := maker.Validate(token)
	require.NoError(t, err)
	assert.Equal(t, "user-123", claims.UserID)
	assert.Equal(t, "admin", claims.Role)
}

func TestMaker_Validate_ExpiredToken(t *testing.T) {
	secret := "test-secret-key"
	maker := NewMaker(secret)

	// Create a token that is already expired
	claims := jwtlib.MapClaims{
		"sub":  "user-123",
		"role": "member",
		"exp":  time.Now().Add(-time.Hour).Unix(),
		"iat":  time.Now().Add(-2 * time.Hour).Unix(),
	}
	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(secret))
	require.NoError(t, err)

	_, err = maker.Validate(tokenStr)
	assert.Error(t, err, "expired token should fail validation")
}

func TestMaker_Validate_InvalidToken(t *testing.T) {
	maker := NewMaker("test-secret-key")

	_, err := maker.Validate("not-a-valid-jwt")
	assert.Error(t, err)
}

func TestMaker_Validate_WrongSecret(t *testing.T) {
	maker1 := NewMaker("secret-one")
	maker2 := NewMaker("secret-two")

	token, err := maker1.Generate("user-123", "member")
	require.NoError(t, err)

	_, err = maker2.Validate(token)
	assert.Error(t, err, "token signed with different secret should fail validation")
}

func TestMaker_Validate_EmptyToken(t *testing.T) {
	maker := NewMaker("test-secret-key")

	_, err := maker.Validate("")
	assert.Error(t, err)
}

func TestMaker_Generate_DifferentUserIDs(t *testing.T) {
	maker := NewMaker("secret")

	token1, err := maker.Generate("user-1", "member")
	require.NoError(t, err)

	token2, err := maker.Generate("user-2", "admin")
	require.NoError(t, err)

	assert.NotEqual(t, token1, token2, "tokens for different users should differ")

	claims1, err := maker.Validate(token1)
	require.NoError(t, err)
	assert.Equal(t, "user-1", claims1.UserID)
	assert.Equal(t, "member", claims1.Role)

	claims2, err := maker.Validate(token2)
	require.NoError(t, err)
	assert.Equal(t, "user-2", claims2.UserID)
	assert.Equal(t, "admin", claims2.Role)
}

func TestMaker_Validate_TamperedToken(t *testing.T) {
	maker := NewMaker("secret")
	token, err := maker.Generate("user-1", "member")
	require.NoError(t, err)

	// Tamper with the token by changing last character
	tampered := token[:len(token)-1] + "X"

	_, err = maker.Validate(tampered)
	assert.Error(t, err, "tampered token should fail validation")
}
