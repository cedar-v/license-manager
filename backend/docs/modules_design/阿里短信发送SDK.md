# 阿里云短信发送SDK集成

## 概述

基于阿里云短信服务，实现验证码发送功能。已集成到项目中，支持注册、登录、密码重置等场景。

## 已实现功能 ✅

- ✅ AccessKey认证方式（推荐）
- ✅ 短信发送API调用
- ✅ 错误处理和响应解析
- ✅ 缓存集成
- ✅ 频率限制
- ✅ 多环境配置支持

## 技术栈

```go
package main

import (
  "context"
  "fmt"
  dysmsapi20170525  "github.com/alibabacloud-go/dysmsapi-20170525/v5/client"
  openapi  "github.com/alibabacloud-go/darabonba-openapi/v2/client"
  util  "github.com/alibabacloud-go/tea-utils/v2/service"
  "github.com/alibabacloud-go/tea/tea"
)


// createSMSClient 创建阿里云SMS客户端 (AccessKey方式)
func (s *smsService) createSMSClient() (*dysmsapi20170525.Client, error) {
  config := &openapi.Config{
    AccessKeyId:     tea.String(s.config.AccessKeyID),
    AccessKeySecret: tea.String(s.config.AccessKeySecret),
  }

  // 设置地域和端点
  if s.config.RegionID != "" {
    config.RegionId = tea.String(s.config.RegionID)
  } else {
    config.RegionId = tea.String("cn-hangzhou")
  }

  if s.config.Endpoint != "" {
    config.Endpoint = tea.String(s.config.Endpoint)
  } else {
    config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
  }

  client, err := dysmsapi20170525.NewClient(config)
  if err != nil {
    return nil, err
  }

  return client, nil
}

// CreateClient 使用凭据方式初始化账号Client (阿里云示例)
// 工程代码建议使用更安全的无AK方式，凭据配置方式请参见：https://help.aliyun.com/document_detail/378661.html。
func CreateClientWithCredentials() (_result *dysmsapi20170525.Client, _err error) {
  // 此方法仅供参考，项目中使用AccessKey方式
  return _result, fmt.Errorf("credentials方式未实现，请使用AccessKey方式")
}

// sendSMS 调用阿里云API发送短信 (项目中的实现)
func (s *smsService) sendSMS(phone, templateCode, code string) error {
  // 创建阿里云SMS客户端
  client, err := s.createSMSClient()
  if err != nil {
    return fmt.Errorf("failed to create SMS client: %w", err)
  }

  // 构建发送短信请求
  sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
    SignName:      tea.String(s.config.SignName),
    TemplateCode:  tea.String(templateCode),
    PhoneNumbers:  tea.String(phone),
    TemplateParam: tea.String(fmt.Sprintf(`{"code":"%s"}`, code)),
  }

  // 创建运行时选项
  runtime := &util.RuntimeOptions{}

  // 发送短信
  resp, err := client.SendSmsWithOptions(sendSmsRequest, runtime)
  if err != nil {
    return fmt.Errorf("failed to send SMS: %w", err)
  }

  // 检查响应结果
  if resp.Body == nil {
    return fmt.Errorf("SMS response body is nil")
  }

  // 检查返回码
  if tea.StringValue(resp.Body.Code) != "OK" {
    return fmt.Errorf("SMS send failed, code: %s, message: %s",
      tea.StringValue(resp.Body.Code),
      tea.StringValue(resp.Body.Message))
  }

  // 记录成功发送的日志
  fmt.Printf("SMS sent successfully, BizId: %s, Code: %s\n",
    tea.StringValue(resp.Body.BizId),
    tea.StringValue(resp.Body.Code))

  return nil
}

// _main 阿里云官方示例代码 (仅供参考)
func _main(args []*string) (_err error) {
  // 项目中不使用此方法，改用AccessKey方式
  return fmt.Errorf("此示例方法已废弃，请查看sendSMS方法")
}


## 配置和使用

### 1. 配置文件

```yaml
# configs/config.yaml
sms:
  enabled: true                           # 启用短信服务
  access_key_id: "your-access-key-id"     # 阿里云AccessKey ID
  access_key_secret: "your-access-key-secret" # 阿里云AccessKey Secret
  sign_name: "惠州顺视智能科技"              # 短信签名
  region_id: "cn-hangzhou"                # 地域ID
  endpoint: "dysmsapi.aliyuncs.com"       # 服务端点

  templates:
    register: "SMS_330275014"      # 注册验证码模板
    reset_pwd: "SMS_330275014"     # 重置密码验证码模板
    login: "SMS_330275014"         # 登录验证码模板
    current_phone: "SMS_330275014" # 当前手机号验证码模板
    new_phone: "SMS_330275014"     # 新手机号验证码模板
```

### 2. 使用示例

```go
// 发送注册验证码
err := smsService.SendVerificationCode(ctx, "+8613800012345", "+86", "register")
if err != nil {
    log.Printf("Failed to send SMS: %v", err)
}

// 验证验证码
valid, err := smsService.VerifyCode(ctx, "+8613800012345", "+86", "123456")
if err != nil {
    log.Printf("Failed to verify code: %v", err)
}
if !valid {
    log.Println("Invalid verification code")
}
```

### 3. 错误处理

项目中的SMS服务包含完善的错误处理：

- **认证错误**: AccessKey无效等
- **参数错误**: 手机号格式、模板不存在等
- **频率限制**: 自动频率控制和缓存
- **网络错误**: 重试和超时处理

### 4. 测试验证

```bash
# 运行SMS测试
go test ./pkg/utils/ -v -run TestSMS

# 查看测试输出，确认API调用正常
```

## 注意事项

1. **AccessKey安全**: 生产环境使用RAM用户AccessKey，不要使用主账号
2. **费用控制**: 阿里云SMS按条收费，注意测试用量
3. **频率限制**: 已实现客户端频率限制，避免被阿里云限制
4. **错误日志**: 重要错误已记录，便于排查问题

## 技术栈版本

- **阿里云SMS SDK**: `github.com/alibabacloud-go/dysmsapi-20170525/v5 v5.4.0`
- **OpenAPI SDK**: `github.com/alibabacloud-go/darabonba-openapi/v2 v2.1.13`
- **Tea Utils**: `github.com/alibabacloud-go/tea-utils/v2 v2.0.9`
- **Tea Core**: `github.com/alibabacloud-go/tea v1.4.0`

---

*此文档基于阿里云SMS服务官方SDK示例，项目中已完整实现并集成。*
