import request from '@/utils/request'
import type { ApiResponse, PageData, Conversation, Message, Wallet, WalletTransaction, Review, Notification } from '@/types'

export const messageApi = {
  listConversations() {
    return request.get<ApiResponse<Conversation[]>>('/messages/conversations')
  },

  getMessages(conversationId: number, params?: { page?: number; page_size?: number }) {
    return request.get<ApiResponse<PageData<Message>>>(`/messages/conversations/${conversationId}`, { params })
  },

  send(data: { conversation_id?: number; receiver_id?: number; content: string; message_type?: string }) {
    return request.post<ApiResponse<Message>>('/messages', data)
  },

  markRead(conversationId: number) {
    return request.put<ApiResponse>(`/messages/conversations/${conversationId}/read`)
  },
}

export const walletApi = {
  getWallet() {
    return request.get<ApiResponse<Wallet>>('/wallet')
  },

  listTransactions(params?: { page?: number; page_size?: number }) {
    return request.get<ApiResponse<PageData<WalletTransaction>>>('/wallet/transactions', { params })
  },
}

export const reviewApi = {
  create(data: { trade_request_id: number; reviewee_id: number; rating: number; comment?: string }) {
    return request.post<ApiResponse<Review>>('/reviews', data)
  },

  listByUser(userId: number, params?: { page?: number; page_size?: number }) {
    return request.get<ApiResponse<PageData<Review>>>(`/reviews/user/${userId}`, { params })
  },
}

export const notificationApi = {
  list(params?: { page?: number; page_size?: number }) {
    return request.get<ApiResponse<PageData<Notification>>>('/notifications', { params })
  },

  unreadCount() {
    return request.get<ApiResponse<{ count: number }>>('/notifications/unread-count')
  },

  markRead(id: number) {
    return request.put<ApiResponse>(`/notifications/${id}/read`)
  },

  markAllRead() {
    return request.put<ApiResponse>('/notifications/read-all')
  },
}
