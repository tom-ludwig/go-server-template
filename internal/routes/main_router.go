package routes

import (
	"com.tom-ludwig/go-server-template/internal/api"
	"com.tom-ludwig/go-server-template/internal/config"
	"com.tom-ludwig/go-server-template/internal/handler"
	"com.tom-ludwig/go-server-template/internal/middleware"
	"com.tom-ludwig/go-server-template/internal/repository"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	oapimiddleware "github.com/oapi-codegen/nethttp-middleware"
)

func NewRouter(cfg *config.Config, queries *repository.Queries) chi.Router {
	r := chi.NewRouter()

	// Core middleware
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(middleware.ColoredLogger(cfg.DebugMode))
	r.Use(chimiddleware.Recoverer)

	// Security headers
	r.Use(middleware.SecurityHeaders)

	// CORS configuration
	corsOptions := cors.Options{
		AllowedOrigins:   cfg.CORSAllowedOrigins,
		AllowedMethods:   cfg.CORSAllowedMethods,
		AllowedHeaders:   cfg.CORSAllowedHeaders,
		ExposedHeaders:   cfg.CORSExposedHeaders,
		AllowCredentials: cfg.CORSAllowCredentials,
		MaxAge:           cfg.CORSMaxAge,
	}
	r.Use(cors.Handler(corsOptions))

	// OpenAPI request validation
	swagger, err := api.GetSwagger()
	if err != nil {
		panic("Failed to load swagger spec")
	}
	r.Use(oapimiddleware.OapiRequestValidator(swagger))

	// Setup handlers
	server := handler.NewServer(queries)
	strictServer := api.NewStrictHandler(server, nil)

	// Register routes using chi handler
	api.HandlerFromMux(strictServer, r)

	return r
}
