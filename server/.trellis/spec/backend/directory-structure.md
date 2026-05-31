# Directory Structure

> Module organization and file layout for the NeoBarter backend (Go + Gin).

---

## Project Layout

```
server/
├── cmd/
│   ├── server/main.go          # HTTP server entry, DI wiring, route registration
│   └── migrate/main.go         # Database migration + seed data
├── internal/
│   ├── config/config.go        # Configuration structs + Viper loader
│   ├── model/                  # GORM model definitions (one file per domain)
│   ├── repository/             # Data access layer (one file per domain)
│   ├── service/                # Business logic layer (one file per domain)
│   ├── handler/                # HTTP handlers (one file per domain)
│   ├── middleware/             # Gin middleware (auth, cors, ratelimit)
│   ├── ws/                     # WebSocket hub and client management
│   └── pkg/                    # Internal utility packages
│       ├── jwt/                # JWT token generation and parsing
│       ├── sms/                # SMS provider interface + mock
│       └── response/           # Unified JSON response helpers
├── config.example.yaml
├── go.mod
└── Dockerfile
```

## Layer Responsibilities

| Layer | Responsibility | May import |
|-------|---------------|------------|
| `model` | Struct definitions, TableName(), constants | Nothing internal |
| `repository` | Database queries, GORM operations, row locking | `model` |
| `service` | Business logic, orchestration, transactions | `repository`, other services |
| `handler` | HTTP parsing, validation, response formatting | `service`, `middleware`, `pkg/response` |
| `middleware` | Cross-cutting (auth, CORS, rate limiting) | `pkg/jwt`, `pkg/response` |

## File Naming

- Model: `internal/model/<domain>.go` (e.g. `user.go`, `wallet.go`)
- Repository: `internal/repository/<domain>_repo.go`
- Service: `internal/service/<domain>_service.go`
- Handler: `internal/handler/<domain>_handler.go`
- Shared service errors: `internal/service/errors.go`
- Common handler helpers: `internal/handler/common.go`

## Adding a New Domain

1. `internal/model/<domain>.go` — struct + `TableName()` + status constants
2. `internal/repository/<domain>_repo.go` — `New<Domain>Repository(db)` constructor
3. `internal/service/<domain>_service.go` — `New<Domain>Service(repo, ...)` constructor
4. `internal/handler/<domain>_handler.go` — `New<Domain>Handler(svc)` constructor
5. `cmd/server/main.go` — instantiate repo → svc → handler, register routes in the `v1` group
6. `cmd/migrate/main.go` — add model to `AutoMigrate()` call

## Forbidden Patterns

- ❌ Handler calling repository directly (bypass service)
- ❌ Repository importing service (circular)
- ❌ Business logic in handler (handlers only parse + respond)
- ❌ Raw SQL in service layer (all queries belong in repository)
