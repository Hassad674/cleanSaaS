package resend

import (
	"context"
	"fmt"
	"time"

	resendgo "github.com/resend/resend-go/v2"

	"github.com/hassad/boilerplateSaaS/backend/internal/port/service"
	"github.com/hassad/boilerplateSaaS/backend/pkg/ctxutil"
)

const defaultFrom = "CleanSaaS <onboarding@resend.dev>"

// defaultCallTimeout is the fallback per-call ceiling used when an EmailService
// is constructed without an explicit timeout.
const defaultCallTimeout = 15 * time.Second

// EmailService implements service.EmailService using Resend.
type EmailService struct {
	client      *resendgo.Client
	callTimeout time.Duration
}

// NewEmailService creates a new Resend email service.
func NewEmailService(apiKey string) *EmailService {
	return NewEmailServiceWithTimeout(apiKey, defaultCallTimeout)
}

// NewEmailServiceWithTimeout builds an EmailService that bounds every Resend API
// call to callTimeout (a ceiling; a nearer caller deadline still wins).
func NewEmailServiceWithTimeout(apiKey string, callTimeout time.Duration) *EmailService {
	return &EmailService{
		client:      NewClient(apiKey),
		callTimeout: callTimeout,
	}
}

func (s *EmailService) Send(ctx context.Context, email service.Email) error {
	ctx, cancel := ctxutil.WithTimeout(ctx, s.callTimeout)
	defer cancel()

	params := &resendgo.SendEmailRequest{
		From:    defaultFrom,
		To:      []string{email.To},
		Subject: email.Subject,
		Html:    email.HTML,
	}

	_, err := s.client.Emails.SendWithContext(ctx, params)
	if err != nil {
		return fmt.Errorf("sending email via resend: %w", err)
	}

	return nil
}

func (s *EmailService) SendTemplate(ctx context.Context, to string, template string, data map[string]string) error {
	subject, body := renderTemplate(template, data)

	return s.Send(ctx, service.Email{
		To:      to,
		Subject: subject,
		HTML:    body,
	})
}
