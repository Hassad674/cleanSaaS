package middleware

import (
	"net/http"

	"github.com/hassad/boilerplateSaaS/backend/internal/handler/dto/response"
)

func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role := UserRoleFromContext(r.Context())
		if role != "admin" {
			response.Error(w, http.StatusForbidden, "admin access required")
			return
		}
		next.ServeHTTP(w, r)
	})
}
