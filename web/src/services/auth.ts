import request from '@/utils/request'
import type { ApiResponse, User } from '@/types'

export const authApi = {
  sendCode(phone: string) {
    return request.post<ApiResponse>('/auth/send-code', { phone })
  },

  login(phone: string, code: string, userType = 'personal') {
    return request.post<ApiResponse<{ token: string; user: User }>>('/auth/login', {
      phone,
      code,
      user_type: userType,
    })
  },
}
