import request from '@/utils/request'
import type { ApiResponse, PageData, Item, Category } from '@/types'

export interface ItemQuery {
  page?: number
  page_size?: number
  category_id?: number
  keyword?: string
  location?: string
  condition?: string
  min_value?: number
  max_value?: number
  sort_by?: string
  mine?: boolean
}

export const itemApi = {
  create(data: Partial<Item>) {
    return request.post<ApiResponse<Item>>('/items', data)
  },

  list(params: ItemQuery) {
    return request.get<ApiResponse<PageData<Item>>>('/items', { params })
  },

  get(id: number) {
    return request.get<ApiResponse<Item>>(`/items/${id}`)
  },

  update(id: number, data: Partial<Item>) {
    return request.put<ApiResponse<Item>>(`/items/${id}`, data)
  },

  delete(id: number) {
    return request.delete<ApiResponse>(`/items/${id}`)
  },

  updateStatus(id: number, status: string) {
    return request.put<ApiResponse>(`/items/${id}/status`, { status })
  },

  listCategories() {
    return request.get<ApiResponse<Category[]>>('/categories')
  },
}
