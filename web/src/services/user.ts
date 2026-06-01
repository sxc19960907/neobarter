import request from '@/utils/request'
import type { ApiResponse, User, UserAddress } from '@/types'

export const userApi = {
  getMe() {
    return request.get<ApiResponse<User>>('/users/me')
  },

  updateMe(data: Partial<User>) {
    return request.put<ApiResponse<User>>('/users/me', data)
  },

  getUser(id: number) {
    return request.get<ApiResponse<User>>(`/users/${id}`)
  },

  listAddresses() {
    return request.get<ApiResponse<UserAddress[]>>('/users/me/addresses')
  },

  createAddress(data: Omit<UserAddress, 'id' | 'user_id'>) {
    return request.post<ApiResponse<UserAddress>>('/users/me/addresses', data)
  },

  updateAddress(id: number, data: Omit<UserAddress, 'id' | 'user_id'>) {
    return request.put<ApiResponse<UserAddress>>(`/users/me/addresses/${id}`, data)
  },

  deleteAddress(id: number) {
    return request.delete<ApiResponse>(`/users/me/addresses/${id}`)
  },

  verifyRealName(data: { real_name: string; id_card: string }) {
    return request.post<ApiResponse>('/users/me/verify-realname', data)
  },

  verifyEnterprise(data: { enterprise_name: string; license_url: string }) {
    return request.post<ApiResponse>('/users/me/verify-enterprise', data)
  },
}
