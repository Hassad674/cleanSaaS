package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/hassad/boilerplateSaaS/backend/internal/handler/dto/response"
	"github.com/hassad/boilerplateSaaS/backend/pkg/jwt"
)

type contextKey string

const userIDKey contextKey = "userID"
const userRoleKey contextKey = "userRole"

func Auth(secret string) func(http.Handler) http.Handler {
	maker := jwt.NewMaker(secret)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				response.Error(w, http.StatusUnauthorized, "missing authorization header")
				return
			}

			token := strings.TrimPrefix(header, "Bearer ")
			if token == header {
				response.Error(w, http.StatusUnauthorized, "invalid authorization format")
				return
			}

			claims, err := maker.Validate(token)
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "invalid token")
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
			ctx = context.WithValue(ctx, userRoleKey, claims.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIDFromContext(ctx context.Context) string {
	v, _ := ctx.Value(userIDKey).(string)
	return v
}

func UserRoleFromContext(ctx context.Context) string {
	v, _ := ctx.Value(userRoleKey).(string)
	return v
}
