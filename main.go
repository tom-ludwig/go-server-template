package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"com.tom-ludwig/go-server-template/internal/api/health"
	"com.tom-ludwig/go-server-template/internal/api/users"
	"com.tom-ludwig/go-server-template/internal/config"
	"com.tom-ludwig/go-server-template/internal/middleware"
	"com.tom-ludwig/go-server-template/internal/repository"
	"com.tom-ludwig/go-server-template/internal/routes"
)

func main() {
	// Setup Logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	err := godotenv.Load()
	if err != nil {
		slog.Warn("Error loading .env file", "error", err)
	}

	// Load configuration
	cfg := config.Load()

	dbpool, err := connectToDatabase(cfg)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer dbpool.Close()
	slog.Info("Successfully connected to database")

	queries := repository.New(dbpool)

	// Initialize JWT auth if OIDC is enabled
	var jwtAuth *middleware.JWTAuth
	if cfg.OIDCEnabled {
		if cfg.OIDCIssuer == "" {
			slog.Error("OIDC_ISSUER must be set when OIDC_ENABLED is true")
			os.Exit(1)
		}
		var err error
		jwtAuth, err = middleware.NewJWTAuth(context.Background(), cfg.OIDCIssuer, cfg.OIDCAudience)
		if err != nil {
			slog.Error("Failed to initialize JWT auth", "error", err)
			os.Exit(1)
		}
		slog.Info("JWT authentication enabled", "issuer", cfg.OIDCIssuer)
	}

	router := routes.NewRouter(cfg, queries, jwtAuth)

	// Print registered routes in debug mode
	if cfg.DebugMode {
		// Add swagger specs here when you create new OpenAPI files
		swaggers := []*openapi3.T{}
		if s, err := health.GetSwagger(); err == nil {
			swaggers = append(swaggers, s)
		}
		if s, err := users.GetSwagger(); err == nil {
			swaggers = append(swaggers, s)
		}
		routes.PrintRoutes(router, swaggers)
	}

	port := fmt.Sprintf(":%s", cfg.Port)
	server := &http.Server{
		Addr:         port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	slog.Info("Server starting", "port", cfg.Port, "debug", cfg.DebugMode)
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}

func connectToDatabase(cfg *config.Config) (*pgxpool.Pool, error) {
	// Create context with timeout for database connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var dsn string
	if !cfg.PGLocal {
		dsn = fmt.Sprintf(
			"host=%s port=%s dbname=%s user=%s sslmode=%s sslcert=%s sslkey=%s sslrootcert=%s password=%s",
			cfg.PGHost, cfg.PGPort, cfg.PGDB, cfg.PGUser, cfg.PGSSLMode,
			cfg.PGTLSCert, cfg.PGTLSKey, cfg.PGSSLRootCert, cfg.PGPassword,
		)
	} else {
		// Local development
		dsn = fmt.Sprintf(
			"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
			cfg.PGHost, cfg.PGPort, cfg.PGDB, cfg.PGUser, cfg.PGPassword, cfg.PGSSLMode,
		)
	}

	// Parse config to validate DSN format, then create pool with context timeout
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database configuration: %w", err)
	}

	// Create pool with timeout context
	return pgxpool.NewWithConfig(ctx, config)
}
