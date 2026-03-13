package resend

import (
	resendgo "github.com/resend/resend-go/v2"
)

// NewClient creates a configured Resend API client.
func NewClient(apiKey string) *resendgo.Client {
	return resendgo.NewClient(apiKey)
}
