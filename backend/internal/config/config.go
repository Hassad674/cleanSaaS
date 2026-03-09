package config

import (
	"log"
	"os"
	"strings"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	FrontendURL string

	// Stripe
	StripeKey           string
	StripeWebhookSecret string

	// Resend
	ResendKey string

	// AI
	ClaudeKey string
	OpenAIKey string
	GeminiKey string

	// Storage (Cloudflare R2)
	R2AccountID  string
	R2AccessKey  string
	R2SecretKey  string
	R2BucketName string
	R2PublicURL  string

	// OAuth
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
}

func Load() *Config {
	cfg := &Config{
		Port:                env("PORT", "8081"),
		DatabaseURL:         cleanDSN(env("DATABASE_URL", "postgres://postgres:postgres@localhost:5433/cleansaas?sslmode=disable")),
		JWTSecret:           env("JWT_SECRET", "dev-secret-change-me"),
		FrontendURL:         env("FRONTEND_URL", "http://localhost:3000"),
		StripeKey:           env("STRIPE_SECRET_KEY", ""),
		StripeWebhookSecret: env("STRIPE_WEBHOOK_SECRET", ""),
		ResendKey:           env("RESEND_API_KEY", ""),
		ClaudeKey:           env("CLAUDE_API_KEY", ""),
		OpenAIKey:           env("OPENAI_API_KEY", ""),
		GeminiKey:           env("GEMINI_API_KEY", ""),
		R2AccountID:         env("R2_ACCOUNT_ID", ""),
		R2AccessKey:         env("R2_ACCESS_KEY", ""),
		R2SecretKey:         env("R2_SECRET_KEY", ""),
		R2BucketName:        env("R2_BUCKET_NAME", ""),
		R2PublicURL:         env("R2_PUBLIC_URL", ""),
		GoogleClientID:      env("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret:  env("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:   env("GOOGLE_REDIRECT_URL", ""),
	}

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	return cfg
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func cleanDSN(dsn string) string {
	return strings.Replace(dsn, "&channel_binding=require", "", 1)
}
