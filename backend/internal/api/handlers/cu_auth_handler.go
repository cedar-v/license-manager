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

type CuAuthHandler struct {
	cuUserService service.CuUserService
}

func NewCuAuthHandler(cuUserService service.CuUserService) *CuAuthHandler {
	return &CuAuthHandler{
		cuUserService: cuUserService,
	}
}

// CuUserRegister 客户用户注册
// @Summary 客户用户注册
// @Description 通过手机号注册客户用户账号
// @Tags 客户用户管理
// @Accept json
// @Produce json
// @Param user body models.CuUserRegisterRequest true "注册信息"
// @Success 200 {object} models.APIResponse{data=object{user=object{id=string,customer_id=string,phone=string,real_name=string,email=string,user_role=string,status=string},token=string}} "成功"
// @Failure 400 {object} models.ErrorResponse "请求参数无效"
// @Failure 409 {object} models.ErrorResponse "手机号已被注册"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/cu/register [post]
func (h *CuAuthHandler) CuUserRegister(c *gin.Context) {
	var req models.CuUserRegisterRequest
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

	// 验证短信验证码（暂时跳过，后面实现短信服务时补上）

	user, err := h.cuUserService.Register(c.Request.Context(), &req)
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

	// 生成JWT Token
	token, err := utils.GenerateCuToken(user.ID, user.CustomerID, user.UserRole, user.Phone)
	if err != nil {
		lang := middleware.GetLanguage(c)
		status, errCode, message := i18n.NewI18nErrorResponse("500001", lang)
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
		Data: gin.H{
			"user": gin.H{
				"id":          user.ID,
				"customer_id": user.CustomerID,
				"phone":       user.Phone,
				"real_name":   user.RealName,
				"email":       user.Email,
				"user_role":   user.UserRole,
				"status":      user.Status,
			},
			"token": token,
		},
	})
}

// CuUserLogin 客户用户登录
// @Summary 客户用户登录
// @Description 通过手机号和密码登录客户用户账号
// @Tags 客户用户管理
// @Accept json
// @Produce json
// @Param user body models.CuUserLoginRequest true "登录信息"
// @Success 200 {object} models.APIResponse{data=object{user=object{id=string,customer_id=string,phone=string,real_name=string,email=string,user_role=string,status=string},token=string}} "成功"
// @Failure 400 {object} models.ErrorResponse "请求参数无效"
// @Failure 401 {object} models.ErrorResponse "手机号或密码错误"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/cu/login [post]
func (h *CuAuthHandler) CuUserLogin(c *gin.Context) {
	var req models.CuUserLoginRequest
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

	// 获取客户端IP
	ip := c.ClientIP()

	user, token, err := h.cuUserService.Login(c.Request.Context(), &req, ip)
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
		Data: gin.H{
			"user": gin.H{
				"id":          user.ID,
				"customer_id": user.CustomerID,
				"phone":       user.Phone,
				"real_name":   user.RealName,
				"email":       user.Email,
				"user_role":   user.UserRole,
				"status":      user.Status,
			},
			"token": token,
		},
	})
}

// CuUserForgotPassword 忘记密码 - 发送验证码
// @Summary 忘记密码 - 发送验证码
// @Description 提交手机号，发送密码重置验证码到手机
// @Tags 客户用户管理
// @Accept json
// @Produce json
// @Param request body models.CuUserForgotPasswordRequest true "忘记密码请求"
// @Success 200 {object} models.APIResponse "验证码发送成功"
// @Failure 400 {object} models.ErrorResponse "请求参数无效"
// @Failure 404 {object} models.ErrorResponse "用户不存在"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/cu/forgot-password [post]
func (h *CuAuthHandler) CuUserForgotPassword(c *gin.Context) {
	var req models.CuUserForgotPasswordRequest
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

	// TODO: 实现短信验证码发送逻辑（在服务层处理）

	err := h.cuUserService.ForgotPassword(c.Request.Context(), &req)
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

// CuUserSendRegisterSms 注册发送验证码
// @Summary 注册发送验证码
// @Description 注册前发送短信验证码到手机
// @Tags 客户用户管理
// @Accept json
// @Produce json
// @Param request body models.CuUserSendRegisterSmsRequest true "注册发送验证码请求"
// @Success 200 {object} models.APIResponse "验证码发送成功"
// @Failure 400 {object} models.ErrorResponse "请求参数无效"
// @Failure 409 {object} models.ErrorResponse "手机号已被注册"
// @Failure 429 {object} models.ErrorResponse "请求过于频繁"
// @Failure 500 {object} models.ErrorResponse "短信发送失败"
// @Router /api/cu/send-register-sms [post]
func (h *CuAuthHandler) CuUserSendRegisterSms(c *gin.Context) {
	var req models.CuUserSendRegisterSmsRequest
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

	err := h.cuUserService.SendRegisterSms(c.Request.Context(), &req)
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

// CuUserResetPassword 重置密码
// @Summary 重置密码
// @Description 通过验证码重置密码
// @Tags 客户用户管理
// @Accept json
// @Produce json
// @Param request body models.CuUserResetPasswordRequest true "重置密码请求"
// @Success 200 {object} models.APIResponse "成功"
// @Failure 400 {object} models.ErrorResponse "请求参数无效"
// @Failure 409 {object} models.ErrorResponse "重置密码失败"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/cu/reset-password [post]
func (h *CuAuthHandler) CuUserResetPassword(c *gin.Context) {
	var req models.CuUserResetPasswordRequest
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

	// 验证短信验证码（暂时跳过）

	err := h.cuUserService.ResetPassword(c.Request.Context(), &req)
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
