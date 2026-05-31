# Error Handling

> How errors are handled in the NeoBarter backend.

---

## Error Types

### Service-level errors (`internal/service/errors.go`)

```go
var (
    ErrForbidden = errors.New("无权执行此操作")
    ErrNotFound  = errors.New("资源不存在")
)
```

Domain-specific errors are returned as plain `errors.New("message")` with user-facing Chinese messages.

### GORM errors

- `gorm.ErrRecordNotFound` — checked in service layer to distinguish "not found" from other DB errors

## Error Propagation

```
Repository → returns error (GORM errors, nil on success)
Service    → wraps/translates errors, returns user-facing messages
Handler    → maps errors to HTTP status codes via response helpers
```

## API Error Response Format

All errors use the unified response structure from `pkg/response`:

```json
{
  "code": 40000,
  "message": "手机号格式错误",
  "data": null
}
```

### HTTP status → error code mapping

| HTTP Status | Code | Usage |
|-------------|------|-------|
| 400 | 40000 | Bad request / validation error |
| 401 | 40100 | Unauthorized (missing/expired token) |
| 403 | 40300 | Forbidden (no permission) |
| 404 | 40400 | Resource not found |
| 429 | 42900 | Rate limited |
| 500 | 50000 | Internal server error |

### Response helper usage in handlers

```go
// Validation error
response.BadRequest(c, "参数错误")

// Auth error
response.Unauthorized(c, "认证已过期，请重新登录")

// Permission error
response.Forbidden(c, "无权修改此物品")

// Not found
response.NotFound(c, "物品不存在")

// Server error (hide internal details from client)
response.ServerError(c, "获取物品列表失败")
```

## Conventions

- Error messages are in Chinese (user-facing)
- Never expose internal error details (stack traces, SQL) to the client
- Service layer errors carry the user-facing message; handler just passes it through
- Use `errors.Is()` for sentinel error comparison
- Use `fmt.Errorf("context: %w", err)` to wrap errors with context

## Forbidden Patterns

- ❌ Returning raw GORM errors to the client
- ❌ Panicking on recoverable errors
- ❌ Ignoring errors with `_` (except fire-and-forget operations like view count increment)
- ❌ English error messages in API responses (keep consistent Chinese UX)
