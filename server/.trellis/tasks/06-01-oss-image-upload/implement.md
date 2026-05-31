# OSS Image Upload - Implementation Plan

## Execution Order

### Phase 1: Storage Abstraction (backend)
- [ ] Create `internal/pkg/storage/storage.go` — Storage interface + helpers (filename gen, type/size validation)
- [ ] Create `internal/pkg/storage/local.go` — LocalProvider (saves to ./uploads, serves via static route)
- [ ] Create `internal/pkg/storage/oss.go` — OSSProvider (aliyun-oss-go-sdk)
- [ ] Factory function to pick provider by config

### Phase 2: Upload API (backend)
- [ ] Create `internal/service/upload_service.go` — wraps Storage, validation logic
- [ ] Create `internal/handler/upload_handler.go` — multipart parsing, POST /upload/image
- [ ] Wire into main.go: init storage, register route with auth + rate limit
- [ ] Serve static /uploads when using local provider
- [ ] Add storage config to config.example.yaml

### Phase 3: Tests (backend)
- [ ] storage_test.go — filename generation, type/size validation, local upload roundtrip
- [ ] Verify go vet + go test pass

### Phase 4: Frontend Integration
- [ ] Create `web/src/services/upload.ts` — upload API
- [ ] Items/Publish.tsx — real Upload with customRequest, collect URLs into form
- [ ] Items/Detail or edit — same upload component for editing
- [ ] Profile avatar upload
- [ ] Loading/error states, preview

### Phase 5: Verify
- [ ] cd server && CGO_ENABLED=1 go test ./...
- [ ] cd web && npx tsc --noEmit && npm run build
- [ ] gitignore uploads/

## Validation Commands

```bash
cd server && go vet ./... && CGO_ENABLED=1 go test ./...
cd web && npm run lint && npx tsc --noEmit && npm run build
```

## Notes

- LocalProvider serves files at /uploads/* — Gin Static route
- Filename: {YYYYMM}/{uuid}.{ext} to avoid collisions and path traversal
- Validate both extension and content-type sniff (http.DetectContentType)
- Frontend: use antd Upload `customRequest` to call our API instead of default behavior
