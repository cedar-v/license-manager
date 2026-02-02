import Axios from './https'

// 企业线索接口类型定义
export interface Lead {
  id: string;
  lead_no: string;
  company_name: string;
  contact_name: string;
  contact_phone: string;
  contact_email: string;
  requirement: string;
  extra_info: string;
  status: 'pending' | 'contacting' | 'completed' | 'rejected';
  follow_up_date: string | null;
  follow_up_record: string;
  internal_note: string;
  created_at: string;
  updated_at: string;
}

// 查询请求参数
export interface LeadQueryRequest {
  page?: number;
  page_size?: number;
  search?: string;
  status?: string;
}

// 通用响应接口
export interface ApiResponse<T = any> {
  code: string;
  message: string;
  timestamp: string;
  data?: T;
}

// 分页数据响应接口
export interface LeadQueryResponse {
  code: string;
  message: string;
  timestamp: string;
  data: {
    leads: Lead[];
    total_count: number;
    page: number;
    page_size: number;
  };
}

/**
 * 获取企业线索列表
 */
export function getLeads(params?: LeadQueryRequest): Promise<LeadQueryResponse> {
  return Axios.get('/api/leads', { params })
}

/**
 * 获取企业线索统计数据 (假设有这个接口，如果没有则前端计算或者等后续补充)
 */
export function getLeadSummary(): Promise<ApiResponse<any>> {
  return Axios.get('/api/leads/summary')
}

/**
 * 获取企业线索详情
 */
export function getLeadDetail(id: string | number): Promise<ApiResponse<Lead>> {
  return Axios.get(`/api/leads/${id}`)
}

/**
 * 更新企业线索
 */
export function updateLead(id: string | number, data: Partial<Lead>): Promise<ApiResponse<Lead>> {
  return Axios.put(`/api/leads/${id}`, data)
}

/**
 * 删除企业线索
 */
export function deleteLead(id: string | number): Promise<ApiResponse<any>> {
  return Axios.delete(`/api/leads/${id}`)
}



