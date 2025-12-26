# Load .env file if it exists (optional, won't fail if missing)
-include .env
export

MIGRATIONS_DIR = migrations
CONTAINER_CMD ?= podman

# Database connection defaults (can be overridden by .env file or environment)
PG_HOST ?= localhost
PG_PORT ?= 5432
PG_DB ?= orbis
PG_USER ?= user
PG_PASSWORD ?= password
PG_SSLMODE ?= disable

# Construct DB_URL from environment variables
DB_URL = postgres://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(PG_DB)?sslmode=$(PG_SSLMODE)

download:
	@echo "Download go.mod dependencies"
	@go mod download
 
install-tools: download
	@echo Installing tools from tools.go
	@cat tools/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go get -tool %

# Create a new migration
# Usage: make new NAME=your_migration_name
new:
ifndef NAME
	$(error NAME is not set. Usage: make new NAME=your_migration_name)
endif
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(NAME)

up:
	@echo "Running migrations with DB_URL: postgres://$(PG_USER)@$(PG_HOST):$(PG_PORT)/$(PG_DB)"
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

down:
	@echo "Rolling back migrations with DB_URL: postgres://$(PG_USER)@$(PG_HOST):$(PG_PORT)/$(PG_DB)"
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down

start-dev-db:
	$(CONTAINER_CMD) compose -f docker-compose.dev.yaml up -d

stop-dev-db:
	$(CONTAINER_CMD) compose -f docker-compose.dev.yaml down

psql:
	$(CONTAINER_CMD) exec -it postgres psql -U $(PG_USER) -d $(PG_DB)

sqlc-gen-code:
	@echo "→ Generating database go code from SQL queries..."
	@go tool sqlc generate

api-gen-code:
	@echo "→ Generating server code from OpenAPI specs..."
	@for spec in docs/*.openapi.yaml; do \
		name=$$(basename "$$spec" .openapi.yaml); \
		echo "  - Generating from $$spec → internal/api/$$name/"; \
		mkdir -p "internal/api/$$name"; \
		go tool oapi-codegen \
			-package "$$name" \
			-generate chi-server,strict-server,models,embedded-spec \
			-o "internal/api/$$name/openapi.gen.go" \
			"$$spec"; \
	done

generate: sqlc-gen-code api-gen-code
	@echo "→ Code generation complete."

build-swagger-docs:
	$(CONTAINER_CMD) run -p 8000:8080 -e SWAGGER_JSON=/docs/openapi.yaml -v $(shell pwd)/docs:/docs swaggerapi/swagger-ui

lint:
	@echo "→ Running golangci-lint..."
	@golangci-lint run

lint-fix:
	@echo "→ Running golangci-lint with auto-fix..."
	@golangci-lint run --fix

.PHONY: new up down start-dev-db psql api-gen-code sqlc-gen-code build-swagger-docs generate download install-tools lint lint-fix

