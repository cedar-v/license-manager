package models

// ========== 管理端发票请求/响应结构体 ==========

// InvoiceRejectRequest 发票驳回请求结构（管理端）
type InvoiceRejectRequest struct {
	RejectReason string `json:"reject_reason" binding:"required,max=500"` // 驳回原因，必填，最大500字符
	Suggestion   string `json:"suggestion" binding:"required,max=500"`    // 修改建议，必填，最大500字符
}

// InvoiceIssueRequest 发票开票请求结构（管理端）
type InvoiceIssueRequest struct {
	InvoiceFileURL string `json:"invoice_file_url" binding:"required,max=500"` // 发票文件URL，必填，最大500字符
	IssuedAt       string `json:"issued_at" binding:"required"`                // 开票时间，必填，ISO 8601格式时间字符串
}

// InvoiceUploadRequest 发票文件上传请求结构（管理端）
type InvoiceUploadRequest struct {
	InvoiceNo string `form:"invoice_no" binding:"required"` // 发票申请号，必填
}

// InvoiceUploadResponse 发票文件上传响应结构（管理端）
type InvoiceUploadResponse struct {
	FileURL string `json:"file_url"` // 上传后的文件URL
}

// InvoiceDetailResponse 发票详情响应结构（包含关联信息，管理端使用）
type InvoiceDetailResponse struct {
	*InvoiceResponse
	OrderPackageName string `json:"order_package_name,omitempty"` // 关联订单的套餐名称
	ApplicantName    string `json:"applicant_name,omitempty"`     // 发票申请人真实姓名
	ApplicantPhone   string `json:"applicant_phone,omitempty"`    // 发票申请人电话号码
	UploaderName     string `json:"uploader_name,omitempty"`      // 发票文件上传人姓名（管理员）
	RejecterName     string `json:"rejecter_name,omitempty"`      // 发票驳回人姓名（管理员）
}
