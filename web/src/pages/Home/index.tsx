import React, { useEffect, useState } from 'react'
import { Row, Col, Input, Select, Card, Tag, Empty, Pagination, Spin } from 'antd'
import { SearchOutlined, EnvironmentOutlined, EyeOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { itemApi, type ItemQuery } from '@/services/item'
import type { Item, Category } from '@/types'

const { Search } = Input

const HomePage: React.FC = () => {
  const navigate = useNavigate()
  const [items, setItems] = useState<Item[]>([])
  const [categories, setCategories] = useState<Category[]>([])
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(false)
  const [query, setQuery] = useState<ItemQuery>({ page: 1, page_size: 12 })

  useEffect(() => {
    loadCategories()
  }, [])

  useEffect(() => {
    loadItems()
  }, [query])

  const loadCategories = async () => {
    try {
      const res = await itemApi.listCategories()
      setCategories(res.data.data)
    } catch { /* ignore */ }
  }

  const loadItems = async () => {
    setLoading(true)
    try {
      const res = await itemApi.list(query)
      setItems(res.data.data.list || [])
      setTotal(res.data.data.total)
    } catch { /* ignore */ }
    setLoading(false)
  }

  const conditionLabel: Record<string, string> = {
    new: '全新', like_new: '几乎全新', good: '良好', fair: '一般',
  }

  return (
    <div>
      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        <Col xs={24} sm={12} md={8}>
          <Search
            placeholder="搜索物品"
            prefix={<SearchOutlined />}
            onSearch={(v) => setQuery({ ...query, keyword: v, page: 1 })}
            allowClear
          />
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Select
            placeholder="选择分类"
            allowClear
            style={{ width: '100%' }}
            onChange={(v) => setQuery({ ...query, category_id: v, page: 1 })}
            options={categories.map((c) => ({ label: c.name, value: c.id }))}
          />
        </Col>
        <Col xs={12} md={4}>
          <Select
            placeholder="成色"
            allowClear
            style={{ width: '100%' }}
            onChange={(v) => setQuery({ ...query, condition: v, page: 1 })}
            options={[
              { label: '全新', value: 'new' },
              { label: '几乎全新', value: 'like_new' },
              { label: '良好', value: 'good' },
              { label: '一般', value: 'fair' },
            ]}
          />
        </Col>
        <Col xs={12} md={4}>
          <Select
            placeholder="排序"
            defaultValue="created_at"
            style={{ width: '100%' }}
            onChange={(v) => setQuery({ ...query, sort_by: v, page: 1 })}
            options={[
              { label: '最新', value: 'created_at' },
              { label: '最热', value: 'view_count' },
              { label: '估值', value: 'estimated_value' },
            ]}
          />
        </Col>
      </Row>

      <Spin spinning={loading}>
        {items.length === 0 ? (
          <Empty description="暂无物品" />
        ) : (
          <Row gutter={[16, 16]}>
            {items.map((item) => (
              <Col xs={24} sm={12} md={8} lg={6} key={item.id}>
                <Card
                  hoverable
                  cover={
                    item.images?.[0] ? (
                      <img alt={item.title} src={item.images[0]} style={{ height: 200, objectFit: 'cover' }} />
                    ) : (
                      <div style={{ height: 200, background: '#f5f5f5', display: 'flex', alignItems: 'center', justifyContent: 'center', color: '#ccc' }}>
                        暂无图片
                      </div>
                    )
                  }
                  onClick={() => navigate(`/items/${item.id}`)}
                >
                  <Card.Meta
                    title={item.title}
                    description={
                      <div>
                        <div style={{ marginBottom: 8 }}>
                          <Tag color="blue">{conditionLabel[item.condition] || item.condition}</Tag>
                          {item.estimated_value && <Tag color="orange">¥{item.estimated_value} 巴特币</Tag>}
                        </div>
                        <div style={{ color: '#999', fontSize: 12 }}>
                          {item.location && <span><EnvironmentOutlined /> {item.location} </span>}
                          <span><EyeOutlined /> {item.view_count}</span>
                        </div>
                      </div>
                    }
                  />
                </Card>
              </Col>
            ))}
          </Row>
        )}
      </Spin>

      {total > 0 && (
        <div style={{ textAlign: 'center', marginTop: 24 }}>
          <Pagination
            current={query.page}
            pageSize={query.page_size}
            total={total}
            onChange={(page) => setQuery({ ...query, page })}
            showSizeChanger={false}
          />
        </div>
      )}
    </div>
  )
}

export default HomePage
