# Elasticsearch Search - Implementation Plan

## Execution Order

### Phase 1: ES Infrastructure
- [ ] Update docker-compose to add IK analyzer plugin to ES container
- [ ] Create ES client wrapper (`internal/pkg/es/client.go`)
- [ ] Define item index mapping (`internal/pkg/es/mapping.go`)
- [ ] Add ES initialization to server startup (create index if not exists)

### Phase 2: Index Sync via RabbitMQ
- [ ] Create RabbitMQ connection wrapper (`internal/pkg/mq/rabbitmq.go`)
- [ ] Define item event messages (create/update/delete)
- [ ] Publish events from ItemService on create/update/delete/status change
- [ ] Create consumer (`cmd/consumer/main.go`) that listens and syncs to ES
- [ ] Add bulk reindex command (`cmd/reindex/main.go`) for initial data load

### Phase 3: Search Service
- [ ] Create SearchRepository (`internal/repository/search_repo.go`)
  - Full-text search with IK analyzer
  - Multi-condition filter (bool query)
  - Sort options (score, created_at, view_count, estimated_value)
  - Highlight on title + description
  - Suggest/completion
- [ ] Create SearchService (`internal/service/search_service.go`)
- [ ] Create SearchHandler (`internal/handler/search_handler.go`)
- [ ] Register route: `GET /v1/search/items`

### Phase 4: Frontend Integration
- [ ] Create search service API (`web/src/services/search.ts`)
- [ ] Update Home page to use search API when keyword is present
- [ ] Add highlight rendering in search results
- [ ] Add search suggestion dropdown (debounced input)

## Validation Commands

```bash
# Backend
cd server && go build ./...

# Verify ES is running
curl http://localhost:9200/_cluster/health

# Verify index exists
curl http://localhost:9200/neobarter_items

# Frontend
cd web && npx tsc --noEmit
```
