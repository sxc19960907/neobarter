// 格式化工具函数

/** 物品成色中文标签 */
export const conditionLabel: Record<string, string> = {
  new: '全新',
  like_new: '几乎全新',
  good: '良好',
  fair: '一般',
}

/** 交易状态中文标签 */
export const tradeStatusLabel: Record<string, string> = {
  pending: '待确认',
  accepted: '已接受',
  rejected: '已拒绝',
  completed: '已完成',
  cancelled: '已取消',
  expired: '已过期',
}

/** 格式化巴特币金额 */
export function formatBarterCoin(value: string | number): string {
  const num = typeof value === 'string' ? parseFloat(value) : value
  if (isNaN(num)) return '0.00'
  return num.toFixed(2)
}

/** 格式化时间（截取到分钟） */
export function formatDateTime(time: string): string {
  if (!time) return ''
  return time.slice(0, 16).replace('T', ' ')
}

/** 校验手机号格式 */
export function isValidPhone(phone: string): boolean {
  return /^1[3-9]\d{9}$/.test(phone)
}
