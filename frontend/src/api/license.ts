import Axios from './https'

// 授权接口类型定义
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
  search?: string;
  customer_id?: string;
  status?: 'active' | 'inactive' | 'expired';
  license_type?: string;
  sort?: 'created_at' | 'updated_at' | 'expiry_date' | 'activation_date';
  order?: 'asc' | 'desc';
  start_date?: string;
  end_date?: string;
}

// 创建授权请求参数
export interface LicenseCreateRequest {
  customer_id: string;
  description: string;
  license_type: string;
  expiry_date: string;
  max_users?: number;
  features: string[];
  ip_restrictions?: string[];
  hardware_id?: string;
}

// 更新授权请求参数
export interface LicenseUpdateRequest {
  description?: string;
  status?: 'active' | 'inactive' | 'expired';
  expiry_date?: string;
  max_users?: number;
  features?: string[];
  ip_restrictions?: string[];
  hardware_id?: string;
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
    list: License[];
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
  return Axios.get('/api/licenses', { params })
}

/**
 * 获取授权详情
 */
export function getLicenseDetail(id: string): Promise<ApiResponse<License>> {
  return Axios.get(`/api/licenses/${id}`)
}

/**
 * 创建授权
 */
export function createLicense(data: LicenseCreateRequest): Promise<ApiResponse<License>> {
  return Axios.post('/api/licenses', data)
}

/**
 * 更新授权
 */
export function updateLicense(id: string, data: LicenseUpdateRequest): Promise<ApiResponse<License>> {
  return Axios.put(`/api/licenses/${id}`, data)
}

/**
 * 删除授权
 */
export function deleteLicense(id: string): Promise<ApiResponse> {
  return Axios.delete(`/api/licenses/${id}`)
}

/**
 * 激活授权
 */
export function activateLicense(id: string): Promise<ApiResponse<License>> {
  return Axios.patch(`/api/licenses/${id}/activate`)
}

/**
 * 停用授权
 */
export function deactivateLicense(id: string): Promise<ApiResponse<License>> {
  return Axios.patch(`/api/licenses/${id}/deactivate`)
}

/**
 * 获取授权统计数据
 */
export function getLicenseStats(): Promise<ApiResponse<LicenseStats>> {
  return Axios.get('/api/licenses/stats')
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
export function renewLicense(id: string, expiry_date: string): Promise<ApiResponse<License>> {
  return Axios.patch(`/api/licenses/${id}/renew`, { expiry_date })
}