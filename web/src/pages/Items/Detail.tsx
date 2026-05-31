import React, { useEffect, useState } from 'react'
import { Descriptions, Tag, Button, Image, Card, Space, message, Modal, Input } from 'antd'
import { SwapOutlined, MessageOutlined, UserOutlined } from '@ant-design/icons'
import { useParams, useNavigate } from 'react-router-dom'
import { itemApi } from '@/services/item'
import { tradeApi } from '@/services/trade'
import { useAuthStore } from '@/stores/auth'
import type { Item } from '@/types'

const ItemDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const { user } = useAuthStore()
  const [item, setItem] = useState<Item | null>(null)
  const [tradeModal, setTradeModal] = useState(false)
  const [tradeMessage, setTradeMessage] = useState('')
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (id) loadItem(Number(id))
  }, [id])

  const loadItem = async (itemId: number) => {
    try {
      const res = await itemApi.get(itemId)
      setItem(res.data.data)
    } catch {
      message.error('物品不存在')
      navigate('/')
    }
  }

  const handleTrade = async () => {
    if (!item) return
    setLoading(true)
    try {
      await tradeApi.create({
        target_item_id: item.id,
        barter_coin_amount: Number(item.estimated_value) || 0,
        message: tradeMessage,
      })
      message.success('交换请求已发送')
      setTradeModal(false)
    } catch (err: unknown) {
      message.error((err as Error).message)
    }
    setLoading(false)
  }

  if (!item) return null

  const conditionLabel: Record<string, string> = {
    new: '全新', like_new: '几乎全新', good: '良好', fair: '一般',
  }

  const isOwner = user?.id === item.user_id

  return (
    <div>
      <Card>
        <div style={{ display: 'flex', gap: 24, flexWrap: 'wrap' }}>
          <div style={{ flex: '0 0 400px', maxWidth: '100%' }}>
            {item.images?.length > 0 ? (
              <Image.PreviewGroup>
                <Image src={item.images[0]} style={{ width: '100%', maxHeight: 400, objectFit: 'cover' }} />
                <div style={{ display: 'flex', gap: 8, marginTop: 8, flexWrap: 'wrap' }}>
                  {item.images.slice(1).map((img, i) => (
                    <Image key={i} src={img} width={80} height={80} style={{ objectFit: 'cover' }} />
                  ))}
                </div>
              </Image.PreviewGroup>
            ) : (
              <div style={{ height: 300, background: '#f5f5f5', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                暂无图片
              </div>
            )}
          </div>

          <div style={{ flex: 1, minWidth: 300 }}>
            <h2>{item.title}</h2>
            <Descriptions column={1} size="small">
              <Descriptions.Item label="估值">
                <span style={{ color: '#f5222d', fontSize: 20, fontWeight: 'bold' }}>
                  {item.estimated_value || '面议'} 巴特币
                </span>
              </Descriptions.Item>
              <Descriptions.Item label="成色">
                <Tag color="blue">{conditionLabel[item.condition]}</Tag>
              </Descriptions.Item>
              <Descriptions.Item label="分类">{item.category?.name}</Descriptions.Item>
              <Descriptions.Item label="位置">{item.location || '未设置'}</Descriptions.Item>
              <Descriptions.Item label="浏览量">{item.view_count}</Descriptions.Item>
              <Descriptions.Item label="发布者">
                <Space>
                  <UserOutlined />
                  {item.user?.nickname}
                </Space>
              </Descriptions.Item>
            </Descriptions>

            {item.want_items?.length > 0 && (
              <div style={{ marginTop: 16 }}>
                <strong>期望交换：</strong>
                {item.want_items.map((w, i) => <Tag key={i}>{w}</Tag>)}
              </div>
            )}

            {!isOwner && (
              <Space style={{ marginTop: 24 }}>
                <Button type="primary" icon={<SwapOutlined />} onClick={() => setTradeModal(true)}>
                  发起交换
                </Button>
                <Button icon={<MessageOutlined />} onClick={() => navigate(`/messages?to=${item.user_id}`)}>
                  联系卖家
                </Button>
              </Space>
            )}
          </div>
        </div>

        {item.description && (
          <div style={{ marginTop: 24 }}>
            <h3>物品描述</h3>
            <p style={{ whiteSpace: 'pre-wrap' }}>{item.description}</p>
          </div>
        )}
      </Card>

      <Modal
        title="发起交换"
        open={tradeModal}
        onOk={handleTrade}
        onCancel={() => setTradeModal(false)}
        confirmLoading={loading}
        okText="发送请求"
      >
        <p>你想交换「{item.title}」，预计需要 {item.estimated_value || 0} 巴特币</p>
        <Input.TextArea
          placeholder="给对方留言（可选）"
          value={tradeMessage}
          onChange={(e) => setTradeMessage(e.target.value)}
          rows={3}
        />
      </Modal>
    </div>
  )
}

export default ItemDetail
