package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/hassad/boilerplateSaaS/backend/internal/app/billing"
	"github.com/hassad/boilerplateSaaS/backend/internal/handler/dto/request"
	"github.com/hassad/boilerplateSaaS/backend/internal/handler/dto/response"
	"github.com/hassad/boilerplateSaaS/backend/internal/handler/middleware"
)

type BillingHandler struct {
	svc *billing.Service
}

func NewBillingHandler(svc *billing.Service) *BillingHandler {
	return &BillingHandler{svc: svc}
}

func (h *BillingHandler) GetPlans(w http.ResponseWriter, r *http.Request) {
	plans, err := h.svc.GetPlans(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to get plans")
		return
	}

	result := make([]response.PlanResponse, len(plans))
	for i, p := range plans {
		result[i] = response.PlanFromDomain(p)
	}

	response.JSON(w, http.StatusOK, result)
}

func (h *BillingHandler) CreateCheckout(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())

	var req request.CheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.PlanID == "" {
		response.Error(w, http.StatusBadRequest, "plan_id is required")
		return
	}

	url, err := h.svc.CreateCheckout(r.Context(), userID, req.PlanID)
	if err != nil {
		response.HandleDomainError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"url": url})
}

func (h *BillingHandler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())

	sub, err := h.svc.GetSubscription(r.Context(), userID)
	if err != nil {
		response.HandleDomainError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, response.SubscriptionFromDomain(sub))
}

func (h *BillingHandler) CancelSubscription(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())

	if err := h.svc.CancelSubscription(r.Context(), userID); err != nil {
		response.HandleDomainError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "subscription will be canceled at period end"})
}

func (h *BillingHandler) CreatePortalSession(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())

	url, err := h.svc.CreatePortalSession(r.Context(), userID)
	if err != nil {
		response.HandleDomainError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"url": url})
}

func (h *BillingHandler) GetInvoices(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())

	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 20
	}

	invoices, total, err := h.svc.GetInvoices(r.Context(), userID, offset, limit)
	if err != nil {
		response.HandleDomainError(w, err)
		return
	}

	result := make([]response.InvoiceResponse, len(invoices))
	for i, inv := range invoices {
		result[i] = response.InvoiceFromDomain(inv)
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"invoices": result,
		"total":    total,
	})
}

func (h *BillingHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	payload, err := io.ReadAll(io.LimitReader(r.Body, 65536))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "failed to read body")
		return
	}

	signature := r.Header.Get("Stripe-Signature")
	if signature == "" {
		response.Error(w, http.StatusBadRequest, "missing stripe signature")
		return
	}

	if err := h.svc.HandleWebhook(r.Context(), payload, signature); err != nil {
		response.Error(w, http.StatusBadRequest, "webhook processing failed")
		return
	}

	w.WriteHeader(http.StatusOK)
}
