# 实名/企业认证（F-U-011 / F-U-012）

## Goal

补上认证提交接口：字段已在 User 模型，但没有提交认证的接口。

## 设计 & 范围声明

**MVP 范围**：提交认证资料 → 存储 → 直接置为已认证。
**不做**：人工审核队列、三方实名核验（公安/芝麻信用等）——依赖运营后台和外部付费接口，超出当前范围。实现和文档如实标注，不假装做审核。

后续接三方核验时，只需把"提交即通过"改为"提交→pending→回调置 verified"。

## Requirements

### 个人实名认证
- POST /users/me/verify-realname：提交 real_name + id_card
- 身份证 18 位格式校验
- 存储后置 real_name_verified=true
- 已认证不可重复提交

### 企业认证
- POST /users/me/verify-enterprise：提交 enterprise_name + license_url
- 仅 user_type=enterprise 可提交
- 存储后置 enterprise_verified=true

### 隐私
- id_card 等敏感字段 json:"-" 不返回（已有）

## Acceptance Criteria

- [ ] 个人可提交实名认证，状态变已认证
- [ ] 身份证格式非法被拒
- [ ] 重复认证被拒
- [ ] 企业用户可提交企业认证，个人用户提交被拒
- [ ] 敏感字段不泄露
- [ ] service 单测覆盖
- [ ] 前端个人中心认证入口
- [ ] 编译 + 测试 + CI 通过

## Constraints

- 不引入三方核验/付费接口
- 不破坏现有 User 结构
