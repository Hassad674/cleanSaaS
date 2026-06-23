package gemini

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/google/generative-ai-go/genai"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain/ai"
	"github.com/hassad/boilerplateSaaS/backend/internal/port/service"
	"github.com/hassad/boilerplateSaaS/backend/pkg/ctxutil"
)

// defaultCallTimeout is the fallback per-call ceiling used when an AIService is
// constructed without an explicit timeout. It applies to non-streaming Chat
// only — Stream is deliberately exempt (see Stream for the rationale).
const defaultCallTimeout = 15 * time.Second

// defaultSystemInstruction is set on every model to provide broad,
// helpful behaviour across all topics.
const defaultSystemInstruction = `You are a helpful, knowledgeable, and friendly AI assistant. You can help with anything: coding, writing, math, science, analysis, brainstorming, creative tasks, general knowledge, and more.

Guidelines:
- Give thorough, accurate, and well-structured answers
- Use markdown formatting when it helps readability (headers, lists, code blocks, bold, etc.)
- When asked about code, provide working examples with explanations
- Be conversational and natural — avoid sounding robotic or overly formal
- If you're unsure about something, say so honestly
- Adapt your response length to the complexity of the question — short for simple questions, detailed for complex ones`

type AIService struct {
	client      *genai.Client
	model       string
	callTimeout time.Duration
}

func NewAIService(client *genai.Client) *AIService {
	return NewAIServiceWithTimeout(client, defaultCallTimeout)
}

// NewAIServiceWithTimeout builds an AIService that bounds non-streaming Chat
// calls to callTimeout (a ceiling; a nearer caller deadline still wins).
// Streaming (Stream) is intentionally NOT bounded by this timeout.
func NewAIServiceWithTimeout(client *genai.Client, callTimeout time.Duration) *AIService {
	return &AIService{
		client:      client,
		model:       "gemini-2.5-flash",
		callTimeout: callTimeout,
	}
}

// configureModel applies shared settings (system instruction, temperature)
// to the given GenerativeModel.
func (s *AIService) configureModel(model *genai.GenerativeModel) {
	model.SystemInstruction = genai.NewUserContent(genai.Text(defaultSystemInstruction))
	temp := float32(0.7)
	model.Temperature = &temp
}

func (s *AIService) Chat(ctx context.Context, messages []ai.Message) (*service.AIResponse, error) {
	ctx, cancel := ctxutil.WithTimeout(ctx, s.callTimeout)
	defer cancel()

	model := s.client.GenerativeModel(s.model)
	s.configureModel(model)

	cs := model.StartChat()
	cs.History = toGeminiHistory(messages[:len(messages)-1])

	lastMsg := messages[len(messages)-1].Content
	resp, err := cs.SendMessage(ctx, genai.Text(lastMsg))
	if err != nil {
		return nil, fmt.Errorf("gemini chat: %w", err)
	}

	content := extractText(resp)
	tokens := 0
	if resp.UsageMetadata != nil {
		tokens = int(resp.UsageMetadata.TotalTokenCount)
	}

	return &service.AIResponse{
		Content: content,
		Model:   ai.ModelGemini,
		Tokens:  tokens,
	}, nil
}

// Stream streams an AI response chunk-by-chunk to writer. It deliberately does
// NOT apply the per-call timeout: a long, legitimate generation can far exceed
// any short ceiling, and clamping it would truncate valid output mid-stream.
// Cancellation is instead governed by the incoming ctx (request lifecycle) and
// the HTTP server's WriteTimeout, so a genuinely dead connection still unwinds.
func (s *AIService) Stream(ctx context.Context, messages []ai.Message, writer io.Writer) error {
	model := s.client.GenerativeModel(s.model)
	s.configureModel(model)

	cs := model.StartChat()
	cs.History = toGeminiHistory(messages[:len(messages)-1])

	lastMsg := messages[len(messages)-1].Content
	iter := cs.SendMessageStream(ctx, genai.Text(lastMsg))

	for {
		resp, err := iter.Next()
		if err != nil {
			if err.Error() == "no more items in iterator" {
				break
			}
			return fmt.Errorf("gemini stream: %w", err)
		}

		text := extractText(resp)
		if text != "" {
			if _, err := fmt.Fprint(writer, text); err != nil {
				return fmt.Errorf("writing stream chunk: %w", err)
			}
			if f, ok := writer.(interface{ Flush() }); ok {
				f.Flush()
			}
		}
	}

	return nil
}

func toGeminiHistory(messages []ai.Message) []*genai.Content {
	var history []*genai.Content
	for _, m := range messages {
		role := "user"
		if m.Role == ai.RoleAssistant {
			role = "model"
		} else if m.Role == ai.RoleSystem {
			continue // system messages handled separately
		}
		history = append(history, &genai.Content{
			Role:  role,
			Parts: []genai.Part{genai.Text(m.Content)},
		})
	}
	return history
}

func extractText(resp *genai.GenerateContentResponse) string {
	if resp == nil || len(resp.Candidates) == 0 {
		return ""
	}
	var text string
	for _, part := range resp.Candidates[0].Content.Parts {
		if t, ok := part.(genai.Text); ok {
			text += string(t)
		}
	}
	return text
}
