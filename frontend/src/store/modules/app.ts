import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { DeviceType } from '@/utils/useDevice'

export const useAppStore = defineStore('app', () => {
  // 状态
  const sidebarCollapsed = ref(false)
  const deviceType = ref<DeviceType>(DeviceType.Desktop)
  const screenWidth = ref(window.innerWidth)
  const screenHeight = ref(window.innerHeight)
  const orientation = ref<'landscape' | 'portrait'>('landscape')
  const isTouchDevice = ref(false)
  const theme = ref<'light' | 'dark' | 'auto'>('light')
  const language = ref<'zh' | 'en' | 'ja'>('zh')
  const loading = ref(false)

  // 计算属性
  const isMobile = computed(() => deviceType.value === DeviceType.Mobile)
  const isTablet = computed(() => deviceType.value === DeviceType.Tablet)
  const isDesktop = computed(() => deviceType.value === DeviceType.Desktop)
  const isSmallScreen = computed(() => screenWidth.value < 640)
  
  const isDark = computed(() => {
    if (theme.value === 'auto') {
      return window.matchMedia('(prefers-color-scheme: dark)').matches
    }
    return theme.value === 'dark'
  })

  // 响应式侧边栏状态
  const shouldCollapseSidebar = computed(() => {
    return isMobile.value || (isTablet.value && orientation.value === 'portrait')
  })

  // 操作
  const toggleSidebar = () => {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  const setSidebarCollapsed = (collapsed: boolean) => {
    sidebarCollapsed.value = collapsed
  }

  const setDeviceInfo = (info: {
    deviceType: DeviceType
    screenWidth: number
    screenHeight: number
    orientation: 'landscape' | 'portrait'
    isTouchDevice: boolean
  }) => {
    deviceType.value = info.deviceType
    screenWidth.value = info.screenWidth
    screenHeight.value = info.screenHeight
    orientation.value = info.orientation
    isTouchDevice.value = info.isTouchDevice
    
    // 自动调整侧边栏状态
    if (shouldCollapseSidebar.value && !sidebarCollapsed.value) {
      setSidebarCollapsed(true)
    }
  }

  // 兼容性方法
  const setMobile = (mobile: boolean) => {
    deviceType.value = mobile ? DeviceType.Mobile : DeviceType.Desktop
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
    deviceType,
    screenWidth,
    screenHeight,
    orientation,
    isTouchDevice,
    theme,
    language,
    loading,
    // 计算属性
    isMobile,
    isTablet,
    isDesktop,
    isSmallScreen,
    isDark,
    shouldCollapseSidebar,
    // 操作
    toggleSidebar,
    setSidebarCollapsed,
    setDeviceInfo,
    setMobile,
    setTheme,
    setLanguage,
    setLoading
  }
})