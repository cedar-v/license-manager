import type { AxiosResponse } from 'axios'
import Axios from './https'
// 授权码(AuthorizationCode)类型 - 供授权管理页面使用
export interface AuthorizationCode {
  id: string;
  code: string;
  customer_id: string;
  customer_name?: string;
  software_id?: string;
  description?: string;
  status: 'normal' | 'locked' | 'expired';
  status_display?: string;
  created_at: string;
  updated_at?: string;
  start_date?: string;
  created_by?: string;
  end_date?: string;
  deployment_type: string;
  deployment_type_display?: string;
  encryption_type?: string;
  encryption_type_display?: string;
  software_version?: string;
  max_activations: number;
  current_activations?: number;
  activated_licenses_count?: number;
  is_locked?: boolean;
  feature_config?: any;
  usage_limits?: any;
  custom_parameters?: any;
  customer_info?: {
    customer_code?: string;
    customer_name?: string;
  }
}

// 许可证（License）类型 - 供其他页面使用（保持兼容）
export interface License {
  id: string;
  license_code: string;
  customer_id: string;
  customer_name: string;
  description: string;
  status: 'active' | 'inactive' | 'expired';
  status_display: string;
  created_at: string;
  updated_at: string;
  expiry_date: string;
  activation_date?: string;
  license_type: string;
  license_type_display: string;
  max_users?: number;
  current_users?: number;
  features: string[];
  license_key: string;
  hardware_id?: string;
  ip_restrictions?: string[];
  created_by: string;
  updated_by: string;
}

// 授权查询请求参数
export interface LicenseQueryRequest {
  page?: number;
  page_size?: number;
  customer_id?: string;
  status?: 'normal' | 'locked' | 'expired';
  sort?: 'created_at' | 'updated_at' | 'code';
  order?: 'asc' | 'desc';
  start_date?: string; // 创建时间开始
  end_date?: string; // 创建时间结束
  expiry_start_date?: string; // 到期时间开始
  expiry_end_date?: string; // 到期时间结束
  code?: string; // 授权码搜索
}

// 创建授权请求参数
export interface AuthorizationCodeCreateRequest {
  customer_id: string;
  software_id?: string;
  description?: string;
  validity_days: number;
  deployment_type: string;
  encryption_type?: string;
  software_version?: string;
  max_activations: number;
  feature_config?: any;
  usage_limits?: any;
  custom_parameters?: any;
}

// 更新授权请求参数
export interface LicenseUpdateRequest {
  software_id?: string;
  description?: string;
  validity_days?: number;
  deployment_type?: string;
  encryption_type?: string;
  software_version?: string;
  max_activations?: number;
  feature_config?: any;
  usage_limits?: any;
  custom_parameters?: any;
  change_type: string;
  reason?: string;
}

// 授权统计数据
export interface LicenseStats {
  total: number;
  active: number;
  inactive: number;
  expired: number;
  expiring_soon: number;
}

// 授权查询响应
export interface LicenseQueryResponse {
  code: string;
  message: string;
  timestamp: string;
  data: {
    list: AuthorizationCode[];
    total: number;
    page: number;
    page_size: number;
    total_pages: number;
  };
}

// 通用响应 - 泛型类型
export interface ApiResponse<T = any> {
  code: string;
  message: string;
  timestamp: string;
  data?: T;
}

/**
 * 查询授权列表
 */
export function getLicenses(params?: LicenseQueryRequest): Promise<LicenseQueryResponse> {
  return Axios.get('/api/v1/authorization-codes', { params })
}
/**
 * 获取授权详情
 */
export function getLicenseDetail(id: string): Promise<ApiResponse<AuthorizationCode>> {
  return Axios.get(`/api/v1/authorization-codes/${id}`)
}

/**
 * 创建授权
 */
export function createLicense(data: AuthorizationCodeCreateRequest): Promise<ApiResponse<AuthorizationCode>> {
  return Axios.post('/api/v1/authorization-codes', data)
}

/**
 * 更新授权
 */
export function updateLicense(id: string, data: LicenseUpdateRequest): Promise<ApiResponse<AuthorizationCode>> {
  return Axios.put(`/api/v1/authorization-codes/${id}`, data)
}

/**
 * 删除授权
 */
export function deleteLicense(id: string): Promise<ApiResponse> {
  return Axios.delete(`/api/v1/authorization-codes/${id}`)
}

/**
 * 锁定/解锁授权码
 */
export function lockAuthorizationCode(
  id: string,
  data: { is_locked: boolean; lock_reason?: string; reason?: string }
): Promise<ApiResponse<AuthorizationCode>> {
  return Axios.put(`/api/v1/authorization-codes/${id}/lock`, data)
}


/**
 * 获取授权统计数据
 */
export function getLicenseStats(): Promise<ApiResponse<LicenseStats>> {
  return Axios.get('/api/v1/authorization-codes/stats')
}

