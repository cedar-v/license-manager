import Axios from './https'

// 套餐接口类型定义
export interface Package {
  id: string | number;
  name: string;
  price: number;
  price_description?: string;
  duration_description: string;
  status: number; // 1: enabled, 2: disabled
  type: string;
  sort_order: number;
  description?: string;
  features?: any;
  remark?: string;
  created_at: string;
  updated_at: string;
}

// 套餐列表响应定义
export interface PackageListResponse {
  packages: Package[];
  total_count: number;
}

// 通用响应接口
export interface ApiResponse<T = any> {
  code: string;
  message: string;
  timestamp: string;
  data?: T;
}

/**
 * 获取套餐列表
 */
export function getPackages(): Promise<ApiResponse<PackageListResponse>> {
  return Axios.get('/api/packages')
}

/**
 * 获取套餐详情
 */
export function getPackageDetail(id: string | number): Promise<ApiResponse<Package>> {
  return Axios.get(`/api/packages/${id}`)
}

/**
 * 更新套餐
 */
export function updatePackage(id: string | number, data: Partial<Package>): Promise<ApiResponse<Package>> {
  return Axios.put(`/api/packages/${id}`, data)
}
