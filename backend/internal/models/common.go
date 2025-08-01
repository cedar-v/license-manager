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
	Code      int    `json:"code"`      // HTTP状态码，如 401, 500
	Error     string `json:"error"`     // 业务错误码，如 AUTH_001, LOGIN_FAILED
	Message   string `json:"message"`   // 错误描述信息
	Timestamp string `json:"timestamp"` // 错误发生时间
}
