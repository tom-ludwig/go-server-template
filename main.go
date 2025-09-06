package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"com.tom-ludwig/go-server-template/internal/repository"
	"com.tom-ludwig/go-server-template/internal/routes"
	"github.com/gin-gonic/gin"
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

	dbpool, err := connectToDatabase()
	if err != nil {
		log.Fatalf("An error occured while connecting to the database: %s", err)
	}
	defer dbpool.Close()
	fmt.Println("Succesfully connected to database.")

	queries := repository.New(dbpool)

	if getEnv("DEBUG_MODE", "false") == "true" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := routes.NewRouter(queries)

	port := fmt.Sprintf(":%s", getEnv("PORT", "8080"))
	fmt.Printf("Running on %s \n", port)
	err = router.Run(port)
	if err != nil {
		log.Fatalf("Error occued while starting router: %s", err)
	}
}

func connectToDatabase() (*pgxpool.Pool, error) {
	devMode := getEnv("PG_LOCAL", "false") == "true"
	pg_host := getEnv("PG_HOST", "localhost")
	pg_port := getEnv("PG_PORT", "5432")
	pg_db := getEnv("PG_DB", "orbis")
	pg_user := getEnv("PG_USER", "user")
	pg_sslmode := getEnv("PG_SSLMODE", "verify-full")

	if !devMode {
		pg_tls_cert := getEnv("PG_CLIENT_CERT", "/certs/tls.crt")
		pg_client_key := getEnv("PG_CLIENT_KEY", "/certs/tls.key")
		pg_ssl_root_cert := getEnv("PG_SSLROOTCERT", "/certs/ca.crt")

		dsn := fmt.Sprintf(
			"host=%s port=%s dbname=%s user=%s sslmode=%s sslcert=%s sslkey=%s sslrootcert=%s",
			pg_host, pg_port, pg_db, pg_user, pg_sslmode, pg_tls_cert, pg_client_key, pg_ssl_root_cert,
		)

		return pgxpool.New(context.Background(), dsn)
	} else {
		pg_password := getEnv("PG_PASSWORD", "password")

		dsn := fmt.Sprintf(
			"host=%s port=%s dbname=%s user=%s password =%s sslmode=%s",
			pg_host, pg_port, pg_db, pg_user, pg_password, pg_sslmode,
		)

		return pgxpool.New(context.Background(), dsn)
	}
}
func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
