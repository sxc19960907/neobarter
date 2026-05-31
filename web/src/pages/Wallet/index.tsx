import React, { useEffect, useState } from 'react'
import { Card, Statistic, Table, Row, Col, Tag } from 'antd'
import { WalletOutlined } from '@ant-design/icons'
import { walletApi } from '@/services/misc'
import type { Wallet, WalletTransaction } from '@/types'

const WalletPage: React.FC = () => {
  const [wallet, setWallet] = useState<Wallet | null>(null)
  const [transactions, setTransactions] = useState<WalletTransaction[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)

  useEffect(() => { loadWallet() }, [])
  useEffect(() => { loadTransactions() }, [page])

  const loadWallet = async () => {
    try {
      const res = await walletApi.getWallet()
      setWallet(res.data.data)
    } catch { /* ignore */ }
  }

  const loadTransactions = async () => {
    try {
      const res = await walletApi.listTransactions({ page, page_size: 10 })
      setTransactions(res.data.data.list || [])
      setTotal(res.data.data.total)
    } catch { /* ignore */ }
  }

  const typeLabel: Record<string, { text: string; color: string }> = {
    reward: { text: '奖励', color: 'green' },
    trade_in: { text: '交易收入', color: 'green' },
    trade_out: { text: '交易支出', color: 'red' },
    deposit: { text: '充值', color: 'blue' },
    withdraw: { text: '提现', color: 'orange' },
  }

  const columns = [
    {
      title: '类型', dataIndex: 'type', key: 'type',
      render: (t: string) => <Tag color={typeLabel[t]?.color}>{typeLabel[t]?.text || t}</Tag>,
    },
    {
      title: '金额', dataIndex: 'amount', key: 'amount',
      render: (v: string) => (
        <span style={{ color: Number(v) >= 0 ? '#52c41a' : '#f5222d' }}>
          {Number(v) >= 0 ? '+' : ''}{v}
        </span>
      ),
    },
    { title: '余额', dataIndex: 'balance_after', key: 'balance' },
    { title: '说明', dataIndex: 'description', key: 'desc' },
    { title: '时间', dataIndex: 'created_at', key: 'time', render: (t: string) => t?.slice(0, 16) },
  ]

  return (
    <div>
      <Card style={{ marginBottom: 24 }}>
        <Row gutter={24}>
          <Col span={8}>
            <Statistic
              title="巴特币余额"
              value={wallet?.balance || '0.00'}
              prefix={<WalletOutlined />}
              precision={2}
            />
          </Col>
          <Col span={8}>
            <Statistic title="累计收入" value={wallet?.total_income || '0.00'} precision={2} valueStyle={{ color: '#52c41a' }} />
          </Col>
          <Col span={8}>
            <Statistic title="累计支出" value={wallet?.total_expense || '0.00'} precision={2} valueStyle={{ color: '#f5222d' }} />
          </Col>
        </Row>
      </Card>

      <Card title="交易流水">
        <Table
          columns={columns}
          dataSource={transactions}
          rowKey="id"
          pagination={{ current: page, total, pageSize: 10, onChange: setPage }}
        />
      </Card>
    </div>
  )
}

export default WalletPage
