package routes

import (
	"com.tom-ludwig/go-server-template/internal/api/health"
	"com.tom-ludwig/go-server-template/internal/api/users"
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

	// Core middleware (applied to all routes)
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

	// Mount Health API (public - no auth required)
	mountHealthAPI(r, queries)

	// Mount Users API (can add auth middleware here)
	mountUsersAPI(r, queries)

	return r
}

// mountHealthAPI mounts health check endpoints (public, no authentication)
func mountHealthAPI(r chi.Router, queries *repository.Queries) {
	healthHandler := handler.NewHealthHandler(queries)
	strictHealthServer := health.NewStrictHandler(healthHandler, nil)

	// OpenAPI request validation for health routes
	healthSwagger, err := health.GetSwagger()
	if err != nil {
		panic("Failed to load health swagger spec")
	}

	r.Group(func(r chi.Router) {
		// Public routes - no auth middleware
		r.Use(oapimiddleware.OapiRequestValidator(healthSwagger))
		health.HandlerFromMux(strictHealthServer, r)
	})
}

// mountUsersAPI mounts user management endpoints (can be protected with auth)
func mountUsersAPI(r chi.Router, queries *repository.Queries) {
	userHandler := handler.NewUserHandler(queries)
	strictUsersServer := users.NewStrictHandler(userHandler, nil)

	// OpenAPI request validation for user routes
	usersSwagger, err := users.GetSwagger()
	if err != nil {
		panic("Failed to load users swagger spec")
	}

	r.Group(func(r chi.Router) {
		// Add authentication/authorization middleware here:
		// r.Use(middleware.AuthRequired)
		// r.Use(middleware.RequireRole("admin"))

		r.Use(oapimiddleware.OapiRequestValidator(usersSwagger))
		users.HandlerFromMux(strictUsersServer, r)
	})
}
