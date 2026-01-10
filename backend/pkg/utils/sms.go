package utils

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"license-manager/internal/config"
)

// 默认短信模板常量
const (
	DefaultTemplateRegister     = "SMS_330275014"
	DefaultTemplateResetPwd     = "SMS_330275015"
	DefaultTemplateLogin        = "SMS_330275016"
	DefaultTemplateCurrentPhone = "SMS_330275017" // 当前手机号验证码
	DefaultTemplateNewPhone     = "SMS_330275018" // 新手机号验证码
)

// 缓存键常量
const (
	SMSCodeKeyPrefix   = "sms:code"
	SMSLimitKeyPrefix  = "sms:limit"
	SMSCodeExpireTime  = 5 * time.Minute
	SMSLimitTimeWindow = time.Hour
	SMSLimitMaxCount   = 5
)

// SMSService 短信服务接口
type SMSService interface {
	SendVerificationCode(ctx context.Context, phone, phoneCountryCode, templateCode string) error
	VerifyCode(ctx context.Context, phone, phoneCountryCode, code string) (bool, error)
	IsRateLimited(ctx context.Context, phone, phoneCountryCode string) bool
	SendRegisterCode(ctx context.Context, phone, phoneCountryCode string) error
	SendResetPwdCode(ctx context.Context, phone, phoneCountryCode string) error
	SendLoginCode(ctx context.Context, phone, phoneCountryCode string) error
	SendCurrentPhoneCode(ctx context.Context, phone, phoneCountryCode string) error
	SendNewPhoneCode(ctx context.Context, phone, phoneCountryCode string) error
}

// smsService 短信服务实现
type smsService struct {
	config    *config.SMSConfig
	templates *config.SMSTemplates
	logger    interface{} // 支持不同类型的logger
	// TODO: 添加Redis客户端和阿里云SMS客户端
	// redis     *redis.Client
	// client    *dysmsapi20170525.Client
}

// NewSMSService 创建SMS服务实例
func NewSMSService(cfg *config.SMSConfig, logger interface{}) (SMSService, error) {
	if !cfg.Enabled {
		return &disabledSMSService{}, nil
	}

	// 验证配置有效性
	if cfg.AccessKeyID == "" || cfg.AccessKeySecret == "" {
		return nil, fmt.Errorf("SMS service configuration is invalid: AccessKeyID and AccessKeySecret are required")
	}

	if cfg.SignName == "" {
		return nil, fmt.Errorf("SMS service configuration is invalid: SignName is required")
	}

	return &smsService{
		config:    cfg,
		templates: &cfg.Templates,
		logger:    logger,
	}, nil
}

// disabledSMSService 禁用的SMS服务
type disabledSMSService struct{}

func (s *disabledSMSService) SendVerificationCode(ctx context.Context, phone, phoneCountryCode, templateCode string) error {
	return fmt.Errorf("SMS service is disabled")
}

func (s *disabledSMSService) VerifyCode(ctx context.Context, phone, phoneCountryCode, code string) (bool, error) {
	return false, fmt.Errorf("SMS service is disabled")
}

func (s *disabledSMSService) IsRateLimited(ctx context.Context, phone, phoneCountryCode string) bool {
	return true
}

func (s *disabledSMSService) SendRegisterCode(ctx context.Context, phone, phoneCountryCode string) error {
	return fmt.Errorf("SMS service is disabled")
}

func (s *disabledSMSService) SendResetPwdCode(ctx context.Context, phone, phoneCountryCode string) error {
	return fmt.Errorf("SMS service is disabled")
}

func (s *disabledSMSService) SendLoginCode(ctx context.Context, phone, phoneCountryCode string) error {
	return fmt.Errorf("SMS service is disabled")
}

func (s *disabledSMSService) SendCurrentPhoneCode(ctx context.Context, phone, phoneCountryCode string) error {
	return fmt.Errorf("SMS service is disabled")
}

func (s *disabledSMSService) SendNewPhoneCode(ctx context.Context, phone, phoneCountryCode string) error {
	return fmt.Errorf("SMS service is disabled")
}

// createSMSClient 创建阿里云短信客户端
// TODO: 实现阿里云SMS客户端创建，需要添加依赖：
// github.com/alibabacloud-go/dysmsapi-20170525/v5/client
// github.com/alibabacloud-go/darabonba-openapi/v2/client
// github.com/alibabacloud-go/tea-utils/v2/service
// github.com/alibabacloud-go/tea/tea
func createSMSClient(accessKeyID, accessKeySecret, regionID, endpoint string) (interface{}, error) {
	// 暂时返回nil，等待依赖添加
	return nil, fmt.Errorf("阿里云SMS客户端未实现，请添加相关依赖")
}

