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

## API Design Standards

- RESTful resource naming: `/v1/items`, `/v1/trades/:id/accept`
- Consistent pagination: `?page=1&page_size=20`
- All list endpoints return `{ list, total, page, page_size }`
- All mutations require JWT auth (except `/auth/send-code` and `/auth/login`)
- Use HTTP verbs correctly: GET (read), POST (create), PUT (update), DELETE (remove)

## Security Checklist

- [ ] All user input validated via `binding:"required"` or manual checks
- [ ] Ownership verified before mutation (user can only modify own resources)
- [ ] Rate limiting on auth endpoints
- [ ] JWT expiration enforced
- [ ] Sensitive fields excluded from JSON response (`json:"-"`)
- [ ] No SQL injection vectors (parameterized queries only)
