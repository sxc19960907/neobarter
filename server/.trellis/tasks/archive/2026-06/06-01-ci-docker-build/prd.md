# CI 集成 Docker 镜像构建与推送

## Goal

在 CI 中增加 Docker 镜像构建验证，并支持推送到镜像仓库（GitHub Container Registry, ghcr.io）。同时修复现有 Dockerfile 中阻碍构建的问题。

## Background / 已知问题（必须先修）

1. **web/Dockerfile 非法 COPY**：`COPY ../deploy/nginx/nginx.conf` 引用了 build context 之外的父目录，Docker 构建会直接失败。需把 nginx 配置纳入 web 的 build context。
2. **server/Dockerfile Go 版本过低**：`go.mod` 要求 `go 1.26.1`，Dockerfile 却用 `golang:1.21-alpine`，构建会因版本不匹配失败。需升级基础镜像。
3. server 镜像构建需要 `docs/` 包（已入库，OK）。

## Requirements

### Dockerfile 修复
- 三个 Dockerfile（server / web / ai-service）都能独立 `docker build` 成功
- server：升级 Go 基础镜像匹配 go.mod
- web：解决跨 context COPY 问题，nginx 配置正确打进镜像
- 多阶段构建保持精简镜像体积

### CI 集成
- 新增 docker 构建 job
- 使用 docker/build-push-action 构建三个镜像
- PR / push 到 main：构建验证（不推送）
- push 到 main：构建并推送到 ghcr.io
- 镜像 tag：`ghcr.io/<owner>/neobarter-{server,web,ai}` + latest / commit SHA
- 路径过滤：仅对应目录变更才构建对应镜像
- 用 GITHUB_TOKEN 登录 ghcr（无需额外 secret）
- GitHub Actions cache 加速构建

## Acceptance Criteria

- [ ] 三个 Dockerfile 本地 docker build 均成功
- [ ] web 镜像内 nginx 配置正确
- [ ] CI docker job 配置完成
- [ ] push 到 main 后镜像推送到 ghcr.io
- [ ] 路径过滤生效
- [ ] 本地验证镜像可构建后再推 CI

## Constraints

- 镜像仓库用 ghcr.io（GITHUB_TOKEN 即可推送）
- 不引入额外付费服务
- 保持现有 docker-compose 本地开发流程可用
