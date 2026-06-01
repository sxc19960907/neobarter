package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Hub 管理所有 WebSocket 连接
type Hub struct {
	mu         sync.RWMutex
	clients    map[int64]*Client // userID -> client
	register   chan *Client
	unregister chan *Client
}

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	userID int64
	send   chan []byte
}

type WSMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[int64]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.userID] = client
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			// 仅当 map 中仍是同一个 client 时才删除（避免同用户重连时误删新连接）
			if existing, ok := h.clients[client.userID]; ok && existing == client {
				delete(h.clients, client.userID)
				close(client.send)
			}
			h.mu.Unlock()
		}
	}
}

// SendToUser 向指定在线用户推送一条已封装的 WSMessage。离线则静默跳过。
func (h *Hub) SendToUser(userID int64, eventType string, payload interface{}) {
	data, err := json.Marshal(WSMessage{Type: eventType, Data: payload})
	if err != nil {
		return
	}
	h.deliver(userID, data)
}

// SendToUsers 向一批在线用户精确推送。用于会话消息：调用方传入会话参与者
// （通常已排除发送者）。离线用户静默跳过，不影响其他人。
func (h *Hub) SendToUsers(userIDs []int64, eventType string, payload interface{}) {
	if len(userIDs) == 0 {
		return
	}
	data, err := json.Marshal(WSMessage{Type: eventType, Data: payload})
	if err != nil {
		return
	}
	for _, uid := range userIDs {
		h.deliver(uid, data)
	}
}

// deliver 将原始字节投递给指定用户的连接（若在线）。send 缓冲满时丢弃，避免阻塞。
func (h *Hub) deliver(userID int64, data []byte) {
	h.mu.RLock()
	client, ok := h.clients[userID]
	h.mu.RUnlock()
	if !ok {
		return
	}
	select {
	case client.send <- data:
	default:
		// 缓冲已满，丢弃本条（连接可能已僵死，由 read/write pump 负责清理）
	}
}

// ServeWS 处理 WebSocket 连接升级
func ServeWS(hub *Hub, c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &Client{
		hub:    hub,
		conn:   conn,
		userID: userID.(int64),
		send:   make(chan []byte, 256),
	}

	hub.register <- client

	go client.writePump()
	go client.readPump()
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(512 * 1024)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.conn.WriteMessage(websocket.TextMessage, message)

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
