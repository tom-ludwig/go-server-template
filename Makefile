# Load .env file if it exists (optional, won't fail if missing)
-include .env
export

MIGRATIONS_DIR = migrations
SCHEMA_FILE = $(MIGRATIONS_DIR)/schema.sql
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

# Declarative schema migration using pg-schema-diff
schema-apply:
	@echo "Applying schema from $(SCHEMA_FILE) to database..."
	@go tool pg-schema-diff apply \
		--from-dsn "$(DB_URL)" \
		--to-dir "$(MIGRATIONS_DIR)" \
		--allow-hazards DELETES_DATA,INDEX_BUILD

schema-diff:
	@echo "Showing schema diff between database and $(SCHEMA_FILE)..."
	@go tool pg-schema-diff diff \
		--from-dsn "$(DB_URL)" \
		--to-dir "$(MIGRATIONS_DIR)"

schema-plan:
	@echo "Planning schema changes (dry-run)..."
	@go tool pg-schema-diff plan \
		--from-dsn "$(DB_URL)" \
		--to-dir "$(MIGRATIONS_DIR)" \
		--disable-plan-validation

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
	$(CONTAINER_CMD) run -p 8000:8080 -e SWAGGER_JSON=/docs/health.openapi.yaml -v $(shell pwd)/docs:/docs swaggerapi/swagger-ui

lint:
	@echo "→ Running golangci-lint..."
	@golangci-lint run

lint-fix:
	@echo "→ Running golangci-lint with auto-fix..."
	@golangci-lint run --fix

.PHONY: new up down schema-apply schema-diff schema-plan start-dev-db psql api-gen-code sqlc-gen-code build-swagger-docs generate download install-tools lint lint-fix

