import React, { useState, useEffect } from 'react'
import { Form, Input, Button, Radio, message, Card } from 'antd'
import { MobileOutlined, LockOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { authApi } from '@/services/auth'
import { useAuthStore } from '@/stores/auth'

const LoginPage: React.FC = () => {
  const navigate = useNavigate()
  const { setAuth, isLoggedIn } = useAuthStore()
  const [countdown, setCountdown] = useState(0)
  const [loading, setLoading] = useState(false)
  const [form] = Form.useForm()

  useEffect(() => {
    if (isLoggedIn) navigate('/')
  }, [isLoggedIn, navigate])

  useEffect(() => {
    if (countdown > 0) {
      const timer = setTimeout(() => setCountdown(countdown - 1), 1000)
      return () => clearTimeout(timer)
    }
  }, [countdown])

  const handleSendCode = async () => {
    const phone = form.getFieldValue('phone')
    if (!phone || phone.length !== 11) {
      message.error('请输入正确的手机号')
      return
    }
    try {
      await authApi.sendCode(phone)
      setCountdown(60)
      message.success('验证码已发送')
    } catch (err: unknown) {
      message.error((err as Error).message)
    }
  }

  const handleLogin = async (values: { phone: string; code: string; user_type: string }) => {
    setLoading(true)
    try {
      const res = await authApi.login(values.phone, values.code, values.user_type)
      const { token, user } = res.data.data
      setAuth(token, user)
      message.success('登录成功')
      navigate('/')
    } catch (err: unknown) {
      message.error((err as Error).message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '100vh', background: '#f0f2f5' }}>
      <Card title="欢迎来到 Easy Barter" style={{ width: 400 }}>
        <Form form={form} onFinish={handleLogin} initialValues={{ user_type: 'personal' }}>
          <Form.Item name="phone" rules={[{ required: true, message: '请输入手机号' }, { len: 11, message: '手机号格式错误' }]}>
            <Input prefix={<MobileOutlined />} placeholder="手机号" maxLength={11} />
          </Form.Item>

          <Form.Item name="code" rules={[{ required: true, message: '请输入验证码' }]}>
            <Input
              prefix={<LockOutlined />}
              placeholder="验证码"
              maxLength={6}
              suffix={
                <Button type="link" disabled={countdown > 0} onClick={handleSendCode} style={{ padding: 0 }}>
                  {countdown > 0 ? `${countdown}s` : '获取验证码'}
                </Button>
              }
            />
          </Form.Item>

          <Form.Item name="user_type" label="用户类型">
            <Radio.Group>
              <Radio value="personal">个人用户</Radio>
              <Radio value="enterprise">企业用户</Radio>
            </Radio.Group>
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading} block>
              登录 / 注册
            </Button>
          </Form.Item>
        </Form>
        <div style={{ textAlign: 'center', color: '#999', fontSize: 12 }}>
          未注册的手机号将自动创建账户
        </div>
      </Card>
    </div>
  )
}

export default LoginPage
