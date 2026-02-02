package handlers

import (
	"net/http"

	"license-manager/internal/api/middleware"
	"license-manager/internal/models"
	"license-manager/internal/service"
	"license-manager/pkg/i18n"

	"github.com/gin-gonic/gin"
)

type LeadHandler struct {
	leadService service.LeadService
}

func NewLeadHandler(leadService service.LeadService) *LeadHandler {
	return &LeadHandler{
		leadService: leadService,
	}
}

// CreateLead 创建线索
// @Summary 创建线索
// @Description 用户提交企业线索，无需鉴权
// @Tags 企业线索
// @Accept json
// @Produce json
// @Param request body models.LeadCreateRequest true "线索创建请求"
// @Success 200 {object} models.APIResponse{data=models.LeadResponse} "成功"
// @Failure 400 {object} models.ErrorResponse "请求参数无效"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/leads [post]
func (h *LeadHandler) CreateLead(c *gin.Context) {
	lang := middleware.GetLanguage(c)

	var req models.LeadCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		status, errCode, message := i18n.NewI18nErrorResponse("900001", lang)
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message + ": " + err.Error(),
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	lead, err := h.leadService.CreateLead(c.Request.Context(), &req)
	if err != nil {
		handleI18nError(c, err, lang)
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:      "000000",
		Message:   i18n.GetI18nErrorMessage("000000", lang),
		Data:      lead.ToResponse(),
		Timestamp: getCurrentTimestamp(),
	})
}

// GetLeads 获取线索列表
// @Summary 获取线索列表
// @Description 管理员获取线索列表，支持分页和筛选
// @Tags 企业线索
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码，默认1" minimum(1)
// @Param page_size query int false "每页条数，默认10，最大100" minimum(1) maximum(100)
// @Param search query string false "关键词检索"
// @Param status query string false "状态筛选"
// @Success 200 {object} models.APIResponse{data=models.LeadListResponse} "成功"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 403 {object} models.ErrorResponse "权限不足"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/leads [get]
func (h *LeadHandler) GetLeads(c *gin.Context) {
	lang := middleware.GetLanguage(c)

	var req models.LeadListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		status, errCode, message := i18n.NewI18nErrorResponse("900001", lang)
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message + ": " + err.Error(),
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	result, err := h.leadService.GetLeadList(c.Request.Context(), &req)
	if err != nil {
		handleI18nError(c, err, lang)
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:      "000000",
		Message:   i18n.GetI18nErrorMessage("000000", lang),
		Data:      result,
		Timestamp: getCurrentTimestamp(),
	})
}

// GetLeadSummary 获取线索汇总
// @Summary 获取线索汇总
// @Description 管理员获取企业线索汇总统计信息
// @Tags 企业线索
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=models.LeadSummaryResponse} "成功"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 403 {object} models.ErrorResponse "权限不足"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/leads/summary [get]
func (h *LeadHandler) GetLeadSummary(c *gin.Context) {
	lang := middleware.GetLanguage(c)

	summary, err := h.leadService.GetLeadSummary(c.Request.Context())
	if err != nil {
		handleI18nError(c, err, lang)
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:      "000000",
		Message:   i18n.GetI18nErrorMessage("000000", lang),
		Data:      summary,
		Timestamp: getCurrentTimestamp(),
	})
}

// GetLead 获取线索详情
// @Summary 获取线索详情
// @Description 根据ID获取线索详细信息
// @Tags 企业线索
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "线索ID"
// @Success 200 {object} models.APIResponse{data=models.LeadResponse} "成功"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 403 {object} models.ErrorResponse "权限不足"
// @Failure 404 {object} models.ErrorResponse "线索不存在"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/leads/{id} [get]
func (h *LeadHandler) GetLead(c *gin.Context) {
	lang := middleware.GetLanguage(c)
	id := c.Param("id")

	lead, err := h.leadService.GetLead(c.Request.Context(), id)
	if err != nil {
		handleI18nError(c, err, lang)
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:      "000000",
		Message:   i18n.GetI18nErrorMessage("000000", lang),
		Data:      lead.ToResponse(),
		Timestamp: getCurrentTimestamp(),
	})
}

