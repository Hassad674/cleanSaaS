package jwt

import (
	"strings"
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

	// Tamper with the first character of the signature segment. The signature
	// is the third dot-delimited segment; flipping a non-final base64url char
	// guarantees the decoded signature bytes change (the final char only
	// carries a few significant bits, so swapping it can be a no-op).
	parts := strings.Split(token, ".")
	require.Len(t, parts, 3)
	sig := parts[2]
	flipped := "A"
	if sig[0] == 'A' {
		flipped = "B"
	}
	parts[2] = flipped + sig[1:]
	tampered := strings.Join(parts, ".")
	require.NotEqual(t, token, tampered)

	_, err = maker.Validate(tampered)
	assert.Error(t, err, "tampered token should fail validation")
}

func TestMaker_Generate_CarriesIssuerAudienceAndTTL(t *testing.T) {
	secret := "a-sufficiently-long-jwt-secret-value-1234"
	ttl := 15 * time.Minute
	maker := NewMakerWithOptions(secret, ttl, "myiss", "myaud")

	tokenStr, err := maker.Generate("user-1", "member")
	require.NoError(t, err)

	parsed, err := jwtlib.Parse(tokenStr, func(_ *jwtlib.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	require.NoError(t, err)
	claims, ok := parsed.Claims.(jwtlib.MapClaims)
	require.True(t, ok)

	assert.Equal(t, "myiss", claims["iss"])
	assert.Equal(t, "myaud", claims["aud"])

	exp, ok := claims["exp"].(float64)
	require.True(t, ok)
	iat, ok := claims["iat"].(float64)
	require.True(t, ok)
	// exp - iat should equal the configured TTL (within a second of jitter).
	assert.InDelta(t, ttl.Seconds(), exp-iat, 2)
}

func TestNewMakerWithOptions_FallsBackToDefaults(t *testing.T) {
	maker := NewMakerWithOptions("secret", 0, "", "")
	assert.Equal(t, defaultAccessTTL, maker.AccessTTL())
}

func TestMaker_AccessTTL_IsShortByDefault(t *testing.T) {
	maker := NewMaker("secret")
	assert.Equal(t, 15*time.Minute, maker.AccessTTL(), "default access TTL should be short-lived")
}

func TestMaker_Validate_RejectsAlgNone(t *testing.T) {
	secret := "a-sufficiently-long-jwt-secret-value-1234"
	maker := NewMaker(secret)

	// Forge an unsigned token (alg=none). The alg-confusion guard must reject it.
	claims := jwtlib.MapClaims{
		"sub":  "user-1",
		"role": "admin",
		"exp":  time.Now().Add(time.Hour).Unix(),
	}
	token := jwtlib.NewWithClaims(jwtlib.SigningMethodNone, claims)
	tokenStr, err := token.SignedString(jwtlib.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	_, err = maker.Validate(tokenStr)
	assert.Error(t, err, "alg=none token must be rejected")
}

func TestMaker_Validate_RejectsNonHMACAlg(t *testing.T) {
	secret := "a-sufficiently-long-jwt-secret-value-1234"
	maker := NewMaker(secret)

	// A token whose header claims a non-HMAC algorithm must be rejected before
	// the secret is ever consulted (defeats RS256->HS256 key-confusion attacks).
	claims := jwtlib.MapClaims{
		"sub":  "user-1",
		"role": "admin",
		"alg":  "RS256",
	}
	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	token.Header["alg"] = "RS256"
	tokenStr, err := token.SignedString([]byte(secret))
	require.NoError(t, err)

	_, err = maker.Validate(tokenStr)
	assert.Error(t, err, "non-HMAC alg header must be rejected")
}

func TestGenerateRefreshToken_RandomAndOpaque(t *testing.T) {
	a, err := GenerateRefreshToken()
	require.NoError(t, err)
	b, err := GenerateRefreshToken()
	require.NoError(t, err)

	assert.NotEmpty(t, a)
	assert.NotEqual(t, a, b, "refresh tokens must be unique")
	// Opaque, not a JWT (no dot-delimited segments to parse).
	assert.NotContains(t, a, ".")
}

func TestHashRefreshToken_StableAndOneWay(t *testing.T) {
	token, err := GenerateRefreshToken()
	require.NoError(t, err)

	h1 := HashRefreshToken(token)
	h2 := HashRefreshToken(token)
	assert.Equal(t, h1, h2, "hashing must be deterministic")
	assert.NotEqual(t, token, h1, "hash must differ from the raw token")
	assert.Len(t, h1, 64, "sha-256 hex digest is 64 chars")
	assert.NotEqual(t, h1, HashRefreshToken(token+"x"), "different inputs hash differently")
}
