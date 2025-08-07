package models

// APIResponse 通用API响应结构
type APIResponse struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp string      `json:"timestamp,omitempty"`
}

// ErrorResponse 错误响应结构
type ErrorResponse struct {
	Code      string `json:"code"`      // 业务错误码，如 "100001", "900001"
	Message   string `json:"message"`   // 错误描述信息
	Timestamp string `json:"timestamp"` // 错误发生时间
}
