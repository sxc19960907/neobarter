import React, { useEffect, useState, useRef } from 'react'
import { List, Input, Button, Avatar, Badge, Empty, Card } from 'antd'
import { SendOutlined, UserOutlined } from '@ant-design/icons'
import { useSearchParams, useNavigate } from 'react-router-dom'
import { messageApi } from '@/services/misc'
import { useAuthStore } from '@/stores/auth'
import type { Conversation, Message } from '@/types'

// 物品卡片消息渲染
const ItemCardBubble: React.FC<{ extraData: string; onClick: (id: number) => void }> = ({ extraData, onClick }) => {
  let card: { item_id: number; title: string; image: string; estimated_value: string; condition: string }
  try {
    card = JSON.parse(extraData)
  } catch {
    return null
  }
  return (
    <Card
      hoverable
      size="small"
      style={{ width: 220, cursor: 'pointer' }}
      onClick={() => onClick(card.item_id)}
      cover={card.image ? <img src={card.image} alt={card.title} style={{ height: 120, objectFit: 'cover' }} /> : undefined}
    >
      <Card.Meta
        title={card.title}
        description={card.estimated_value ? `${card.estimated_value} 巴特币` : '面议'}
      />
    </Card>
  )
}

const MessagesPage: React.FC = () => {
  const { user } = useAuthStore()
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const [conversations, setConversations] = useState<Conversation[]>([])
  const [activeConv, setActiveConv] = useState<number | null>(null)
  const [messages, setMessages] = useState<Message[]>([])
  const [inputValue, setInputValue] = useState('')
  const messagesEndRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    loadConversations()
    const toUser = searchParams.get('to')
    if (toUser) {
      // 如果有 to 参数，发送一条空消息来创建会话（或直接打开）
    }
  }, [searchParams])

  useEffect(() => {
    if (activeConv) loadMessages(activeConv)
  }, [activeConv])

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages])

  const loadConversations = async () => {
    try {
      const res = await messageApi.listConversations()
      setConversations(res.data.data || [])
    } catch { /* ignore */ }
  }

  const loadMessages = async (convId: number) => {
    try {
      const res = await messageApi.getMessages(convId, { page: 1, page_size: 50 })
      setMessages((res.data.data.list || []).reverse())
      messageApi.markRead(convId)
    } catch { /* ignore */ }
  }

  const handleSend = async () => {
    if (!inputValue.trim() || !activeConv) return
    try {
      await messageApi.send({ conversation_id: activeConv, content: inputValue })
      setInputValue('')
      loadMessages(activeConv)
    } catch { /* ignore */ }
  }

  return (
    <div style={{ display: 'flex', height: 'calc(100vh - 180px)', gap: 16 }}>
      {/* 会话列表 */}
      <div style={{ width: 280, borderRight: '1px solid #f0f0f0', overflowY: 'auto' }}>
        <List
          dataSource={conversations}
          renderItem={(conv) => (
            <List.Item
              onClick={() => setActiveConv(conv.conversation_id)}
              style={{
                cursor: 'pointer',
                padding: '12px 16px',
                background: activeConv === conv.conversation_id ? '#e6f7ff' : undefined,
              }}
            >
              <List.Item.Meta
                avatar={<Badge count={conv.unread_count} size="small"><Avatar icon={<UserOutlined />} /></Badge>}
                title={conv.user?.nickname || `会话 ${conv.conversation_id}`}
              />
            </List.Item>
          )}
          locale={{ emptyText: <Empty description="暂无会话" /> }}
        />
      </div>

      {/* 聊天区域 */}
      <div style={{ flex: 1, display: 'flex', flexDirection: 'column' }}>
        {activeConv ? (
          <>
            <div style={{ flex: 1, overflowY: 'auto', padding: 16 }}>
              {messages.map((msg) => (
                <div
                  key={msg.id}
                  style={{
                    display: 'flex',
                    justifyContent: msg.sender_id === user?.id ? 'flex-end' : 'flex-start',
                    marginBottom: 12,
                  }}
                >
                  {msg.message_type === 'item_card' && msg.extra_data ? (
                    <ItemCardBubble extraData={msg.extra_data} onClick={(id) => navigate(`/items/${id}`)} />
                  ) : (
                    <div
                      style={{
                        maxWidth: '70%',
                        padding: '8px 12px',
                        borderRadius: 8,
                        background: msg.sender_id === user?.id ? '#1890ff' : '#f0f0f0',
                        color: msg.sender_id === user?.id ? '#fff' : '#000',
                      }}
                    >
                      {msg.content}
                    </div>
                  )}
                </div>
              ))}
              <div ref={messagesEndRef} />
            </div>
            <div style={{ display: 'flex', gap: 8, padding: '8px 0' }}>
              <Input
                value={inputValue}
                onChange={(e) => setInputValue(e.target.value)}
                onPressEnter={handleSend}
                placeholder="输入消息..."
              />
              <Button type="primary" icon={<SendOutlined />} onClick={handleSend}>
                发送
              </Button>
            </div>
          </>
        ) : (
          <div style={{ flex: 1, display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
            <Empty description="选择一个会话开始聊天" />
          </div>
        )}
      </div>
    </div>
  )
}

export default MessagesPage
