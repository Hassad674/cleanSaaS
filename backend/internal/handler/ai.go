package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	appai "github.com/hassad/boilerplateSaaS/backend/internal/app/ai"
	"github.com/hassad/boilerplateSaaS/backend/internal/handler/dto/request"
	"github.com/hassad/boilerplateSaaS/backend/internal/handler/dto/response"
	"github.com/hassad/boilerplateSaaS/backend/internal/handler/middleware"
)

type AIHandler struct {
	svc *appai.Service
}

func NewAIHandler(svc *appai.Service) *AIHandler {
	return &AIHandler{svc: svc}
}

func (h *AIHandler) CreateConversation(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())

	var req request.CreateConversationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		req.Title = "New conversation"
	}

	conv, err := h.svc.CreateConversation(r.Context(), userID, req.Title)
	if err != nil {
		response.HandleDomainError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, response.ConversationFromDomain(conv))
}

func (h *AIHandler) ListConversations(w http.ResponseWriter, r *http.Request) {
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

	convos, total, err := h.svc.ListConversations(r.Context(), userID, offset, limit)
	if err != nil {
		response.HandleDomainError(w, err)
		return
	}

	items := make([]response.ConversationResponse, len(convos))
	for i, c := range convos {
		items[i] = response.ConversationFromDomain(c)
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"conversations": items,
		"total":         total,
		"page":          page,
		"limit":         limit,
	})
}

func (h *AIHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	convID := chi.URLParam(r, "id")

	conv, err := h.svc.GetConversation(r.Context(), userID, convID)
	if err != nil {
		response.HandleDomainError(w, err)
		return
	}

	msgs := make([]response.MessageResponse, len(conv.Messages))
	for i, m := range conv.Messages {
		msgs[i] = response.MessageResponse{
			Role:    string(m.Role),
			Content: m.Content,
		}
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"messages": msgs,
	})
}

func (h *AIHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	convID := chi.URLParam(r, "id")

	var req request.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Content == "" {
		response.Error(w, http.StatusBadRequest, "content is required")
		return
	}

	if len(req.Content) > 32000 {
		response.Error(w, http.StatusBadRequest, "message too long (max 32000 characters)")
		return
	}

	reply, err := h.svc.SendMessage(r.Context(), userID, convID, req.Content)
	if err != nil {
		response.HandleDomainError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, response.MessageResponse{
		Role:    "assistant",
		Content: reply,
	})
}

func (h *AIHandler) StreamMessage(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	convID := chi.URLParam(r, "id")

	var req request.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Content == "" {
		response.Error(w, http.StatusBadRequest, "content is required")
		return
	}

	if len(req.Content) > 32000 {
		response.Error(w, http.StatusBadRequest, "message too long (max 32000 characters)")
		return
	}

	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		response.Error(w, http.StatusInternalServerError, "streaming not supported")
		return
	}

	sseWriter := &sseWriter{w: w, flusher: flusher}

	if err := h.svc.StreamMessage(r.Context(), userID, convID, req.Content, sseWriter); err != nil {
		// If headers already sent, write error as SSE event.
		// Never expose internal error details to the client.
		fmt.Fprintf(w, "event: error\ndata: an error occurred while processing your message\n\n")
		flusher.Flush()
		return
	}

	// Send done event
	fmt.Fprintf(w, "event: done\ndata: [DONE]\n\n")
	flusher.Flush()
}

func (h *AIHandler) DeleteConversation(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	convID := chi.URLParam(r, "id")

	if err := h.svc.DeleteConversation(r.Context(), userID, convID); err != nil {
		response.HandleDomainError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// sseWriter wraps http.ResponseWriter to format output as SSE data events
type sseWriter struct {
	w       http.ResponseWriter
	flusher http.Flusher
}

func (s *sseWriter) Write(p []byte) (n int, err error) {
	n, err = fmt.Fprintf(s.w, "data: %s\n\n", string(p))
	s.flusher.Flush()
	return
}
