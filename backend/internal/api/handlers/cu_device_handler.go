package handlers

import (
	"net/http"

	"license-manager/internal/api/middleware"
	"license-manager/internal/models"
	"license-manager/internal/service"
	"license-manager/pkg/i18n"
	"license-manager/pkg/utils"

	"github.com/gin-gonic/gin"
)

type CuDeviceHandler struct {
	cuDeviceService service.CuDeviceService
}

func NewCuDeviceHandler(cuDeviceService service.CuDeviceService) *CuDeviceHandler {
	return &CuDeviceHandler{
		cuDeviceService: cuDeviceService,
	}
}

// GetDevices 获取设备列表
// @Summary 获取设备列表
// @Description 获取当前登录用户的设备列表，支持分页和筛选
// @Tags 客户设备管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码，默认1"
// @Param page_size query int false "每页数量，默认20，最大100"
// @Param device_name query string false "设备名称模糊搜索"
// @Param authorization_code_id query string false "按授权码ID筛选设备"
// @Success 200 {object} models.APIResponse{data=models.DeviceListResponse} "获取成功"
// @Failure 400 {object} models.ErrorResponse "请求参数无效"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/cu/devices [get]
func (h *CuDeviceHandler) GetDevices(c *gin.Context) {
	// 从JWT Token中获取用户ID和客户ID
	claims := c.MustGet("cu_user").(*utils.CuClaims)

	// 解析查询参数
	var req models.DeviceListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		lang := middleware.GetLanguage(c)
		status, errCode, message := i18n.NewI18nErrorResponse("900001", lang)
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message + ": " + err.Error(),
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	// 调用服务层
	result, err := h.cuDeviceService.GetDeviceList(c.Request.Context(), claims.CustomerID, &req)
	if err != nil {
		// err已经是i18n.I18nError，直接使用
		i18nErr, ok := err.(*i18n.I18nError)
		if ok {
			c.JSON(i18nErr.HttpCode, models.ErrorResponse{
				Code:      i18nErr.Code,
				Message:   i18nErr.Message,
				Timestamp: getCurrentTimestamp(),
			})
			return
		}
		// 如果不是i18n错误，使用通用错误
		lang := middleware.GetLanguage(c)
		status, errCode, message := i18n.NewI18nErrorResponse("900004", lang, err.Error())
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message + ": " + err.Error(),
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	lang := middleware.GetLanguage(c)
	successMessage := i18n.GetI18nErrorMessage("000000", lang)
	c.JSON(http.StatusOK, models.APIResponse{
		Code:    "000000",
		Message: successMessage,
		Data:    result,
	})
}

// UnbindDevice 解绑设备
// @Summary 解绑设备
// @Description 解绑指定的设备，物理删除许可证记录
// @Tags 客户设备管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "许可证ID"
// @Success 200 {object} models.APIResponse "解绑成功"
// @Failure 400 {object} models.ErrorResponse "请求参数无效"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 403 {object} models.ErrorResponse "无权限操作此设备"
// @Failure 404 {object} models.ErrorResponse "设备不存在"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/cu/devices/{id} [delete]
func (h *CuDeviceHandler) UnbindDevice(c *gin.Context) {
	// 从JWT Token中获取用户ID和客户ID
	claims := c.MustGet("cu_user").(*utils.CuClaims)

	// 获取路径参数
	licenseID := c.Param("id")
	if licenseID == "" {
		lang := middleware.GetLanguage(c)
		status, errCode, message := i18n.NewI18nErrorResponse("900001", lang)
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message + ": license id is required",
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	// 调用服务层
	err := h.cuDeviceService.UnbindDevice(c.Request.Context(), claims.CustomerID, licenseID)
	if err != nil {
		// err已经是i18n.I18nError，直接使用
		i18nErr, ok := err.(*i18n.I18nError)
		if ok {
			c.JSON(i18nErr.HttpCode, models.ErrorResponse{
				Code:      i18nErr.Code,
				Message:   i18nErr.Message,
				Timestamp: getCurrentTimestamp(),
			})
			return
		}
		// 如果不是i18n错误，使用通用错误
		lang := middleware.GetLanguage(c)
		status, errCode, message := i18n.NewI18nErrorResponse("900004", lang, err.Error())
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message + ": " + err.Error(),
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	lang := middleware.GetLanguage(c)
	successMessage := i18n.GetI18nErrorMessage("000000", lang)
	c.JSON(http.StatusOK, models.APIResponse{
		Code:    "000000",
		Message: successMessage,
		Data:    nil,
	})
}

// GetDeviceSummary 获取设备汇总统计
// @Summary 获取设备汇总统计
// @Description 获取当前登录用户的设备汇总统计信息（总数、在线数、离线数）
// @Tags 客户设备管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=models.DeviceSummaryResponse} "获取成功"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/cu/devices/summary [get]
func (h *CuDeviceHandler) GetDeviceSummary(c *gin.Context) {
	// 从JWT Token中获取用户ID和客户ID
	claims := c.MustGet("cu_user").(*utils.CuClaims)

	// 调用服务层
	result, err := h.cuDeviceService.GetDeviceSummary(c.Request.Context(), claims.CustomerID)
	if err != nil {
		// err已经是i18n.I18nError，直接使用
		i18nErr, ok := err.(*i18n.I18nError)
		if ok {
			c.JSON(i18nErr.HttpCode, models.ErrorResponse{
				Code:      i18nErr.Code,
				Message:   i18nErr.Message,
				Timestamp: getCurrentTimestamp(),
			})
			return
		}
		// 如果不是i18n错误，使用通用错误
		lang := middleware.GetLanguage(c)
		status, errCode, message := i18n.NewI18nErrorResponse("900004", lang, err.Error())
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message + ": " + err.Error(),
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	lang := middleware.GetLanguage(c)
	successMessage := i18n.GetI18nErrorMessage("000000", lang)
	c.JSON(http.StatusOK, models.APIResponse{
		Code:    "000000",
		Message: successMessage,
		Data:    result,
	})
}
