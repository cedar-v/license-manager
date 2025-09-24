/*
 * @Author: 13895237362 2205451508@qq.com
 * @Date: 2025-09-23 17:21:03
 * @LastEditors: 13895237362 2205451508@qq.com
 * @LastEditTime: 2025-09-23 17:30:24
 * @FilePath: \frontend\src\api\dashboard.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import Axios from './https/index'

// 最近授权数据类型定义
export interface RecentAuthorizationItem {
  id: number
  customer: string
  description: string
  status: number // 1代表有效，0代表失效
  expiry: string
  createTime: string
}

// API响应类型
export interface ApiResponse<T> {
  code: string
  message: string
  data: T
}

// 获取最近授权列表 - 暂时使用授权码接口
export const getRecentAuthorizations = (): Promise<ApiResponse<RecentAuthorizationItem[]>> => {
  return Axios.get('/api/v1/authorization-codes')
}

// 获取仪表盘统计数据（如果需要的话）
export const getDashboardStats = (): Promise<ApiResponse<any>> => {
  return Axios.get('/api/v1/dashboard/stats')
}