// UpdateLead 更新线索
// @Summary 更新线索
// @Description 更新线索信息，用于线索跟进
// @Tags 企业线索
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "线索ID"
// @Param request body models.LeadUpdateRequest true "线索更新请求"
// @Success 200 {object} models.APIResponse{data=models.LeadResponse} "成功"
// @Failure 400 {object} models.ErrorResponse "请求参数无效"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 403 {object} models.ErrorResponse "权限不足"
// @Failure 404 {object} models.ErrorResponse "线索不存在"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/leads/{id} [put]
func (h *LeadHandler) UpdateLead(c *gin.Context) {
	lang := middleware.GetLanguage(c)
	id := c.Param("id")

	var req models.LeadUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		status, errCode, message := i18n.NewI18nErrorResponse("900001", lang)
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message + ": " + err.Error(),
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	lead, err := h.leadService.UpdateLead(c.Request.Context(), id, &req)
	if err != nil {
		handleI18nError(c, err, lang)
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:      "000000",
		Message:   i18n.GetI18nErrorMessage("000000", lang),
		Data:      lead.ToResponse(),
		Timestamp: getCurrentTimestamp(),
	})
}

// DeleteLead 删除线索
// @Summary 删除线索
// @Description 删除指定线索
// @Tags 企业线索
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "线索ID"
// @Success 200 {object} models.APIResponse{} "成功"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 403 {object} models.ErrorResponse "权限不足"
// @Failure 404 {object} models.ErrorResponse "线索不存在"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/leads/{id} [delete]
func (h *LeadHandler) DeleteLead(c *gin.Context) {
	lang := middleware.GetLanguage(c)
	id := c.Param("id")

	if err := h.leadService.DeleteLead(c.Request.Context(), id); err != nil {
		handleI18nError(c, err, lang)
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:      "000000",
		Message:   i18n.GetI18nErrorMessage("000000", lang),
		Data:      nil,
		Timestamp: getCurrentTimestamp(),
	})
}

// UpdateLeadStatus 更新线索状态
// @Summary 更新线索状态
// @Description 快速更新线索状态
// @Tags 企业线索
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "线索ID"
// @Param status query string true "状态: pending/contacted/converted/invalid"
// @Success 200 {object} models.APIResponse{data=models.LeadResponse} "成功"
// @Failure 400 {object} models.ErrorResponse "请求参数无效"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 403 {object} models.ErrorResponse "权限不足"
// @Failure 404 {object} models.ErrorResponse "线索不存在"
// @Router /api/leads/{id}/status [put]
func (h *LeadHandler) UpdateLeadStatus(c *gin.Context) {
	lang := middleware.GetLanguage(c)
	id := c.Param("id")

	status := c.Query("status")
	if status == "" {
		status, errCode, message := i18n.NewI18nErrorResponse("900001", lang)
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message,
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	// 验证状态值
	validStatuses := map[string]bool{
		string(models.LeadStatusPending):   true,
		string(models.LeadStatusContacted): true,
		string(models.LeadStatusConverted): true,
		string(models.LeadStatusInvalid):   true,
	}
	if !validStatuses[status] {
		statusCode, errCode, message := i18n.NewI18nErrorResponse("900001", lang)
		c.JSON(statusCode, models.ErrorResponse{
			Code:      errCode,
			Message:   message,
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	req := &models.LeadUpdateRequest{
		Status: status,
	}

	lead, err := h.leadService.UpdateLead(c.Request.Context(), id, req)
	if err != nil {
		handleI18nError(c, err, lang)
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:      "000000",
		Message:   i18n.GetI18nErrorMessage("000000", lang),
		Data:      lead.ToResponse(),
		Timestamp: getCurrentTimestamp(),
	})
}
