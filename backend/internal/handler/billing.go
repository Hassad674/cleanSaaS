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

// DemoCheckout handles POST /demo/billing/checkout
// Creates a real Stripe Checkout Session without requiring authentication.
// Used by the public demo page to demonstrate the billing flow.
func (h *BillingHandler) DemoCheckout(w http.ResponseWriter, r *http.Request) {
	var req request.DemoCheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.PlanID == "" {
		response.Error(w, http.StatusBadRequest, "plan_id is required")
		return
	}
	if req.SuccessURL == "" || req.CancelURL == "" {
		response.Error(w, http.StatusBadRequest, "success_url and cancel_url are required")
		return
	}

	url, err := h.svc.DemoCheckout(r.Context(), req.PlanID, req.SuccessURL, req.CancelURL)
	if err != nil {
		response.HandleDomainError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"url": url})
}

// DemoSession handles GET /demo/billing/session?session_id=xxx
// Retrieves a completed Stripe Checkout Session and returns plan details.
func (h *BillingHandler) DemoSession(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		response.Error(w, http.StatusBadRequest, "session_id is required")
		return
	}

	info, err := h.svc.GetDemoSession(r.Context(), sessionID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to retrieve session")
		return
	}

	response.JSON(w, http.StatusOK, info)
}

// DemoPortal handles POST /demo/billing/portal
// Creates a Stripe Billing Portal session for a demo customer so they can
// manage their subscription (upgrade, cancel, update payment method).
func (h *BillingHandler) DemoPortal(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CustomerID string `json:"customer_id"`
		ReturnURL  string `json:"return_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.CustomerID == "" || req.ReturnURL == "" {
		response.Error(w, http.StatusBadRequest, "customer_id and return_url are required")
		return
	}

	url, err := h.svc.DemoPortalSession(r.Context(), req.CustomerID, req.ReturnURL)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to create portal session")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"url": url})
}

func (h *BillingHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	// 400 (bad request) is reserved for genuinely malformed requests — a body
	// we can't read or a missing signature. Stripe treats 4xx as "do not
	// retry", so we must NOT return 400 for transient/processing failures.
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

	// Processing errors (invalid signature, DB failures, downstream errors) are
	// returned as 500 so Stripe retries the delivery. Combined with the
	// idempotency check in the service, retries are safe and legitimate events
	// are never permanently dropped on a transient failure.
	if err := h.svc.HandleWebhook(r.Context(), payload, signature); err != nil {
		response.Error(w, http.StatusInternalServerError, "webhook processing failed")
		return
	}

	w.WriteHeader(http.StatusOK)
}
