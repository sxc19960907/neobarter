# Logging Guidelines

> Logging conventions for the NeoBarter backend.

---

## Library

- **Development**: GORM's built-in logger (`logger.Info` mode) for SQL queries
- **Application**: Go standard `log` package (current state)
- **Future**: Structured logging with `zap` or `zerolog` (not yet implemented)

## Log Levels

| Level | When to use |
|-------|-------------|
| `log.Printf` | Server startup, configuration loaded, important state changes |
| `log.Fatalf` | Unrecoverable startup errors (DB connection failed, config missing) |
| `fmt.Printf` | Development-only debug output (e.g. mock SMS codes) |

## GORM SQL Logging

```go
// Debug mode: log all SQL queries
gormLogger = logger.Default.LogMode(logger.Info)

// Release mode: silent
gormLogger = logger.Default.LogMode(logger.Silent)
```

Controlled by `server.mode` in config.yaml.

## What to Log

- Server startup and port binding
- Database connection success/failure
- Migration execution results
- Authentication failures (for security monitoring)
- Transaction settlement operations (audit trail)
- WebSocket connection/disconnection events

## What NOT to Log

- ❌ User passwords or tokens
- ❌ Full request bodies containing PII (phone numbers, ID cards)
- ❌ SMS verification codes in production
- ❌ Database connection strings with passwords
- ❌ JWT secret values

## Future Direction

When the project grows, migrate to structured logging:
- Use `zap` for JSON-formatted logs
- Add request ID middleware for trace correlation
- Log to stdout (container-friendly), collect with Prometheus/Loki
