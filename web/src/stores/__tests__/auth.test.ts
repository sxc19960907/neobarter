import { describe, it, expect, beforeEach, vi } from 'vitest'
import { useAuthStore } from '../auth'
import type { User } from '@/types'

// Mock userApi
vi.mock('@/services/user', () => ({
  userApi: {
    getMe: vi.fn(),
  },
}))

const mockUser: User = {
  id: 1,
  phone: '13800138000',
  nickname: '测试用户',
  avatar_url: '',
  user_type: 'personal',
  status: 'active',
  credit_score: 100,
  real_name_verified: false,
  enterprise_name: '',
  enterprise_verified: false,
  location: '',
  bio: '',
  last_login_at: '',
  created_at: '',
}

describe('useAuthStore', () => {
  beforeEach(() => {
    localStorage.clear()
    useAuthStore.setState({ token: null, user: null, isLoggedIn: false })
  })

  it('setAuth 应保存 token 和用户信息', () => {
    useAuthStore.getState().setAuth('test-token', mockUser)

    const state = useAuthStore.getState()
    expect(state.token).toBe('test-token')
    expect(state.user).toEqual(mockUser)
    expect(state.isLoggedIn).toBe(true)
    expect(localStorage.getItem('token')).toBe('test-token')
  })

  it('logout 应清除 token 和用户信息', () => {
    useAuthStore.getState().setAuth('test-token', mockUser)
    useAuthStore.getState().logout()

    const state = useAuthStore.getState()
    expect(state.token).toBeNull()
    expect(state.user).toBeNull()
    expect(state.isLoggedIn).toBe(false)
    expect(localStorage.getItem('token')).toBeNull()
  })
})
