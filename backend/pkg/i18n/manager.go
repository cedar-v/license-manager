package i18n

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/yaml.v2"
)

// SystemInfo 语言包系统信息
type SystemInfo struct {
	Name   string `yaml:"name"`
	Locale string `yaml:"locale"`
}

// LanguageData 语言包数据结构
type LanguageData struct {
	System       SystemInfo             `yaml:"system"`
	Errors       map[string]interface{} `yaml:"errors"` // 灵活结构：支持直接字符串和嵌套对象
	Enums        map[string]interface{} `yaml:"enums"`  // 枚举值多语言映射
	DefaultError string                 `yaml:"default_error"`
	flatErrors   map[string]string      // 扁平化的错误码映射（内部使用）
	flatEnums    map[string]string      // 扁平化的枚举值映射（内部使用）
}

// Manager 多语言管理器
type Manager struct {
	defaultLang string
	loadedLangs map[string]*LanguageData
	configPath  string
	mutex       sync.RWMutex
}

// I18nManager 多语言管理器接口
type I18nManager interface {
	LoadLanguage(lang string) error
	GetErrorMessage(code, lang string) string
	GetEnumMessage(enumType, enumValue, lang string) string
	GetAllEnums(lang string) map[string]map[string]string
	GetEnumsByType(enumType, lang string) map[string]string
	SupportedLanguages() []string
	SetDefaultLanguage(lang string)
	IsLanguageSupported(lang string) bool
}

// NewManager 创建新的多语言管理器
func NewManager(configPath, defaultLang string) *Manager {
	return &Manager{
		defaultLang: defaultLang,
		loadedLangs: make(map[string]*LanguageData),
		configPath:  configPath,
	}
}

// LoadLanguage 加载指定语言包
func (m *Manager) LoadLanguage(lang string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 检查是否已加载
	if _, exists := m.loadedLangs[lang]; exists {
		return nil
	}

	// 构建语言文件路径
	filePath := filepath.Join(m.configPath, fmt.Sprintf("%s.yaml", lang))

	// 读取文件内容
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read language file %s: %v", filePath, err)
	}

	// 解析YAML
	var langData LanguageData
	if err := yaml.Unmarshal(data, &langData); err != nil {
		return fmt.Errorf("failed to parse language file %s: %v", filePath, err)
	}

	// 扁平化错误码映射，便于快速查找
	langData.flatErrors = make(map[string]string)
	langData.flatEnums = make(map[string]string)

	// 处理错误码映射
	for key, value := range langData.Errors {
		switch v := value.(type) {
		case string:
			// 直接定义的错误码（如 "000000": "成功"）
			langData.flatErrors[key] = v
		case map[interface{}]interface{}:
			// 嵌套模块的错误码（如 auth: {"100001": "xxx"}）
			for subKey, subValue := range v {
				if code, ok := subKey.(string); ok {
					if msg, ok := subValue.(string); ok {
						langData.flatErrors[code] = msg
					}
				}
			}
		}
	}

	// 处理枚举值映射
	for enumType, enumData := range langData.Enums {
		if enumMap, ok := enumData.(map[interface{}]interface{}); ok {
			for enumKey, enumValue := range enumMap {
				if key, ok := enumKey.(string); ok {
					if value, ok := enumValue.(string); ok {
						// 使用 "enumType.enumKey" 作为扁平化的键
						flatKey := fmt.Sprintf("%s.%s", enumType, key)
						langData.flatEnums[flatKey] = value
					}
				}
			}
		}
	}

	// 存储加载的语言数据
	m.loadedLangs[lang] = &langData

	return nil
}

// GetErrorMessage 获取指定语言的错误信息
func (m *Manager) GetErrorMessage(code, lang string) string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// 首先尝试获取请求语言的错误信息
	if langData, exists := m.loadedLangs[lang]; exists {
		if message, found := langData.flatErrors[code]; found {
			return message
		}
		// 如果错误码不存在，返回该语言的默认错误信息
		if langData.DefaultError != "" {
			return langData.DefaultError
		}
	}

	// 如果请求语言不存在或错误码不存在，降级到默认语言
	if lang != m.defaultLang {
		if langData, exists := m.loadedLangs[m.defaultLang]; exists {
			if message, found := langData.flatErrors[code]; found {
				return message
			}
			// 返回默认语言的默认错误信息
			if langData.DefaultError != "" {
				return langData.DefaultError
			}
		}
	}

	// 最终降级：返回硬编码的默认错误信息
	return "Unknown error"
}

// GetEnumMessage 获取指定语言的枚举值显示文本
func (m *Manager) GetEnumMessage(enumType, enumValue, lang string) string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// 构建查找键
	flatKey := fmt.Sprintf("%s.%s", enumType, enumValue)

	// 首先尝试获取请求语言的枚举值
	if langData, exists := m.loadedLangs[lang]; exists {
		if message, found := langData.flatEnums[flatKey]; found {
			return message
		}
	}

	// 如果请求语言不存在，降级到默认语言
	if lang != m.defaultLang {
		if langData, exists := m.loadedLangs[m.defaultLang]; exists {
			if message, found := langData.flatEnums[flatKey]; found {
				return message
			}
		}
	}

	// 最终降级：返回原始枚举值
	return enumValue
}

