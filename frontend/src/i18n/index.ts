/*
 * @Author: 13895237362 2205451508@qq.com
 * @Date: 2025-07-29 14:57:17
 * @LastEditors: 13895237362 2205451508@qq.com
 * @LastEditTime: 2025-08-05 15:21:49
 * @FilePath: /vue-demo3.0/src/i18n/index.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { createI18n } from 'vue-i18n'
import en from './locales/en.json'
import zh from './locales/zh.json'
import ja from './locales/ja.json'

// Element Plus 语言包
import elementEnLocale from 'element-plus/es/locale/lang/en'
import elementZhLocale from 'element-plus/es/locale/lang/zh-cn'
import elementJaLocale from 'element-plus/es/locale/lang/ja'

// 分离管理：仅使用自定义语言包
const messages = {
  en,
  zh, 
  ja
}

// Element Plus 语言包映射
export const elementLocales = {
  en: elementEnLocale,
  zh: elementZhLocale,
  ja: elementJaLocale
}

// 支持的语言类型
export type SupportedLocale = 'en' | 'zh' | 'ja'

// 语言代码标准映射
export const LANGUAGE_MAP = {
  'zh': { code: 'zh', iso: 'zh-CN', name: '中文' },
  'en': { code: 'en', iso: 'en-US', name: 'English' }, 
  'ja': { code: 'ja', iso: 'ja-JP', name: '日本語' }
} as const

// 检测浏览器语言并返回支持的语言
function detectBrowserLanguage(): SupportedLocale {
  // 1. 优先使用本地存储中的语言偏好（用户显式选择）
  const savedLanguage = localStorage.getItem('userLanguage') as SupportedLocale
  if (savedLanguage && Object.keys(LANGUAGE_MAP).includes(savedLanguage)) {
    return savedLanguage
  }

  // 2. 检测浏览器语言设置
  const browserLanguages = navigator.languages || [navigator.language]
  
  for (const lang of browserLanguages) {
    const langCode = lang.toLowerCase()
    
    // 精确匹配
    if (langCode === 'zh-cn' || langCode === 'zh') return 'zh'
    if (langCode === 'en-us' || langCode === 'en') return 'en'
    if (langCode === 'ja-jp' || langCode === 'ja') return 'ja'
    
    // 前缀匹配
    if (langCode.startsWith('zh')) return 'zh'
    if (langCode.startsWith('en')) return 'en'
    if (langCode.startsWith('ja')) return 'ja'
  }
  
  // 3. 默认回退到中文
  return 'zh'
}

// 创建 i18n 实例
const i18n = createI18n({
  legacy: false, // 使用 Composition API 模式
  locale: detectBrowserLanguage(), // 动态检测语言
  fallbackLocale: 'en', // 备用语言
  messages
})

// 设置语言的函数
export function setI18nLanguage(locale: SupportedLocale) {
  // 验证语言代码
  if (!Object.keys(LANGUAGE_MAP).includes(locale)) {
    console.warn(`Unsupported locale: ${locale}, falling back to zh`)
    locale = 'zh'
  }
  
  // 更新 i18n 当前语言
  i18n.global.locale.value = locale
  
  // 设置 HTML 文档的 lang 属性为标准 ISO 代码
  const isoCode = LANGUAGE_MAP[locale].iso
  document.documentElement.setAttribute('lang', isoCode)
  
  // 保存用户语言偏好（用于下次访问和 Accept-Language 头）
  localStorage.setItem('userLanguage', locale)
  
  // 触发自定义事件，通知其他组件语言已变更
  window.dispatchEvent(new CustomEvent('language-change', { 
    detail: { locale, isoCode, name: LANGUAGE_MAP[locale].name } 
  }))
  
  return locale
}

// 获取当前语言的 ISO 代码
export function getCurrentLanguageISO(): string {
  const currentLocale = i18n.global.locale.value as SupportedLocale
  return LANGUAGE_MAP[currentLocale]?.iso || 'zh-CN'
}

// 获取当前Element Plus语言包
export function getCurrentElementLocale() {
  const currentLocale = i18n.global.locale.value as SupportedLocale
  return elementLocales[currentLocale] || elementLocales.zh
}

export default i18n