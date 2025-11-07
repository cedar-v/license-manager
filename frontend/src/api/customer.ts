import Axios from './https'

// 授权统计信息类型
export interface AuthorizationStats {
  active_licenses?: number;           // 已激活许可证数量
  expired_auth_codes?: number;        // 已过期授权码数量
  expired_licenses?: number;          // 已过期许可证数量
  expiring_soon_auth_codes?: number;  // 30日内即将到期授权码数量
  inactive_licenses?: number;         // 未激活许可证数量
  total_auth_codes?: number;          // 总授权码数量
  total_licenses?: number;            // 总许可证数量
}

// 客户接口类型定义 - 完全使用接口提供的字段名称
export interface Customer {
  id: string;
  customer_code: string;
  customer_name: string;
  customer_type: string;
  customer_type_display: string;
  contact_person: string;
  contact_title?: string;
  email: string;
  phone?: string;
  address?: string;
  customer_level: string;
  customer_level_display: string;
  status: string;
  status_display: string;
  company_size?: string;
  company_size_display?: string;
  description?: string;
  created_at: string;
  updated_at: string;
  created_by: string;
  updated_by: string;
  authorization_stats?: AuthorizationStats; // 授权统计信息
}

// 客户查询请求参数
export interface CustomerQueryRequest {
  page?: number;
  page_size?: number;
  search?: string;
  customer_name?: string;
  customer_type?: 'individual' | 'enterprise' | 'government' | 'education';
  customer_level?: string;
  status?: 'active' | 'disabled';
  sort?: 'created_at' | 'updated_at' | 'customer_name' | 'customer_code';
  order?: 'asc' | 'desc';
}

// 创建客户请求参数
export interface CustomerCreateRequest {
  customer_name: string;
  customer_type: 'individual' | 'enterprise' | 'government' | 'education';
  contact_person: string;
  contact_title?: string;
  email?: string;
  phone?: string;
  address?: string;
  customer_level: string;
  status: 'active' | 'disabled';
  company_size?: string;
  description?: string;
}

// 更新客户请求参数
export interface CustomerUpdateRequest {
  customer_name?: string;
  customer_type?: 'individual' | 'enterprise' | 'government' | 'education';
  contact_person?: string;
  contact_title?: string;
  email?: string;
  phone?: string;
  address?: string;
  customer_level?: string;
  status?: 'active' | 'disabled';
  company_size?: string;
  description?: string;
}

// 状态更新请求参数
export interface CustomerStatusUpdateRequest {
  status: 'active' | 'disabled';
}

// 客户查询响应
export interface CustomerQueryResponse {
  code: string;
  message: string;
  timestamp: string;
  data: {
    list: Customer[];
    total: number;
    page: number;
    page_size: number;
    total_pages: number;
  };
}

// 客户详情响应
// 通用响应 - 泛型类型
export interface ApiResponse<T = any> {
  code: string;
  message: string;
  timestamp: string;
  data?: T;
}

/**
 * 查询客户列表
 */
export function getCustomers(params?: CustomerQueryRequest): Promise<CustomerQueryResponse> {
  return Axios.get('/api/customers', { params })
}

/**
 * 获取客户详情
 */
export function getCustomerDetail(id: string): Promise<ApiResponse<Customer>> {
  return Axios.get(`/api/customers/${id}`)
}

/**
 * 创建客户
 */
export function createCustomer(data: CustomerCreateRequest): Promise<ApiResponse<Customer>> {
  return Axios.post('/api/customers', data)
}

/**
 * 更新客户
 */
export function updateCustomer(id: string, data: CustomerUpdateRequest): Promise<ApiResponse<Customer>> {
  return Axios.put(`/api/customers/${id}`, data)
}

/**
 * 删除客户
 */
export function deleteCustomer(id: string): Promise<ApiResponse> {
  return Axios.delete(`/api/customers/${id}`)
}

/**
 * 切换客户状态
 */
export function toggleCustomerStatus(id: string, status: 'active' | 'disabled'): Promise<ApiResponse<Customer>> {
  return Axios.patch(`/api/customers/${id}/status`, { status })
}