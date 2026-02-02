import Axios from './https'

// 发票详情接口类型定义
export interface Invoice {
  id: string;
  invoice_no: string;
  applicant_name: string;
  applicant_phone: string;
  receiver_email: string;
  order_no: string;
  order_id: string;
  order_package_name: string;
  created_at: string;
  updated_at: string;
  status: 'pending' | 'success' | 'rejected';
  status_display: string;
  amount: number;
  title: string;
  taxpayer_id?: string;
  content?: string;
  remark?: string;
  invoice_type: string;
  invoice_type_display: string;
  // 审核/上传相关
  reject_reason?: string;
  rejected_at?: string;
  rejected_by?: string;
  uploaded_at?: string;
  uploaded_by?: string;
  uploader_name?: string;
  invoice_file_url?: string;
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
    invoices: Invoice[];
    total_count: number;
    page: number;
    page_size: number;
    total_pages?: number;
  };
}


// 发票统计数据接口响应
export interface InvoiceSummaryResponse {
  code: string;
  message: string;
  timestamp: string;
  data: {
    total_count: number;
    pending_count: number;
    completed_count: number;
    rejected_count: number;
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
export function uploadInvoice(data: FormData): Promise<ApiResponse> {
  return Axios.post('/api/v1/invoices/upload', data, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

/**
 * 驳回发票申请
 */
export function rejectInvoice(id: string, data: any): Promise<ApiResponse> {
  return Axios.post(`/api/v1/invoices/${id}/reject`, data)
}

/**
 * 管理员发票开票
 */
export function issueInvoice(id: string, data: { invoice_file_url: string, issued_at: string }): Promise<ApiResponse> {
  return Axios.post(`/api/v1/invoices/${id}/issue`, data)
}

