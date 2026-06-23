package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

// defaultJWTSecret is the local-development fallback. It is intentionally weak and
// MUST never be used outside development — Validate() refuses to boot if it is.
const defaultJWTSecret = "dev-secret-change-me"

// minJWTSecretLen is the minimum acceptable JWT secret length (bytes) outside dev.
const minJWTSecretLen = 32

type Config struct {
	AppEnv      string // "development" | "test" | "production" | "staging" | ...
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
		AppEnv:              env("APP_ENV", "development"),
		Port:                env("PORT", "8081"),
		DatabaseURL:         cleanDSN(env("DATABASE_URL", "postgres://postgres:postgres@localhost:5433/cleansaas?sslmode=disable")),
		JWTSecret:           env("JWT_SECRET", defaultJWTSecret),
		FrontendURL:         env("FRONTEND_URL", "http://localhost:3010"),
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

	if err := cfg.Validate(); err != nil {
		log.Fatalf("invalid configuration: %v", err)
	}

	return cfg
}

// IsDevelopment reports whether the app is running in a non-production context where
// insecure defaults are tolerated (development or test).
func (c *Config) IsDevelopment() bool {
	return c.AppEnv == "" || c.AppEnv == "development" || c.AppEnv == "dev" || c.AppEnv == "test"
}

// Validate enforces production-safety invariants. Outside development it FAILS CLOSED on
// insecure configuration (default/short JWT secret) rather than silently booting with it.
// In development it only warns, so local setup stays frictionless.
func (c *Config) Validate() error {
	if c.DatabaseURL == "" {
		return errors.New("DATABASE_URL is required")
	}

	if weak := c.weakJWTSecretReason(); weak != "" {
		if c.IsDevelopment() {
			log.Printf("WARNING: %s — acceptable for local development ONLY; set a strong JWT_SECRET before deploying", weak)
		} else {
			return fmt.Errorf("refusing to start in APP_ENV=%q: %s", c.AppEnv, weak)
		}
	}

	return nil
}

// weakJWTSecretReason returns a non-empty reason if the JWT secret is unsafe, else "".
func (c *Config) weakJWTSecretReason() string {
	switch {
	case c.JWTSecret == "":
		return "JWT_SECRET is empty"
	case c.JWTSecret == defaultJWTSecret:
		return "JWT_SECRET is the built-in default value"
	case len(c.JWTSecret) < minJWTSecretLen:
		return fmt.Sprintf("JWT_SECRET is too short (%d bytes; need at least %d)", len(c.JWTSecret), minJWTSecretLen)
	default:
		return ""
	}
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
