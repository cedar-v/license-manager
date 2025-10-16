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
  end_date?: string;
  deployment_type: string;
  deployment_type_display?: string;
  encryption_type?: string;
  encryption_type_display?: string;
  software_version?: string;
  max_activations: number;
  current_activations?: number;
  is_locked?: boolean;
  feature_config?: any;
  usage_limits?: any;
  custom_parameters?: any;
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