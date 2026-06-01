# 端到端冒烟测试结果

**执行时间**: 2026-06-01
**环境**: docker-compose 起 PostgreSQL + Redis（本地），`go run` 跑 migrate + server。ES/MQ 留空走降级（不影响核心流程）。

## 验证通过的流程（11 步全绿）

1. ✅ Swagger UI `/swagger/index.html` 返回 200（首次运行时验证）
2. ✅ 公开接口 `/v1/categories` 返回种子分类数据
3. ✅ 用户A 手机验证码注册登录，JWT 签发，自动生成昵称
4. ✅ 注册赠送 100 巴特币，钱包 + reward 流水正确
5. ✅ 用户B 注册 + 发布物品（估值30）
6. ✅ 物品列表（关联发布者）+ 详情浏览量自增
7. ✅ **核心交易闭环**：A 发起交换 → B 收到通知 → B 接受 → 完成
8. ✅ **巴特币结算**（真实 PG 事务）：A 100→70，B 100→130，物品 active→traded
9. ✅ 互评 + 信用积分 100→102（好评+2），评价带评价人
10. ✅ 消息系统：自动建会话、未读计数、消息历史
11. ✅ 图片上传（真实 multipart）+ 静态 URL 访问 200
12. ✅ 安全边界：无/错 token→401，越权删除→403，重复评价拦截，错误验证码拦截

服务运行期间**无 error/panic**。

## 发现的小问题（非阻断，单独跟进）

### 问题1: /v1/categories 无需鉴权即可访问
- `ListCategories` 路由注册在 authorized 组的代码块内，但用了 `v1.GET` 而非 `authorized.GET`，导致它实际挂在了无鉴权的 v1 组上。
- 影响：分类列表本就适合公开，功能无害，但与代码组织意图不符（看起来在鉴权块内）。
- 位置：cmd/server/main.go 路由注册处。
- 处置：建议显式移到公开路由区，或改为 authorized.GET。低优先级。

## docker-compose 打包层面待改进（本次未用 compose 起 server，故未验证）

- server/web 容器互联用的是 localhost，容器内需改为服务名（postgres/redis/server）。
- server 容器需要 config.yaml（gitignore），compose 应提供容器专用配置或用环境变量覆盖。
- 这些属于"用 compose 一键起全栈"的完善，核心应用逻辑已验证无误。
