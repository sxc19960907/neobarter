# AI 推荐服务（Python + FastAPI）

## Goal

实现独立的 AI 微服务，提供物品智能推荐、物品估值建议、图片分类识别功能，供 Go 后端通过 HTTP 调用。

## Requirements

### 功能模块
1. **智能匹配推荐** — 基于用户浏览/交易历史 + 物品特征，推荐可能感兴趣的物品
2. **物品估值建议** — 根据物品分类、成色、历史交易数据，给出参考估值
3. **相似物品推荐** — 给定一个物品，推荐相似的可交换物品

### 技术栈
- Python 3.11+
- FastAPI（HTTP 服务）
- scikit-learn（协同过滤、TF-IDF 相似度）
- 连接 PostgreSQL 读取物品/交易数据

### API 设计
- `POST /recommend/items` — 为用户推荐物品（输入 user_id，输出 item_id 列表）
- `POST /estimate/value` — 物品估值（输入 category_id + condition + title）
- `POST /similar/items` — 相似物品（输入 item_id，输出相似 item_id 列表）

### 集成方式
- Go 后端通过 HTTP 调用 AI 服务
- AI 服务独立部署，端口 8081
- docker-compose 中加入 ai-service 容器

## Acceptance Criteria

- [ ] FastAPI 服务启动正常，健康检查通过
- [ ] 推荐 API 返回合理的物品列表（非空、不重复、不包含用户自己的物品）
- [ ] 估值 API 返回合理的价格范围
- [ ] 相似物品 API 返回相关度较高的结果
- [ ] Go 后端能成功调用 AI 服务
- [ ] docker-compose 中 ai-service 容器正常运行
