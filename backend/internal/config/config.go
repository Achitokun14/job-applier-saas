package config

import (
	"log"
	"os"
	"strings"
)

// NOTE: This package still uses stdlib "log" for fatal startup errors that
// occur before zerolog is initialised.  All runtime logging should use the
// zerolog instance created in main.

type Config struct {
	DatabaseURL                  string
	JWTSecret                    string
	JWTExpiry                    string
	PythonServiceURL             string
	CORSAllowedOrigins           []string
	RedisURL                     string
	EncryptionKey                string
	AppEnv                       string
	StripeSecretKey              string
	StripeWebhookSecret          string
	StripePriceProMonthly        string
	StripePriceEnterpriseMonthly string
	SentryDSN                    string
}

func Load() *Config {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "sqlite:./data/jobapplier.db"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-change-in-production"
	}

	jwtExpiry := os.Getenv("JWT_EXPIRY")
	if jwtExpiry == "" {
		jwtExpiry = "24h"
	}

	pythonURL := os.Getenv("PYTHON_SERVICE_URL")
	if pythonURL == "" {
		pythonURL = "http://localhost:8001"
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	encryptionKey := os.Getenv("ENCRYPTION_KEY")

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "development"
	}

	stripeSecretKey := os.Getenv("STRIPE_SECRET_KEY")
	stripeWebhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	stripePriceProMonthly := os.Getenv("STRIPE_PRICE_PRO_MONTHLY")
	stripePriceEnterpriseMonthly := os.Getenv("STRIPE_PRICE_ENTERPRISE_MONTHLY")

	sentryDSN := os.Getenv("SENTRY_DSN")

	var corsOrigins []string
	corsRaw := os.Getenv("CORS_ALLOWED_ORIGINS")
	if corsRaw != "" {
		for _, origin := range strings.Split(corsRaw, ",") {
			trimmed := strings.TrimSpace(origin)
			if trimmed != "" {
				corsOrigins = append(corsOrigins, trimmed)
			}
		}
	}

	// Validate security settings in production
	if appEnv == "production" {
		if jwtSecret == "default-secret-change-in-production" {
			log.Fatal("FATAL: JWT_SECRET must be changed from the default value in production")
		}
		if len(jwtSecret) < 32 {
			log.Fatal("FATAL: JWT_SECRET must be at least 32 characters long in production")
		}
	}

	return &Config{
		DatabaseURL:                  dbURL,
		JWTSecret:                    jwtSecret,
		JWTExpiry:                    jwtExpiry,
		PythonServiceURL:             pythonURL,
		CORSAllowedOrigins:           corsOrigins,
		RedisURL:                     redisURL,
		EncryptionKey:                encryptionKey,
		AppEnv:                       appEnv,
		StripeSecretKey:              stripeSecretKey,
		StripeWebhookSecret:          stripeWebhookSecret,
		StripePriceProMonthly:        stripePriceProMonthly,
		StripePriceEnterpriseMonthly: stripePriceEnterpriseMonthly,
		SentryDSN:                    sentryDSN,
	}
}
