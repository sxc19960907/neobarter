// API 通用响应
export interface ApiResponse<T = unknown> {
  code: number
  message: string
  data: T
}

export interface PageData<T> {
  list: T[]
  total: number
  page: number
  page_size: number
}

// 用户
export interface User {
  id: number
  phone: string
  nickname: string
  avatar_url: string
  user_type: 'personal' | 'enterprise'
  status: string
  credit_score: number
  real_name_verified: boolean
  enterprise_name: string
  enterprise_verified: boolean
  location: string
  bio: string
  last_login_at: string
  created_at: string
}

export interface UserAddress {
  id: number
  user_id: number
  name: string
  phone: string
  province: string
  city: string
  district: string
  detail: string
  is_default: boolean
}

// 钱包
export interface Wallet {
  id: number
  user_id: number
  balance: string
  frozen_balance: string
  total_income: string
  total_expense: string
}

export interface WalletTransaction {
  id: number
  wallet_id: number
  type: string
  amount: string
  balance_after: string
  reference_type: string
  reference_id: number
  description: string
  created_at: string
}

// 物品
export interface Category {
  id: number
  name: string
  parent_id: number | null
  icon: string
  sort_order: number
}

export interface Item {
  id: number
  user_id: number
  title: string
  description: string
  category_id: number
  estimated_value: string
  condition: string
  images: string[]
  video_url: string
  status: string
  location: string
  view_count: number
  want_items: string[]
  created_at: string
  updated_at: string
  user?: User
  category?: Category
}

// 交易
export interface TradeRequest {
  id: number
  initiator_id: number
  target_user_id: number
  target_item_id: number
  offered_item_id: number | null
  barter_coin_amount: string
  status: string
  message: string
  reject_reason: string
  expired_at: string
  completed_at: string
  created_at: string
  initiator?: User
  target_user?: User
  target_item?: Item
  offered_item?: Item
}

// 消息
export interface Conversation {
  id: number
  conversation_id: number
  user_id: number
  unread_count: number
  last_read_at: string
  user?: User
}

export interface Message {
  id: number
  conversation_id: number
  sender_id: number
  content: string
  message_type: string
  extra_data: string | null
  created_at: string
  sender?: User
}

// 评价
export interface Review {
  id: number
  trade_request_id: number
  reviewer_id: number
  reviewee_id: number
  rating: number
  comment: string
  created_at: string
  reviewer?: User
}

// 通知
export interface Notification {
  id: number
  user_id: number
  type: string
  title: string
  content: string
  reference_type: string
  reference_id: number
  is_read: boolean
  created_at: string
}
