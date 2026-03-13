package resend

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWelcomeEmail(t *testing.T) {
	subject, body := WelcomeEmail("John")
	assert.Equal(t, "Welcome to CleanSaaS!", subject)
	assert.Contains(t, body, "John")
	assert.Contains(t, body, "Welcome")
}

func TestVerificationEmail(t *testing.T) {
	link := "https://example.com/verify?token=abc123"
	subject, body := VerificationEmail("Jane", link)
	assert.Equal(t, "Verify your email address", subject)
	assert.Contains(t, body, "Jane")
	assert.Contains(t, body, link)
}

func TestPasswordResetEmail(t *testing.T) {
	link := "https://example.com/reset?token=xyz789"
	subject, body := PasswordResetEmail("Bob", link)
	assert.Equal(t, "Reset your password", subject)
	assert.Contains(t, body, "Bob")
	assert.Contains(t, body, link)
}

func TestRenderTemplate_Welcome(t *testing.T) {
	subject, body := renderTemplate("welcome", map[string]string{"name": "Alice"})
	assert.Equal(t, "Welcome to CleanSaaS!", subject)
	assert.Contains(t, body, "Alice")
}

func TestRenderTemplate_Unknown(t *testing.T) {
	subject, body := renderTemplate("unknown", map[string]string{"message": "Hello"})
	assert.Equal(t, "CleanSaaS Notification", subject)
	assert.Contains(t, body, "Hello")
}
