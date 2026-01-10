package utils

import (
	"context"
	"testing"

	"license-manager/internal/config"
	"license-manager/pkg/cache"
)

func TestSMSServiceInitialization(t *testing.T) {
	// 创建内存缓存用于测试
	memCache := cache.NewMemoryCache(100)

	// 创建SMS配置
	smsConfig := &config.SMSConfig{
		Enabled:         false, // 测试禁用状态
		AccessKeyID:     "test-key-id",
		AccessKeySecret: "test-key-secret",
		SignName:        "Test Sign",
		RegionID:        "cn-hangzhou",
		Endpoint:        "dysmsapi.aliyuncs.com",
	}

	// 测试创建SMS服务
	smsService, err := NewSMSService(smsConfig, memCache, nil)
	if err != nil {
		t.Fatalf("Failed to create SMS service: %v", err)
	}

	if smsService == nil {
		t.Fatal("SMS service should not be nil")
	}

	// 测试禁用状态下的方法调用
	ctx := context.Background()

	// 这些方法在禁用状态下应该返回错误
	err = smsService.SendVerificationCode(ctx, "+8613800012345", "+86", "test-template")
	if err == nil {
		t.Error("SendVerificationCode should return error when SMS is disabled")
	}

	_, err = smsService.VerifyCode(ctx, "+8613800012345", "+86", "123456")
	if err == nil {
		t.Error("VerifyCode should return error when SMS is disabled")
	}

	exists := smsService.IsRateLimited(ctx, "+8613800012345", "+86")
	if !exists { // 禁用状态下应该总是限流
		t.Error("IsRateLimited should return true when SMS is disabled")
	}
}

func TestSMSEnabledService(t *testing.T) {
	// 创建内存缓存用于测试
	memCache := cache.NewMemoryCache(100)

	// 创建SMS配置（启用状态）
	smsConfig := &config.SMSConfig{
		Enabled:         true,
		AccessKeyID:     "test-key-id",
		AccessKeySecret: "test-key-secret",
		SignName:        "Test Sign",
		RegionID:        "cn-hangzhou",
		Endpoint:        "dysmsapi.aliyuncs.com",
		Templates: config.SMSTemplates{
			Register:     "SMS_111111111",
			ResetPwd:     "SMS_222222222",
			Login:        "SMS_333333333",
			CurrentPhone: "SMS_444444444",
			NewPhone:     "SMS_555555555",
		},
	}

	// 测试创建SMS服务
	smsService, err := NewSMSService(smsConfig, memCache, nil)
	if err != nil {
		t.Fatalf("Failed to create SMS service: %v", err)
	}

	if smsService == nil {
		t.Fatal("SMS service should not be nil")
	}

	// 测试模板配置是否正确应用
	// 注意：由于smsService是内部类型，我们通过接口方法间接验证
	// 在实际应用中，SendCurrentPhoneCode会使用配置中的current_phone模板

	ctx := context.Background()

	// 测试频率限制（应该允许第一次发送）
	isLimited := smsService.IsRateLimited(ctx, "+8613800012345", "+86")
	if isLimited {
		t.Error("First SMS should not be rate limited")
	}

	// 测试发送登录验证码（会因为无效的AccessKey而失败，但不应该panic）
	err = smsService.SendLoginCode(ctx, "+8613800012345", "+86")
	if err == nil {
		t.Log("SMS send succeeded (unexpected with test credentials)")
	} else {
		t.Logf("SMS send failed as expected: %v", err)
	}
}
