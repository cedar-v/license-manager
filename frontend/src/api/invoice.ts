import Axios from './https'

// 发票详情接口类型定义
export interface Invoice {
  id: string;
  invoiceNo: string;
  user: string;
  userType: string;
  userTypeLabel: string;
  orderNo: string;
  time: string;
  status: 'pending' | 'success' | 'rejected' | 'completed';
  statusLabel: string;
  amount: number;
  invoiceTitle: string;
  taxId?: string;
  content?: string;
  remark?: string;
  // 详情字段
  applyTime?: string;
  phone?: string;
  email?: string;
  invoiceType?: string;
  rejectReason?: string;
  rejectTime?: string;
  rejectUser?: string;
  finishTime?: string;
  approveUser?: string;
  approveTime?: string;
  fileName?: string;
  fileSize?: string;
}

// 发票查询请求参数
export interface InvoiceQueryRequest {
  page?: number;
  page_size?: number;
  keyword?: string;
  status?: string;
  start_date?: string;
  end_date?: string;
}

// 通用响应接口
export interface ApiResponse<T = any> {
  code: string;
  message: string;
  timestamp: string;
  data?: T;
}

// 分页数据响应接口
export interface InvoiceQueryResponse {
  code: string;
  message: string;
  timestamp: string;
  data: {
    list: Invoice[];
    total: number;
    page: number;
    page_size: number;
    total_pages: number;
  };
}

// 发票统计数据接口响应
export interface InvoiceSummaryResponse {
  code: string;
  message: string;
  timestamp: string;
  data: {
    total: number;
    pending: number;
    completed: number;
    rejected: number;
  };
}

/**
 * 查询发票列表
 */
export function getInvoices(params?: InvoiceQueryRequest): Promise<InvoiceQueryResponse> {
  return Axios.get('/api/v1/invoices', { params })
}

/**
 * 获取发票统计数据
 */
export function getInvoiceSummary(): Promise<InvoiceSummaryResponse> {
  return Axios.get('/api/v1/invoices/summary')
}


/**
 * 获取发票详情
 */
export function getInvoiceDetail(id: string): Promise<ApiResponse<Invoice>> {
  return Axios.get(`/api/v1/invoices/${id}`)
}

/**
 * 上传发票文件
 */
export function uploadInvoice(id: string, data: any): Promise<ApiResponse> {
  return Axios.post(`/api/v1/invoices/${id}/upload`, data)
}

/**
 * 驳回发票申请
 */
export function rejectInvoice(id: string, data: any): Promise<ApiResponse> {
  return Axios.post(`/api/v1/invoices/${id}/reject`, data)
}
