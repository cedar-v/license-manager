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

// TestSMSEnabledService 测试已删除 - 避免在CI/CD环境中调用真实的外部API
// 如需测试SMS功能，请在本地开发环境中进行手动测试
