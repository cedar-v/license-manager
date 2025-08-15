<template>
  <div class="app-layout" :class="layoutClasses">
    <!-- 移动端遮罩层 -->
    <div 
      v-if="appStore.isMobile && !appStore.sidebarCollapsed" 
      class="layout-overlay"
      @click="appStore.setSidebarCollapsed(true)"
    ></div>
    
    <!-- 侧边栏 -->
    <Sidebar 
      :app-name="props.appName"
      :nav-items="navItems"
      @nav-click="handleNavClick"
    />
    
    <!-- 主内容区域 -->
    <div class="layout-main" :class="mainClasses">
      <!-- 顶部导航 -->
      <NavContent />
      
      <!-- 页面内容 -->
      <div class="layout-content">
        <div class="content-container">
          <slot />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/store/modules/app'
import Sidebar from './Sidebar.vue'
import NavContent from './NavContent.vue'

// 导航项和面包屑项接口定义
interface NavItem {
  id: string
  label: string
  href: string
  icon?: string
  active?: boolean
  children?: NavItem[]
}


// 组件 Props
interface Props {
  appName?: string
  pageTitle?: string
}

const props = withDefaults(defineProps<Props>(), {
  appName: 'Cedar-V',
  pageTitle: ''
})

// 使用国际化
const { t } = useI18n()

// 默认导航配置
const defaultNavItems = computed(() => [
  { id: "dashboard", label: t('navigation.menu.dashboard'), href: "/dashboard", icon: "dashboard" },
  { id: "customers", label: t('navigation.menu.customers'), href: "/customers", icon: "customers" },
  { id: "licenses", label: t('navigation.menu.licenses'), href: "/licenses", icon: "licenses" },
  { id: "roles", label: t('navigation.menu.roles'), href: "/roles", icon: "roles" },
  { id: "users", label: t('navigation.menu.users'), href: "/users", icon: "users" }
])

// 使用 store 和路由
const appStore = useAppStore()
const route = useRoute()
const router = useRouter()

// 计算当前激活的导航项
const navItems = computed(() => {
  return defaultNavItems.value.map(item => ({
    ...item,
    active: route.path === item.href
  }))
})

// 定义组件事件
const emit = defineEmits<{
  navClick: [item: NavItem, event: Event]
}>()


// 处理导航点击
const handleNavClick = (item: NavItem, event: Event) => {
  router.push(item.href)
  emit('navClick', item, event)
}

// 计算类名
const layoutClasses = computed(() => ({
  'app-layout--mobile': appStore.isMobile,
  'app-layout--sidebar-collapsed': appStore.sidebarCollapsed
}))

const mainClasses = computed(() => ({
  'layout-main--mobile': appStore.isMobile,
  'layout-main--sidebar-collapsed': appStore.sidebarCollapsed
}))

// 响应式设备检测 - 简化版本
const checkResponsive = () => {
  const width = window.innerWidth
  const isMobile = width <= 768      // 768px 及以下为移动端
  const isTablet = width > 768 && width <= 1024  // 769-1024px 为平板
  const isDesktop = width > 1024  // 1025px+ 桌面端（包含2K、4K）
  
  appStore.setMobile(isMobile)
  
  // 移动端和小平板自动折叠侧边栏
  if (isMobile || isTablet) {
    appStore.setSidebarCollapsed(true)
  } else if (isDesktop) {
    // 桌面端默认展开侧边栏（rem会自动缩放）
    appStore.setSidebarCollapsed(false)
  }
  
  // 调试信息：显示当前根字体大小
  if (isDesktop) {
    const rootFontSize = parseFloat(getComputedStyle(document.documentElement).fontSize)
    console.log(`桌面屏幕: ${width}px, 根字体: ${rootFontSize}px`)
    
    // 更新页面调试信息
    document.body.setAttribute('data-font-size', `${rootFontSize}px (${width}px)`)
    
    // 额外调试：检测屏幕类型
    if (width >= 3840) {
      console.log('🖥️ 4K屏幕检测')
    } else if (width >= 2560) {
      console.log('🖥️ 2K屏幕检测')
    } else {
      console.log('🖥️ 1080p屏幕检测')
    }
  }
}

