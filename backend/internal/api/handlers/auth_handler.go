package handlers

import (
	"net/http"
	"time"

	"license-manager/internal/api/middleware"
	"license-manager/internal/models"
	"license-manager/internal/service"
	"license-manager/pkg/i18n"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login 用户登录
// @Summary 用户登录
// @Description 管理员用户登录接口
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "登录请求参数"
// @Success 200 {object} models.LoginResponse "登录成功"
// @Failure 400 {object} models.ErrorResponse "请求参数无效"
// @Failure 401 {object} models.ErrorResponse "用户名或密码错误"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/v1/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		lang := middleware.GetLanguage(c)
		status, errCode, message := i18n.NewI18nErrorResponse("900001", lang)
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message,
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	// 获取客户端IP
	clientIP := c.ClientIP()

	data, err := h.authService.Login(c.Request.Context(), &req, clientIP)
	if err != nil {
		errorCode := "100003" // 登录失败
		if err.Error() == "配置未初始化" {
			errorCode = "900004" // 服务器内部错误
		}

		lang := middleware.GetLanguage(c)
		status, errCode, message := i18n.NewI18nErrorResponse(errorCode, lang, err.Error())
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message,
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	lang := middleware.GetLanguage(c)
	successMessage := i18n.GetErrorMessage("000000", lang)
	c.JSON(http.StatusOK, models.LoginResponse{
		Code:    "000000",
		Message: successMessage,
		Data:    data,
	})
}

// Logout 用户登出
// @Summary 用户登出
// @Description 用户登出接口
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse "登出成功"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Router /api/v1/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// 对于JWT无状态认证，客户端删除Token即可
	// 这里可以记录登出日志或将Token加入黑名单（可选实现）

	lang := middleware.GetLanguage(c)
	successMessage := i18n.GetErrorMessage("000000", lang)
	c.JSON(http.StatusOK, models.APIResponse{
		Code:    "000000",
		Message: successMessage,
	})
}

// RefreshToken 刷新Token
// @Summary 刷新Token
// @Description 刷新用户Token
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse "刷新成功"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// 从请求头获取Token
	token := extractTokenFromHeader(c)
	if token == "" {
		lang := middleware.GetLanguage(c)
		status, errCode, message := i18n.NewI18nErrorResponse("100004", lang)
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message,
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	newToken, err := h.authService.RefreshToken(token)
	if err != nil {
		// 根据错误类型选择合适的认证错误码
		errorCode := "100001" // token过期，可以刷新
		if err.Error() == "token无效" {
			errorCode = "100002" // token无效，需要重新登录
		}

		lang := middleware.GetLanguage(c)
		status, errCode, message := i18n.NewI18nErrorResponse(errorCode, lang, err.Error())
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message,
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	lang := middleware.GetLanguage(c)
	successMessage := i18n.GetErrorMessage("000000", lang)
	c.JSON(http.StatusOK, models.APIResponse{
		Code:    "000000",
		Message: successMessage,
		Data: map[string]interface{}{
			"token": newToken,
		},
	})
}

// getCurrentTimestamp 获取当前时间戳
func getCurrentTimestamp() string {
	return time.Now().Format(time.RFC3339)
}

// extractTokenFromHeader 从请求头提取Token
func extractTokenFromHeader(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return ""
}