// GetAllEnums 获取指定语言的所有枚举值
func (m *Manager) GetAllEnums(lang string) map[string]map[string]string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	result := make(map[string]map[string]string)

	// 首先尝试获取请求语言的枚举值
	if langData, exists := m.loadedLangs[lang]; exists {
		result = m.buildEnumMap(langData)
		if len(result) > 0 {
			return result
		}
	}

	// 如果请求语言不存在，降级到默认语言
	if lang != m.defaultLang {
		if langData, exists := m.loadedLangs[m.defaultLang]; exists {
			result = m.buildEnumMap(langData)
		}
	}

	return result
}

// GetEnumsByType 获取指定类型的枚举值
func (m *Manager) GetEnumsByType(enumType, lang string) map[string]string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// 首先尝试获取请求语言的枚举值
	if langData, exists := m.loadedLangs[lang]; exists {
		if enumMap := m.getEnumMapByType(langData, enumType); len(enumMap) > 0 {
			return enumMap
		}
	}

	// 如果请求语言不存在，降级到默认语言
	if lang != m.defaultLang {
		if langData, exists := m.loadedLangs[m.defaultLang]; exists {
			return m.getEnumMapByType(langData, enumType)
		}
	}

	return make(map[string]string)
}

// buildEnumMap 构建完整的枚举映射
func (m *Manager) buildEnumMap(langData *LanguageData) map[string]map[string]string {
	result := make(map[string]map[string]string)

	for enumType, enumData := range langData.Enums {
		if enumMap, ok := enumData.(map[interface{}]interface{}); ok {
			typeMap := make(map[string]string)
			for enumKey, enumValue := range enumMap {
				if key, ok := enumKey.(string); ok {
					if value, ok := enumValue.(string); ok {
						typeMap[key] = value
					}
				}
			}
			if len(typeMap) > 0 {
				result[enumType] = typeMap
			}
		}
	}

	return result
}

// getEnumMapByType 获取指定类型的枚举映射
func (m *Manager) getEnumMapByType(langData *LanguageData, enumType string) map[string]string {
	result := make(map[string]string)

	if enumData, exists := langData.Enums[enumType]; exists {
		if enumMap, ok := enumData.(map[interface{}]interface{}); ok {
			for enumKey, enumValue := range enumMap {
				if key, ok := enumKey.(string); ok {
					if value, ok := enumValue.(string); ok {
						result[key] = value
					}
				}
			}
		}
	}

	return result
}

// SupportedLanguages 返回已加载的语言列表
func (m *Manager) SupportedLanguages() []string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	languages := make([]string, 0, len(m.loadedLangs))
	for lang := range m.loadedLangs {
		languages = append(languages, lang)
	}
	return languages
}

// SetDefaultLanguage 设置默认语言
func (m *Manager) SetDefaultLanguage(lang string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.defaultLang = lang
}

// IsLanguageSupported 检查语言是否被支持
func (m *Manager) IsLanguageSupported(lang string) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	_, exists := m.loadedLangs[lang]
	return exists
}

// GetLanguageInfo 获取语言包信息
func (m *Manager) GetLanguageInfo(lang string) *SystemInfo {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if langData, exists := m.loadedLangs[lang]; exists {
		return &langData.System
	}
	return nil
}

// parseAcceptLanguage 解析 Accept-Language 头部
func ParseAcceptLanguage(acceptLang string) string {
	if acceptLang == "" {
		return ""
	}

	// 简单解析 Accept-Language 头部，获取权重最高的语言
	// 例：zh-CN,zh;q=0.9,en;q=0.8
	languages := strings.Split(acceptLang, ",")
	if len(languages) > 0 {
		// 取第一个语言（权重最高）
		firstLang := strings.TrimSpace(languages[0])
		// 移除权重信息
		if idx := strings.Index(firstLang, ";"); idx != -1 {
			firstLang = firstLang[:idx]
		}
		return strings.TrimSpace(firstLang)
	}

	return ""
}

// 全局管理器实例（单例模式）
var globalManager *Manager
var once sync.Once

// InitGlobalManager 初始化全局管理器
func InitGlobalManager(configPath, defaultLang string) error {
	var initErr error
	once.Do(func() {
		globalManager = NewManager(configPath, defaultLang)

		// 预加载默认语言
		if err := globalManager.LoadLanguage(defaultLang); err != nil {
			initErr = fmt.Errorf("failed to load default language %s: %v", defaultLang, err)
		}
	})
	return initErr
}

// GetGlobalManager 获取全局管理器实例
func GetGlobalManager() *Manager {
	return globalManager
}

// GetErrorMessage 全局函数：获取错误信息
func GetErrorMessage(code, lang string) string {
	if globalManager == nil {
		return "I18n not initialized"
	}
	return globalManager.GetErrorMessage(code, lang)
}

// LoadLanguage 全局函数：加载语言
func LoadLanguage(lang string) error {
	if globalManager == nil {
		return fmt.Errorf("I18n not initialized")
	}
	return globalManager.LoadLanguage(lang)
}

// GetEnumMessage 全局函数：获取枚举值显示文本
func GetEnumMessage(enumType, enumValue, lang string) string {
	if globalManager == nil {
		return enumValue
	}
	return globalManager.GetEnumMessage(enumType, enumValue, lang)
}

// GetAllEnums 全局函数：获取所有枚举值
func GetAllEnums(lang string) map[string]map[string]string {
	if globalManager == nil {
		return make(map[string]map[string]string)
	}
	return globalManager.GetAllEnums(lang)
}

// GetEnumsByType 全局函数：获取指定类型的枚举值
func GetEnumsByType(enumType, lang string) map[string]string {
	if globalManager == nil {
		return make(map[string]string)
	}
	return globalManager.GetEnumsByType(enumType, lang)
}
