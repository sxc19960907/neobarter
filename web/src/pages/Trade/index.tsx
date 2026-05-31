import React, { useEffect, useState } from 'react'
import { Table, Tag, Button, Space, Tabs, message, Modal, Input } from 'antd'
import { tradeApi } from '@/services/trade'
import { useAuthStore } from '@/stores/auth'
import type { TradeRequest } from '@/types'

const TradeList: React.FC = () => {
  const { user } = useAuthStore()
  const [trades, setTrades] = useState<TradeRequest[]>([])
  const [loading, setLoading] = useState(false)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [status, setStatus] = useState('')
  const [rejectModal, setRejectModal] = useState<number | null>(null)
  const [rejectReason, setRejectReason] = useState('')

  useEffect(() => { loadTrades() }, [page, status])

  const loadTrades = async () => {
    setLoading(true)
    try {
      const res = await tradeApi.list({ status, page, page_size: 10 })
      setTrades(res.data.data.list || [])
      setTotal(res.data.data.total)
    } catch { /* ignore */ }
    setLoading(false)
  }

  const handleAccept = async (id: number) => {
    try {
      await tradeApi.accept(id)
      message.success('已接受')
      loadTrades()
    } catch (err: unknown) { message.error((err as Error).message) }
  }

  const handleReject = async () => {
    if (!rejectModal || !rejectReason) return
    try {
      await tradeApi.reject(rejectModal, rejectReason)
      message.success('已拒绝')
      setRejectModal(null)
      setRejectReason('')
      loadTrades()
    } catch (err: unknown) { message.error((err as Error).message) }
  }

  const handleComplete = async (id: number) => {
    try {
      await tradeApi.complete(id)
      message.success('交易完成')
      loadTrades()
    } catch (err: unknown) { message.error((err as Error).message) }
  }

  const handleCancel = async (id: number) => {
    try {
      await tradeApi.cancel(id)
      message.success('已取消')
      loadTrades()
    } catch (err: unknown) { message.error((err as Error).message) }
  }

  const statusMap: Record<string, { color: string; text: string }> = {
    pending: { color: 'orange', text: '待确认' },
    accepted: { color: 'blue', text: '已接受' },
    rejected: { color: 'red', text: '已拒绝' },
    completed: { color: 'green', text: '已完成' },
    cancelled: { color: 'default', text: '已取消' },
    expired: { color: 'default', text: '已过期' },
  }

  const columns = [
    { title: '目标物品', key: 'target', render: (_: unknown, r: TradeRequest) => r.target_item?.title || '-' },
    { title: '发起方', key: 'initiator', render: (_: unknown, r: TradeRequest) => r.initiator?.nickname || '-' },
    { title: '巴特币', dataIndex: 'barter_coin_amount', key: 'amount' },
    {
      title: '状态', dataIndex: 'status', key: 'status',
      render: (s: string) => <Tag color={statusMap[s]?.color}>{statusMap[s]?.text || s}</Tag>,
    },
    {
      title: '操作', key: 'action',
      render: (_: unknown, r: TradeRequest) => {
        const isTarget = r.target_user_id === user?.id
        const isInitiator = r.initiator_id === user?.id
        return (
          <Space>
            {r.status === 'pending' && isTarget && (
              <>
                <Button size="small" type="primary" onClick={() => handleAccept(r.id)}>接受</Button>
                <Button size="small" danger onClick={() => setRejectModal(r.id)}>拒绝</Button>
              </>
            )}
            {r.status === 'pending' && isInitiator && (
              <Button size="small" onClick={() => handleCancel(r.id)}>取消</Button>
            )}
            {r.status === 'accepted' && (
              <Button size="small" type="primary" onClick={() => handleComplete(r.id)}>确认完成</Button>
            )}
          </Space>
        )
      },
    },
  ]

  return (
    <div>
      <h2>交易管理</h2>
      <Tabs
        activeKey={status}
        onChange={(k) => { setStatus(k); setPage(1) }}
        items={[
          { key: '', label: '全部' },
          { key: 'pending', label: '待确认' },
          { key: 'accepted', label: '进行中' },
          { key: 'completed', label: '已完成' },
        ]}
      />
      <Table
        columns={columns}
        dataSource={trades}
        rowKey="id"
        loading={loading}
        pagination={{ current: page, total, pageSize: 10, onChange: setPage }}
      />
      <Modal
        title="拒绝原因"
        open={!!rejectModal}
        onOk={handleReject}
        onCancel={() => setRejectModal(null)}
        okText="确认拒绝"
      >
        <Input.TextArea
          value={rejectReason}
          onChange={(e) => setRejectReason(e.target.value)}
          placeholder="请填写拒绝原因"
          rows={3}
        />
      </Modal>
    </div>
  )
}

export default TradeList
