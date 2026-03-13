package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	appblog "github.com/hassad/boilerplateSaaS/backend/internal/app/blog"
	"github.com/hassad/boilerplateSaaS/backend/internal/app/user"
	"github.com/hassad/boilerplateSaaS/backend/internal/handler/dto/response"
)

type AdminHandler struct {
	userSvc *user.Service
	blogSvc *appblog.Service
}

func NewAdminHandler(userSvc *user.Service, blogSvc *appblog.Service) *AdminHandler {
	return &AdminHandler{userSvc: userSvc, blogSvc: blogSvc}
}

func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	users, total, err := h.userSvc.ListUsers(r.Context(), search, offset, limit)
	if err != nil {
		response.HandleDomainError(w, err)
		return
	}

	items := make([]response.UserResponse, len(users))
	for i, u := range users {
		items[i] = response.UserFromDomain(u)
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"users": items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *AdminHandler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")

	var req struct {
		Role string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Role != "member" && req.Role != "admin" {
		response.Error(w, http.StatusBadRequest, "role must be 'member' or 'admin'")
		return
	}

	u, err := h.userSvc.UpdateRole(r.Context(), userID, req.Role)
	if err != nil {
		response.HandleDomainError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, response.UserFromDomain(u))
}

func (h *AdminHandler) DashboardStats(w http.ResponseWriter, r *http.Request) {
	userCount, err := h.userSvc.CountUsers(r.Context())
	if err != nil {
		response.HandleDomainError(w, err)
		return
	}

	posts, postCount, err := h.blogSvc.ListAll(r.Context(), "", "", 0, 0)
	_ = posts

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"total_users": userCount,
		"total_posts": postCount,
	})
}
