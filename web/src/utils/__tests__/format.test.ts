import { describe, it, expect } from 'vitest'
import {
  conditionLabel,
  tradeStatusLabel,
  formatBarterCoin,
  formatDateTime,
  isValidPhone,
} from '../format'

describe('formatBarterCoin', () => {
  it('应格式化字符串金额为两位小数', () => {
    expect(formatBarterCoin('100')).toBe('100.00')
    expect(formatBarterCoin('99.5')).toBe('99.50')
  })

  it('应格式化数字金额', () => {
    expect(formatBarterCoin(50)).toBe('50.00')
  })

  it('无效值应返回 0.00', () => {
    expect(formatBarterCoin('abc')).toBe('0.00')
  })
})

describe('formatDateTime', () => {
  it('应截取到分钟', () => {
    expect(formatDateTime('2026-05-31T10:30:45.123Z')).toBe('2026-05-31 10:30')
  })

  it('空字符串应返回空', () => {
    expect(formatDateTime('')).toBe('')
  })
})

describe('isValidPhone', () => {
  it('合法手机号应返回 true', () => {
    expect(isValidPhone('13800138000')).toBe(true)
    expect(isValidPhone('19912345678')).toBe(true)
  })

  it('非法手机号应返回 false', () => {
    expect(isValidPhone('12345')).toBe(false)
    expect(isValidPhone('10000000000')).toBe(false)
    expect(isValidPhone('1380013800')).toBe(false) // 10位
    expect(isValidPhone('abcdefghijk')).toBe(false)
  })
})

describe('labels', () => {
  it('成色标签映射正确', () => {
    expect(conditionLabel.new).toBe('全新')
    expect(conditionLabel.like_new).toBe('几乎全新')
  })

  it('交易状态标签映射正确', () => {
    expect(tradeStatusLabel.pending).toBe('待确认')
    expect(tradeStatusLabel.completed).toBe('已完成')
  })
})
