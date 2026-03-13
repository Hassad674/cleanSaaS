package handler

import (
	"encoding/json"
	"net/http"

	"github.com/hassad/boilerplateSaaS/backend/internal/app/user"
	"github.com/hassad/boilerplateSaaS/backend/internal/handler/dto/request"
	"github.com/hassad/boilerplateSaaS/backend/internal/handler/dto/response"
	"github.com/hassad/boilerplateSaaS/backend/internal/handler/middleware"
)

type UserHandler struct {
	svc *user.Service
}

func NewUserHandler(svc *user.Service) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())

	u, err := h.svc.GetProfile(r.Context(), userID)
	if err != nil {
		response.HandleDomainError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, response.UserFromDomain(u))
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())

	var req request.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	u, err := h.svc.UpdateProfile(r.Context(), userID, req.Name, req.AvatarURL)
	if err != nil {
		response.HandleDomainError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, response.UserFromDomain(u))
}

func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())

	var req request.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.OldPassword == "" || req.NewPassword == "" {
		response.Error(w, http.StatusBadRequest, "old_password and new_password are required")
		return
	}

	if len(req.NewPassword) < 8 {
		response.Error(w, http.StatusBadRequest, "new password must be at least 8 characters")
		return
	}

	if err := h.svc.ChangePassword(r.Context(), userID, req.OldPassword, req.NewPassword); err != nil {
		response.HandleDomainError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "password updated"})
}

func (h *UserHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())

	if err := h.svc.DeleteAccount(r.Context(), userID); err != nil {
		response.HandleDomainError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
