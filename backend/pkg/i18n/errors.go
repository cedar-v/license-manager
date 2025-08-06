package i18n

import (
	"fmt"
)

// ErrorI18n 错误信息国际化封装
type ErrorI18n struct {
	manager I18nManager
}

// NewErrorI18n 创建新的错误国际化实例
func NewErrorI18n(manager I18nManager) *ErrorI18n {
	return &ErrorI18n{
		manager: manager,
	}
}

// NewI18nErrorResponse 创建多语言错误响应
// 扩展现有错误响应函数，支持多语言
// 返回：HTTP状态码、错误码、本地化错误消息
func NewI18nErrorResponse(code, lang string, customMessage ...string) (int, string, string) {
	// 获取本地化错误信息
	message := GetI18nErrorMessage(code, lang)
	
	// 支持自定义消息覆盖
	if len(customMessage) > 0 && customMessage[0] != "" {
		message = customMessage[0]
	}
	
	// 获取对应的HTTP状态码
	httpStatus := getHTTPStatusByCode(code)
	
	return httpStatus, code, message
}

// GetI18nErrorMessage 获取多语言错误信息
func GetI18nErrorMessage(code, lang string) string {
	if globalManager == nil {
		// 如果多语言系统未初始化，降级到默认错误处理
		return getDefaultErrorMessage(code)
	}
	
	return globalManager.GetErrorMessage(code, lang)
}

// GetI18nErrorMessageWithManager 使用指定管理器获取多语言错误信息
func GetI18nErrorMessageWithManager(manager I18nManager, code, lang string) string {
	if manager == nil {
		return getDefaultErrorMessage(code)
	}
	
	return manager.GetErrorMessage(code, lang)
}

// getDefaultErrorMessage 获取默认错误信息（作为降级方案）
func getDefaultErrorMessage(code string) string {
	// 这里可以集成原有的 ErrorMessages 映射作为降级方案
	defaultMessages := map[string]string{
		"000000": "Success",
		"100001": "Authentication expired",
		"100002": "Invalid authentication",
		"100003": "Invalid username or password",
		"100004": "Authentication required",
		"100005": "Insufficient permissions",
		"200001": "Customer not found",
		"200002": "Customer already exists",
		"200003": "Customer name is required",
		"200004": "Invalid customer information format",
		"300001": "Invalid license key",
		"300002": "License expired",
		"300003": "License key is required",
		"300004": "License key already used",
		"900001": "Invalid request parameters",
		"900002": "Resource not found",
		"900003": "Resource conflict",
		"900004": "Internal server error",
	}
	
	if msg, ok := defaultMessages[code]; ok {
		return msg
	}
	return "Unknown error"
}

// getHTTPStatusByCode 根据业务错误码获取对应的 HTTP 状态码
// 保持与原有错误处理逻辑一致
func getHTTPStatusByCode(code string) int {
	const (
		StatusOK                  = 200
		StatusBadRequest          = 400
		StatusUnauthorized        = 401
		StatusForbidden           = 403
		StatusNotFound            = 404
		StatusConflict            = 409
		StatusInternalServerError = 500
	)
	
	switch code {
	case "000000": // 成功
		return StatusOK
	case "900001": // 请求参数无效
		return StatusBadRequest
	case "100001", "100002", "100003", "100004": // 认证相关错误
		return StatusUnauthorized
	case "100005": // 权限不足
		return StatusForbidden
	case "900002", "200001", "300001": // 资源不存在
		return StatusNotFound
	case "900003", "200002", "300004": // 资源冲突
		return StatusConflict
	case "900004": // 服务器内部错误
		return StatusInternalServerError
	default:
		// 根据错误码前缀判断
		if len(code) >= 2 {
			prefix := code[:2]
			switch prefix {
			case "10": // 认证模块错误，默认401
				return StatusUnauthorized
			case "20": // 客户模块错误，默认404
				return StatusNotFound
			case "30": // 授权模块错误，默认404
				return StatusNotFound
			case "90": // 系统错误，默认500
				return StatusInternalServerError
			}
		}
		return StatusInternalServerError
	}
}

// I18nError 多语言错误结构体
type I18nError struct {
	Code     string
	Lang     string
	Message  string
	HttpCode int
}

// Error 实现 error 接口
func (e *I18nError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewI18nError 创建多语言错误
func NewI18nError(code, lang string, customMessage ...string) *I18nError {
	httpCode, _, message := NewI18nErrorResponse(code, lang, customMessage...)
	
	return &I18nError{
		Code:     code,
		Lang:     lang,
		Message:  message,
		HttpCode: httpCode,
	}
}

// 便利函数：常用错误码的快速创建方法

// NewAuthError 创建认证错误
func NewAuthError(lang string, code ...string) *I18nError {
	errorCode := "100001" // 默认认证过期
	if len(code) > 0 && code[0] != "" {
		errorCode = code[0]
	}
	return NewI18nError(errorCode, lang)
}

// NewValidationError 创建验证错误
func NewValidationError(lang string) *I18nError {
	return NewI18nError("900001", lang)
}

// NewNotFoundError 创建资源不存在错误
func NewNotFoundError(lang string) *I18nError {
	return NewI18nError("900002", lang)
}

// NewConflictError 创建资源冲突错误
func NewConflictError(lang string) *I18nError {
	return NewI18nError("900003", lang)
}

// NewInternalError 创建内部服务器错误
func NewInternalError(lang string) *I18nError {
	return NewI18nError("900004", lang)
}

// 批量错误处理函数

// ValidateI18nSupport 验证多语言支持状态
func ValidateI18nSupport() error {
	if globalManager == nil {
		return fmt.Errorf("i18n manager not initialized")
	}
	
	supportedLangs := globalManager.SupportedLanguages()
	if len(supportedLangs) == 0 {
		return fmt.Errorf("no languages loaded")
	}
	
	return nil
}

// GetSupportedLanguages 获取支持的语言列表
func GetSupportedLanguages() []string {
	if globalManager == nil {
		return []string{}
	}
	return globalManager.SupportedLanguages()
}

// IsLanguageSupported 检查语言是否支持
func IsLanguageSupported(lang string) bool {
	if globalManager == nil {
		return false
	}
	return globalManager.IsLanguageSupported(lang)
}