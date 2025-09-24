package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"license-manager/internal/api/middleware"
	"license-manager/internal/models"
	"license-manager/internal/service"
	"license-manager/pkg/i18n"

	"github.com/gin-gonic/gin"
)

type LicenseHandler struct {
	licenseService service.LicenseService
}

// NewLicenseHandler 创建许可证处理器
func NewLicenseHandler(licenseService service.LicenseService) *LicenseHandler {
	return &LicenseHandler{
		licenseService: licenseService,
	}
}

// GetLicenseList 查询许可证列表
// @Summary 查询许可证列表
// @Description 分页查询许可证列表，支持筛选
// @Tags 许可证管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码，默认1" minimum(1)
// @Param page_size query int false "每页条数，默认20，最大100" minimum(1) maximum(100)
// @Param authorization_code_id query string false "授权码ID筛选"
// @Param customer_id query string false "客户ID筛选"
// @Param status query string false "状态筛选" Enums(active, inactive, revoked)
// @Param is_online query bool false "在线状态筛选"
// @Param sort query string false "排序字段，默认created_at" Enums(created_at, updated_at, activated_at, last_heartbeat)
// @Param order query string false "排序方向，默认desc" Enums(asc, desc)
// @Success 200 {object} models.APIResponse{data=models.LicenseListResponse} "查询成功"
// @Failure 400 {object} models.ErrorResponse "请求参数无效"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/v1/licenses [get]
func (h *LicenseHandler) GetLicenseList(c *gin.Context) {
	var req models.LicenseListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		lang := middleware.GetLanguage(c)
		status, errCode, message := i18n.NewI18nErrorResponse("900001", lang)
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message + ": " + err.Error(),
			Timestamp: time.Now().Format(time.RFC3339),
		})
		return
	}

	// 设置语言到Context中
	ctx := middleware.WithLanguage(c.Request.Context(), c)

	data, err := h.licenseService.GetLicenseList(ctx, &req)
	if err != nil {
		// 错误已经在Service层完全包装好了，直接使用
		var i18nErr *i18n.I18nError
		if errors.As(err, &i18nErr) {
			c.JSON(i18nErr.HttpCode, models.ErrorResponse{
				Code:      i18nErr.Code,
				Message:   i18nErr.Message,
				Timestamp: time.Now().Format(time.RFC3339),
			})
		} else {
			// 兜底：理论上不应该到这里，因为Service层应该返回I18nError
			lang := middleware.GetLanguage(c)
			status, errCode, message := i18n.NewI18nErrorResponse("900004", lang)
			c.JSON(status, models.ErrorResponse{
				Code:      errCode,
				Message:   message,
				Timestamp: time.Now().Format(time.RFC3339),
			})
		}
		return
	}

	lang := middleware.GetLanguage(c)
	successMessage := i18n.GetErrorMessage("000000", lang)
	c.JSON(http.StatusOK, models.APIResponse{
		Code:    "000000",
		Message: successMessage,
		Data:    data,
	})
}

// GetLicense 获取许可证详情
// @Summary 获取许可证详情
// @Description 根据许可证ID获取许可证详细信息
// @Tags 许可证管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "许可证ID"
// @Success 200 {object} models.APIResponse{data=models.LicenseDetailResponse} "查询成功"
// @Failure 400 {object} models.ErrorResponse "请求参数无效"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 404 {object} models.ErrorResponse "许可证不存在"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/v1/licenses/{id} [get]
func (h *LicenseHandler) GetLicense(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		lang := middleware.GetLanguage(c)
		status, errCode, message := i18n.NewI18nErrorResponse("900001", lang)
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message,
			Timestamp: time.Now().Format(time.RFC3339),
		})
		return
	}

	// 设置语言到Context中
	ctx := middleware.WithLanguage(c.Request.Context(), c)

	data, err := h.licenseService.GetLicense(ctx, id)
	if err != nil {
		// 错误已经在Service层完全包装好了，直接使用
		var i18nErr *i18n.I18nError
		if errors.As(err, &i18nErr) {
			c.JSON(i18nErr.HttpCode, models.ErrorResponse{
				Code:      i18nErr.Code,
				Message:   i18nErr.Message,
				Timestamp: time.Now().Format(time.RFC3339),
			})
		} else {
			// 兜底：理论上不应该到这里，因为Service层应该返回I18nError
			lang := middleware.GetLanguage(c)
			status, errCode, message := i18n.NewI18nErrorResponse("900004", lang)
			c.JSON(status, models.ErrorResponse{
				Code:      errCode,
				Message:   message,
				Timestamp: time.Now().Format(time.RFC3339),
			})
		}
		return
	}

	lang := middleware.GetLanguage(c)
	successMessage := i18n.GetErrorMessage("000000", lang)
	c.JSON(http.StatusOK, models.APIResponse{
		Code:    "000000",
		Message: successMessage,
		Data:    data,
	})
}

