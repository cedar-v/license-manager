/*
 * @Author: 13895237362 2205451508@qq.com
 * @Date: 2025-09-23 17:21:03
 * @LastEditors: 13895237362 2205451508@qq.com
 * @LastEditTime: 2025-09-24 16:18:40
 * @FilePath: \frontend\src\api\dashboard.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import Axios from './https/index'

// 最近授权数据类型定义

// API响应类型
export interface RecentAuthorizationItem {
  id: number
  customer_name: string
  description: string
  status: number // 1代表有效，0代表失效
  end_date: string
  created_at: string
}

// 请求格式
export interface ApiResponse<T> {
  code: string
  message: string
  data: T
}

// 最近授权响应类型
export interface RecentAuthorizationResponse {
  list: RecentAuthorizationItem[]
  total: number
}

// 获取最近授权列表 
export const getRecentAuthorizations = (): Promise<ApiResponse<RecentAuthorizationResponse>> => {
  return Axios.get('api/v1/dashboard/recent-authorizations')
}

// 授权趋势数据类型定义
export interface TrendDataItem {
  date: string
  total_authorizations: number
}

export interface AuthorizationTrendResponse {
  'trend_data': TrendDataItem[]
}

// 获取授权趋势数据
export const getAuthorizationTrend = (params: {
  type:string,
  start_date: string
  end_date: string
}): Promise<ApiResponse<AuthorizationTrendResponse>> => {
  return Axios.get('/api/v1/dashboard/authorization-trend', { params })
}

// 获取仪表盘统计数据（新接口，返回卡片区数据）
export const getOverviewStats = () => {
  return Axios.get('/api/v1/stats/overview')
}