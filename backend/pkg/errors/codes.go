package errors

// HTTP 状态码常量
const (
	StatusOK                  = 200
	StatusBadRequest          = 400
	StatusUnauthorized        = 401
	StatusForbidden           = 403
	StatusNotFound            = 404
	StatusConflict            = 409
	StatusInternalServerError = 500
)

// ErrorMessages 业务错误码到错误信息的映射
// 六位数字格式：前2位模块码 + 后4位错误号
var ErrorMessages = map[string]string{
	"000000": "成功",

	// 10xxxx - 认证模块
	"100001": "认证已过期",
	"100002": "认证无效", 
	"100003": "用户名或密码错误",
	"100004": "缺少认证信息",
	"100005": "权限不足",

	// 20xxxx - 客户模块
	"200001": "客户不存在",
	"200002": "客户已存在",
	"200003": "客户名称不能为空",
	"200004": "客户信息格式无效",
	"200005": "客户列表查询失败",

	// 90xxxx - 系统通用错误
	"900001": "请求参数无效",
	"900002": "资源不存在",
	"900003": "资源冲突", 
	"900004": "服务器内部错误",
}

// GetErrorMessage 获取错误码对应的错误信息
func GetErrorMessage(code string) string {
	if msg, ok := ErrorMessages[code]; ok {
		return msg
	}
	return "未知错误"
}

// NewErrorResponse 创建标准错误响应，返回 HTTP 状态码、错误码、错误消息
// 大部分业务错误返回 200，只有协议级错误才返回对应 HTTP 状态码
func NewErrorResponse(code string, customMessage ...string) (int, string, string) {
	httpStatus := getHTTPStatusByCode(code)
	message := GetErrorMessage(code)
	
	if len(customMessage) > 0 && customMessage[0] != "" {
		message = customMessage[0]
	}

	return httpStatus, code, message
}

// getHTTPStatusByCode 根据业务错误码获取对应的 HTTP 状态码
func getHTTPStatusByCode(code string) int {
	switch code {
	case "000000": // 成功
		return StatusOK
	case "900001": // 请求参数无效
		return StatusBadRequest
	case "100001", "100002", "100003", "100004": // 认证相关错误
		return StatusUnauthorized
	case "100005": // 权限不足
		return StatusForbidden
	case "200001", "900002": // 资源不存在
		return StatusNotFound
	case "200002", "900003": // 资源冲突
		return StatusConflict
	case "200003", "200004": // 客户信息无效
		return StatusBadRequest
	case "200005", "900004": // 服务器内部错误
		return StatusOK // 根据设计理念，大部分业务错误返回200，通过响应体区分
	default:
		return StatusInternalServerError
	}
}
