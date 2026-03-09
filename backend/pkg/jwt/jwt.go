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

func NewMaker(secret string) *Maker {
	return &Maker{secret: []byte(secret)}
}

func (m *Maker) Generate(userID, role string) (string, error) {
	claims := jwtlib.MapClaims{
		"sub":  userID,
		"role": role,
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
		"iat":  time.Now().Unix(),
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

	return &Claims{
		UserID: claims["sub"].(string),
		Role:   claims["role"].(string),
	}, nil
}
