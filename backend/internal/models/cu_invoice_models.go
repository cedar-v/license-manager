package models

// ========== 客户端发票请求/响应结构体 ==========

// InvoiceCreateRequest 客户端申请发票请求结构
// Content 非必填。
type InvoiceCreateRequest struct {
	OrderID       string  `json:"order_id" binding:"required"`                                           // 订单ID，必填
	InvoiceType   string  `json:"invoice_type" binding:"required,oneof=personal enterprise vat_special"` // 发票类型，必填，personal-个人/enterprise-企业普票/vat_special-增值税专用发票
	ReceiverEmail string  `json:"receiver_email" binding:"required,email"`                               // 收票邮箱，必填，格式为邮箱地址
	Title         string  `json:"title" binding:"required,max=200"`                                      // 发票抬头，必填，最多200字符
	TaxpayerID    *string `json:"taxpayer_id,omitempty"`                                                 // 纳税人识别号，企业发票必填，个人发票可为空
	Content       string  `json:"content" binding:"omitempty,max=200"`                                   // 开票内容，非必填，最多200字符
	Remark        *string `json:"remark,omitempty"`                                                      // 备注，可选
}

// InvoiceListRequest 发票列表查询请求结构
// 注：apply_date 支持 YYYY-MM-DD 或 RFC3339 时间范围 "start,end"。
type InvoiceListRequest struct {
	Page       int    `form:"page" binding:"omitempty,min=1"`              // 页码，默认1
	PageSize   int    `form:"page_size" binding:"omitempty,min=1,max=100"` // 每页条数，默认10，最大100
	Status     string `form:"status" binding:"omitempty"`                  // 状态筛选，pending-待处理/issued-已开票/rejected-已驳回
	Search     string `form:"search" binding:"omitempty"`                  // 搜索关键字，支持发票号或订单号模糊匹配
	ApplyDate  string `form:"apply_date" binding:"omitempty"`              // 申请日期筛选，格式YYYY-MM-DD 或 RFC3339范围start,end
	CustomerID string `form:"customer_id" binding:"omitempty"`             // 客户ID筛选（管理端）
}

// InvoiceListResponse 发票列表响应结构
type InvoiceListResponse struct {
	Invoices   []*InvoiceResponse `json:"invoices"`    // 发票列表
	TotalCount int64              `json:"total_count"` // 总记录数
	Page       int                `json:"page"`        // 当前页码
	PageSize   int                `json:"page_size"`   // 每页条数
}

// InvoiceSummaryResponse 发票汇总统计响应结构
type InvoiceSummaryResponse struct {
	TotalCount    int64 `json:"total_count"`    // 发票总数
	PendingCount  int64 `json:"pending_count"`  // 待处理发票数量
	IssuedCount   int64 `json:"issued_count"`   // 已开票发票数量
	RejectedCount int64 `json:"rejected_count"` // 已驳回发票数量
}
