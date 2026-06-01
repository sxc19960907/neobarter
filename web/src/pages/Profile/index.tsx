import React, { useEffect, useState } from 'react'
import { Card, Form, Input, Button, Avatar, message, Descriptions, Upload, Tag, Modal } from 'antd'
import { UserOutlined } from '@ant-design/icons'
import type { UploadProps } from 'antd'
import { useAuthStore } from '@/stores/auth'
import { userApi } from '@/services/user'
import { uploadApi } from '@/services/upload'
import { creditLevelMap } from '@/utils/format'

const ProfilePage: React.FC = () => {
  const { user, fetchUser } = useAuthStore()
  const [form] = Form.useForm()
  const [verifyForm] = Form.useForm()
  const [avatarUploading, setAvatarUploading] = useState(false)
  const [verifyModal, setVerifyModal] = useState(false)

  useEffect(() => {
    if (user) {
      form.setFieldsValue({
        nickname: user.nickname,
        bio: user.bio,
        location: user.location,
      })
    }
  }, [user, form])

  // 头像上传
  const handleAvatarUpload: UploadProps['customRequest'] = async (options) => {
    const file = options.file as File
    setAvatarUploading(true)
    try {
      const res = await uploadApi.uploadImage(file)
      await userApi.updateMe({ avatar_url: res.data.data.url })
      await fetchUser()
      message.success('头像更新成功')
      options.onSuccess?.(res.data.data)
    } catch (err) {
      message.error('头像上传失败')
      options.onError?.(err as Error)
    }
    setAvatarUploading(false)
  }

  const beforeAvatarUpload = (file: File) => {
    const isValidType = ['image/jpeg', 'image/png', 'image/webp'].includes(file.type)
    if (!isValidType) {
      message.error('仅支持 jpg/png/webp 格式')
      return Upload.LIST_IGNORE
    }
    const isLt5M = file.size / 1024 / 1024 < 5
    if (!isLt5M) {
      message.error('头像大小不能超过 5MB')
      return Upload.LIST_IGNORE
    }
    return true
  }

  const handleSave = async (values: { nickname: string; bio: string; location: string }) => {
    try {
      await userApi.updateMe(values)
      await fetchUser()
      message.success('保存成功')
    } catch (err: unknown) {
      message.error((err as Error).message)
    }
  }

  const handleVerify = async (values: { real_name: string; id_card: string }) => {
    try {
      await userApi.verifyRealName(values)
      await fetchUser()
      message.success('实名认证成功')
      setVerifyModal(false)
      verifyForm.resetFields()
    } catch (err: unknown) {
      message.error((err as Error).message)
    }
  }

  return (
    <div style={{ maxWidth: 600 }}>
      <Card title="个人信息" style={{ marginBottom: 24 }}>
        <div style={{ textAlign: 'center', marginBottom: 24 }}>
          <Upload
            showUploadList={false}
            customRequest={handleAvatarUpload}
            beforeUpload={beforeAvatarUpload}
            accept="image/*"
          >
            <div style={{ cursor: 'pointer', display: 'inline-block' }}>
              <Avatar size={80} src={user?.avatar_url} icon={<UserOutlined />} />
              <div style={{ marginTop: 8, color: '#1890ff', fontSize: 12 }}>
                {avatarUploading ? '上传中...' : '点击更换头像'}
              </div>
            </div>
          </Upload>
        </div>
        <Descriptions column={1} size="small" style={{ marginBottom: 24 }}>
          <Descriptions.Item label="手机号">{user?.phone}</Descriptions.Item>
          <Descriptions.Item label="用户类型">{user?.user_type === 'personal' ? '个人用户' : '企业用户'}</Descriptions.Item>
          <Descriptions.Item label="信用积分">
            {user?.credit_score}
            {user?.credit_level && (
              <Tag color={creditLevelMap[user.credit_level]?.color} style={{ marginLeft: 8 }}>
                {creditLevelMap[user.credit_level]?.text || user.credit_level}
              </Tag>
            )}
          </Descriptions.Item>
          <Descriptions.Item label="实名认证">
            {user?.real_name_verified ? (
              <Tag color="green">已认证</Tag>
            ) : (
              <Button size="small" type="link" onClick={() => setVerifyModal(true)}>去认证</Button>
            )}
          </Descriptions.Item>
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

      <Modal
        title="实名认证"
        open={verifyModal}
        onOk={() => verifyForm.submit()}
        onCancel={() => setVerifyModal(false)}
        okText="提交认证"
      >
        <Form form={verifyForm} layout="vertical" onFinish={handleVerify}>
          <Form.Item name="real_name" label="真实姓名" rules={[{ required: true, message: '请填写真实姓名' }]}>
            <Input placeholder="请输入真实姓名" />
          </Form.Item>
          <Form.Item
            name="id_card"
            label="身份证号"
            rules={[
              { required: true, message: '请填写身份证号' },
              { pattern: /^\d{17}[\dXx]$/, message: '身份证号格式不正确' },
            ]}
          >
            <Input placeholder="18位身份证号" maxLength={18} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default ProfilePage
