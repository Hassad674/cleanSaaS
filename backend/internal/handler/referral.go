package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	appreferral "github.com/hassad/boilerplateSaaS/backend/internal/app/referral"
	"github.com/hassad/boilerplateSaaS/backend/internal/handler/dto/request"
	"github.com/hassad/boilerplateSaaS/backend/internal/handler/dto/response"
	"github.com/hassad/boilerplateSaaS/backend/internal/handler/middleware"
)

type ReferralHandler struct {
	svc *appreferral.Service
}

func NewReferralHandler(svc *appreferral.Service) *ReferralHandler {
	return &ReferralHandler{svc: svc}
}

func (h *ReferralHandler) GetCode(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())

	code, err := h.svc.GetOrCreateCode(r.Context(), userID)
	if err != nil {
		response.HandleDomainError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"code": code})
}

func (h *ReferralHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())

	stats, err := h.svc.GetStats(r.Context(), userID)
	if err != nil {
		response.HandleDomainError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, stats)
}

func (h *ReferralHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	referrals, total, err := h.svc.ListReferrals(r.Context(), userID, offset, limit)
	if err != nil {
		response.HandleDomainError(w, err)
		return
	}

	items := make([]response.ReferralResponse, len(referrals))
	for i, ref := range referrals {
		items[i] = response.ReferralFromDomain(ref)
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"referrals": items,
		"total":     total,
		"page":      page,
		"limit":     limit,
	})
}

func (h *ReferralHandler) Apply(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())

	var req request.ApplyReferralRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Code == "" {
		response.Error(w, http.StatusBadRequest, "referral code is required")
		return
	}

	if err := h.svc.ApplyReferral(r.Context(), userID, req.Code); err != nil {
		response.HandleDomainError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "referral applied successfully"})
}
