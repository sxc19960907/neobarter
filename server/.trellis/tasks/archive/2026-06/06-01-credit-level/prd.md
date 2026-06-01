# 信用等级体系（F-R-010）

## Goal

基于已有的 credit_score 计算用户信用等级（普通/银牌/金牌/钻石），在用户信息接口返回。

## 设计

信用等级是 credit_score 的派生值，**不单独存储字段**（避免冗余/不一致），计算函数 + 序列化时带出。

等级阈值（基于初始 100 分）：
- 普通 normal：< 120
- 银牌 silver：120 - 199
- 金牌 gold：200 - 499
- 钻石 diamond：>= 500

## Requirements

- model 层提供 `CreditLevel(score int) string` 计算函数 + 常量
- User 序列化时附带 `credit_level` 字段（GORM 计算字段或 MarshalJSON）
- GetMe / GetUser 返回的用户信息包含等级
- 前端个人中心展示信用等级

## Acceptance Criteria

- [ ] 各分数段返回正确等级（边界值 120/200/500）
- [ ] /users/me 和 /users/:id 返回 credit_level
- [ ] 计算函数有单测（含边界）
- [ ] 前端展示等级
- [ ] 编译 + 测试 + CI 通过

## Constraints

- 不新增数据库字段（纯派生）
- 不破坏现有 User JSON 结构（新增字段）
