import React from 'react'
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { ConfigProvider } from 'antd'
import zhCN from 'antd/locale/zh_CN'
import { useAuthStore } from '@/stores/auth'
import AppLayout from '@/components/Layout'
import LoginPage from '@/pages/Auth/Login'
import HomePage from '@/pages/Home'
import ItemDetail from '@/pages/Items/Detail'
import PublishItem from '@/pages/Items/Publish'
import MyItems from '@/pages/Items/MyItems'
import TradeList from '@/pages/Trade'
import MessagesPage from '@/pages/Messages'
import WalletPage from '@/pages/Wallet'
import ProfilePage from '@/pages/Profile'
import NotificationsPage from '@/pages/Notifications'

// 路由守卫
const ProtectedRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { isLoggedIn } = useAuthStore()
  if (!isLoggedIn) return <Navigate to="/login" replace />
  return <>{children}</>
}

const App: React.FC = () => {
  return (
    <ConfigProvider locale={zhCN}>
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route
            path="/"
            element={
              <ProtectedRoute>
                <AppLayout />
              </ProtectedRoute>
            }
          >
            <Route index element={<HomePage />} />
            <Route path="items/:id" element={<ItemDetail />} />
            <Route path="items/publish" element={<PublishItem />} />
            <Route path="items/mine" element={<MyItems />} />
            <Route path="trades" element={<TradeList />} />
            <Route path="messages" element={<MessagesPage />} />
            <Route path="wallet" element={<WalletPage />} />
            <Route path="profile" element={<ProfilePage />} />
            <Route path="notifications" element={<NotificationsPage />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </ConfigProvider>
  )
}

export default App
