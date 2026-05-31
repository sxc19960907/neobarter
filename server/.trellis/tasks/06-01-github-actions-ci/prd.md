# 配置 GitHub Actions CI 流水线

## Goal

配置 GitHub Actions，在每次 push 和 PR 时自动运行后端、前端、AI 服务的测试、构建和静态检查，保证主分支始终可构建、测试通过。

## Requirements

### CI 任务（三个独立 job，并行运行）

1. **后端 CI（Go）**
   - `go vet` 静态检查
   - `go build ./...` 编译
   - `go test ./...`（CGO_ENABLED=1，因为 SQLite 测试需要）
   - 缓存 Go module

2. **前端 CI（Node）**
   - `npm ci` 安装依赖
   - `npm run lint` ESLint 检查
   - `tsc --noEmit` 类型检查
   - `npm run test` Vitest 单元测试
   - `npm run build` 生产构建
   - 缓存 npm

3. **AI 服务 CI（Python）**
   - `pip install -r requirements.txt`
   - 语法检查（py_compile 或 ruff）
   - 缓存 pip

### 触发条件
- push 到 main 分支
- 针对 main 的 PR
- 路径过滤：只有相关目录变更才触发对应 job（server/ → 后端 job，web/ → 前端 job，ai-service/ → AI job）

### 约束
- 使用 GitHub-hosted runner（ubuntu-latest）
- Go 1.21、Node 18、Python 3.11
- 国内 goproxy 在 CI 中不需要（GitHub runner 在海外）

## Acceptance Criteria

- [ ] `.github/workflows/ci.yml` 创建完成
- [ ] push 后三个 job 都能在 GitHub Actions 中运行
- [ ] 后端测试在 CI 中通过
- [ ] 前端 lint + 测试 + 构建在 CI 中通过
- [ ] AI 服务依赖安装 + 语法检查通过
- [ ] 路径过滤生效（改前端不触发后端 job）
- [ ] README 添加 CI 状态徽章
