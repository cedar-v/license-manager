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

// 检测浏览器语言并返回支持的语言
function detectBrowserLanguage(): SupportedLocale {
  // 优先使用本地存储中的语言偏好
  const savedLanguage = localStorage.getItem('userLanguage') as SupportedLocale
  if (savedLanguage && ['en', 'zh', 'ja'].includes(savedLanguage)) {
    return savedLanguage
  }

  // 获取浏览器语言
  const browserLanguage = navigator.language.toLowerCase()
  
  // 匹配浏览器语言到支持的语言
  if (browserLanguage.startsWith('en')) {
    return 'en'
  } else if (browserLanguage.startsWith('ja')) {
    return 'ja'
  } else if (browserLanguage.startsWith('zh')) {
    return 'zh'
  }
  
  // 默认回退到中文
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
  i18n.global.locale.value = locale
  // 设置 HTML 文档的 lang 属性
  document.documentElement.setAttribute('lang', locale)
  return locale
}

// 获取当前Element Plus语言包
export function getCurrentElementLocale() {
  const currentLocale = i18n.global.locale.value as SupportedLocale
  return elementLocales[currentLocale] || elementLocales.zh
}

export default i18n