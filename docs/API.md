# API 接口文档

## 基础信息

- Base URL: `https://api.neobarter.com/v1`
- 认证方式: Bearer Token (JWT)
- 响应格式: JSON

## 通用响应格式

成功:
```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

错误:
```json
{
  "code": 40001,
  "message": "手机号格式错误",
  "data": null
}
```

---

## 一、认证模块

### 1.1 发送验证码
**POST** `/auth/send-code`

### 1.2 登录/注册
**POST** `/auth/login`

---

## 二、用户模块

### 2.1 获取当前用户信息
**GET** `/users/me`

### 2.2 更新用户信息
**PUT** `/users/me`

### 2.3 获取用户公开信息
**GET** `/users/:id`

---

## 三、物品模块

### 3.1 发布物品
**POST** `/items`

### 3.2 获取物品列表
**GET** `/items?page=1&page_size=20&category_id=1`

### 3.3 获取物品详情
**GET** `/items/:id`

### 3.4 更新物品信息
**PUT** `/items/:id`

### 3.5 删除物品
**DELETE** `/items/:id`

---

## 四、交易模块

### 4.1 发起交换请求
**POST** `/trades`

### 4.2 获取交易列表
**GET** `/trades?status=pending`

### 4.3 接受交换
**PUT** `/trades/:id/accept`

### 4.4 拒绝交换
**PUT** `/trades/:id/reject`

### 4.5 完成交易
**PUT** `/trades/:id/complete`

---

## 五、消息模块

### 5.1 获取会话列表
**GET** `/messages/conversations`

### 5.2 获取消息历史
**GET** `/messages/:conversationId`

### 5.3 发送消息
**POST** `/messages`

### 5.4 WebSocket 接入
**连接地址**: `wss://api.neobarter.com/ws/messages`

---

*接口版本：v1.0*  
*更新时间：2026-05-31*
