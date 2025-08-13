import { ref, computed, onMounted, onUnmounted } from 'vue'

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
}

export function useDevice() {
  const screenWidth = ref(window.innerWidth)
  
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
  
  // 监听窗口变化
  const handleResize = () => {
    screenWidth.value = window.innerWidth
  }
  
  onMounted(() => {
    window.addEventListener('resize', handleResize)
  })
  
  onUnmounted(() => {
    window.removeEventListener('resize', handleResize)
  })
  
  return {
    screenWidth,
    deviceType,
    isMobile,
    isTablet,
    isDesktop
  }
}