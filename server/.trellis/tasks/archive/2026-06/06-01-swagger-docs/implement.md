# Swagger Docs - Implementation Plan

## Execution Order

### Phase 1: Setup
- [ ] Install swag CLI (go install github.com/swaggo/swag/cmd/swag)
- [ ] Add deps: swaggo/gin-swagger, swaggo/files, swaggo/swag
- [ ] Add global API annotations to cmd/server/main.go (title, version, securityDefinitions)

### Phase 2: Annotate handlers
- [ ] auth_handler.go — SendCode, Login
- [ ] user_handler.go — GetMe, UpdateMe, GetUser, address CRUD, wallet
- [ ] item_handler.go — Create, List, Get, Update, Delete, UpdateStatus, ListCategories
- [ ] trade_handler.go — Create, List, Get, Accept, Reject, Complete, Cancel
- [ ] message_handler.go — conversations, messages, send, markRead
- [ ] review_handler.go — Create, ListByUser, notifications
- [ ] search_handler.go — Search, Suggest
- [ ] upload_handler.go — UploadImage

Note: define a shared response wrapper type for swag (response.Response) so @Success can reference it.

### Phase 3: Generate + Wire UI
- [ ] Run swag init -g cmd/server/main.go -o docs
- [ ] Import docs package + gin-swagger in main.go
- [ ] Register /swagger/*any route (gate behind non-release mode optionally)
- [ ] go build verify

### Phase 4: Tooling + Docs
- [ ] Add Makefile with `docs` target (swag init)
- [ ] Update GETTING_STARTED / README with how to view + regenerate
- [ ] Verify CI passes (docs committed, build works)

## Validation Commands

```bash
cd server
swag init -g cmd/server/main.go -o docs --parseDependency
go build ./...
CGO_ENABLED=1 go test ./...
# Manual: run server, open http://localhost:8080/swagger/index.html
```

## Notes

- swag CLI needs PATH to include $(go env GOPATH)/bin
- Use --parseDependency so swag resolves types from other packages (model, response)
- Annotate with example response type response.Response{data=model.X}
- docs/docs.go is generated — commit it so CI build works without running swag
