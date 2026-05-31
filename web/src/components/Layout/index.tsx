import React from 'react'
import { Layout, Menu, Badge, Avatar, Dropdown, Space } from 'antd'
import {
  HomeOutlined,
  ShoppingOutlined,
  SwapOutlined,
  MessageOutlined,
  WalletOutlined,
  BellOutlined,
  UserOutlined,
  LogoutOutlined,
} from '@ant-design/icons'
import { Outlet, useNavigate, useLocation } from 'react-router-dom'
import { useAuthStore } from '@/stores/auth'

const { Header, Content, Sider } = Layout

const AppLayout: React.FC = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const { user, logout } = useAuthStore()

  const menuItems = [
    { key: '/', icon: <HomeOutlined />, label: '首页' },
    { key: '/items/mine', icon: <ShoppingOutlined />, label: '我的物品' },
    { key: '/trades', icon: <SwapOutlined />, label: '交易管理' },
    { key: '/messages', icon: <MessageOutlined />, label: '消息' },
    { key: '/wallet', icon: <WalletOutlined />, label: '钱包' },
  ]

  const userMenuItems = [
    { key: 'profile', icon: <UserOutlined />, label: '个人中心' },
    { key: 'logout', icon: <LogoutOutlined />, label: '退出登录' },
  ]

  const handleUserMenu = ({ key }: { key: string }) => {
    if (key === 'logout') {
      logout()
      navigate('/login')
    } else if (key === 'profile') {
      navigate('/profile')
    }
  }

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', padding: '0 24px' }}>
        <div style={{ color: '#fff', fontSize: 20, fontWeight: 'bold', cursor: 'pointer' }} onClick={() => navigate('/')}>
          新易巴特
        </div>
        <Space size="large">
          <Badge count={0} size="small">
            <BellOutlined style={{ color: '#fff', fontSize: 18, cursor: 'pointer' }} onClick={() => navigate('/notifications')} />
          </Badge>
          <Dropdown menu={{ items: userMenuItems, onClick: handleUserMenu }} placement="bottomRight">
            <Space style={{ cursor: 'pointer', color: '#fff' }}>
              <Avatar src={user?.avatar_url} icon={<UserOutlined />} size="small" />
              <span>{user?.nickname || '用户'}</span>
            </Space>
          </Dropdown>
        </Space>
      </Header>
      <Layout>
        <Sider
          breakpoint="lg"
          collapsedWidth="0"
          style={{ background: '#fff' }}
        >
          <Menu
            mode="inline"
            selectedKeys={[location.pathname]}
            items={menuItems}
            onClick={({ key }) => navigate(key)}
            style={{ height: '100%', borderRight: 0 }}
          />
        </Sider>
        <Content style={{ padding: 24, minHeight: 280 }}>
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  )
}

export default AppLayout
