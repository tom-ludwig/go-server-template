# Go-Server Template

A template for a Go server with OpenAPI, PostgreSQL, and database migrations.

[api-fiddle](https://api.api-fiddle.com/v1/public/resources/oas_api_3_1/tommys-organization-wzn/yellow-sheep-yssp)

# Requiremnets

- sqlc: `brew install sqlc`
- golang-migrate: `brew install golang-migrate`

## Migrations

New: `make new NAME=xxx`

Up: `make up`

Down: `make down`

Edit `DB_URL` in the Makefile. Migrations live in `migrations/`.

## DB Schema

Use sqlc to generate go boilerplate code: `sqlc generate`

## API Spec

OpenAPI is used to define the API spec. The spec is located in `docs/openapi.yaml`.
Generate the server boilerplate code using `oapi-codegen`:

```bash
make api-gen-code
```

which runs:

```bash
go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=oapi-codegen.yaml docs/openapi.yaml
```

All handlers have to implemented by the `Server` struct in `internal/handler/handlers.go`.
To split the handlers into multiple files, create new files in the same package and add either directly implement the methods or create new structs that will be embedded by the `Server` and implement the methods there.

Example:

```go
type ClientHandler struct { }
func (s *ClientHandler) GetClient(_ context.Context, _ api.GetClientRequestObject) (api.GetClientResponseObject, error) {
	return api.GetClient200JSONResponse{
		Name: "Client Name",
	}, nil
}

type AdminHandler struct {
    DB *dbpool.Pool
}
func (s *AdminHandler) GetAdmins(_ context.Context, _ api.GetAdminRequestObject) (api.GetAdminResponseObject, error) {
    // use s.DB to query the database
	return api.GetAdmin200JSONResponse{
		Name: "Client Name",
	}, nil
}

type Server struct {
    ClientHandler
    AdminHandler
}

// compile-time check
var _ api.StrictServerInterface = (*Server)(nil)

func NewServer(db *dbpool.Pool) *Server {
    return &Server{
        ClientHandler: ClientHandler{},
        AdminHandler: AdminHandler{DB: db},
    }
}
```

## Get Started

- Install the Requirements
- Copy the example `.env.example` file and adjust to your use case: `cp .env.example .env`
- Perpare local development (start postgres): `make start-dev-db`
- Run migrations: `make up`
- Generate the API and DB code: `make generate`
- Start the server: `make run`
- Test the server:

```bash
curl -X POST "http://localhost:8080/user" -H "Content-Type: application/json" -d '{"first_name": "John", "last_name": "Doe", "email": "john.doe@example.com"}'
curl "http://localhost:8080/user?user_id=<uuid>"
```
