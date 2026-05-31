import React, { useEffect, useState } from 'react'
import { Table, Tag, Button, Space, Popconfirm, message } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { itemApi } from '@/services/item'
import type { Item } from '@/types'

const MyItems: React.FC = () => {
  const navigate = useNavigate()
  const [items, setItems] = useState<Item[]>([])
  const [loading, setLoading] = useState(false)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)

  useEffect(() => { loadItems() }, [page])

  const loadItems = async () => {
    setLoading(true)
    try {
      const res = await itemApi.list({ mine: true, page, page_size: 10 })
      setItems(res.data.data.list || [])
      setTotal(res.data.data.total)
    } catch { /* ignore */ }
    setLoading(false)
  }

  const handleToggleStatus = async (item: Item) => {
    const newStatus = item.status === 'active' ? 'inactive' : 'active'
    try {
      await itemApi.updateStatus(item.id, newStatus)
      message.success(newStatus === 'active' ? '已上架' : '已下架')
      loadItems()
    } catch (err: unknown) {
      message.error((err as Error).message)
    }
  }

  const handleDelete = async (id: number) => {
    try {
      await itemApi.delete(id)
      message.success('已删除')
      loadItems()
    } catch (err: unknown) {
      message.error((err as Error).message)
    }
  }

  const statusMap: Record<string, { color: string; text: string }> = {
    active: { color: 'green', text: '上架中' },
    inactive: { color: 'default', text: '已下架' },
    traded: { color: 'blue', text: '已交易' },
    deleted: { color: 'red', text: '已删除' },
  }

  const columns = [
    { title: '标题', dataIndex: 'title', key: 'title' },
    {
      title: '估值', dataIndex: 'estimated_value', key: 'value',
      render: (v: string) => v ? `${v} 巴特币` : '-',
    },
    {
      title: '状态', dataIndex: 'status', key: 'status',
      render: (s: string) => <Tag color={statusMap[s]?.color}>{statusMap[s]?.text || s}</Tag>,
    },
    { title: '浏览量', dataIndex: 'view_count', key: 'views' },
    {
      title: '操作', key: 'action',
      render: (_: unknown, record: Item) => (
        <Space>
          <Button size="small" onClick={() => navigate(`/items/${record.id}`)}>查看</Button>
          {record.status !== 'traded' && record.status !== 'deleted' && (
            <>
              <Button size="small" onClick={() => handleToggleStatus(record)}>
                {record.status === 'active' ? '下架' : '上架'}
              </Button>
              <Popconfirm title="确定删除？" onConfirm={() => handleDelete(record.id)}>
                <Button size="small" danger>删除</Button>
              </Popconfirm>
            </>
          )}
        </Space>
      ),
    },
  ]

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <h2>我的物品</h2>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => navigate('/items/publish')}>
          发布物品
        </Button>
      </div>
      <Table
        columns={columns}
        dataSource={items}
        rowKey="id"
        loading={loading}
        pagination={{ current: page, total, pageSize: 10, onChange: setPage }}
      />
    </div>
  )
}

export default MyItems
