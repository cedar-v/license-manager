import Axios from './https'

// 原始枚举项接口（来自后端）
export interface RawEnumItem {
  key: string
  display: string
}

// 原始枚举数据接口（来自后端）
export interface RawEnumData {
  type: string
  items: RawEnumItem[]
}

// 枚举响应接口
export interface EnumResponse {
  code: string
  message: string
  timestamp?: string
  data: RawEnumData
}

/**
 * 根据枚举类型获取对应的多语言显示值
 * @param type 枚举类型，如: customer_type, customer_level, status, company_size 等
 */
export function getEnumOptions(type: string): Promise<EnumResponse> {
  return Axios.get(`/api/enums/${type}`)
}

/**
 * 获取客户类型枚举
 */
export function getCustomerTypeEnums(): Promise<EnumResponse> {
  return getEnumOptions('customer_type')
}

/**
 * 获取客户等级枚举
 */
export function getCustomerLevelEnums(): Promise<EnumResponse> {
  return getEnumOptions('customer_level')
}

/**
 * 获取状态枚举
 */
export function getStatusEnums(): Promise<EnumResponse> {
  return getEnumOptions('customer_status')
}

/**
 * 获取企业规模枚举
 */
export function getCompanySizeEnums(): Promise<EnumResponse> {
  return getEnumOptions('company_size')
}

/**
 * 获取授权码状态枚举
 */
export function getAuthorizationStatusEnums(): Promise<EnumResponse> {
  return getEnumOptions('authorization_status')
}