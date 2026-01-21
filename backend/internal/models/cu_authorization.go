package models

// CuAuthorizationCodeListRequest 用户端授权码列表查询请求
type CuAuthorizationCodeListRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`              // 页码，默认1
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"` // 每页数量，默认10，最大100
	Status   string `form:"status" binding:"omitempty,oneof=normal locked expired"`
	Search   string `form:"search" binding:"omitempty"` // 授权码模糊匹配（code LIKE %search%）
}

// CuAuthorizationCodeListItem 用户端授权码列表项
type CuAuthorizationCodeListItem struct {
	ID                   string `json:"id"`
	Code                 string `json:"code"`
	Status               string `json:"status"`
	StatusDisplay        string `json:"status_display,omitempty"`
	MaxActivations       int    `json:"max_activations"`
	CurrentActivations   int    `json:"current_activations"`
	RemainingActivations int    `json:"remaining_activations"`
	CreatedAt            string `json:"created_at"`
	EndDate              string `json:"end_date"`
}

// CuAuthorizationCodeListResponse 用户端授权码列表响应
type CuAuthorizationCodeListResponse struct {
	List       []CuAuthorizationCodeListItem `json:"list"`
	Total      int64                         `json:"total"`
	Page       int                           `json:"page"`
	PageSize   int                           `json:"page_size"`
	TotalPages int                           `json:"total_pages"`
}

// CuAuthorizationCodeSummaryResponse 用户端授权信息统计响应
type CuAuthorizationCodeSummaryResponse struct {
	TotalCount             int64 `json:"total_count"`
	ExpiredCount           int64 `json:"expired_count"`
	ValidCount             int64 `json:"valid_count"`
	ValidMaxActivationsSum int64 `json:"valid_max_activations_sum"`
}
