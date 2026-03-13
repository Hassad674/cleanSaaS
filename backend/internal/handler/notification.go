package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	appnotif "github.com/hassad/boilerplateSaaS/backend/internal/app/notification"
	"github.com/hassad/boilerplateSaaS/backend/internal/handler/dto/response"
	"github.com/hassad/boilerplateSaaS/backend/internal/handler/middleware"
)

type NotificationHandler struct {
	svc *appnotif.Service
}

func NewNotificationHandler(svc *appnotif.Service) *NotificationHandler {
	return &NotificationHandler{svc: svc}
}

func (h *NotificationHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	unreadOnly := r.URL.Query().Get("unread") == "true"

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	notifications, total, err := h.svc.List(r.Context(), userID, unreadOnly, offset, limit)
	if err != nil {
		response.HandleDomainError(w, err)
		return
	}

	items := make([]response.NotificationResponse, len(notifications))
	for i, n := range notifications {
		items[i] = response.NotificationFromDomain(n)
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"notifications": items,
		"total":         total,
		"page":          page,
		"limit":         limit,
	})
}

func (h *NotificationHandler) UnreadCount(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())

	count, err := h.svc.UnreadCount(r.Context(), userID)
	if err != nil {
		response.HandleDomainError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]int{"count": count})
}

func (h *NotificationHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	notifID := chi.URLParam(r, "id")

	if err := h.svc.MarkAsRead(r.Context(), userID, notifID); err != nil {
		response.HandleDomainError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *NotificationHandler) MarkAllRead(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())

	if err := h.svc.MarkAllAsRead(r.Context(), userID); err != nil {
		response.HandleDomainError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
