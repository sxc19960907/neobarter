import request from '@/utils/request'
import type { ApiResponse } from '@/types'

export const uploadApi = {
  /** 上传单张图片，返回可访问 URL */
  uploadImage(file: File) {
    const formData = new FormData()
    formData.append('file', file)
    return request.post<ApiResponse<{ url: string }>>('/upload/image', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
  },
}
