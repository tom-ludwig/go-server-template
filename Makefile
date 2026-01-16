# Load .env file if it exists (optional, won't fail if missing)
-include .env
export

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
SCHEMA_FILE ?= migrations/schema.sql
OLD_SCHEMA  ?= migrations/current-schema.sql

TERM_BOLD  := $(shell tput bold)
TERM_RESET := $(shell tput sgr0)
TERM_GREEN := $(shell tput setaf 2)
TERM_RED   := $(shell tput setaf 1)

download:
	@echo "Download go.mod dependencies"
	@go mod download
 
install-tools: download
	@echo Installing tools from tools.go
	@cat tools/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go get -tool %

plan: ## Preview changes (Dry Run)
	@echo "$(TERM_BOLD)» Checking for schema changes...$(TERM_RESET)"
	@PGSSLMODE=disable PGPASSWORD=$(PG_PASSWORD) LOG_LEVEL=INFO go tool psqldef -h $(PG_HOST) -p $(PG_PORT) -U $(PG_USER) $(PG_DB) --dry-run < $(SCHEMA_FILE)

apply: ## Apply changes to DB
	@echo "$(TERM_BOLD)» Applying schema changes...$(TERM_RESET)"
	@PGSSLMODE=disable PGPASSWORD=$(PG_PASSWORD) LOG_LEVEL=INFO go tool psqldef --config sqldef.yaml -h $(PG_HOST) -p $(PG_PORT) -U $(PG_USER) $(PG_DB) < $(SCHEMA_FILE) && \
	echo "$(TERM_GREEN)✔ Schema applied successfully.$(TERM_RESET)"

diff: ## Compare two local files 
	@echo "$(TERM_BOLD)» Comparing '$(OLD_SCHEMA)' and '$(SCHEMA_FILE)'...$(TERM_RESET)"
	@if [ ! -f "$(OLD_SCHEMA)" ]; then \
		echo "$(TERM_RED)✘ Error:$(TERM_RESET) File '$(OLD_SCHEMA)' not found."; \
		exit 1; \
	fi
	@if [ ! -f "$(SCHEMA_FILE)" ]; then \
		echo "$(TERM_RED)✘ Error:$(TERM_RESET) File '$(SCHEMA_FILE)' not found."; \
		exit 1; \
	fi
	@LOG_LEVEL=INFO go tool psqldef $(OLD_SCHEMA) < $(SCHEMA_FILE)

dump: ## Dump current DB schema to migrations/old-schema.sql (Schema Only)
	@echo "$(TERM_BOLD)» Dumping schema from DB...$(TERM_RESET)"
	@mkdir -p migrations
	@PGSSLMODE=disable PGPASSWORD=$(PG_PASSWORD) pg_dump -h $(PG_HOST) -p $(PG_PORT) -U $(PG_USER) --schema-only --no-owner --no-privileges $(PG_DB) > migrations/pg_dump.sql
	@echo "$(TERM_GREEN)✔ Dumped to $(OLD_SCHEMA)$(TERM_RESET)"

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

.PHONY: new up down start-dev-db psql api-gen-code sqlc-gen-code build-swagger-docs generate download install-tools lint lint-fix

