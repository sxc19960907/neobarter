# 集成 Swagger API 文档自动生成

## Goal

为 Go 后端集成 Swagger（OpenAPI），通过代码注解自动生成交互式 API 文档，开发者可在浏览器中查看所有接口、参数、响应结构并在线调试。

## Requirements

### 技术方案
- `swaggo/swag`：解析 Go 注解生成 OpenAPI 2.0 spec（docs/swagger.json|yaml）
- `swaggo/gin-swagger` + `swaggo/files`：提供交互式 Swagger UI
- UI 路由：`GET /swagger/index.html`

### 注解覆盖范围
- 全局 API 信息（标题、版本、描述、base path、安全定义 Bearer Token）
- 所有 handler 方法添加 swag 注解：
  - `@Summary` / `@Description` / `@Tags`
  - `@Param`（路径参数、query、body）
  - `@Success` / `@Failure`（响应结构）
  - `@Security`（需要鉴权的接口）
- 请求/响应 DTO 结构体可被 swag 识别

### 工程化
- `make docs` 或脚本：运行 `swag init` 重新生成文档
- 生成的 `docs/` 包加入版本控制（CI 不强制重新生成，但提供命令）
- release 模式下可选择关闭 Swagger UI（安全考虑）

## Acceptance Criteria

- [ ] `swag init` 能成功生成 docs（无解析错误）
- [ ] 启动服务后访问 `/swagger/index.html` 能看到完整 API 文档
- [ ] 认证、用户、物品、交易、钱包、消息、搜索、上传等模块接口都有文档
- [ ] Bearer Token 鉴权在 UI 中可配置并能在线调试
- [ ] 响应结构（统一 Response 格式）正确展示
- [ ] go build / go test / CI 通过
- [ ] README 说明如何查看和重新生成文档

## Constraints

- swag 注解写在 handler 方法上方
- 不破坏现有 API 行为
- docs 目录在 server 内（server/docs/）
