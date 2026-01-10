package cache

import (
	"strings"
)

// KeyBuilder 缓存键构建器
type KeyBuilder struct {
	prefix    string
	separator string
}

// NewKeyBuilder 创建键构建器
func NewKeyBuilder(prefix string) *KeyBuilder {
	return &KeyBuilder{
		prefix:    prefix,
		separator: ":",
	}
}

// Build 构建缓存键
func (k *KeyBuilder) Build(parts ...string) string {
	if k.prefix != "" {
		parts = append([]string{k.prefix}, parts...)
	}
	return strings.Join(parts, k.separator)
}

// User 构建用户相关键
func (k *KeyBuilder) User(userID string) string {
	return k.Build("user", userID)
}

// Session 构建会话相关键
func (k *KeyBuilder) Session(sessionID string) string {
	return k.Build("session", sessionID)
}

// SMSCode 构建短信验证码键
func (k *KeyBuilder) SMSCode(phone string) string {
	return k.Build("sms", "code", phone)
}

// SMSLimit 构建短信频率限制键
func (k *KeyBuilder) SMSLimit(phone string) string {
	return k.Build("sms", "limit", phone)
}

// Config 构建配置相关键
func (k *KeyBuilder) Config(key string) string {
	return k.Build("config", key)
}

// Lock 构建分布式锁键
func (k *KeyBuilder) Lock(resource string) string {
	return k.Build("lock", resource)
}

// Counter 构建计数器键
func (k *KeyBuilder) Counter(name string) string {
	return k.Build("counter", name)
}

// 预定义键模式常量
const (
	UserPrefix      = "user"
	SessionPrefix   = "session"
	SMSCodePrefix   = "sms:code"
	SMSLimitPrefix  = "sms:limit"
	ConfigPrefix    = "config"
	LockPrefix      = "lock"
	CounterPrefix   = "counter"
)
