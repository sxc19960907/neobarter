# WebSocket 按会话参与者精确推送

## Goal

修复 `Hub.SendToConversation` 当前"广播给所有在线用户"的实现，改为只推送给该会话的实际参与者（排除发送者），消除消息串台和隐私泄露风险。

## Background / 问题

`internal/ws/hub.go` 的 `SendToConversation` 注释里写着"简化实现：实际应查询会话参与者"，当前是 `for uid := range h.clients` 广播给所有在线用户。后果：用户A和B的私聊消息会被推送给在线的用户C。

## 设计约束

- ws 是基础设施层，**不应反向依赖 repository/service**（避免分层倒置）。
- 方案：Hub 只提供"按 userID 列表精确投递"的能力（`SendToUsers`）；由 MessageService 查出参与者后调用。
- 会话参与者信息来自 `conversation_participants` 表。

## Requirements

- 新增 `Hub.SendToUsers(userIDs []int64, eventType string, payload interface{})`——只推给指定在线用户。
- 新增 repository 方法查会话参与者 userID 列表（排除指定用户）。
- MessageService.Send：发消息后查出"除发送者外的参与者"，调 `SendToUsers` 精确推送。
- 删除/废弃旧的 `SendToConversation` 广播逻辑。
- 推送消息结构保持 `{type:"new_message", data:<message>}`，前端已按此消费。
- 清理 Hub 里未使用的 `broadcast` channel（如确认无用）。

## Acceptance Criteria

- [ ] 私聊消息只推给会话双方，不泄露给无关在线用户
- [ ] ws 包不 import repository/service（分层不倒置）
- [ ] 发送者自己不收到自己刚发的消息推送
- [ ] 离线参与者不报错（静默跳过）
- [ ] Hub 推送逻辑有单元测试覆盖
- [ ] 编译 + 测试 + CI 通过

## Constraints

- 不改前端消费协议（type=new_message）
- 不引入新的并发 bug（Hub 的 map 读写加锁）
