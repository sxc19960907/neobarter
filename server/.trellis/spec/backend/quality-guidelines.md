# Quality Guidelines

> Code quality standards for the NeoBarter backend.

---

## Code Style

- Follow standard Go conventions (`gofmt`, `go vet`)
- Use `golangci-lint` for static analysis (to be configured)
- Exported types and functions must have doc comments
- Keep functions under 50 lines where possible

## Required Patterns

### Constructor pattern for all layers

```go
func NewItemService(itemRepo *repository.ItemRepository) *ItemService {
    return &ItemService{itemRepo: itemRepo}
}
```

### Unified response in all handlers

```go
response.Success(c, data)
response.SuccessPage(c, list, total, page, pageSize)
response.BadRequest(c, "message")
```

### Permission checks in service layer

```go
func (s *ItemService) Delete(id int64, userID int64) error {
    item, err := s.itemRepo.GetByID(id)
    if err != nil { return err }
    if item.UserID != userID { return ErrForbidden }
    // ...
}
```

### Request struct with binding tags

```go
type CreateItemReq struct {
    Title     string `json:"title" binding:"required"`
    Condition string `json:"condition" binding:"required"`
}
```

## Forbidden Patterns

- ❌ Global mutable state (use dependency injection)
- ❌ `init()` functions with side effects
- ❌ Naked goroutines without error handling
- ❌ `interface{}` when a concrete type is known
- ❌ Ignoring context cancellation in long operations
- ❌ Hardcoded configuration values (use config.yaml)
- ❌ `fmt.Sprintf` for SQL queries (SQL injection risk)

## Testing Requirements

- Unit tests for service layer business logic (especially wallet transfers)
- Integration tests for critical flows (auth, trade completion)
- Test file naming: `<file>_test.go` in the same package
- Use table-driven tests for multiple scenarios
- Mock external dependencies (SMS, OSS) via interfaces

### Test conventions (established)

- **Test framework**: `github.com/stretchr/testify` (`assert` + `require`)
- **In-memory DB**: service-layer tests use `gorm.io/driver/sqlite` with `:memory:` — fast, no external dependency. Run with `CGO_ENABLED=1`.
- **decimal comparison**: NEVER use `assert.Equal` on `decimal.Decimal` — internal `exp` differs between `NewFromFloat(100)` (exp=2) and DB-loaded values (exp=0). Use `expected.Equal(actual)` wrapped in `assert.True`.
- **Optional deps as nil**: services with optional dependencies (e.g. `ItemService`'s MQ publisher) accept `nil` in tests to skip side effects.
- **Test setup helpers**: `setupTestDB(t)` migrates only the tables under test; `createTestUser(t, db, phone)` for fixtures.

Run tests:
```bash
cd server && CGO_ENABLED=1 go test ./...
```

## API Design Standards

- RESTful resource naming: `/v1/items`, `/v1/trades/:id/accept`
- Consistent pagination: `?page=1&page_size=20`
- All list endpoints return `{ list, total, page, page_size }`
- All mutations require JWT auth (except `/auth/send-code` and `/auth/login`)
- Use HTTP verbs correctly: GET (read), POST (create), PUT (update), DELETE (remove)

### Swagger / OpenAPI docs

- Every handler method MUST carry swag annotations (`@Summary`, `@Tags`, `@Param`, `@Success`, `@Security BearerAuth`, `@Router`). See `auth_handler.go` for the canonical pattern.
- `@Router` paths omit the `/v1` base path prefix (e.g. `/items/{id}` not `/v1/items/{id}`).
- Custom field types that swag can't introspect need a `swaggertype` struct tag:
  - `decimal.Decimal` → `swaggertype:"string"`
  - `pq.StringArray` → `swaggertype:"array,string"`
- **Generation gotcha**: `swag init --parseDependency` crashes on Go 1.26 stdlib (`reflect.Value.Len on zero Value`). Use `--parseInternal` ONLY (resolves internal/ packages) plus the `swaggertype` tags above for third-party types. Command: `make docs`.
- The generated `docs/` package IS committed (so CI `go build` works without running swag).
- Swagger UI served at `/swagger/index.html`, gated behind non-`release` mode.

## Docker / Deployment

- Each service has its own Dockerfile + `.dockerignore`; build context is the service dir (server/ web/ ai-service/).
- **Never `COPY` from outside the build context** (e.g. `COPY ../deploy/...` is illegal in Docker). Files needed in an image must live inside that service's dir — `web/nginx.conf` is the in-image nginx config (the `deploy/nginx/` one is for compose/host use).
- **server Dockerfile Go version must match `go.mod`** (currently `golang:1.26-alpine`). A stale base image (1.21) fails the build.
- server Dockerfile sets `ENV GOPROXY=https://goproxy.cn,direct` so in-container `go mod download` works from CN networks; `direct` fallback keeps overseas CI working.
- **CN base-image pulls**: configure OrbStack/Docker daemon `registry-mirrors` (`~/.orbstack/config/docker.json`) when docker.io times out, then `orb restart docker`.
- CI builds/pushes images to `ghcr.io/<owner>/neobarter-{server,web,ai}` via `docker/build-push-action`; PR builds only, push-to-main pushes. Auth via `GITHUB_TOKEN` (needs `packages: write` job permission).

### Config & full-stack compose

- **Config precedence**: env vars (`NEOBARTER_` prefix, nested keys joined by `_`, e.g. `NEOBARTER_DATABASE_HOST`) > `config.yaml` > built-in defaults (`config.setDefaults()`). The config file is OPTIONAL — containers run on env vars alone.
- viper `AutomaticEnv` only binds keys that have a registered default; that's why `setDefaults()` registers every connection key. **Slices/special types (`elasticsearch.addresses`, `rabbitmq.url`) need manual `os.Getenv` post-processing** — AutomaticEnv doesn't split comma lists into `[]string`.
- server image builds three binaries (`server`/`migrate`/`consumer`) into `/app`; one image (`neobarter-server:local`) is reused by all three compose services via different `command:`.
- compose orchestration order: infra `service_healthy` → `migrate` (one-shot, `service_completed_successfully`) → `server`/`consumer`. Use a YAML anchor (`&server-env`/`*server-env`) to DRY the connection env across the three.
- **ES index creation must degrade gracefully**: the IK analyzer (`ik_smart`) requires the `analysis-ik` plugin which the stock ES image lacks. `EnsureIndex` tries the IK mapping, falls back to a built-in `standard`-analyzer mapping (`ItemMappingFallback`) on failure — otherwise consumer crash-loops. CJK precision suffers (single-char tokenization) until the plugin is installed.
- Swagger UI is gated behind non-`release` mode, so it 404s under compose (which sets `NEOBARTER_SERVER_MODE=release`) — that's intended, not a bug.

## Security Checklist

- [ ] All user input validated via `binding:"required"` or manual checks
- [ ] Ownership verified before mutation (user can only modify own resources)
- [ ] Rate limiting on auth endpoints
- [ ] JWT expiration enforced
- [ ] Sensitive fields excluded from JSON response (`json:"-"`)
- [ ] No SQL injection vectors (parameterized queries only)
