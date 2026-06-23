package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/hassad/boilerplateSaaS/backend/internal/handler/dto/response"
	"github.com/hassad/boilerplateSaaS/backend/pkg/jwt"
	"github.com/hassad/boilerplateSaaS/backend/pkg/orgctx"
)

type contextKey string

const userIDKey contextKey = "userID"
const userRoleKey contextKey = "userRole"

// OrgResolver resolves and authorizes a caller's active organization. Given the
// authenticated userID and the org claim from the token (which may be empty), it
// returns the org id to use for the request. Implementations must verify the user
// is a member of any explicitly-requested org and fall back to the user's default
// organization otherwise. Returning ("", nil) means "no active org" (the tenant
// request path will then fail closed at the database).
//
// It is injected (not imported) so the middleware stays in the handler layer
// without depending on a repository adapter — the composition root wires it.
type OrgResolver func(ctx context.Context, userID, tokenOrgID string) (string, error)

// Auth authenticates the bearer token and stores userID + role in the context.
// It does not resolve a tenant org; use AuthWithOrg for tenant-scoped routes.
func Auth(secret string) func(http.Handler) http.Handler {
	return AuthWithOrg(secret, nil)
}

// AuthWithOrg authenticates the bearer token, stores userID + role, and resolves
// the caller's active organization into the context (via pkg/orgctx) so the
// org-scoped database path can enforce RLS. When resolve is nil it behaves like
// Auth. A resolver error is treated as unauthorized — we never serve a tenant
// request without a known, authorized org.
func AuthWithOrg(secret string, resolve OrgResolver) func(http.Handler) http.Handler {
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

			activeOrg := claims.OrgID
			if resolve != nil {
				orgID, err := resolve(ctx, claims.UserID, claims.OrgID)
				if err != nil {
					response.Error(w, http.StatusUnauthorized, "invalid organization")
					return
				}
				activeOrg = orgID
				if orgID != "" {
					ctx = orgctx.WithOrgID(ctx, orgID)
				}
			} else if claims.OrgID != "" {
				ctx = orgctx.WithOrgID(ctx, claims.OrgID)
			}

			// Attribute the access log to the caller (read by StructuredLogging
			// when the request completes). No-op if logging is not in the chain.
			setLogUser(ctx, claims.UserID, activeOrg)

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
