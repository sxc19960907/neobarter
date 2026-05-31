import request from '@/utils/request'
import type { ApiResponse } from '@/types'

export interface SearchResultItem {
  id: number
  title: string
  description: string
  category_id: number
  category_name: string
  estimated_value: number
  condition: string
  images: string[]
  location: string
  view_count: number
  user_nickname: string
  created_at: string
  highlight?: {
    title?: string[]
    description?: string[]
  }
}

export interface SearchResponse {
  items: SearchResultItem[]
  total: number
  page: number
  page_size: number
}

export interface SearchParams {
  keyword?: string
  category_id?: number
  condition?: string
  min_value?: number
  max_value?: number
  location?: string
  sort_by?: string
  page?: number
  page_size?: number
}

export const searchApi = {
  search(params: SearchParams) {
    return request.get<ApiResponse<SearchResponse>>('/search/items', { params })
  },

  suggest(q: string) {
    return request.get<ApiResponse<string[]>>('/search/suggest', { params: { q } })
  },
}
