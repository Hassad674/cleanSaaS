package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPassword_GeneratesHash(t *testing.T) {
	hashed, err := Password("mypassword")
	require.NoError(t, err)
	assert.NotEmpty(t, hashed)
	assert.NotEqual(t, "mypassword", hashed, "hash should differ from plaintext")
}

func TestPassword_DifferentHashesForSamePassword(t *testing.T) {
	hash1, err := Password("samepassword")
	require.NoError(t, err)

	hash2, err := Password("samepassword")
	require.NoError(t, err)

	assert.NotEqual(t, hash1, hash2, "bcrypt should produce different hashes due to unique salts")
}

func TestCheck_CorrectPassword(t *testing.T) {
	hashed, err := Password("correctpassword")
	require.NoError(t, err)

	assert.True(t, Check("correctpassword", hashed))
}

func TestCheck_WrongPassword(t *testing.T) {
	hashed, err := Password("correctpassword")
	require.NoError(t, err)

	assert.False(t, Check("wrongpassword", hashed))
}

func TestCheck_EmptyPassword(t *testing.T) {
	hashed, err := Password("password123")
	require.NoError(t, err)

	assert.False(t, Check("", hashed))
}

func TestCheck_EmptyHash(t *testing.T) {
	assert.False(t, Check("password", ""))
}

func TestCheck_InvalidHash(t *testing.T) {
	assert.False(t, Check("password", "not-a-bcrypt-hash"))
}

func TestPassword_EmptyPassword(t *testing.T) {
	// bcrypt can hash empty strings
	hashed, err := Password("")
	require.NoError(t, err)
	assert.NotEmpty(t, hashed)
	assert.True(t, Check("", hashed))
}

func TestPassword_LongPassword(t *testing.T) {
	// bcrypt rejects passwords over 72 bytes
	longPw := ""
	for len(longPw) < 100 {
		longPw += "a"
	}
	_, err := Password(longPw)
	assert.Error(t, err, "bcrypt should reject passwords exceeding 72 bytes")
}

func TestPassword_Max72Bytes(t *testing.T) {
	// bcrypt accepts passwords up to 72 bytes
	pw72 := ""
	for len(pw72) < 72 {
		pw72 += "a"
	}
	hashed, err := Password(pw72)
	require.NoError(t, err)
	assert.True(t, Check(pw72, hashed))
}
