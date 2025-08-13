/**
 * 语言管理工具类
 * 提供统一的语言切换、检测和 Accept-Language 头管理
 */

import { setI18nLanguage, LANGUAGE_MAP, type SupportedLocale } from '@/i18n'
import { useAppStore } from '@/store/modules/app'

// 导出类型供其他模块使用
export type { SupportedLocale }

export class LanguageManager {
  private static instance: LanguageManager
  
  public static getInstance(): LanguageManager {
    if (!LanguageManager.instance) {
      LanguageManager.instance = new LanguageManager()
    }
    return LanguageManager.instance
  }

  /**
   * 切换语言并同步所有相关状态
   * @param locale 目标语言代码
   */
  public changeLanguage(locale: SupportedLocale): void {
    try {
      // 1. 更新 i18n
      setI18nLanguage(locale)
      
      // 2. 更新 app store
      const appStore = useAppStore()
      appStore.setLanguage(locale)
      
      // 3. 记录日志
      console.log(`Language changed to: ${locale} (${LANGUAGE_MAP[locale].iso})`)
      
    } catch (error) {
      console.error('Failed to change language:', error)
      throw error
    }
  }

  /**
   * 获取当前语言信息
   */
  public getCurrentLanguage(): {
    code: SupportedLocale
    iso: string
    name: string
  } {
    const code = (localStorage.getItem('userLanguage') as SupportedLocale) || 'zh'
    const languageInfo = LANGUAGE_MAP[code]
    
    return {
      code,
      iso: languageInfo.iso,
      name: languageInfo.name
    }
  }

  /**
   * 生成标准的 Accept-Language 头
   * 格式: 主语言;q=1.0, 备用语言;q=0.8, *;q=0.5
   */
  public generateAcceptLanguageHeader(): string {
    const currentLang = this.getCurrentLanguage()
    const primaryLang = currentLang.iso
    const fallbackLang = primaryLang !== 'en-US' ? 'en-US' : 'zh-CN'
    
    return `${primaryLang};q=1.0, ${fallbackLang};q=0.8, *;q=0.5`
  }

  /**
   * 获取所有支持的语言列表
   */
  public getSupportedLanguages(): Array<{
    code: SupportedLocale
    iso: string
    name: string
  }> {
    return Object.entries(LANGUAGE_MAP).map(([code, info]) => ({
      code: code as SupportedLocale,
      iso: info.iso,
      name: info.name
    }))
  }

  /**
   * 检测浏览器首选语言
   */
  public detectBrowserLanguage(): SupportedLocale {
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
    
    return 'zh' // 默认中文
  }

  /**
   * 初始化语言设置
   * 优先级：localStorage > 浏览器语言 > 默认中文
   */
  public initializeLanguage(): SupportedLocale {
    const savedLanguage = localStorage.getItem('userLanguage') as SupportedLocale
    
    if (savedLanguage && Object.keys(LANGUAGE_MAP).includes(savedLanguage)) {
      this.changeLanguage(savedLanguage)
      return savedLanguage
    }
    
    const browserLanguage = this.detectBrowserLanguage()
    this.changeLanguage(browserLanguage)
    return browserLanguage
  }

  /**
   * 监听语言变更事件
   */
  public onLanguageChange(callback: (event: CustomEvent) => void): void {
    window.addEventListener('language-change', callback as EventListener)
  }

  /**
   * 移除语言变更监听器
   */
  public offLanguageChange(callback: (event: CustomEvent) => void): void {
    window.removeEventListener('language-change', callback as EventListener)
  }
}

// 导出单例实例
export const languageManager = LanguageManager.getInstance()

// 导出便捷函数
export const changeLanguage = (locale: SupportedLocale) => languageManager.changeLanguage(locale)
export const getCurrentLanguage = () => languageManager.getCurrentLanguage()
export const generateAcceptLanguageHeader = () => languageManager.generateAcceptLanguageHeader()