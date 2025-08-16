package models

// APIResponse 通用API响应结构
type APIResponse struct {
	Code      string      `json:"code"`
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

// 枚举相关结构

// EnumItem 枚举项结构
type EnumItem struct {
	Key     string `json:"key"`     // 枚举值
	Display string `json:"display"` // 多语言显示文本
}

// EnumTypeResponse 单个枚举类型响应
type EnumTypeResponse struct {
	Type  string     `json:"type"`  // 枚举类型名称
	Items []EnumItem `json:"items"` // 枚举项列表
}

// EnumListResponse 枚举列表响应
type EnumListResponse struct {
	Enums []EnumTypeResponse `json:"enums"` // 所有枚举类型
}
