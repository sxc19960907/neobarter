import React, { useEffect } from 'react'
import { Card, Form, Input, Button, Avatar, message, Descriptions } from 'antd'
import { UserOutlined } from '@ant-design/icons'
import { useAuthStore } from '@/stores/auth'
import { userApi } from '@/services/user'

const ProfilePage: React.FC = () => {
  const { user, fetchUser } = useAuthStore()
  const [form] = Form.useForm()

  useEffect(() => {
    if (user) {
      form.setFieldsValue({
        nickname: user.nickname,
        bio: user.bio,
        location: user.location,
      })
    }
  }, [user, form])

  const handleSave = async (values: { nickname: string; bio: string; location: string }) => {
    try {
      await userApi.updateMe(values)
      await fetchUser()
      message.success('保存成功')
    } catch (err: unknown) {
      message.error((err as Error).message)
    }
  }

  return (
    <div style={{ maxWidth: 600 }}>
      <Card title="个人信息" style={{ marginBottom: 24 }}>
        <div style={{ textAlign: 'center', marginBottom: 24 }}>
          <Avatar size={80} src={user?.avatar_url} icon={<UserOutlined />} />
        </div>
        <Descriptions column={1} size="small" style={{ marginBottom: 24 }}>
          <Descriptions.Item label="手机号">{user?.phone}</Descriptions.Item>
          <Descriptions.Item label="用户类型">{user?.user_type === 'personal' ? '个人用户' : '企业用户'}</Descriptions.Item>
          <Descriptions.Item label="信用积分">{user?.credit_score}</Descriptions.Item>
          <Descriptions.Item label="实名认证">{user?.real_name_verified ? '已认证' : '未认证'}</Descriptions.Item>
        </Descriptions>

        <Form form={form} layout="vertical" onFinish={handleSave}>
          <Form.Item name="nickname" label="昵称" rules={[{ max: 50, message: '昵称不超过50字' }]}>
            <Input placeholder="设置昵称" />
          </Form.Item>
          <Form.Item name="bio" label="个人简介">
            <Input.TextArea placeholder="介绍一下自己" rows={3} />
          </Form.Item>
          <Form.Item name="location" label="所在地">
            <Input placeholder="例如：北京市" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit">保存修改</Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  )
}

export default ProfilePage
