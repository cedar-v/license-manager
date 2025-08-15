/**
 * 数据字典工具函数
 * 用于管理系统中的枚举类型和选项数据
 * 支持多语言国际化
 */

import i18n from '@/i18n'

// 选项接口定义
export interface DictionaryOption {
  label: string
  value: string | number
}

// 获取 i18n t 函数
const getT = () => {
  return i18n.global.t
}

// 客户类型字典
export const getCustomerTypeOptions = (): DictionaryOption[] => {
  const t = getT()
  return [
    { label: t('dictionary.common.all'), value: '' },
    { label: t('dictionary.customerType.individual'), value: 'individual' },
    { label: t('dictionary.customerType.enterprise'), value: 'enterprise' },
    { label: t('dictionary.customerType.government'), value: 'government' },
    { label: t('dictionary.customerType.education'), value: 'education' }
  ]
}

// 保持向后兼容
export const CUSTOMER_TYPE_OPTIONS: DictionaryOption[] = getCustomerTypeOptions()

// 客户等级字典
export const getCustomerLevelOptions = (): DictionaryOption[] => {
  const t = getT()
  return [
    { label: t('dictionary.common.all'), value: '' },
    { label: t('dictionary.customerLevel.normal'), value: 'normal' },
    { label: t('dictionary.customerLevel.vip'), value: 'vip' },
    { label: t('dictionary.customerLevel.enterprise'), value: 'enterprise' },
    { label: t('dictionary.customerLevel.strategic'), value: 'strategic' }
  ]
}

// 保持向后兼容
export const CUSTOMER_LEVEL_OPTIONS: DictionaryOption[] = getCustomerLevelOptions()

// 通用状态字典
export const getStatusOptions = (): DictionaryOption[] => {
  const t = getT()
  return [
    { label: t('dictionary.common.all'), value: '' },
    { label: t('dictionary.status.active'), value: 'active' },
    { label: t('dictionary.status.disabled'), value: 'disabled' }
  ]
}

// 保持向后兼容
export const STATUS_OPTIONS: DictionaryOption[] = getStatusOptions()

// 字典映射函数 - 根据值获取标签
export const getDictionaryLabel = (options: DictionaryOption[], value: string | number): string => {
  const t = getT()
  const option = options.find(item => item.value === value)
  return option?.label || t('dictionary.common.unknown')
}

// 客户类型映射函数
export const getCustomerTypeLabel = (value: string): string => {
  const t = getT()
  if (!value) return t('dictionary.common.all')
  return t(`dictionary.customerType.${value}`) || t('dictionary.common.unknown')
}

// 客户等级映射函数
export const getCustomerLevelLabel = (value: string): string => {
  const t = getT()
  if (!value) return t('dictionary.common.all')
  return t(`dictionary.customerLevel.${value}`) || t('dictionary.common.unknown')
}

// 状态映射函数
export const getStatusLabel = (value: string): string => {
  const t = getT()
  if (!value) return t('dictionary.common.all')
  return t(`dictionary.status.${value}`) || t('dictionary.common.unknown')
}

// 获取所有字典选项的函数（支持国际化）
export const getCustomerDictionaries = () => ({
  customerType: getCustomerTypeOptions(),
  customerLevel: getCustomerLevelOptions(),
  status: getStatusOptions()
})

// 保持向后兼容
export const CUSTOMER_DICTIONARIES = {
  customerType: CUSTOMER_TYPE_OPTIONS,
  customerLevel: CUSTOMER_LEVEL_OPTIONS,
  status: STATUS_OPTIONS
} as const