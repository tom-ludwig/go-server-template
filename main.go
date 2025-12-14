package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"com.tom-ludwig/go-server-template/internal/config"
	"com.tom-ludwig/go-server-template/internal/repository"
	"com.tom-ludwig/go-server-template/internal/routes"
	"com.tom-ludwig/go-server-template/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	// Setup Logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	// Load configuration
	cfg := config.Load()

	dbpool, err := connectToDatabase(cfg)
	if err != nil {
		log.Fatalf("An error occurred while connecting to the database: %s", err)
	}
	defer dbpool.Close()
	fmt.Println("Successfully connected to database.")

	queries := repository.New(dbpool)

	router := routes.NewRouter(cfg, queries)

	// Print registered routes in debug mode
	if cfg.DebugMode {
		utils.PrintRoutes(router)
	}

	port := fmt.Sprintf(":%s", cfg.Port)
	server := &http.Server{
		Addr:    port,
		Handler: router,
	}

	fmt.Printf("Running on %s (Debug: %v)\n", port, cfg.DebugMode)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Error occurred while starting server: %s", err)
	}
}

func connectToDatabase(cfg *config.Config) (*pgxpool.Pool, error) {
	if !cfg.PGLocal {
		dsn := fmt.Sprintf(
			"host=%s port=%s dbname=%s user=%s sslmode=%s sslcert=%s sslkey=%s sslrootcert=%s",
			cfg.PGHost, cfg.PGPort, cfg.PGDB, cfg.PGUser, cfg.PGSSLMode,
			cfg.PGTLSCert, cfg.PGTLSKey, cfg.PGSSLRootCert,
		)

		return pgxpool.New(context.Background(), dsn)
	} else {
		dsn := fmt.Sprintf(
			"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
			cfg.PGHost, cfg.PGPort, cfg.PGDB, cfg.PGUser, cfg.PGPassword, cfg.PGSSLMode,
		)

		return pgxpool.New(context.Background(), dsn)
	}
}
