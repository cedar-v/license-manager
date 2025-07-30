package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"license-manager/internal/models"
	"license-manager/internal/service"
)

type SystemHandler struct {
	systemService service.SystemService
}

// NewSystemHandler 创建系统处理器
func NewSystemHandler(systemService service.SystemService) *SystemHandler {
	return &SystemHandler{
		systemService: systemService,
	}
}

// HealthCheck 健康检测接口
// @Summary 健康检测
// @Description 系统健康状态检查
// @Tags 系统
// @Accept json
// @Produce json
// @Success 200 {object} models.HealthResponse "健康状态"
// @Router /health [get]
func (h *SystemHandler) HealthCheck(c *gin.Context) {
	response := h.systemService.GetHealthStatus()
	c.JSON(http.StatusOK, response)
}

// GetSystemInfo 获取系统信息（需要认证）
// @Summary 获取系统信息
// @Description 获取详细的系统信息（需要管理员权限）
// @Tags 系统
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse "系统信息"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 403 {object} models.ErrorResponse "权限不足"
// @Router /api/v1/admin/system/info [get]
func (h *SystemHandler) GetSystemInfo(c *gin.Context) {
	systemInfo := h.systemService.GetSystemInfo()

	c.JSON(http.StatusOK, models.APIResponse{
		Code:    http.StatusOK,
		Message: "获取成功",
		Data:    systemInfo,
	})
}