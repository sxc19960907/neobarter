# 集成 Elasticsearch 搜索服务

## Goal

为物品模块接入 Elasticsearch，实现高性能全文搜索、分词、高亮、自动补全，替代当前 PostgreSQL ILIKE 查询。

## Requirements

### 核心功能
1. **物品索引同步** — 物品创建/更新/删除时自动同步到 ES 索引
2. **全文搜索** — 支持中文分词（IK 分析器），搜索标题+描述
3. **多条件筛选** — 分类、成色、估值范围、地区，与全文搜索组合
4. **搜索结果排序** — 相关度、最新、最热、估值
5. **搜索建议/自动补全** — 输入时提示热门搜索词
6. **搜索高亮** — 命中关键词高亮返回

### 技术要求
- ES 8.x，使用官方 Go 客户端 `github.com/elastic/go-elasticsearch/v8`
- 索引名：`neobarter_items`
- 通过 RabbitMQ 异步同步（物品变更 → 发消息 → consumer 写 ES）
- 搜索 API 独立于现有物品列表 API（新增 `/v1/search/items`）

### 前端对接
- 首页搜索框调用新的搜索 API
- 搜索结果页展示高亮片段

## Acceptance Criteria

- [ ] 物品创建后 5 秒内可被搜索到
- [ ] 中文分词正确（"苹果手机" 能搜到 "iPhone 苹果"）
- [ ] 多条件组合筛选正常工作
- [ ] 搜索响应时间 < 200ms（1000条数据量级）
- [ ] 物品删除/下架后从搜索结果中移除
- [ ] 前端搜索页面正常展示高亮结果

## Notes

- 开发环境 ES 已在 docker-compose 中配置（端口 9200）
- IK 分词器需要作为 ES 插件安装（Dockerfile 或 docker-compose 配置）
