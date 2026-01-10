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

type CuProfileHandler struct {
	cuUserService service.CuUserService
}

func NewCuProfileHandler(cuUserService service.CuUserService) *CuProfileHandler {
	return &CuProfileHandler{
		cuUserService: cuUserService,
	}
}

// GetCuUserProfile 获取用户个人资料
// @Summary 获取用户个人资料
// @Description 获取当前登录用户的个人资料信息
// @Tags 客户用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=models.CuUserResponse} "获取成功"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 404 {object} models.ErrorResponse "用户不存在"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/cu/profile [get]
func (h *CuProfileHandler) GetCuUserProfile(c *gin.Context) {
	// 从JWT Token中获取用户ID
	claims := c.MustGet("cu_user").(*utils.CuClaims)

	user, err := h.cuUserService.GetProfile(c.Request.Context(), claims.UserID)
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
	successMessage := i18n.GetErrorMessage("000000", lang)
	c.JSON(http.StatusOK, models.APIResponse{
		Code:    "000000",
		Message: successMessage,
		Data:    user.ToResponse(),
	})
}

// UpdateCuUserProfile 更新用户个人资料
// @Summary 更新用户个人资料
// @Description 更新当前登录用户的个人资料信息
// @Tags 客户用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param profile body models.CuUserProfileUpdateRequest true "个人资料信息"
// @Success 200 {object} models.APIResponse "成功"
// @Failure 400 {object} models.ErrorResponse "请求参数无效"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 409 {object} models.ErrorResponse "个人资料更新失败"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/cu/profile [put]
func (h *CuProfileHandler) UpdateCuUserProfile(c *gin.Context) {
	// 从JWT Token中获取用户ID
	claims := c.MustGet("cu_user").(*utils.CuClaims)

	var req models.CuUserProfileUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		lang := middleware.GetLanguage(c)
		status, errCode, message := i18n.NewI18nErrorResponse("900001", lang)
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message + ": " + err.Error(),
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	err := h.cuUserService.UpdateProfile(c.Request.Context(), claims.UserID, &req)
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
	successMessage := i18n.GetErrorMessage("000000", lang)
	c.JSON(http.StatusOK, models.APIResponse{
		Code:    "000000",
		Message: successMessage,
		Data:    nil,
	})
}

// CuUserSendCurrentPhoneSms 发送当前手机号验证码
// @Summary 发送当前手机号验证码
// @Description 发送当前手机号验证码，用于手机号更新时的身份验证
// @Tags 客户用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse "验证码发送成功"
// @Failure 400 {object} models.ErrorResponse "请求参数无效"
// @Failure 401 {object} models.ErrorResponse "未授权"
// @Failure 404 {object} models.ErrorResponse "用户不存在"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/cu/profile/send-current-phone-sms [post]
func (h *CuProfileHandler) CuUserSendCurrentPhoneSms(c *gin.Context) {
	claims := c.MustGet("cu_user").(*utils.CuClaims)

	err := h.cuUserService.SendCurrentPhoneSms(c.Request.Context(), claims.UserID)
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
		Data:    nil,
	})
}

// CuUserSendNewPhoneSms 发送新手机号验证码
// @Summary 发送新手机号验证码
// @Description 发送新手机号验证码，用于手机号更新时的可用性验证
// @Tags 客户用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CuUserSendNewPhoneSmsRequest true "新手机号信息"
// @Success 200 {object} models.APIResponse "验证码发送成功"
// @Failure 400 {object} models.ErrorResponse "请求参数无效"
// @Failure 401 {object} models.ErrorResponse "未授权"
// @Failure 409 {object} models.ErrorResponse "手机号已被注册"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/cu/profile/send-new-phone-sms [post]
func (h *CuProfileHandler) CuUserSendNewPhoneSms(c *gin.Context) {

	var req models.CuUserSendNewPhoneSmsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		lang := middleware.GetLanguage(c)
		status, errCode, message := i18n.NewI18nErrorResponse("900001", lang)
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message + ": " + err.Error(),
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	err := h.cuUserService.SendNewPhoneSms(c.Request.Context(), &req)
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
		Data:    nil,
	})
}

// UpdateCuUserPhone 更新用户手机号
// @Summary 更新用户手机号
// @Description 更新当前登录用户的手机号，需要验证新旧手机号
// @Tags 客户用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param phone body models.CuUserPhoneUpdateRequest true "手机号更新信息"
// @Success 200 {object} models.APIResponse "成功"
// @Failure 400 {object} models.ErrorResponse "请求参数无效"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 409 {object} models.ErrorResponse "手机号更新失败"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/cu/profile/phone [put]
func (h *CuProfileHandler) UpdateCuUserPhone(c *gin.Context) {
	// 从JWT Token中获取用户ID
	claims := c.MustGet("cu_user").(*utils.CuClaims)

	var req models.CuUserPhoneUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		lang := middleware.GetLanguage(c)
		status, errCode, message := i18n.NewI18nErrorResponse("900001", lang)
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message + ": " + err.Error(),
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	// 获取用户的当前手机号信息（从数据库）
	user, err := h.cuUserService.GetProfile(c.Request.Context(), claims.UserID)
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
			Message:   message,
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	// 验证当前手机号验证码
	valid, err := h.cuUserService.VerifyPhoneCode(c.Request.Context(), user.Phone, user.PhoneCountryCode, req.CurrentSmsCode, "current_phone")
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
			Message:   message,
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	if !valid {
		lang := middleware.GetLanguage(c)
		status, errCode, message := i18n.NewI18nErrorResponse("500012", lang) // 验证码错误或已过期
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message,
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	// 验证新手机号验证码
	newPhoneCountryCode := req.NewPhoneCountryCode
	if newPhoneCountryCode == "" {
		newPhoneCountryCode = "+86"
	}

	valid, err = h.cuUserService.VerifyPhoneCode(c.Request.Context(), req.NewPhone, newPhoneCountryCode, req.NewSmsCode, "new_phone")
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
			Message:   message,
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	if !valid {
		lang := middleware.GetLanguage(c)
		status, errCode, message := i18n.NewI18nErrorResponse("500012", lang) // 验证码错误或已过期
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message,
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	err = h.cuUserService.UpdatePhone(c.Request.Context(), claims.UserID, &req)
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
	successMessage := i18n.GetErrorMessage("000000", lang)
	c.JSON(http.StatusOK, models.APIResponse{
		Code:    "000000",
		Message: successMessage,
		Data:    nil,
	})
}

// ChangeCuUserPassword 修改用户密码
// @Summary 修改用户密码
// @Description 修改当前登录用户的密码，需要提供旧密码
// @Tags 客户用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param password body models.CuUserChangePasswordRequest true "密码修改信息"
// @Success 200 {object} models.APIResponse "成功"
// @Failure 400 {object} models.ErrorResponse "请求参数无效"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 409 {object} models.ErrorResponse "密码修改失败"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/cu/profile/password [put]
func (h *CuProfileHandler) ChangeCuUserPassword(c *gin.Context) {
	// 从JWT Token中获取用户ID
	claims := c.MustGet("cu_user").(*utils.CuClaims)

	var req models.CuUserChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		lang := middleware.GetLanguage(c)
		status, errCode, message := i18n.NewI18nErrorResponse("900001", lang)
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message + ": " + err.Error(),
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	err := h.cuUserService.ChangePassword(c.Request.Context(), claims.UserID, req.OldPassword, req.NewPassword)
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
		Data:    nil,
	})
}
