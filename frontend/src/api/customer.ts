import Axios from './https'

// 客户接口类型定义 - 完全使用接口提供的字段名称
export interface Customer {
  id: string;
  customer_code: string;
  customer_name: string;
  customer_type: string;
  contact_person: string;
  email: string;
  customer_level: string;
  status: string;
  created_at: string;
}

// 客户查询请求参数
export interface CustomerQueryRequest {
  page?: number;
  page_size?: number;
  search?: string;
  customer_type?: 'individual' | 'enterprise' | 'government' | 'education';
  customer_level?: 'normal' | 'vip' | 'enterprise' | 'strategic';
  status?: 'active' | 'disabled';
}

// 客户查询响应
export interface CustomerQueryResponse {
  code: number;
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
export interface CustomerDetailResponse {
  code: number;
  message: string;
  timestamp: string;
  data: Customer;
}

// 客户创建/更新请求
export interface CustomerRequest {
  customer_code?: string;
  customer_name: string;
  customer_type: 'individual' | 'enterprise' | 'government' | 'education';
  contact_person: string;
  email: string;
  customer_level: 'normal' | 'vip' | 'enterprise' | 'strategic';
  status: 'active' | 'disabled';
}

// 通用响应
export interface ApiResponse {
  code: number;
  message: string;
  timestamp: string;
  data?: any;
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
export function getCustomerDetail(id: string): Promise<CustomerDetailResponse> {
  return Axios.get(`/api/customers/${id}`)
}

/**
 * 创建客户
 */
export function createCustomer(data: CustomerRequest): Promise<ApiResponse> {
  return Axios.post('/api/customers', data)
}

/**
 * 更新客户
 */
export function updateCustomer(id: string, data: CustomerRequest): Promise<ApiResponse> {
  return Axios.put(`/api/customers/${id}`, data)
}

/**
 * 删除客户
 */
export function deleteCustomer(id: string): Promise<ApiResponse> {
  return Axios.delete(`/api/customers/${id}`)
}

/**
 * 禁用/启用客户
 */
export function toggleCustomerStatus(id: string, status: 'active' | 'disabled'): Promise<ApiResponse> {
  return Axios.patch(`/api/customers/${id}/status`, { status })
}