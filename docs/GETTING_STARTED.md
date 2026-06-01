# 开发指南

## 环境要求

- Go 1.21+
- Node.js 18+
- PostgreSQL 15+
- Redis 7+
- Docker (可选)

## 快速开始

### 1. 克隆项目

```bash
cd /Users/timmy/PycharmProjects/ppt-master/projects/neobarter
```

### 2. 后端开发

```bash
cd server
go mod download
cp config.example.yaml config.yaml
go run cmd/migrate/main.go
go run cmd/server/main.go
```

### API 文档（Swagger）

后端启动后（非 release 模式），访问交互式 API 文档：

```
http://localhost:8080/swagger/index.html
```

修改 handler 注解后，重新生成文档：

```bash
cd server
# 首次需安装 swag CLI：go install github.com/swaggo/swag/cmd/swag@v1.16.3
make docs   # 等价于 swag init -g cmd/server/main.go -o docs --parseInternal
```

> 在 Swagger UI 右上角 Authorize 中填入 `Bearer {token}` 即可在线调试需要登录的接口。

### 3. 前端开发

```bash
cd web
npm install
npm run dev
```

### 4. 使用 Docker（一键启动全栈）

```bash
docker compose up -d --build   # 构建并启动全部 9 个服务
docker compose ps              # 查看状态
docker compose logs -f server  # 跟踪后端日志
docker compose down            # 停止
```

启动顺序由 compose 编排：基础设施（postgres/redis/es/rabbitmq）healthy
→ `migrate` 自动迁移（一次性）→ `server` / `consumer` 启动。

访问入口：
- 前端：http://localhost:3000 （nginx 反代 API 到后端）
- 后端 API：http://localhost:8080/v1
- AI 服务：http://localhost:8081

**配置说明**：容器内服务通过环境变量连接（无需 config.yaml）。
配置优先级：环境变量（前缀 `NEOBARTER_`，嵌套键用下划线，如
`NEOBARTER_DATABASE_HOST=postgres`）> config.yaml > 内置默认值。
本地 `go run` 开发仍可用 config.yaml。

> 注：Swagger UI 仅在非 release 模式开放；compose 默认 release 模式，
> 如需文档调试，本地 `go run` 或设 `NEOBARTER_SERVER_MODE=debug`。
> ES 中文分词需 analysis-ik 插件，未安装时自动降级到内置 standard 分析器。

### 5. 预构建镜像（GitHub Container Registry）

每次 push 到 main，CI 会自动构建并推送镜像到 ghcr.io：

```bash
docker pull ghcr.io/sxc19960907/neobarter-server:latest
docker pull ghcr.io/sxc19960907/neobarter-web:latest
docker pull ghcr.io/sxc19960907/neobarter-ai:latest
```

> 国内构建镜像时，daemon 可配置 registry-mirrors 加速基础镜像拉取；server 镜像内已设置 `GOPROXY=https://goproxy.cn,direct`。

## 项目结构

```
neobarter/
+-- web/                    # 前端应用
|   +-- src/
|       +-- components/     # 通用组件
|       +-- pages/          # 页面
|       +-- services/       # API调用
|       +-- stores/         # 状态管理
+-- server/                 # 后端服务
|   +-- cmd/                # 入口
|   +-- internal/           # 内部实现
|       +-- handler/        # HTTP处理器
|       +-- service/        # 业务逻辑
|       +-- repository/     # 数据访问
+-- ai-service/            # AI服务
+-- deploy/                # 部署配置
+-- docs/                  # 文档
```

## 开发规范

### Git 提交规范

```
<type>(<scope>): <subject>

类型:
- feat: 新功能
- fix: 修复
- docs: 文档
- style: 格式
- refactor: 重构
- test: 测试
- chore: 构建/工具
```

---

*更新时间：2026-05-31*
