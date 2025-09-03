import { ref, computed, onMounted, onUnmounted } from 'vue'
import { debounce } from '@/utils/performance'

// 设备类型枚举
export enum DeviceType {
  Mobile = 'mobile',
  Tablet = 'tablet', 
  Desktop = 'desktop'
}

// 断点配置
const BREAKPOINTS = {
  mobile: 768,
  tablet: 1024
} as const

export function useDevice() {
  const screenWidth = ref(window.innerWidth)
  const screenHeight = ref(window.innerHeight)
  
  // 设备类型判断
  const deviceType = computed(() => {
    if (screenWidth.value < BREAKPOINTS.mobile) {
      return DeviceType.Mobile
    } else if (screenWidth.value < BREAKPOINTS.tablet) {
      return DeviceType.Tablet
    } else {
      return DeviceType.Desktop
    }
  })
  
  // 便捷判断
  const isMobile = computed(() => deviceType.value === DeviceType.Mobile)
  const isTablet = computed(() => deviceType.value === DeviceType.Tablet)  
  const isDesktop = computed(() => deviceType.value === DeviceType.Desktop)
  
  // 屏幕方向
  const orientation = computed(() => {
    return screenWidth.value > screenHeight.value ? 'landscape' : 'portrait'
  })
  
  // 是否触摸设备
  const isTouchDevice = computed(() => {
    return 'ontouchstart' in window || navigator.maxTouchPoints > 0
  })
  
  // 是否小屏幕（适用于密集布局调整）
  const isSmallScreen = computed(() => screenWidth.value < 640)
  
  // 监听窗口变化（使用防抖优化性能，避免频繁更新）
  const handleResize = debounce(() => {
    screenWidth.value = window.innerWidth
    screenHeight.value = window.innerHeight
  }, 150)
  
  onMounted(() => {
    // 初始化屏幕尺寸
    screenWidth.value = window.innerWidth
    screenHeight.value = window.innerHeight
    
    window.addEventListener('resize', handleResize)
  })
  
  onUnmounted(() => {
    window.removeEventListener('resize', handleResize)
  })
  
  return {
    screenWidth,
    screenHeight,
    deviceType,
    isMobile,
    isTablet,
    isDesktop,
    orientation,
    isTouchDevice,
    isSmallScreen,
    BREAKPOINTS
  }
}