// CreateLicense 手动添加许可证
// @Summary 手动添加许可证
// @Description 为指定授权码手动创建许可证
// @Tags 许可证管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.LicenseCreateRequest true "许可证创建信息"
// @Success 201 {object} models.APIResponse{data=models.License} "创建成功"
// @Failure 400 {object} models.ErrorResponse "请求参数无效"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 404 {object} models.ErrorResponse "授权码不存在"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/v1/licenses [post]
func (h *LicenseHandler) CreateLicense(c *gin.Context) {
	var req models.LicenseCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		lang := middleware.GetLanguage(c)
		status, errCode, message := i18n.NewI18nErrorResponse("900001", lang)
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message + ": " + err.Error(),
			Timestamp: time.Now().Format(time.RFC3339),
		})
		return
	}

	// 设置语言到Context中
	ctx := middleware.WithLanguage(c.Request.Context(), c)

	data, err := h.licenseService.CreateLicense(ctx, &req)
	if err != nil {
		// 错误已经在Service层完全包装好了，直接使用
		var i18nErr *i18n.I18nError
		if errors.As(err, &i18nErr) {
			c.JSON(i18nErr.HttpCode, models.ErrorResponse{
				Code:      i18nErr.Code,
				Message:   i18nErr.Message,
				Timestamp: time.Now().Format(time.RFC3339),
			})
		} else {
			// 兜底：理论上不应该到这里，因为Service层应该返回I18nError
			lang := middleware.GetLanguage(c)
			status, errCode, message := i18n.NewI18nErrorResponse("900004", lang)
			c.JSON(status, models.ErrorResponse{
				Code:      errCode,
				Message:   message,
				Timestamp: time.Now().Format(time.RFC3339),
			})
		}
		return
	}

	lang := middleware.GetLanguage(c)
	successMessage := i18n.GetErrorMessage("000000", lang)
	c.JSON(http.StatusCreated, models.APIResponse{
		Code:    "000000",
		Message: successMessage,
		Data:    data,
	})
}

// RevokeLicense 撤销许可证
// @Summary 撤销许可证
// @Description 撤销指定的许可证，撤销后无法再使用
// @Tags 许可证管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "许可证ID"
// @Param request body models.LicenseRevokeRequest true "撤销原因"
// @Success 200 {object} models.APIResponse{data=models.License} "撤销成功"
// @Failure 400 {object} models.ErrorResponse "请求参数无效"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 404 {object} models.ErrorResponse "许可证不存在"
// @Failure 409 {object} models.ErrorResponse "许可证已被撤销"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/v1/licenses/{id}/revoke [put]
func (h *LicenseHandler) RevokeLicense(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		lang := middleware.GetLanguage(c)
		status, errCode, message := i18n.NewI18nErrorResponse("900001", lang)
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message,
			Timestamp: time.Now().Format(time.RFC3339),
		})
		return
	}

	var req models.LicenseRevokeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		lang := middleware.GetLanguage(c)
		status, errCode, message := i18n.NewI18nErrorResponse("900001", lang)
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message + ": " + err.Error(),
			Timestamp: time.Now().Format(time.RFC3339),
		})
		return
	}

	// 设置语言到Context中
	ctx := middleware.WithLanguage(c.Request.Context(), c)

	data, err := h.licenseService.RevokeLicense(ctx, id, &req)
	if err != nil {
		// 错误已经在Service层完全包装好了，直接使用
		var i18nErr *i18n.I18nError
		if errors.As(err, &i18nErr) {
			c.JSON(i18nErr.HttpCode, models.ErrorResponse{
				Code:      i18nErr.Code,
				Message:   i18nErr.Message,
				Timestamp: time.Now().Format(time.RFC3339),
			})
		} else {
			// 兜底：理论上不应该到这里，因为Service层应该返回I18nError
			lang := middleware.GetLanguage(c)
			status, errCode, message := i18n.NewI18nErrorResponse("900004", lang)
			c.JSON(status, models.ErrorResponse{
				Code:      errCode,
				Message:   message,
				Timestamp: time.Now().Format(time.RFC3339),
			})
		}
		return
	}

	lang := middleware.GetLanguage(c)
	successMessage := i18n.GetErrorMessage("000000", lang)
	c.JSON(http.StatusOK, models.APIResponse{
		Code:    "000000",
		Message: successMessage,
		Data:    data,
	})
}

