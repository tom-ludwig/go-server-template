package routes

import (
	"log/slog"
	"os"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	oapimiddleware "github.com/oapi-codegen/nethttp-middleware"

	"com.tom-ludwig/go-server-template/internal/api/health"
	"com.tom-ludwig/go-server-template/internal/api/users"
	"com.tom-ludwig/go-server-template/internal/config"
	"com.tom-ludwig/go-server-template/internal/handler"
	"com.tom-ludwig/go-server-template/internal/middleware"
	"com.tom-ludwig/go-server-template/internal/repository"
)

func NewRouter(cfg *config.Config, queries *repository.Queries, jwtAuth *middleware.JWTAuth) chi.Router {
	r := chi.NewRouter()

	// Core middleware (applied to all routes)
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(middleware.RequestLogger(cfg.DebugMode))
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

	// Mount Health API (public)
	mountHealthAPI(r, queries)

	// Mount Users API (protected with JWT auth if enabled)
	mountUsersAPI(r, queries, jwtAuth)

	return r
}

// mountHealthAPI mounts health check endpoints
func mountHealthAPI(r chi.Router, queries *repository.Queries) {
	healthHandler := handler.NewHealthHandler(queries)
	strictHealthServer := health.NewStrictHandler(healthHandler, nil)

	healthSwagger, err := health.GetSwagger()
	if err != nil {
		slog.Error("Failed to load health swagger spec", "error", err)
		os.Exit(1)
	}

	r.Group(func(r chi.Router) {
		r.Use(oapimiddleware.OapiRequestValidator(healthSwagger))
		health.HandlerFromMux(strictHealthServer, r)
	})
}

// mountUsersAPI mounts user management endpoints
func mountUsersAPI(r chi.Router, queries *repository.Queries, jwtAuth *middleware.JWTAuth) {
	userHandler := handler.NewUserHandler(queries)
	strictUsersServer := users.NewStrictHandler(userHandler, nil)

	usersSwagger, err := users.GetSwagger()
	if err != nil {
		slog.Error("Failed to load users swagger spec", "error", err)
		os.Exit(1)
	}

	r.Group(func(r chi.Router) {
		// Configure request validator with JWT authentication function
		var validatorOptions oapimiddleware.Options
		if jwtAuth != nil {
			validatorOptions = oapimiddleware.Options{
				Options: openapi3filter.Options{
					AuthenticationFunc: jwtAuth.OAPIMiddleware(jwtAuth),
				},
			}
		}
		r.Use(oapimiddleware.OapiRequestValidatorWithOptions(usersSwagger, &validatorOptions))

		// Add JWT authentication if enabled
		if jwtAuth != nil {
			r.Use(jwtAuth.Middleware)
			// r.Use(middleware.RequireScope("read:users"))
			// r.Use(middleware.RequireRole("groups", "admin"))
		}

		users.HandlerFromMux(strictUsersServer, r)
	})
}
