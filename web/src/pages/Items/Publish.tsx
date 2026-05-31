import React, { useState, useEffect } from 'react'
import { Form, Input, Select, InputNumber, Upload, Button, message, Card } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { itemApi } from '@/services/item'
import type { Category } from '@/types'

const { TextArea } = Input

const PublishItem: React.FC = () => {
  const navigate = useNavigate()
  const [categories, setCategories] = useState<Category[]>([])
  const [loading, setLoading] = useState(false)
  const [form] = Form.useForm()

  useEffect(() => {
    itemApi.listCategories().then((res) => setCategories(res.data.data))
  }, [])

  const handleSubmit = async (values: Record<string, unknown>) => {
    setLoading(true)
    try {
      await itemApi.create({
        title: values.title as string,
        description: values.description as string,
        category_id: values.category_id as number,
        estimated_value: String(values.estimated_value || 0),
        condition: values.condition as string,
        images: [],
        want_items: values.want_items ? (values.want_items as string).split(',').map((s: string) => s.trim()) : [],
        location: values.location as string,
      })
      message.success('物品发布成功')
      navigate('/items/mine')
    } catch (err: unknown) {
      message.error((err as Error).message)
    }
    setLoading(false)
  }

  return (
    <Card title="发布物品">
      <Form form={form} layout="vertical" onFinish={handleSubmit} style={{ maxWidth: 600 }}>
        <Form.Item name="title" label="物品标题" rules={[{ required: true, message: '请输入标题' }, { max: 50, message: '标题不超过50字' }]}>
          <Input placeholder="请输入物品标题" maxLength={50} showCount />
        </Form.Item>

        <Form.Item name="description" label="物品描述">
          <TextArea placeholder="详细描述物品的情况" rows={4} />
        </Form.Item>

        <Form.Item name="category_id" label="物品分类" rules={[{ required: true, message: '请选择分类' }]}>
          <Select placeholder="选择分类" options={categories.map((c) => ({ label: c.name, value: c.id }))} />
        </Form.Item>

        <Form.Item name="condition" label="物品成色" rules={[{ required: true, message: '请选择成色' }]}>
          <Select
            placeholder="选择成色"
            options={[
              { label: '全新', value: 'new' },
              { label: '几乎全新', value: 'like_new' },
              { label: '良好', value: 'good' },
              { label: '一般', value: 'fair' },
            ]}
          />
        </Form.Item>

        <Form.Item name="estimated_value" label="估值（巴特币）">
          <InputNumber min={0} precision={2} style={{ width: '100%' }} placeholder="物品估值" />
        </Form.Item>

        <Form.Item name="images" label="物品图片">
          <Upload listType="picture-card" maxCount={9} beforeUpload={() => false}>
            <div>
              <PlusOutlined />
              <div style={{ marginTop: 8 }}>上传</div>
            </div>
          </Upload>
        </Form.Item>

        <Form.Item name="location" label="所在地">
          <Input placeholder="例如：北京市朝阳区" />
        </Form.Item>

        <Form.Item name="want_items" label="期望交换的物品">
          <Input placeholder="多个用逗号分隔，如：手机,平板,耳机" />
        </Form.Item>

        <Form.Item>
          <Button type="primary" htmlType="submit" loading={loading} block>
            发布物品
          </Button>
        </Form.Item>
      </Form>
    </Card>
  )
}

export default PublishItem
