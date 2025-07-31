/*
 * @Author: 13895237362 2205451508@qq.com
 * @Date: 2025-07-29 14:57:17
 * @LastEditors: 13895237362 2205451508@qq.com
 * @LastEditTime: 2025-07-29 15:51:27
 * @FilePath: /vue-demo3.0/src/i18n/index.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { createI18n, type Composer } from 'vue-i18n'
import { nextTick } from 'vue'

type SupportedLocale = 'zh-CN' | 'en-US' | 'ja-JP'

const DEFAULT_LOCALE: SupportedLocale = 'zh-CN'
const SUPPORTED_LOCALES: SupportedLocale[] = ['zh-CN', 'en-US', 'ja-JP']

const getStartingLocale = (): SupportedLocale => {
  const storedLocale = localStorage.getItem('locale') as SupportedLocale | null
  if (storedLocale && SUPPORTED_LOCALES.includes(storedLocale)) {
    return storedLocale
  }

  const browserLocale = navigator.language as SupportedLocale
  if (SUPPORTED_LOCALES.includes(browserLocale)) {
    return browserLocale
  }

  const simplifiedLocale = browserLocale.split('-')[0]
  const matchedLocale = SUPPORTED_LOCALES.find(locale => 
    locale.startsWith(simplifiedLocale)
  )

  return matchedLocale || DEFAULT_LOCALE
}

const i18n = createI18n({
  legacy: false,
  locale: getStartingLocale(),
  fallbackLocale: DEFAULT_LOCALE,
  globalInjection: true,
  silentTranslationWarn: true,
  messages: {},
  datetimeFormats: {},
  numberFormats: {}
})

// 按需语言包
export async function loadLocaleMessages(locale: SupportedLocale): Promise<Record<string, string>> {
  try {
    // 修改这里的路径为平级引用
    const messages = await import(`./locales/${locale}.json`)
    ;(i18n.global as Composer).setLocaleMessage(locale, messages.default)
    return messages.default
  } catch (error) {
    console.error(`Failed to load locale messages for ${locale}`, error)
    if (locale !== DEFAULT_LOCALE) {
      return loadLocaleMessages(DEFAULT_LOCALE)
    }
    throw error
  }
}

// 设置语言
export async function setI18nLanguage(locale: SupportedLocale): Promise<void> {
  if (!SUPPORTED_LOCALES.includes(locale)) {
    locale = DEFAULT_LOCALE
  }

  if (!i18n.global.availableLocales.includes(locale)) {
    await loadLocaleMessages(locale)
  }

  i18n.global.locale.value = locale
  document.querySelector('html')?.setAttribute('lang', locale)
  localStorage.setItem('locale', locale)

  return nextTick()
}

// 初始化
setI18nLanguage(i18n.global.locale.value as SupportedLocale)

export default i18n