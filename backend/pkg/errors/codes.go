package errors

// HTTP状态码常量定义
const (
	// 成功
	StatusOK = 200

	// 客户端错误
	StatusBadRequest   = 400
	StatusUnauthorized = 401
	StatusForbidden    = 403
	StatusNotFound     = 404
	StatusConflict     = 409

	// 服务器错误
	StatusInternalServerError = 500
)

// 业务错误码常量定义 - 极简实用设计
const (
	// 通用错误码 - 大部分场景直接用HTTP状态码即可
	CodeBadRequest  = "BAD_REQUEST"  // 400 - 请求参数错误
	CodeForbidden   = "FORBIDDEN"    // 403 - 权限不足
	CodeNotFound    = "NOT_FOUND"    // 404 - 资源不存在
	CodeConflict    = "CONFLICT"     // 409 - 资源冲突
	CodeServerError = "SERVER_ERROR" // 500 - 服务器错误

	// 认证相关错误码 - 前端需要区分处理
	CodeAuthExpired = "AUTH_EXPIRED" // 401 - token过期，需要刷新
	CodeAuthInvalid = "AUTH_INVALID" // 401 - token无效，需要重新登录
	CodeLoginFailed = "LOGIN_FAILED" // 401 - 登录失败
	CodeAuthMissing = "AUTH_MISSING" // 401 - 缺少token
)

// 错误码描述映射
var ErrorMessages = map[string]string{
	CodeBadRequest:  "请求参数无效",
	CodeForbidden:   "权限不足",
	CodeNotFound:    "资源不存在",
	CodeConflict:    "资源冲突",
	CodeServerError: "服务器内部错误",

	// 认证相关
	CodeAuthExpired: "认证已过期",
	CodeAuthInvalid: "认证无效",
	CodeLoginFailed: "用户名或密码错误",
	CodeAuthMissing: "缺少认证信息",
}

// GetErrorMessage 获取错误码对应的错误信息
func GetErrorMessage(code string) string {
	if msg, ok := ErrorMessages[code]; ok {
		return msg
	}
	return "未知错误"
}

// NewErrorResponse 创建标准错误响应 - 返回HTTP状态码、错误码、错误消息
func NewErrorResponse(code string, customMessage ...string) (int, string, string) {
	// 自动匹配HTTP状态码
	var httpStatus int
	switch code {
	case CodeBadRequest:
		httpStatus = StatusBadRequest
	case CodeAuthExpired, CodeAuthInvalid, CodeLoginFailed, CodeAuthMissing:
		httpStatus = StatusUnauthorized
	case CodeForbidden:
		httpStatus = StatusForbidden
	case CodeNotFound:
		httpStatus = StatusNotFound
	case CodeConflict:
		httpStatus = StatusConflict
	default:
		httpStatus = StatusInternalServerError
		code = CodeServerError
	}

	// 获取错误消息
	message := GetErrorMessage(code)
	if len(customMessage) > 0 && customMessage[0] != "" {
		message = customMessage[0]
	}

	return httpStatus, code, message
}
