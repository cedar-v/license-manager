import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useAppStore = defineStore('app', () => {
  // 状态
  const sidebarCollapsed = ref(false)
  const isMobile = ref(false)
  const theme = ref<'light' | 'dark' | 'auto'>('light')
  const language = ref<'zh' | 'en' | 'ja'>('zh')
  const loading = ref(false)

  // 计算属性
  const isDark = computed(() => {
    if (theme.value === 'auto') {
      return window.matchMedia('(prefers-color-scheme: dark)').matches
    }
    return theme.value === 'dark'
  })

  // 操作
  const toggleSidebar = () => {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  const setSidebarCollapsed = (collapsed: boolean) => {
    sidebarCollapsed.value = collapsed
  }

  const setMobile = (mobile: boolean) => {
    isMobile.value = mobile
  }

  const setTheme = (newTheme: 'light' | 'dark' | 'auto') => {
    theme.value = newTheme
    document.documentElement.setAttribute('data-theme', newTheme)
  }

  const setLanguage = (lang: 'zh' | 'en' | 'ja') => {
    language.value = lang
    // 同时更新 localStorage 以保持一致性
    localStorage.setItem('userLanguage', lang)
  }

  const setLoading = (isLoading: boolean) => {
    loading.value = isLoading
  }

  return {
    // 状态
    sidebarCollapsed,
    isMobile,
    theme,
    language,
    loading,
    // 计算属性
    isDark,
    // 操作
    toggleSidebar,
    setSidebarCollapsed,
    setMobile,
    setTheme,
    setLanguage,
    setLoading
  }
})