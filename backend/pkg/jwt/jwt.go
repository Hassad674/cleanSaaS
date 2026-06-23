package jwt

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

// Default access-token configuration. These keep the zero-argument NewMaker
// constructor secure-by-default while NewMakerWithOptions allows overrides.
const (
	defaultAccessTTL = 15 * time.Minute
	defaultIssuer    = "cleansaas"
	defaultAudience  = "cleansaas"
)

// refreshTokenBytes is the entropy of an opaque refresh token before encoding
// (32 bytes = 256 bits). The server never stores the token itself, only its
// SHA-256 hash, so it cannot be recovered from the database.
const refreshTokenBytes = 32

type Claims struct {
	UserID string
	Role   string
	// OrgID is the caller's active organization (the tenant). It is optional in the
	// token: an older token without it, or a non-tenant token (e.g. the WebSocket
	// upgrade), simply yields an empty OrgID, and the request path then resolves the
	// user's default organization. Carrying it lets the common case skip a lookup.
	OrgID string
}

type Maker struct {
	secret    []byte
	accessTTL time.Duration
	issuer    string
	audience  string
}

// NewMaker creates a JWT maker with secure defaults (15m access TTL,
// iss/aud = "cleansaas"). The secret must be at least 32 characters for HS256
// security. In production, always use a strong random secret.
func NewMaker(secret string) *Maker {
	return NewMakerWithOptions(secret, defaultAccessTTL, defaultIssuer, defaultAudience)
}

// NewMakerWithOptions creates a JWT maker with a configurable access-token TTL,
// issuer and audience. Empty/zero values fall back to the secure defaults so
// callers can override only what they need.
func NewMakerWithOptions(secret string, accessTTL time.Duration, issuer, audience string) *Maker {
	if len(secret) < 32 {
		// Allow short secrets in dev but log a warning
		fmt.Println("WARNING: JWT_SECRET is shorter than 32 characters — insecure for production")
	}
	if accessTTL <= 0 {
		accessTTL = defaultAccessTTL
	}
	if issuer == "" {
		issuer = defaultIssuer
	}
	if audience == "" {
		audience = defaultAudience
	}
	return &Maker{
		secret:    []byte(secret),
		accessTTL: accessTTL,
		issuer:    issuer,
		audience:  audience,
	}
}

// AccessTTL returns the configured access-token lifetime.
func (m *Maker) AccessTTL() time.Duration {
	return m.accessTTL
}

// Generate issues a short-lived HS256 access token carrying sub + role + iss +
// aud + standard time claims. Used where no tenant context applies (e.g. the
// WebSocket upgrade token).
func (m *Maker) Generate(userID, role string) (string, error) {
	return m.GenerateWithOrg(userID, role, "")
}

// GenerateWithOrg issues an access token that additionally carries the caller's
// active organization id as the "org" claim. The request path reads it to scope
// every tenant query. An empty orgID omits the claim, falling back to Generate's
// behavior.
func (m *Maker) GenerateWithOrg(userID, role, orgID string) (string, error) {
	now := time.Now()
	claims := jwtlib.MapClaims{
		"sub":  userID,
		"role": role,
		"iss":  m.issuer,
		"aud":  m.audience,
		"exp":  now.Add(m.accessTTL).Unix(),
		"iat":  now.Unix(),
		"nbf":  now.Unix(), // Not valid before issuance time
	}
	if orgID != "" {
		claims["org"] = orgID
	}

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *Maker) Validate(tokenStr string) (*Claims, error) {
	token, err := jwtlib.Parse(tokenStr, func(token *jwtlib.Token) (interface{}, error) {
		// Pin the algorithm to HMAC to defeat alg-confusion attacks (e.g. a
		// token forged with alg=none or an RS256->HS256 key swap).
		if _, ok := token.Method.(*jwtlib.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwtlib.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	sub, ok := claims["sub"].(string)
	if !ok || sub == "" {
		return nil, fmt.Errorf("invalid token: missing subject")
	}

	role, ok := claims["role"].(string)
	if !ok || role == "" {
		return nil, fmt.Errorf("invalid token: missing role")
	}

	// org is optional — absent in legacy/non-tenant tokens.
	orgID, _ := claims["org"].(string)

	return &Claims{
		UserID: sub,
		Role:   role,
		OrgID:  orgID,
	}, nil
}

// GenerateRefreshToken returns a high-entropy, URL-safe opaque refresh token.
// It is NOT a JWT: it carries no claims and is meaningless without the matching
// database row, which makes per-token revocation trivial. The caller persists
// only HashRefreshToken(token); the raw value is returned to the client once.
func GenerateRefreshToken() (string, error) {
	b := make([]byte, refreshTokenBytes)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generating refresh token: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// HashRefreshToken returns the lowercase hex SHA-256 of an opaque refresh token.
// The server stores and looks up tokens by this hash so a database leak never
// exposes usable refresh tokens.
func HashRefreshToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
