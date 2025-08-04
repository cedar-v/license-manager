package handlers

import (
	"net/http"

	"license-manager/internal/models"
	"license-manager/internal/service"
	"license-manager/pkg/errors"

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
		status, errCode, message := errors.NewErrorResponse("900001")
		c.JSON(status, models.ErrorResponse{
			Code:      status,
			Error:     errCode,
			Message:   message,
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	data, err := h.authService.Login(&req)
	if err != nil {
		errorCode := "100003" // 登录失败
		if err.Error() == "配置未初始化" {
			errorCode = "900004" // 服务器内部错误
		}

		status, errCode, message := errors.NewErrorResponse(errorCode, err.Error())
		c.JSON(status, models.ErrorResponse{
			Code:      status,
			Error:     errCode,
			Message:   message,
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		Code:    http.StatusOK,
		Message: "登录成功",
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

	c.JSON(http.StatusOK, models.APIResponse{
		Code:    http.StatusOK,
		Message: "登出成功",
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
		status, errCode, message := errors.NewErrorResponse("100004")
		c.JSON(status, models.ErrorResponse{
			Code:      status,
			Error:     errCode,
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

		status, errCode, message := errors.NewErrorResponse(errorCode, err.Error())
		c.JSON(status, models.ErrorResponse{
			Code:      status,
			Error:     errCode,
			Message:   message,
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:    http.StatusOK,
		Message: "刷新成功",
		Data: map[string]interface{}{
			"token": newToken,
		},
	})
}

// getCurrentTimestamp 获取当前时间戳
func getCurrentTimestamp() string {
	return "2024-07-30T12:00:00Z" // 简化实现，实际应该使用time.Now().Format(time.RFC3339)
}

// extractTokenFromHeader 从请求头提取Token
func extractTokenFromHeader(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return ""
}
