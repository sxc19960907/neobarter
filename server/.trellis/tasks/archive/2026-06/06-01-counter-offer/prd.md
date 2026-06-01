# 交易反向提议（F-T-009）

## Goal

目标用户收到交换请求后，除接受/拒绝外，可"反向提议"——要求用发起方的其他物品交换或调整巴特币，发起方再决定接受/拒绝。

## 设计

单记录扩展（不新建交易），新增状态 countered + 还价字段。

状态流转：
- pending →(B 反向提议)→ countered
- countered →(A 接受)→ accepted（套用还价条件后走原结算）
- countered →(A 拒绝)→ rejected
- countered →(A 取消)→ cancelled

新增字段（trade_requests）：
- counter_item_id *int64：B 要求 A 改用的物品（可空=只改巴特币）
- counter_coin_amount decimal：B 期望的巴特币
- counter_message string：还价留言

接受 countered 时：把 counter_* 落到生效字段（offered_item_id / barter_coin_amount），状态置 accepted，Complete 结算逻辑不变。

## Requirements

- B（target）对 pending 交易发起反向提议 → countered，通知 A
- A（initiator）接受 countered → 套用还价条件 → accepted，通知 B
- A 拒绝 countered → rejected，通知 B
- 权限：只有 target 能 counter，只有 initiator 能响应
- 状态校验：只能对 pending 反向提议，只能对 countered 响应

## Acceptance Criteria

- [ ] B 反向提议后状态 countered，A 收到通知
- [ ] A 接受后还价条件生效（巴特币/物品被替换），状态 accepted
- [ ] A 拒绝后状态 rejected
- [ ] 非法状态/越权被拦截
- [ ] service 单测覆盖（反提议/接受套用/拒绝/越权）
- [ ] 前端交易页支持反向提议入口
- [ ] 编译 + 测试 + CI 通过

## Constraints

- 不破坏现有 pending→accepted→completed 主流程
- 结算逻辑(Complete)不变，靠"接受时落地还价条件"复用
