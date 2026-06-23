package response

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
	"github.com/hassad/boilerplateSaaS/backend/internal/domain/user"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

// AuthResponse is returned by login/register/refresh. `token` is kept as a
// backward-compatible alias for `access_token` (both carry the same value) so
// existing frontend clients keep working; new clients should read
// `access_token` + `refresh_token`.
type AuthResponse struct {
	Token        string       `json:"token"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

type UserResponse struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	AvatarURL     string `json:"avatar_url"`
	Role          string `json:"role"`
	EmailVerified bool   `json:"email_verified"`
}

func UserFromDomain(u *user.User) UserResponse {
	return UserResponse{
		ID:            u.ID,
		Email:         u.Email,
		Name:          u.Name,
		AvatarURL:     u.AvatarURL,
		Role:          string(u.Role),
		EmailVerified: u.EmailVerified,
	}
}

func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, ErrorResponse{Error: message})
}

func HandleDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		Error(w, http.StatusNotFound, "not found")
	case errors.Is(err, domain.ErrAlreadyExists):
		Error(w, http.StatusConflict, "already exists")
	case errors.Is(err, domain.ErrUnauthorized):
		Error(w, http.StatusUnauthorized, "unauthorized")
	case errors.Is(err, domain.ErrForbidden):
		Error(w, http.StatusForbidden, "forbidden")
	case errors.Is(err, domain.ErrValidation):
		Error(w, http.StatusBadRequest, "validation error")
	case errors.Is(err, domain.ErrExpiredToken):
		Error(w, http.StatusBadRequest, "token expired")
	case errors.Is(err, domain.ErrInvalidToken):
		Error(w, http.StatusBadRequest, "invalid token")
	default:
		Error(w, http.StatusInternalServerError, "internal error")
	}
}
