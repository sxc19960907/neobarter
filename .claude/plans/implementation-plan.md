# NeoBarter 实现计划

## 第一步：数据库重新设计

基于 SRS 需求和巴特币结算模型，重新设计完整 schema：

### 新增/修改的表

| 表名 | 说明 |
|------|------|
| `users` | 优化字段，加入密码哈希（预留）、状态字段 |
| `user_addresses` | 收货地址（独立表，支持多地址） |
| `wallets` | 巴特币钱包，每用户一个，记录余额 |
| `wallet_transactions` | 钱包流水（充值/消费/奖励/退款），完整审计链 |
| `categories` | 保持不变 |
| `items` | 保持基本不变 |
| `trade_requests` | 加入 `target_user_id` 冗余字段 |
| `conversations` | 新增，会话元数据 |
| `conversation_participants` | 新增，会话参与者 |
| `messages` | 重构，关联 conversation |
| `reviews` | 保持不变 |
| `notifications` | 新增，通知系统 |

### 巴特币流转设计

- 注册赠送初始巴特币 → `wallet_transactions` 记录 type=`reward`
- 物品交易时：买方扣减、卖方增加 → type=`trade`
- 后续接入 USDC：充值 type=`deposit`，提现 type=`withdraw`

---

## 第二步：后端项目骨架（Go + Gin + GORM）

```
server/
├── cmd/
│   ├── server/main.go          # HTTP 服务入口
│   └── migrate/main.go         # 数据库迁移
├── internal/
│   ├── config/config.go        # 配置加载（Viper）
│   ├── model/                  # GORM 模型定义
│   ├── repository/             # 数据访问层
│   ├── service/                # 业务逻辑层
│   ├── handler/                # HTTP 处理器
│   ├── middleware/             # JWT鉴权、限流、CORS
│   ├── ws/                     # WebSocket 消息推送
│   └── pkg/                    # 工具包（SMS、OSS、JWT）
├── config.example.yaml
├── go.mod
└── Dockerfile
```

### 模块实现顺序

1. **基础设施**：配置、数据库连接、Redis连接、中间件
2. **认证模块**：发送验证码、登录/注册、JWT签发
3. **用户模块**：个人信息CRUD、地址管理
4. **钱包模块**：余额查询、流水记录、内部转账（WalletService 抽象层）
5. **物品模块**：发布、列表、详情、搜索、上下架
6. **交易模块**：发起交换、确认/拒绝、完成交易、巴特币结算
7. **消息模块**：会话管理、消息收发、WebSocket推送
8. **评价模块**：评分、评论、信用积分计算
9. **通知模块**：站内通知、未读计数
10. **搜索服务**：Elasticsearch 索引同步、全文搜索

---

## 第三步：前端项目骨架（React + TypeScript + Ant Design）

```
web/
├── src/
│   ├── components/             # 通用组件
│   │   ├── Layout/
│   │   ├── ItemCard/
│   │   └── ChatBubble/
│   ├── pages/
│   │   ├── Home/               # 首页推荐
│   │   ├── Auth/               # 登录注册
│   │   ├── Items/              # 物品浏览/发布
│   │   ├── Trade/              # 交易管理
│   │   ├── Messages/           # 消息中心
│   │   ├── Profile/            # 个人中心
│   │   └── Wallet/             # 钱包
│   ├── services/               # API 调用封装
│   ├── stores/                 # Zustand 状态管理
│   ├── hooks/                  # 自定义 hooks
│   ├── types/                  # TypeScript 类型定义
│   ├── utils/                  # 工具函数
│   ├── App.tsx
│   └── main.tsx
├── package.json
├── vite.config.ts
├── tsconfig.json
└── Dockerfile
```

---

## 第四步：部署配置

```
deploy/
├── docker-compose.yml          # 本地开发全套环境
├── nginx/nginx.conf            # 反向代理配置
└── k8s/                        # K8s 部署清单（后续）
```

docker-compose 包含：PostgreSQL、Redis、Elasticsearch、RabbitMQ、后端服务、前端服务

---

## 实施节奏

由于代码量较大，按以下批次推进：

**批次 1**：数据库 schema + 后端骨架 + 配置 + docker-compose
**批次 2**：认证模块 + 用户模块 + 钱包模块（后端）
**批次 3**：物品模块 + 交易模块（后端）
**批次 4**：消息模块 + 评价模块 + 通知模块（后端）
**批次 5**：前端项目骨架 + 认证/用户页面
**批次 6**：前端物品/交易/消息/钱包页面
**批次 7**：Elasticsearch 搜索 + AI 推荐服务骨架

先从批次 1 开始。
