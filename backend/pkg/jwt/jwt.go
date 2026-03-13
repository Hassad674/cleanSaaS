package jwt

import (
	"fmt"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string
	Role   string
}

type Maker struct {
	secret []byte
}

// NewMaker creates a JWT maker. The secret must be at least 32 characters
// for HS256 security. In production, always use a strong random secret.
func NewMaker(secret string) *Maker {
	if len(secret) < 32 {
		// Allow short secrets in dev but log a warning
		fmt.Println("WARNING: JWT_SECRET is shorter than 32 characters — insecure for production")
	}
	return &Maker{secret: []byte(secret)}
}

func (m *Maker) Generate(userID, role string) (string, error) {
	now := time.Now()
	claims := jwtlib.MapClaims{
		"sub":  userID,
		"role": role,
		"exp":  now.Add(24 * time.Hour).Unix(),
		"iat":  now.Unix(),
		"nbf":  now.Unix(), // Not valid before issuance time
	}

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *Maker) Validate(tokenStr string) (*Claims, error) {
	token, err := jwtlib.Parse(tokenStr, func(token *jwtlib.Token) (interface{}, error) {
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

	return &Claims{
		UserID: sub,
		Role:   role,
	}, nil
}