// SendVerificationCode 发送验证码
func (s *smsService) SendVerificationCode(ctx context.Context, phone, phoneCountryCode, templateCode string) error {
	// 检查频率限制
	if s.IsRateLimited(ctx, phone, phoneCountryCode) {
		return fmt.Errorf("request too frequent")
	}

	// 生成6位随机验证码
	code := generateVerificationCode()

	// 构造完整手机号
	fullPhone := phoneCountryCode + phone
	if phoneCountryCode == "" {
		fullPhone = "+86" + phone
	}

	// TODO: 发送短信 - 需要阿里云SMS SDK
	err := s.sendSMS(fullPhone, templateCode, code)
	if err != nil {
		fmt.Printf("Failed to send SMS to %s: %v", fullPhone, err)
		return fmt.Errorf("failed to send SMS: %w", err)
	}

	// TODO: 缓存验证码 - 需要Redis客户端
	// cacheKey := fmt.Sprintf("%s:%s:%s", SMSCodeKeyPrefix, phoneCountryCode, phone)
	// err = s.redis.Set(ctx, cacheKey, code, SMSCodeExpireTime).Err()
	// if err != nil {
	//     fmt.Printf("Failed to cache verification code: %v", err)
	// }

	// 记录发送日志
	fmt.Printf("SMS sent successfully to %s with template %s, code: %s", fullPhone, templateCode, code)

	return nil
}

// sendSMS 调用阿里云API发送短信
func (s *smsService) sendSMS(phone, templateCode, code string) error {
	// TODO: 实现阿里云SMS发送，需要添加相关依赖
	// 暂时返回错误，提示未实现
	return fmt.Errorf("阿里云SMS服务未实现，请添加相关依赖并配置有效的AccessKey")
}

// VerifyCode 验证验证码
func (s *smsService) VerifyCode(ctx context.Context, phone, phoneCountryCode, code string) (bool, error) {
	// TODO: 实现验证码验证，需要Redis客户端
	// 暂时返回验证失败，需要实现Redis缓存
	return false, fmt.Errorf("验证码验证服务未实现，需要配置Redis")
}

// IsRateLimited 检查是否被频率限制
func (s *smsService) IsRateLimited(ctx context.Context, phone, phoneCountryCode string) bool {
	// TODO: 实现频率限制检查，需要Redis客户端
	// 暂时允许所有请求通过，需要实现Redis缓存
	return false
}

// generateVerificationCode 生成6位随机验证码
func generateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())
	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	return code
}

// GetTemplateCode 获取模板代码：优先使用配置，否则使用默认值
func (s *smsService) GetTemplateCode(templateType string) string {
	switch templateType {
	case "register":
		if s.templates.Register != "" {
			return s.templates.Register
		}
		return DefaultTemplateRegister
	case "reset_pwd":
		if s.templates.ResetPwd != "" {
			return s.templates.ResetPwd
		}
		return DefaultTemplateResetPwd
	case "login":
		if s.templates.Login != "" {
			return s.templates.Login
		}
		return DefaultTemplateLogin
	case "current_phone":
		// 当前手机号模板暂不支持配置，使用默认值
		return DefaultTemplateCurrentPhone
	case "new_phone":
		// 新手机号模板暂不支持配置，使用默认值
		return DefaultTemplateNewPhone
	default:
		return ""
	}
}

// SendRegisterCode 发送注册验证码
func (s *smsService) SendRegisterCode(ctx context.Context, phone, phoneCountryCode string) error {
	templateCode := s.GetTemplateCode("register")
	return s.SendVerificationCode(ctx, phone, phoneCountryCode, templateCode)
}

// SendResetPwdCode 发送重置密码验证码
func (s *smsService) SendResetPwdCode(ctx context.Context, phone, phoneCountryCode string) error {
	templateCode := s.GetTemplateCode("reset_pwd")
	return s.SendVerificationCode(ctx, phone, phoneCountryCode, templateCode)
}

// SendLoginCode 发送登录验证码
func (s *smsService) SendLoginCode(ctx context.Context, phone, phoneCountryCode string) error {
	templateCode := s.GetTemplateCode("login")
	return s.SendVerificationCode(ctx, phone, phoneCountryCode, templateCode)
}

// SendCurrentPhoneCode 发送当前手机号验证码
func (s *smsService) SendCurrentPhoneCode(ctx context.Context, phone, phoneCountryCode string) error {
	templateCode := s.GetTemplateCode("current_phone")
	return s.SendVerificationCode(ctx, phone, phoneCountryCode, templateCode)
}

// SendNewPhoneCode 发送新手机号验证码
func (s *smsService) SendNewPhoneCode(ctx context.Context, phone, phoneCountryCode string) error {
	templateCode := s.GetTemplateCode("new_phone")
	return s.SendVerificationCode(ctx, phone, phoneCountryCode, templateCode)
}
