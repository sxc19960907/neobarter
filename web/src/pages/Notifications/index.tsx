import React, { useEffect, useState } from 'react'
import { List, Badge, Button, Empty, Tag } from 'antd'
import { notificationApi } from '@/services/misc'
import type { Notification } from '@/types'

const NotificationsPage: React.FC = () => {
  const [notifications, setNotifications] = useState<Notification[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)

  useEffect(() => { loadNotifications() }, [page])

  const loadNotifications = async () => {
    try {
      const res = await notificationApi.list({ page, page_size: 20 })
      setNotifications(res.data.data.list || [])
      setTotal(res.data.data.total)
    } catch { /* ignore */ }
  }

  const handleMarkAllRead = async () => {
    await notificationApi.markAllRead()
    loadNotifications()
  }

  const typeColor: Record<string, string> = {
    trade_request: 'blue',
    trade_accepted: 'green',
    trade_rejected: 'red',
    message: 'purple',
    system: 'default',
  }

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
        <h2>通知</h2>
        <Button onClick={handleMarkAllRead}>全部已读</Button>
      </div>
      <List
        dataSource={notifications}
        pagination={{ current: page, total, pageSize: 20, onChange: setPage }}
        renderItem={(item) => (
          <List.Item>
            <List.Item.Meta
              avatar={<Badge dot={!item.is_read} />}
              title={
                <span>
                  <Tag color={typeColor[item.type]}>{item.type}</Tag>
                  {item.title}
                </span>
              }
              description={
                <div>
                  <div>{item.content}</div>
                  <div style={{ color: '#999', fontSize: 12, marginTop: 4 }}>{item.created_at?.slice(0, 16)}</div>
                </div>
              }
            />
          </List.Item>
        )}
        locale={{ emptyText: <Empty description="暂无通知" /> }}
      />
    </div>
  )
}

export default NotificationsPage
