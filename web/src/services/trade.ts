import request from '@/utils/request'
import type { ApiResponse, PageData, TradeRequest } from '@/types'

export const tradeApi = {
  create(data: {
    target_item_id: number
    offered_item_id?: number
    barter_coin_amount?: number
    message?: string
  }) {
    return request.post<ApiResponse<TradeRequest>>('/trades', data)
  },

  list(params: { status?: string; page?: number; page_size?: number }) {
    return request.get<ApiResponse<PageData<TradeRequest>>>('/trades', { params })
  },

  get(id: number) {
    return request.get<ApiResponse<TradeRequest>>(`/trades/${id}`)
  },

  accept(id: number) {
    return request.put<ApiResponse>(`/trades/${id}/accept`)
  },

  reject(id: number, reason: string) {
    return request.put<ApiResponse>(`/trades/${id}/reject`, { reason })
  },

  complete(id: number) {
    return request.put<ApiResponse>(`/trades/${id}/complete`)
  },

  cancel(id: number) {
    return request.put<ApiResponse>(`/trades/${id}/cancel`)
  },

  counter(id: number, data: { counter_item_id?: number; counter_coin_amount?: number; message?: string }) {
    return request.put<ApiResponse>(`/trades/${id}/counter`, data)
  },

  acceptCounter(id: number) {
    return request.put<ApiResponse>(`/trades/${id}/counter/accept`)
  },

  rejectCounter(id: number, reason: string) {
    return request.put<ApiResponse>(`/trades/${id}/counter/reject`, { reason })
  },
}
