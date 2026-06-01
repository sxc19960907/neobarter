# CI Docker Build - Implementation Plan

## Execution Order

### Phase 1: Fix Dockerfiles (verify each with local docker build)
- [ ] server/Dockerfile: bump golang:1.21-alpine -> golang:1.26-alpine (match go.mod)
- [ ] web/Dockerfile: move nginx conf into web build context
      - Create web/nginx.conf (copy of deploy/nginx but adjusted for in-image use)
      - COPY web/nginx.conf instead of ../deploy/nginx/nginx.conf
- [ ] ai-service/Dockerfile: verify builds (likely OK)
- [ ] Local verify: docker build each image successfully

### Phase 2: CI docker job
- [ ] Add docker build jobs to .github/workflows/ci.yml (or new docker.yml)
- [ ] Use docker/setup-buildx-action + docker/build-push-action
- [ ] Login to ghcr.io with GITHUB_TOKEN (only on push to main)
- [ ] push=false for PR, push=true for main push
- [ ] Path filter: reuse changes job outputs (backend->server image, frontend->web image, ai->ai image)
- [ ] gha cache for layers
- [ ] permissions: packages: write

### Phase 3: Verify
- [ ] Local: all 3 images build
- [ ] Push, watch CI, confirm images appear in ghcr
- [ ] Update README/GETTING_STARTED with image pull instructions

## Validation Commands

```bash
# Local build verification
cd server && docker build -t neobarter-server:test .
cd web && docker build -t neobarter-web:test .
cd ai-service && docker build -t neobarter-ai:test .

# Verify web nginx config baked in
docker run --rm neobarter-web:test cat /etc/nginx/conf.d/default.conf
```

## Notes

- ghcr image names must be lowercase: ghcr.io/sxc19960907/neobarter-server
- docker/build-push-action context per service dir
- Need `permissions: packages: write` and `contents: read` at job level
- Cache: cache-from type=gha, cache-to type=gha,mode=max
