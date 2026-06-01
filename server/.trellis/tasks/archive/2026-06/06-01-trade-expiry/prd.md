# 交易超时自动过期（F-T-004）

## Goal

补上交易超时处理：发起交换时已设 expired_at（24h），但当前没有任何逻辑把超时未处理的 pending 交易置为 expired，导致它们永远停在 pending。

## Requirements

- 后台定时任务（goroutine + ticker）周期扫描：status=pending AND expired_at < now() 的交易批量置为 expired。
- 扫描周期：每 1 分钟一次。
- 过期时通知发起方。
- 任务随 server 启动，优雅停止（context 取消）。

## Acceptance Criteria

- [ ] 超时的 pending 交易被自动置为 expired
- [ ] 不影响 accepted/completed/rejected 的交易
- [ ] 定时任务随 server 启动、随关闭停止
- [ ] service 有单测覆盖过期逻辑
- [ ] 端到端：造一条 expired_at 已过的 pending 交易，扫描后变 expired
- [ ] 编译 + 测试 + CI 通过

## Constraints

- 批量 UPDATE，避免逐条
- 扫描走 status 索引
- 不阻塞主流程
