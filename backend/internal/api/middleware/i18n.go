package middleware

import (
	"context"
	"strings"

	"license-manager/pkg/i18n"
	pkgcontext "license-manager/pkg/context"

	"github.com/gin-gonic/gin"
)

const (
	// LanguageKey 存储语言信息的上下文键
	LanguageKey = "language"

	// DefaultLanguage 默认语言
	DefaultLanguage = "zh-CN"
)

// I18nConfig 多语言中间件配置
type I18nConfig struct {
	Enable       bool     `yaml:"enable"`        // 是否启用多语言
	DefaultLang  string   `yaml:"default_lang"`  // 默认语言
	SupportLangs []string `yaml:"support_langs"` // 支持的语言列表
}

// I18nMiddleware 多语言检测中间件
func I18nMiddleware(config *I18nConfig) gin.HandlerFunc {
	// 如果未启用多语言，使用默认语言
	if config == nil || !config.Enable {
		return func(c *gin.Context) {
			c.Set(LanguageKey, config.DefaultLang)
			c.Next()
		}
	}

	return func(c *gin.Context) {
		lang := detectLanguage(c, config)
		c.Set(LanguageKey, lang)
		c.Next()
	}
}

// detectLanguage 检测客户端语言偏好
// 按优先级顺序：URL参数 > Accept-Language头 > 用户配置 > 系统默认
func detectLanguage(c *gin.Context, config *I18nConfig) string {
	// 1. 检查URL参数 ?lang=en-US (最高优先级)
	if lang := c.Query("lang"); lang != "" {
		if isLanguageSupported(lang, config.SupportLangs) {
			return lang
		}
	}

	// 2. 检查HTTP头 Accept-Language
	if acceptLang := c.GetHeader("Accept-Language"); acceptLang != "" {
		if lang := parseAcceptLanguageHeader(acceptLang, config.SupportLangs); lang != "" {
			return lang
		}
	}

	// 3. 检查用户配置（如果有用户身份验证）
	// 这里可以从用户数据库或JWT token中获取用户语言偏好
	if userLang := getUserLanguagePreference(c); userLang != "" {
		if isLanguageSupported(userLang, config.SupportLangs) {
			return userLang
		}
	}

	// 4. 返回系统默认语言
	return config.DefaultLang
}

// parseAcceptLanguageHeader 解析Accept-Language头部
// 例：zh-CN,zh;q=0.9,en;q=0.8,ja;q=0.7
func parseAcceptLanguageHeader(acceptLang string, supportedLangs []string) string {
	if acceptLang == "" {
		return ""
	}

	// 解析Accept-Language头部
	languages := strings.Split(acceptLang, ",")

	// 按权重排序的语言偏好
	for _, langWithQ := range languages {
		// 移除权重信息，例如 "zh-CN;q=0.9" -> "zh-CN"
		lang := strings.TrimSpace(langWithQ)
		if idx := strings.Index(lang, ";"); idx != -1 {
			lang = lang[:idx]
		}

		lang = strings.TrimSpace(lang)

		// 检查是否支持该语言
		if isLanguageSupported(lang, supportedLangs) {
			return lang
		}

		// 尝试匹配语言的主要部分，例如 "zh-CN" -> "zh"
		if idx := strings.Index(lang, "-"); idx != -1 {
			primaryLang := lang[:idx]
			// 查找支持的语言中是否有相同主要语言
			for _, supported := range supportedLangs {
				if strings.HasPrefix(supported, primaryLang+"-") {
					return supported
				}
			}
		}
	}

	return ""
}

// getUserLanguagePreference 获取用户语言偏好（从用户配置或JWT中）
func getUserLanguagePreference(c *gin.Context) string {
	// 这里可以实现从JWT token或用户数据库中获取用户语言偏好
	// 示例实现：从JWT Claims中获取语言偏好

	// 1. 从JWT token中获取（如果有认证中间件）
	if userID, exists := c.Get("user_id"); exists && userID != nil {
		// 这里可以查询数据库获取用户的语言偏好
		// userLang := getUserLanguageFromDB(userID.(string))
		// return userLang
		_ = userID // 暂时未实现
	}

	// 2. 从自定义头部获取
	if userLang := c.GetHeader("X-User-Language"); userLang != "" {
		return userLang
	}

	return ""
}

// isLanguageSupported 检查语言是否在支持列表中
func isLanguageSupported(lang string, supportedLangs []string) bool {
	if len(supportedLangs) == 0 {
		// 如果没有配置支持语言列表，检查i18n管理器
		return i18n.IsLanguageSupported(lang)
	}

	for _, supported := range supportedLangs {
		if supported == lang {
			return true
		}
	}
	return false
}

// GetLanguage 从gin.Context中获取语言设置
func GetLanguage(c *gin.Context) string {
	if lang, exists := c.Get(LanguageKey); exists {
		if langStr, ok := lang.(string); ok {
			return langStr
		}
	}
	return DefaultLanguage
}

// SetLanguage 设置当前请求的语言
func SetLanguage(c *gin.Context, lang string) {
	c.Set(LanguageKey, lang)
}

// 便利函数：创建不同配置的中间件

// DefaultI18nMiddleware 使用默认配置的多语言中间件
func DefaultI18nMiddleware() gin.HandlerFunc {
	config := &I18nConfig{
		Enable:       true,
		DefaultLang:  "zh-CN",
		SupportLangs: []string{"zh-CN", "en-US", "ja-JP"},
	}
	return I18nMiddleware(config)
}

// SimpleI18nMiddleware 简化的多语言中间件，仅支持指定语言
func SimpleI18nMiddleware(defaultLang string, supportedLangs ...string) gin.HandlerFunc {
	config := &I18nConfig{
		Enable:       true,
		DefaultLang:  defaultLang,
		SupportLangs: supportedLangs,
	}
	return I18nMiddleware(config)
}

// DisableI18nMiddleware 禁用多语言的中间件，总是使用默认语言
func DisableI18nMiddleware(defaultLang string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(LanguageKey, defaultLang)
		c.Next()
	}
}

// I18nResponseHelper 多语言响应辅助函数
type I18nResponseHelper struct {
	Context *gin.Context
}

// NewI18nResponseHelper 创建多语言响应辅助器
func NewI18nResponseHelper(c *gin.Context) *I18nResponseHelper {
	return &I18nResponseHelper{Context: c}
}

// GetLang 获取当前请求的语言
func (h *I18nResponseHelper) GetLang() string {
	return GetLanguage(h.Context)
}

// ErrorResponse 创建多语言错误响应
func (h *I18nResponseHelper) ErrorResponse(code string, customMessage ...string) (int, string, string) {
	lang := h.GetLang()
	return i18n.NewI18nErrorResponse(code, lang, customMessage...)
}

// Context相关工具函数

// WithLanguage 将语言信息添加到Context中
func WithLanguage(ctx context.Context, c *gin.Context) context.Context {
	lang := GetLanguage(c)
	return pkgcontext.WithLanguage(ctx, lang)
}
