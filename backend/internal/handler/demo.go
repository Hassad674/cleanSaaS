package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	domainai "github.com/hassad/boilerplateSaaS/backend/internal/domain/ai"
	"github.com/hassad/boilerplateSaaS/backend/internal/handler/dto/response"
	"github.com/hassad/boilerplateSaaS/backend/internal/port/service"
)

// DemoHandler provides a public (no-auth) AI chat endpoint for the landing-page
// demo. It streams responses from the configured AI provider without persisting
// conversations to the database.
type DemoHandler struct {
	ai service.AIService
}

func NewDemoHandler(ai service.AIService) *DemoHandler {
	return &DemoHandler{ai: ai}
}

type demoChatRequest struct {
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
}

// StreamChat handles POST /demo/ai/chat
// It accepts a list of messages (conversation history) and streams the AI
// response back as Server-Sent Events, identical in format to the
// authenticated /ai/conversations/{id}/stream endpoint.
func (h *DemoHandler) StreamChat(w http.ResponseWriter, r *http.Request) {
	var req demoChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || len(req.Messages) == 0 {
		response.Error(w, http.StatusBadRequest, "messages array is required")
		return
	}

	// Enforce a reasonable limit on history length for the demo
	if len(req.Messages) > 30 {
		response.Error(w, http.StatusBadRequest, "too many messages (max 30 for demo)")
		return
	}

	// Validate last message content length
	lastMsg := req.Messages[len(req.Messages)-1]
	if lastMsg.Content == "" {
		response.Error(w, http.StatusBadRequest, "last message content is required")
		return
	}
	if len(lastMsg.Content) > 4000 {
		response.Error(w, http.StatusBadRequest, "message too long (max 4000 characters for demo)")
		return
	}

	// Convert to domain messages
	messages := make([]domainai.Message, 0, len(req.Messages))
	for _, m := range req.Messages {
		role := domainai.RoleUser
		switch m.Role {
		case "assistant":
			role = domainai.RoleAssistant
		case "system":
			role = domainai.RoleSystem
		}
		messages = append(messages, domainai.Message{
			Role:    role,
			Content: m.Content,
		})
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

	sseW := &demoSSEWriter{w: w, flusher: flusher}

	if err := h.ai.Stream(r.Context(), messages, sseW); err != nil {
		fmt.Fprintf(w, "event: error\ndata: %s\n\n", err.Error())
		flusher.Flush()
		return
	}

	// Send done event
	fmt.Fprintf(w, "event: done\ndata: [DONE]\n\n")
	flusher.Flush()
}

type demoSSEWriter struct {
	w       http.ResponseWriter
	flusher http.Flusher
}

func (s *demoSSEWriter) Write(p []byte) (n int, err error) {
	n, err = fmt.Fprintf(s.w, "data: %s\n\n", string(p))
	s.flusher.Flush()
	return
}
