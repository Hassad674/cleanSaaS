package resend

import (
	"context"
	"fmt"

	resendgo "github.com/resend/resend-go/v2"

	"github.com/hassad/boilerplateSaaS/backend/internal/port/service"
)

const defaultFrom = "CleanSaaS <onboarding@resend.dev>"

// EmailService implements service.EmailService using Resend.
type EmailService struct {
	client *resendgo.Client
}

// NewEmailService creates a new Resend email service.
func NewEmailService(apiKey string) *EmailService {
	return &EmailService{
		client: NewClient(apiKey),
	}
}

func (s *EmailService) Send(ctx context.Context, email service.Email) error {
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
