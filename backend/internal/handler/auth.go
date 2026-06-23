package handler

import (
	"encoding/json"
	"net/http"

	"github.com/hassad/boilerplateSaaS/backend/internal/app/auth"
	"github.com/hassad/boilerplateSaaS/backend/internal/handler/dto/request"
	"github.com/hassad/boilerplateSaaS/backend/internal/handler/dto/response"
	"github.com/hassad/boilerplateSaaS/backend/internal/handler/middleware"
	"github.com/hassad/boilerplateSaaS/backend/pkg/validate"
)

type AuthHandler struct {
	svc *auth.Service
}

func NewAuthHandler(svc *auth.Service) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req request.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" || req.Name == "" {
		response.Error(w, http.StatusBadRequest, "email, name and password are required")
		return
	}

	if !validate.Email(req.Email) {
		response.Error(w, http.StatusBadRequest, "invalid email format")
		return
	}

	if len(req.Password) < 8 {
		response.Error(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}

	if len(req.Name) > 200 {
		response.Error(w, http.StatusBadRequest, "name is too long")
		return
	}

	u, access, refresh, err := h.svc.Register(r.Context(), req.Email, req.Name, req.Password)
	if err != nil {
		response.HandleDomainError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, response.AuthResponse{
		Token:        access,
		AccessToken:  access,
		RefreshToken: refresh,
		User:         response.UserFromDomain(u),
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req request.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	u, access, refresh, err := h.svc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		response.HandleDomainError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, response.AuthResponse{
		Token:        access,
		AccessToken:  access,
		RefreshToken: refresh,
		User:         response.UserFromDomain(u),
	})
}

// Refresh exchanges a valid refresh token for a new access/refresh token pair
// (the presented refresh token is rotated out and revoked).
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req request.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.RefreshToken == "" {
		response.Error(w, http.StatusBadRequest, "refresh_token is required")
		return
	}

	access, refresh, u, err := h.svc.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		response.HandleDomainError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, response.AuthResponse{
		Token:        access,
		AccessToken:  access,
		RefreshToken: refresh,
		User:         response.UserFromDomain(u),
	})
}

// Logout revokes the presented refresh token. It is idempotent — an unknown or
// already-revoked token still returns success.
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req request.LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.RefreshToken == "" {
		response.Error(w, http.StatusBadRequest, "refresh_token is required")
		return
	}

	if err := h.svc.Logout(r.Context(), req.RefreshToken); err != nil {
		response.HandleDomainError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"message": "logged out successfully",
	})
}

func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req request.ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" {
		response.Error(w, http.StatusBadRequest, "email is required")
		return
	}

	// Always return success to avoid leaking user existence
	if err := h.svc.ForgotPassword(r.Context(), req.Email, ""); err != nil {
		response.HandleDomainError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"message": "if an account with that email exists, a reset link has been sent",
	})
}

func (h *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	var req request.VerifyEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Token == "" {
		response.Error(w, http.StatusBadRequest, "token is required")
		return
	}

	if err := h.svc.VerifyEmail(r.Context(), req.Token); err != nil {
		response.HandleDomainError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"message": "email verified successfully",
	})
}

func (h *AuthHandler) ResendVerification(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	if userID == "" {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := h.svc.ResendVerification(r.Context(), userID); err != nil {
		response.HandleDomainError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"message": "verification email sent",
	})
}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req request.ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Token == "" || req.Password == "" {
		response.Error(w, http.StatusBadRequest, "token and password are required")
		return
	}

	if len(req.Password) < 8 {
		response.Error(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}

	if err := h.svc.ResetPassword(r.Context(), req.Token, req.Password); err != nil {
		response.HandleDomainError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"message": "password has been reset successfully",
	})
}
