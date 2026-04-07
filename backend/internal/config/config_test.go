package config_test

import (
	"os"
	"testing"

	"posiflora/backend/internal/config"
)

func TestLoad_Defaults(t *testing.T) {
	for _, key := range []string{"DATABASE_URL", "PORT", "TELEGRAM_ENABLED", "ALLOWED_ORIGINS", "LOCALE", "LOCALES_DIR"} {
		orig, exists := os.LookupEnv(key)
		if err := os.Unsetenv(key); err != nil {
			t.Fatalf("unsetenv %s: %v", key, err)
		}
		if exists {
			k, v := key, orig
			t.Cleanup(func() {
				if err := os.Setenv(k, v); err != nil {
					t.Errorf("restore env %s: %v", k, err)
				}
			})
		}
	}

	cfg := config.Load()

	if cfg.Port != "8080" {
		t.Errorf("expected port 8080, got %q", cfg.Port)
	}
	if cfg.TelegramEnabled {
		t.Error("expected TelegramEnabled=false by default")
	}
	if cfg.AllowedOrigins != "http://localhost:5173" {
		t.Errorf("unexpected AllowedOrigins: %q", cfg.AllowedOrigins)
	}
	if cfg.Locale != "ru" {
		t.Errorf("expected locale ru, got %q", cfg.Locale)
	}
	if cfg.LocalesDir != "locales" {
		t.Errorf("expected localesDir locales, got %q", cfg.LocalesDir)
	}
}

func TestLoad_FromEnv(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("TELEGRAM_ENABLED", "true")
	t.Setenv("LOCALE", "en")
	t.Setenv("ALLOWED_ORIGINS", "https://example.com")

	cfg := config.Load()

	if cfg.Port != "9090" {
		t.Errorf("expected port 9090, got %q", cfg.Port)
	}
	if !cfg.TelegramEnabled {
		t.Error("expected TelegramEnabled=true")
	}
	if cfg.Locale != "en" {
		t.Errorf("expected locale en, got %q", cfg.Locale)
	}
	if cfg.AllowedOrigins != "https://example.com" {
		t.Errorf("unexpected AllowedOrigins: %q", cfg.AllowedOrigins)
	}
}

func TestLoad_TelegramEnabled_False(t *testing.T) {
	t.Setenv("TELEGRAM_ENABLED", "false")
	cfg := config.Load()
	if cfg.TelegramEnabled {
		t.Error("expected TelegramEnabled=false when env is 'false'")
	}
}

func TestLoad_TelegramEnabled_InvalidValue(t *testing.T) {
	t.Setenv("TELEGRAM_ENABLED", "yes")
	cfg := config.Load()
	if cfg.TelegramEnabled {
		t.Error("expected TelegramEnabled=false for non-'true' value")
	}
}
