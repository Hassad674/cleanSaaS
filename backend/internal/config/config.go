package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
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

	// RedisURL enables the shared, multi-instance-safe components (distributed
	// rate limiter and cross-instance WebSocket fan-out). When empty (the default),
	// the app falls back to single-instance in-memory behavior so it still works
	// out-of-the-box with no Redis. Format: "redis://[:password@]host:port[/db]".
	RedisURL string

	// JWT lifecycle
	AccessTokenTTL  time.Duration // short-lived access token (default 15m)
	RefreshTokenTTL time.Duration // long-lived refresh token (default 720h / 30d)
	JWTIssuer       string        // "iss" claim (default "cleansaas")
	JWTAudience     string        // "aud" claim (default "cleansaas")

	// Timeout/cancellation discipline. These bound work that would otherwise run
	// unbounded: hung external API calls and background jobs without a request
	// deadline. They impose a CEILING only — a nearer caller deadline always wins.
	ExternalCallTimeout time.Duration // per-call ceiling for Stripe/Resend/R2/Gemini-Chat (default 15s)
	JobTimeout          time.Duration // per-invocation ceiling for scheduler jobs (default 30s)
	DBQueryTimeout      time.Duration // default ceiling for background DB ops with no deadline (default 15s)

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

	// Observability
	OTELExporterEndpoint string // OTLP/HTTP collector endpoint; empty disables tracing (default "")
	OTELServiceName      string // service.name on emitted spans (default "cleansaas-backend")
	MetricsEnabled       bool   // expose Prometheus /metrics + record HTTP metrics (default true)
}

func Load() *Config {
	cfg := &Config{
		AppEnv:              env("APP_ENV", "development"),
		Port:                env("PORT", "8081"),
		DatabaseURL:         cleanDSN(env("DATABASE_URL", "postgres://postgres:postgres@localhost:5433/cleansaas?sslmode=disable")),
		JWTSecret:           env("JWT_SECRET", defaultJWTSecret),
		FrontendURL:         env("FRONTEND_URL", "http://localhost:3010"),
		RedisURL:            env("REDIS_URL", ""),
		AccessTokenTTL:      envDuration("ACCESS_TOKEN_TTL", 15*time.Minute),
		RefreshTokenTTL:     envDuration("REFRESH_TOKEN_TTL", 720*time.Hour),
		JWTIssuer:           env("JWT_ISSUER", "cleansaas"),
		JWTAudience:         env("JWT_AUDIENCE", "cleansaas"),
		ExternalCallTimeout: envDuration("EXTERNAL_CALL_TIMEOUT", 15*time.Second),
		JobTimeout:          envDuration("JOB_TIMEOUT", 30*time.Second),
		DBQueryTimeout:      envDuration("DB_QUERY_TIMEOUT", 15*time.Second),
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

		OTELExporterEndpoint: env("OTEL_EXPORTER_OTLP_ENDPOINT", ""),
		OTELServiceName:      env("OTEL_SERVICE_NAME", "cleansaas-backend"),
		MetricsEnabled:       envBool("METRICS_ENABLED", true),
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

// envDuration reads a Go duration string (e.g. "15m", "720h") from the
// environment, falling back to the provided default when unset or unparseable.
func envDuration(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		log.Printf("WARNING: invalid duration for %s=%q, using default %s", key, v, fallback)
		return fallback
	}
	return d
}

// envBool reads a boolean ("1", "t", "true", "0", "false", ...) from the
// environment, falling back to the provided default when unset or unparseable.
func envBool(key string, fallback bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		log.Printf("WARNING: invalid bool for %s=%q, using default %v", key, v, fallback)
		return fallback
	}
	return b
}

func cleanDSN(dsn string) string {
	return strings.Replace(dsn, "&channel_binding=require", "", 1)
}
