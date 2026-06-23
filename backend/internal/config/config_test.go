package config

import (
	"strings"
	"testing"
	"time"
)

const strongSecret = "a-sufficiently-long-random-secret-value-123456" // >= 32 bytes

func TestConfig_Validate_RequiresDatabaseURL(t *testing.T) {
	c := &Config{AppEnv: "development", DatabaseURL: "", JWTSecret: strongSecret}
	if err := c.Validate(); err == nil {
		t.Fatal("expected error when DATABASE_URL is empty")
	}
}

func TestConfig_Validate_DevAllowsWeakSecret(t *testing.T) {
	cases := []string{"", defaultJWTSecret, "short"}
	for _, secret := range cases {
		c := &Config{AppEnv: "development", DatabaseURL: "postgres://x", JWTSecret: secret}
		if err := c.Validate(); err != nil {
			t.Errorf("development should tolerate weak secret %q, got error: %v", secret, err)
		}
	}
}

func TestConfig_Validate_ProdRejectsWeakSecret(t *testing.T) {
	cases := map[string]string{
		"empty":   "",
		"default": defaultJWTSecret,
		"short":   "too-short",
	}
	for name, secret := range cases {
		t.Run(name, func(t *testing.T) {
			c := &Config{AppEnv: "production", DatabaseURL: "postgres://x", JWTSecret: secret}
			err := c.Validate()
			if err == nil {
				t.Fatalf("production must reject %s JWT secret", name)
			}
			if !strings.Contains(err.Error(), "refusing to start") {
				t.Errorf("unexpected error message: %v", err)
			}
		})
	}
}

func TestConfig_Validate_ProdAcceptsStrongSecret(t *testing.T) {
	c := &Config{AppEnv: "production", DatabaseURL: "postgres://x", JWTSecret: strongSecret}
	if err := c.Validate(); err != nil {
		t.Errorf("production should accept a strong secret, got: %v", err)
	}
}

func TestConfig_IsDevelopment(t *testing.T) {
	dev := []string{"", "development", "dev", "test"}
	for _, e := range dev {
		if !(&Config{AppEnv: e}).IsDevelopment() {
			t.Errorf("AppEnv %q should be development-like", e)
		}
	}
	prod := []string{"production", "staging", "prod"}
	for _, e := range prod {
		if (&Config{AppEnv: e}).IsDevelopment() {
			t.Errorf("AppEnv %q should NOT be development-like", e)
		}
	}
}

func TestConfig_Load_TimeoutDefaults(t *testing.T) {
	t.Setenv("JWT_SECRET", strongSecret)
	t.Setenv("APP_ENV", "test")
	// Ensure the timeout vars are unset so defaults apply.
	t.Setenv("EXTERNAL_CALL_TIMEOUT", "")
	t.Setenv("JOB_TIMEOUT", "")
	t.Setenv("DB_QUERY_TIMEOUT", "")

	cfg := Load()
	if cfg.ExternalCallTimeout != 15*time.Second {
		t.Errorf("ExternalCallTimeout default = %s, want 15s", cfg.ExternalCallTimeout)
	}
	if cfg.JobTimeout != 30*time.Second {
		t.Errorf("JobTimeout default = %s, want 30s", cfg.JobTimeout)
	}
	if cfg.DBQueryTimeout != 15*time.Second {
		t.Errorf("DBQueryTimeout default = %s, want 15s", cfg.DBQueryTimeout)
	}
}

func TestConfig_Load_TimeoutOverrides(t *testing.T) {
	t.Setenv("JWT_SECRET", strongSecret)
	t.Setenv("APP_ENV", "test")
	t.Setenv("EXTERNAL_CALL_TIMEOUT", "5s")
	t.Setenv("JOB_TIMEOUT", "45s")
	t.Setenv("DB_QUERY_TIMEOUT", "8s")

	cfg := Load()
	if cfg.ExternalCallTimeout != 5*time.Second {
		t.Errorf("ExternalCallTimeout = %s, want 5s", cfg.ExternalCallTimeout)
	}
	if cfg.JobTimeout != 45*time.Second {
		t.Errorf("JobTimeout = %s, want 45s", cfg.JobTimeout)
	}
	if cfg.DBQueryTimeout != 8*time.Second {
		t.Errorf("DBQueryTimeout = %s, want 8s", cfg.DBQueryTimeout)
	}
}

func TestConfig_cleanDSN_StripsChannelBinding(t *testing.T) {
	in := "postgres://u:p@host/db?sslmode=require&channel_binding=require"
	out := cleanDSN(in)
	if strings.Contains(out, "channel_binding") {
		t.Errorf("channel_binding should be stripped, got: %s", out)
	}
}
