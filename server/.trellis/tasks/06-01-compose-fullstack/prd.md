# 打通 docker-compose 全栈一键启动

## Goal

让 `docker compose up` 一条命令拉起整个系统（数据库 + 缓存 + 搜索 + 队列 + 后端 + 消费者 + AI + 前端），各服务正确互联，自动迁移，开箱即用。

## Background / 已知问题

1. server 容器读 `config.yaml`（被 gitignore 且写的是 localhost），容器内 localhost 连不到其他容器。
2. config 加载不支持环境变量覆盖——无法在不改镜像的前提下适配容器网络。
3. 数据库迁移没有在 server 启动前自动执行。
4. consumer（ES 索引同步）未作为服务编排。
5. ES IK 分词器插件目录为空，索引创建会失败（需处理或降级）。

## Requirements

### 配置支持环境变量覆盖
- viper 绑定环境变量（前缀 NEOBARTER_，`.` → `_`），如 `NEOBARTER_DATABASE_HOST=postgres`
- config 文件变为可选：文件不存在时仅靠环境变量也能启动
- 不破坏现有本地 `go run` + config.yaml 的开发方式

### compose 编排
- server 容器通过环境变量连服务名（postgres/redis/elasticsearch/rabbitmq）
- 启动顺序：基础设施 healthy → migrate（一次性）→ server / consumer
- migrate 作为一次性 job（run once 后退出）
- consumer 作为常驻服务
- ai-service 已有，确认环境变量连 postgres
- web 通过 nginx 反代 server（已配 upstream server:8080）
- server 镜像内置默认 config 或纯环境变量启动

### 验证
- `docker compose up -d` 后所有服务 healthy
- 通过 web(3000) 或 server(8080) 真实走一遍核心流程
- 搜索走 ES（若 IK 不可用，记录降级）

## Acceptance Criteria

- [ ] config 支持环境变量覆盖，文件可选
- [ ] `docker compose up -d --build` 全部服务起来且 server 健康
- [ ] 迁移自动执行（categories 种子存在）
- [ ] 通过容器化的 server 完成一次注册→发布→交易→结算
- [ ] consumer 连上 ES + MQ（或合理降级）
- [ ] 本地 go run 方式仍可用（不回归）
- [ ] 编译 + 测试 + CI 通过

## Constraints

- 不引入额外付费服务
- server 镜像保持单一，靠环境变量适配不同环境
- config.yaml 仍 gitignore（容器不依赖它）