// DownloadLicenseFile 下载许可证文件
// @Summary 下载许可证文件
// @Description 下载加密的许可证文件，用于客户端软件激活
// @Tags 许可证管理
// @Accept json
// @Produce application/octet-stream
// @Security BearerAuth
// @Param id path string true "许可证ID"
// @Success 200 {file} binary "许可证文件"
// @Failure 400 {object} models.ErrorResponse "请求参数无效"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 404 {object} models.ErrorResponse "许可证不存在"
// @Failure 409 {object} models.ErrorResponse "许可证已被撤销"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/v1/licenses/{id}/download [get]
func (h *LicenseHandler) DownloadLicenseFile(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		lang := middleware.GetLanguage(c)
		status, errCode, message := i18n.NewI18nErrorResponse("900001", lang)
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message,
			Timestamp: time.Now().Format(time.RFC3339),
		})
		return
	}

	// 设置语言到Context中
	ctx := middleware.WithLanguage(c.Request.Context(), c)

	fileData, fileName, err := h.licenseService.GenerateLicenseFile(ctx, id)
	if err != nil {
		// 错误已经在Service层完全包装好了，直接使用
		var i18nErr *i18n.I18nError
		if errors.As(err, &i18nErr) {
			c.JSON(i18nErr.HttpCode, models.ErrorResponse{
				Code:      i18nErr.Code,
				Message:   i18nErr.Message,
				Timestamp: time.Now().Format(time.RFC3339),
			})
		} else {
			// 兜底：理论上不应该到这里，因为Service层应该返回I18nError
			lang := middleware.GetLanguage(c)
			status, errCode, message := i18n.NewI18nErrorResponse("900004", lang)
			c.JSON(status, models.ErrorResponse{
				Code:      errCode,
				Message:   message,
				Timestamp: time.Now().Format(time.RFC3339),
			})
		}
		return
	}

	// 设置响应头，指示这是一个文件下载
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	// 返回文件内容
	c.Data(http.StatusOK, "application/octet-stream", fileData)
}

// ActivateLicense 激活许可证
// @Summary 激活许可证
// @Description 客户端使用授权码激活软件，获取许可证文件
// @Tags 许可证激活
// @Accept json
// @Produce json
// @Param request body models.ActivateRequest true "激活请求"
// @Success 200 {object} models.APIResponse{data=models.ActivateResponse} "激活成功"
// @Failure 400 {object} models.ErrorResponse "请求参数无效"
// @Failure 404 {object} models.ErrorResponse "授权码不存在"
// @Failure 409 {object} models.ErrorResponse "授权码已锁定或已过期"
// @Failure 429 {object} models.ErrorResponse "激活数量已达上限"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/v1/activate [post]
func (h *LicenseHandler) ActivateLicense(c *gin.Context) {
	var req models.ActivateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		lang := middleware.GetLanguage(c)
		status, errCode, message := i18n.NewI18nErrorResponse("900001", lang)
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message + ": " + err.Error(),
			Timestamp: time.Now().Format(time.RFC3339),
		})
		return
	}

	// 设置语言到Context中
	ctx := middleware.WithLanguage(c.Request.Context(), c)

	// 获取客户端IP
	clientIP := c.ClientIP()

	data, err := h.licenseService.ActivateLicense(ctx, &req, clientIP)
	if err != nil {
		// 错误已经在Service层完全包装好了，直接使用
		var i18nErr *i18n.I18nError
		if errors.As(err, &i18nErr) {
			c.JSON(i18nErr.HttpCode, models.ErrorResponse{
				Code:      i18nErr.Code,
				Message:   i18nErr.Message,
				Timestamp: time.Now().Format(time.RFC3339),
			})
		} else {
			// 兜底：理论上不应该到这里，因为Service层应该返回I18nError
			lang := middleware.GetLanguage(c)
			status, errCode, message := i18n.NewI18nErrorResponse("900004", lang)
			c.JSON(status, models.ErrorResponse{
				Code:      errCode,
				Message:   message,
				Timestamp: time.Now().Format(time.RFC3339),
			})
		}
		return
	}

	lang := middleware.GetLanguage(c)
	successMessage := i18n.GetErrorMessage("000000", lang)
	c.JSON(http.StatusOK, models.APIResponse{
		Code:    "000000",
		Message: successMessage,
		Data:    data,
	})
}

