MIGRATIONS_DIR = migrations
DB_URL = "postgres://user:password@localhost:5432/orbis?sslmode=disable"
CONTAINER_CMD = podman

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
	migrate -path $(MIGRATIONS_DIR) -database $(DB_URL) up

down:
	migrate -path $(MIGRATIONS_DIR) -database $(DB_URL) down

start-dev-db:
	$(CONTAINER_CMD) compose -f docker-compose.dev.yaml up

psql:
	$(CONTAINER_CMD) exec -it postgres psql -U user -d orbis

sqlc-gen-code:
	@echo "→ Generating database go code from SQL queries..."
	@go tool sqlc generate

api-gen-code:
	@echo "→ Generating server code from OpenAPI spec..."
	@go tool oapi-codegen --config=oapi-codegen.yaml docs/openapi.yaml

generate: sqlc-gen-code api-gen-code
	@echo "→ Code generation complete."

build-swagger-docs:
	$(CONTAINER_CMD) run -p 8000:8080 -e SWAGGER_JSON=/docs/openapi.yaml -v $(shell pwd)/docs:/docs swaggerapi/swagger-ui

.PHONY: new up down start-dev-db psql api-gen-code sqlc-gen-code build-swagger-docs generate download install-tools

