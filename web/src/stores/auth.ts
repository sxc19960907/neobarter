import { create } from 'zustand'
import type { User } from '@/types'
import { userApi } from '@/services/user'

interface AuthState {
  token: string | null
  user: User | null
  isLoggedIn: boolean
  setAuth: (token: string, user: User) => void
  logout: () => void
  fetchUser: () => Promise<void>
}

export const useAuthStore = create<AuthState>((set) => ({
  token: localStorage.getItem('token'),
  user: null,
  isLoggedIn: !!localStorage.getItem('token'),

  setAuth: (token, user) => {
    localStorage.setItem('token', token)
    set({ token, user, isLoggedIn: true })
  },

  logout: () => {
    localStorage.removeItem('token')
    set({ token: null, user: null, isLoggedIn: false })
  },

  fetchUser: async () => {
    try {
      const res = await userApi.getMe()
      set({ user: res.data.data })
    } catch {
      localStorage.removeItem('token')
      set({ token: null, user: null, isLoggedIn: false })
    }
  },
}))