// 生命周期
onMounted(() => {
  checkResponsive()
  window.addEventListener('resize', checkResponsive)
})

onUnmounted(() => {
  window.removeEventListener('resize', checkResponsive)
})
</script>

<style lang="scss" scoped>
.app-layout {
  position: relative;
  width: 100vw;
  height: 100vh;
  background: #F5F7FA;
  overflow: hidden;
  
}

// 遮罩层
.layout-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 1998;
  backdrop-filter: blur(2px);
}

// 主内容区域 - 桌面端使用vw单位适配2K/4K
.layout-main {
  margin-left: 14.58vw; /* 280px/1920 = 14.58vw */
  height: 100vh;
  display: flex;
  flex-direction: column;
  transition: all 0.3s ease;
  
  &--sidebar-collapsed {
    margin-left: 3.33vw; /* 64px/1920 = 3.33vw */
  }
  
  &--mobile {
    margin-left: 0;
  }
}

// 页面内容
.layout-content {
  flex: 1;
  padding-top: 4.17vw; /* 80px/1920 = 4.17vw */
  overflow-y: auto;
  position: relative;
  
  /* 滚动条样式 */
  &::-webkit-scrollbar {
    width: 6px;
  }
  
  &::-webkit-scrollbar-track {
    background: #f1f1f1;
    border-radius: 3px;
  }
  
  &::-webkit-scrollbar-thumb {
    background: rgba(0, 0, 0, 0.2);
    border-radius: 3px;
    
    &:hover {
      background: rgba(0, 0, 0, 0.3);
    }
  }
}

.content-container {
  min-height: calc(100vh - 80px); /* 默认使用固定像素的最小高度 */
  padding: 24px;
  width: 100%;
  margin: 0;
  box-sizing: border-box;
}

// 响应式设计 - 统一断点系统
// 平板以下：768px 及以下为移动端，使用固定px单位
@media (max-width: 1024px) {
  .layout-main {
    margin-left: 0; /* 移动端和平板从左边缘开始 */
  }
  
  .layout-content {
    padding-top: 80px; /* 移动端使用固定高度 */
  }
  
  /* 移动端继承基础样式的min-height设置 */
}

@media (max-width: 768px) {
  .content-container {
    padding: 16px;
  }
}

// 小屏手机：480px 及以下
@media (max-width: 480px) {
  .content-container {
    padding: 12px;
  }
}

/* 平板：769px - 1024px 之间 */
@media (min-width: 769px) and (max-width: 1024px) {
  .layout-main {
    margin-left: 64px; /* 平板显示折叠侧边栏 */
  }
}

/* 桌面端：使用vw单位统一适配2K/4K，使用flex布局充满高度 */
@media (min-width: 1025px) {
  .layout-main {
    margin-left: 14.58vw; /* 280px/1920 = 14.58vw */
    
    &--sidebar-collapsed {
      margin-left: 3.33vw; /* 64px/1920 = 3.33vw */
    }
  }
  
  .layout-content {
    padding-top: 4.17vw; /* 80px/1920 = 4.17vw */
  }
  
  .content-container {
    height: calc(100vh - 4.17vw); /* 精确计算可用高度：视口高度减去顶部导航栏高度 */
    padding: 1.25vw; /* 24px/1920 = 1.25vw */
    width: 100%; /* 充满整个屏幕 */
    margin: 0;
    box-sizing: border-box;
    display: flex; /* 桌面端使用flex布局传递高度给子组件 */
    flex-direction: column;
  }
}

// 动画效果
.app-layout * {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

// 打印样式
@media print {
  .app-layout {
    background: white;
  }
  
  .layout-main {
    margin-left: 0;
  }
  
  .layout-content {
    padding-top: 0;
    overflow: visible;
  }
}

// 高对比度模式支持
@media (prefers-contrast: high) {
  .app-layout {
    background: white;
  }
}

// 减少动画模式
@media (prefers-reduced-motion: reduce) {
  .app-layout *,
  .layout-main,
  .layout-overlay {
    transition: none !important;
    animation: none !important;
  }
}
</style>