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
	ID                   string `json:"id"`                       // 授权码ID
	Code                 string `json:"code"`                     // 授权码
	Status               string `json:"status"`                   // 状态
	StatusDisplay        string `json:"status_display,omitempty"` // 状态显示
	MaxActivations       int    `json:"max_activations"`          // 最大激活数量
	CurrentActivations   int    `json:"current_activations"`      // 当前激活数量
	RemainingActivations int    `json:"remaining_activations"`    // 剩余激活数量
	CreatedAt            string `json:"created_at"`               // 创建时间
	EndDate              string `json:"end_date"`                 // 到期时间
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
	TotalCount             int64 `json:"total_count"`               // 总授权码数量
	ExpiredCount           int64 `json:"expired_count"`             // 已过期授权码数量
	ValidCount             int64 `json:"valid_count"`               // 有效授权码数量
	ValidMaxActivationsSum int64 `json:"valid_max_activations_sum"` // 有效授权码的 `max_activations` 之和
}
