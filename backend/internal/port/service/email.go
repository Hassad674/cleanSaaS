package service

import "context"

type Email struct {
	To      string
	Subject string
	HTML    string
}

type EmailService interface {
	Send(ctx context.Context, email Email) error
	SendTemplate(ctx context.Context, to string, template string, data map[string]string) error
}