// Heartbeat 心跳检测
// @Summary 心跳检测
// @Description 客户端定期发送心跳，更新在线状态和使用数据
// @Tags 许可证激活
// @Accept json
// @Produce json
// @Param request body models.HeartbeatRequest true "心跳请求"
// @Success 200 {object} models.APIResponse{data=models.HeartbeatResponse} "心跳成功"
// @Failure 400 {object} models.ErrorResponse "请求参数无效"
// @Failure 404 {object} models.ErrorResponse "许可证不存在"
// @Failure 409 {object} models.ErrorResponse "许可证已被撤销"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/v1/heartbeat [post]
func (h *LicenseHandler) Heartbeat(c *gin.Context) {
	var req models.HeartbeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		lang := middleware.GetLanguage(c)
		status, errCode, message := i18n.NewI18nErrorResponse("900001", lang)
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message + ": " + err.Error(),
			Timestamp: time.Now().Format(time.RFC3339),
		})
		return
	}

	// 设置语言到Context中
	ctx := middleware.WithLanguage(c.Request.Context(), c)

	// 获取客户端IP
	clientIP := c.ClientIP()

	data, err := h.licenseService.Heartbeat(ctx, &req, clientIP)
	if err != nil {
		// 错误已经在Service层完全包装好了，直接使用
		var i18nErr *i18n.I18nError
		if errors.As(err, &i18nErr) {
			c.JSON(i18nErr.HttpCode, models.ErrorResponse{
				Code:      i18nErr.Code,
				Message:   i18nErr.Message,
				Timestamp: time.Now().Format(time.RFC3339),
			})
		} else {
			// 兜底：理论上不应该到这里，因为Service层应该返回I18nError
			lang := middleware.GetLanguage(c)
			status, errCode, message := i18n.NewI18nErrorResponse("900004", lang)
			c.JSON(status, models.ErrorResponse{
				Code:      errCode,
				Message:   message,
				Timestamp: time.Now().Format(time.RFC3339),
			})
		}
		return
	}

	lang := middleware.GetLanguage(c)
	successMessage := i18n.GetErrorMessage("000000", lang)
	c.JSON(http.StatusOK, models.APIResponse{
		Code:    "000000",
		Message: successMessage,
		Data:    data,
	})
}

// GetStatsOverview 获取授权概览统计
// @Summary 获取授权概览统计
// @Description 获取授权码、许可证的总体统计信息
// @Tags 统计分析
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=models.StatsOverviewResponse} "查询成功"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/v1/stats/overview [get]
func (h *LicenseHandler) GetStatsOverview(c *gin.Context) {
	// 设置语言到Context中
	ctx := middleware.WithLanguage(c.Request.Context(), c)

	data, err := h.licenseService.GetStatsOverview(ctx)
	if err != nil {
		// 错误已经在Service层完全包装好了，直接使用
		var i18nErr *i18n.I18nError
		if errors.As(err, &i18nErr) {
			c.JSON(i18nErr.HttpCode, models.ErrorResponse{
				Code:      i18nErr.Code,
				Message:   i18nErr.Message,
				Timestamp: time.Now().Format(time.RFC3339),
			})
		} else {
			// 兜底：理论上不应该到这里，因为Service层应该返回I18nError
			lang := middleware.GetLanguage(c)
			status, errCode, message := i18n.NewI18nErrorResponse("900004", lang)
			c.JSON(status, models.ErrorResponse{
				Code:      errCode,
				Message:   message,
				Timestamp: time.Now().Format(time.RFC3339),
			})
		}
		return
	}

	lang := middleware.GetLanguage(c)
	successMessage := i18n.GetErrorMessage("000000", lang)
	c.JSON(http.StatusOK, models.APIResponse{
		Code:    "000000",
		Message: successMessage,
		Data:    data,
	})
}