/**
 * 根据客户获取授权列表
 */
export function getLicensesByCustomer(customerId: string, params?: Omit<LicenseQueryRequest, 'customer_id'>): Promise<LicenseQueryResponse> {
  return Axios.get(`/api/customers/${customerId}/licenses`, { params })
}

/**
 * 续期授权
 */
export function renewLicense(id: string, expiry_date: string): Promise<ApiResponse<AuthorizationCode>> {
  return Axios.patch(`/api/v1/authorization-codes/${id}/renew`, { expiry_date })
}

/**
 * 激活授权
 */
export function activateLicense(id: string): Promise<ApiResponse<AuthorizationCode>> {
  return Axios.patch(`/api/v1/authorization-codes/${id}/activate`)
}

/**
 * 停用授权
 */
export function deactivateLicense(id: string): Promise<ApiResponse<AuthorizationCode>> {
  return Axios.patch(`/api/v1/authorization-codes/${id}/deactivate`)
}

// 许可证设备信息类型（与后端 LicenseListItem 对应）
export interface LicenseDevice {
  id: string
  license_key: string
  authorization_code_id: string
  authorization_code: string
  customer_name: string
  hardware_fingerprint: string
  status: 'active' | 'inactive' | 'revoked'
  status_display?: string
  is_online: boolean
  is_online_display?: string
  activation_ip?: string
  last_online_ip?: string
  activated_at?: string
  last_heartbeat?: string,
  device_info?: {
    machine_code?: string    // 机器序列号
    mac_address?: string     // MAC地址
    cpu_id?: string          // CPU ID
    [key: string]: any
  }
}

// 许可证详情类型（与后端 LicenseDetailResponse 对应）
export interface LicenseDetail {
  id: string
  license_key: string
  authorization_code_id: string
  authorization_code: string
  customer_id: string
  customer_name: string
  hardware_fingerprint: string
  device_info?: {
    machine_code?: string
    mac_address?: string
    cpu_id?: string
    [key: string]: any
  }
  activation_ip?: string
  status: string
  status_display?: string
  is_online: boolean
  is_online_display?: string
  activated_at?: string
  last_heartbeat?: string
  last_online_ip?: string
  config_updated_at?: string
  usage_data?: any
  created_at: string
  updated_at: string
}

// 许可证列表查询请求
export interface LicenseListQueryRequest {
  authorization_code_id?: string
  page?: number
  page_size?: number
  status?: 'active' | 'inactive' | 'expired'
}

// 许可证列表响应
export interface LicenseListResponse {
  code: string
  message: string
  timestamp: string
  data: {
    list: LicenseDevice[]
    total: number
    page: number
    page_size: number
    total_pages: number
  }
}

/**
 * 获取许可证列表（已激活的设备列表）
 */
export function getLicenseDevices(params: LicenseListQueryRequest): Promise<LicenseListResponse> {
  return Axios.get('/api/v1/licenses', { params })
}

/**
 * 获取许可证详情
 */
export function getLicenseDeviceDetail(id: string): Promise<ApiResponse<LicenseDetail>> {
  return Axios.get(`/api/v1/licenses/${id}`)
}

/**
 * 下载许可证文件
 */
export function downloadLicenseFile(id: string): Promise<AxiosResponse<Blob>> {
  return Axios.get(`/api/v1/licenses/${id}/download`, {
    responseType: 'blob'
  })
}

// 创建许可证请求（手动绑定设备）
export interface LicenseDeviceCreateRequest {
  authorization_code_id: string
  hardware_fingerprint: string
  device_info?: Record<string, any>
  activation_ip?: string
}

/**
 * 手动创建许可证
 */
export function createLicenseDevice(
  data: LicenseDeviceCreateRequest
): Promise<ApiResponse<LicenseDevice>> {
  return Axios.post('/api/v1/licenses', data)
}

// 授权变更历史相关类型定义
export interface AuthorizationChangeItem {
  id: string
  change_type: string
  change_type_display?: string
  operator_id: string
  operator_name?: string
  reason?: string
  created_at: string
}

export interface AuthorizationChangeListRequest {
  page?: number
  page_size?: number
  change_type?: string
  operator_id?: string
  start_date?: string
  end_date?: string
  sort?: 'created_at' | 'change_type'
  order?: 'asc' | 'desc'
}

export interface AuthorizationChangeListResponse {
  code: string
  message: string
  timestamp: string
  data: {
    list: AuthorizationChangeItem[]
    total: number
    page: number
    page_size: number
    total_pages: number
  }
}

/**
 * 获取授权变更历史列表
 */
export function getAuthorizationChanges(
  id: string,
  params?: AuthorizationChangeListRequest
): Promise<AuthorizationChangeListResponse> {
  return Axios.get(`/api/v1/authorization-codes/${id}/changes`, { params })
}