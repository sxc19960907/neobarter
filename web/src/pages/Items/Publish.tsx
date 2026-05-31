import React, { useState, useEffect } from 'react'
import { Form, Input, Select, InputNumber, Upload, Button, message, Card } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import type { UploadFile, UploadProps } from 'antd'
import { useNavigate } from 'react-router-dom'
import { itemApi } from '@/services/item'
import { uploadApi } from '@/services/upload'
import type { Category } from '@/types'

const { TextArea } = Input

const PublishItem: React.FC = () => {
  const navigate = useNavigate()
  const [categories, setCategories] = useState<Category[]>([])
  const [loading, setLoading] = useState(false)
  const [fileList, setFileList] = useState<UploadFile[]>([])
  const [form] = Form.useForm()

  useEffect(() => {
    itemApi.listCategories().then((res) => setCategories(res.data.data))
  }, [])

  // 自定义上传：调用后端上传接口
  const customUpload: UploadProps['customRequest'] = async (options) => {
    const { file, onSuccess, onError } = options
    try {
      const res = await uploadApi.uploadImage(file as File)
      onSuccess?.(res.data.data)
    } catch (err) {
      message.error('图片上传失败')
      onError?.(err as Error)
    }
  }

  const handleChange: UploadProps['onChange'] = ({ fileList: newList }) => {
    setFileList(newList)
  }

  // 校验文件类型和大小
  const beforeUpload = (file: File) => {
    const isValidType = ['image/jpeg', 'image/png', 'image/webp', 'image/gif'].includes(file.type)
    if (!isValidType) {
      message.error('仅支持 jpg/png/webp/gif 格式')
      return Upload.LIST_IGNORE
    }
    const isLt5M = file.size / 1024 / 1024 < 5
    if (!isLt5M) {
      message.error('图片大小不能超过 5MB')
      return Upload.LIST_IGNORE
    }
    return true
  }

  const handleSubmit = async (values: Record<string, unknown>) => {
    // 收集已成功上传的图片 URL
    const images = fileList
      .filter((f) => f.status === 'done')
      .map((f) => (f.response as { url: string })?.url)
      .filter(Boolean)

    setLoading(true)
    try {
      await itemApi.create({
        title: values.title as string,
        description: values.description as string,
        category_id: values.category_id as number,
        estimated_value: String(values.estimated_value || 0),
        condition: values.condition as string,
        images,
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

        <Form.Item label="物品图片" extra="最多9张，单张不超过5MB">
          <Upload
            listType="picture-card"
            maxCount={9}
            fileList={fileList}
            customRequest={customUpload}
            beforeUpload={beforeUpload}
            onChange={handleChange}
            accept="image/*"
          >
            {fileList.length < 9 && (
              <div>
                <PlusOutlined />
                <div style={{ marginTop: 8 }}>上传</div>
              </div>
            )}
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
