package ws

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestClient 构造一个不带真实 websocket 连接的 client，仅用 send channel 验证投递。
func newTestClient(h *Hub, userID int64) *Client {
	return &Client{hub: h, userID: userID, send: make(chan []byte, 8)}
}

// register 直接把 client 放进 map（绕过 Run goroutine，便于同步测试）。
func (h *Hub) registerForTest(c *Client) {
	h.mu.Lock()
	h.clients[c.userID] = c
	h.mu.Unlock()
}

func recv(t *testing.T, c *Client) (WSMessage, bool) {
	select {
	case data := <-c.send:
		var m WSMessage
		require.NoError(t, json.Unmarshal(data, &m))
		return m, true
	case <-time.After(100 * time.Millisecond):
		return WSMessage{}, false
	}
}

func TestSendToUsers_OnlyTargets(t *testing.T) {
	h := NewHub()
	a := newTestClient(h, 1)
	b := newTestClient(h, 2)
	c := newTestClient(h, 3) // 无关用户
	h.registerForTest(a)
	h.registerForTest(b)
	h.registerForTest(c)

	// 推给 1 和 2（模拟会话参与者，排除了发送者）
	h.SendToUsers([]int64{1, 2}, "new_message", map[string]string{"content": "hi"})

	m1, ok1 := recv(t, a)
	assert.True(t, ok1, "用户1应收到")
	assert.Equal(t, "new_message", m1.Type)

	_, ok2 := recv(t, b)
	assert.True(t, ok2, "用户2应收到")

	// 用户3(无关)不应收到 —— 这正是修复广播 bug 的核心
	_, ok3 := recv(t, c)
	assert.False(t, ok3, "无关用户3不应收到")
}

func TestSendToUsers_OfflineSkipped(t *testing.T) {
	h := NewHub()
	a := newTestClient(h, 1)
	h.registerForTest(a)

	// 用户2离线（未注册），不应 panic
	assert.NotPanics(t, func() {
		h.SendToUsers([]int64{1, 2}, "new_message", "x")
	})

	_, ok := recv(t, a)
	assert.True(t, ok, "在线用户1仍应收到")
}

func TestSendToUsers_Empty(t *testing.T) {
	h := NewHub()
	assert.NotPanics(t, func() {
		h.SendToUsers(nil, "new_message", "x")
		h.SendToUsers([]int64{}, "new_message", "x")
	})
}

func TestSendToUser_Single(t *testing.T) {
	h := NewHub()
	a := newTestClient(h, 1)
	h.registerForTest(a)

	h.SendToUser(1, "notification", map[string]int{"count": 5})
	m, ok := recv(t, a)
	require.True(t, ok)
	assert.Equal(t, "notification", m.Type)

	// 离线用户不报错
	assert.NotPanics(t, func() { h.SendToUser(999, "notification", "x") })
}
