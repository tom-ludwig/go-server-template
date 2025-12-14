package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	// Server
	Port      string
	DebugMode bool

	// Database
	PGHost        string
	PGPort        string
	PGDB          string
	PGUser        string
	PGPassword    string
	PGSSLMode     string
	PGLocal       bool
	PGTLSCert     string
	PGTLSKey      string
	PGSSLRootCert string

	// CORS
	CORSAllowedOrigins   []string
	CORSAllowedMethods   []string
	CORSAllowedHeaders   []string
	CORSExposedHeaders   []string
	CORSAllowCredentials bool
	CORSMaxAge           int
}

func Load() *Config {
	return &Config{
		// Server
		Port:      getEnv("PORT", "8080"),
		DebugMode: getEnvBool("DEBUG_MODE", false),

		// Database
		PGHost:        getEnv("PG_HOST", "localhost"),
		PGPort:        getEnv("PG_PORT", "5432"),
		PGDB:          getEnv("PG_DB", "orbis"),
		PGUser:        getEnv("PG_USER", "user"),
		PGPassword:    getEnv("PG_PASSWORD", "password"),
		PGSSLMode:     getEnv("PG_SSLMODE", "verify-full"),
		PGLocal:       getEnvBool("PG_LOCAL", false),
		PGTLSCert:     getEnv("PG_CLIENT_CERT", "/certs/tls.crt"),
		PGTLSKey:      getEnv("PG_CLIENT_KEY", "/certs/tls.key"),
		PGSSLRootCert: getEnv("PG_SSLROOTCERT", "/certs/ca.crt"),

		// CORS - Default to permissive for development, override in production
		CORSAllowedOrigins:   getEnvSlice("CORS_ALLOWED_ORIGINS", []string{"*"}),
		CORSAllowedMethods:   getEnvSlice("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}),
		CORSAllowedHeaders:   getEnvSlice("CORS_ALLOWED_HEADERS", []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"}),
		CORSExposedHeaders:   getEnvSlice("CORS_EXPOSED_HEADERS", []string{"Link"}),
		CORSAllowCredentials: getEnvBool("CORS_ALLOW_CREDENTIALS", true),
		CORSMaxAge:           getEnvInt("CORS_MAX_AGE", 300),
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if v, ok := os.LookupEnv(key); ok {
		if parsed, err := strconv.ParseBool(v); err == nil {
			return parsed
		}
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v, ok := os.LookupEnv(key); ok {
		if parsed, err := strconv.Atoi(v); err == nil {
			return parsed
		}
	}
	return fallback
}

func getEnvSlice(key string, fallback []string) []string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		// Split by comma and trim whitespace
		parts := []string{}
		for _, part := range strings.Split(v, ",") {
			if trimmed := strings.TrimSpace(part); trimmed != "" {
				parts = append(parts, trimmed)
			}
		}
		if len(parts) > 0 {
			return parts
		}
	}
	return fallback
}
