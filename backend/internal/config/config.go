package config

import "os"

type Config struct {
	DatabaseURL     string
	Port            string
	TelegramEnabled bool
	AllowedOrigins  string
	Locale          string
	LocalesDir      string
}

func Load() Config {
	return Config{
		DatabaseURL:     getEnv("DATABASE_URL", "postgres://posiflora:posiflora@localhost:5432/posiflora?sslmode=disable"),
		Port:            getEnv("PORT", "8080"),
		TelegramEnabled: getEnv("TELEGRAM_ENABLED", "false") == "true",
		AllowedOrigins:  getEnv("ALLOWED_ORIGINS", "http://localhost:5173"),
		Locale:          getEnv("LOCALE", "ru"),
		LocalesDir:      getEnv("LOCALES_DIR", "locales"),
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
