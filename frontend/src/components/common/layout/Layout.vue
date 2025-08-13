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

// 默认导航配置
const defaultNavItems: NavItem[] = [
  { id: "dashboard", label: "仪表盘", href: "/dashboard", icon: "dashboard" },
  { id: "customers", label: "客户管理", href: "/customers", icon: "customers" },
  { id: "licenses", label: "授权管理", href: "/licenses", icon: "licenses" },
  { id: "roles", label: "角色权限", href: "/roles", icon: "roles" },
  { id: "users", label: "系统用户", href: "/users", icon: "users" }
]

// 使用 store 和路由
const appStore = useAppStore()
const route = useRoute()
const router = useRouter()

// 计算当前激活的导航项
const navItems = computed(() => {
  return defaultNavItems.map(item => ({
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

// 响应式设备检测 - 统一断点系统
const checkResponsive = () => {
  const width = window.innerWidth
  const isMobile = width <= 768      // 768px 及以下为移动端
  const isTablet = width > 768 && width <= 1024  // 769-1024px 为平板
  
  appStore.setMobile(isMobile)
  
  // 移动端和小平板自动折叠侧边栏
  if (isMobile || isTablet) {
    appStore.setSidebarCollapsed(true)
  } else if (width > 1024) {
    // 桌面端默认展开侧边栏
    appStore.setSidebarCollapsed(false)
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
  z-index: 1999;
  backdrop-filter: blur(2px);
}

// 主内容区域
.layout-main {
  margin-left: 280px;
  height: 100vh;
  display: flex;
  flex-direction: column;
  transition: all 0.3s ease;
  
  &--sidebar-collapsed {
    margin-left: 64px;
  }
  
  &--mobile {
    margin-left: 0;
  }
}

// 页面内容
.layout-content {
  flex: 1;
  padding-top: 80px;
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
  min-height: calc(100vh - 80px);
  padding: 24px;
  max-width: 1600px;
  margin: 0 auto;
}

// 响应式设计 - 统一断点系统
// 平板以下：768px 及以下为移动端
@media (max-width: 768px) {
  .layout-main {
    margin-left: 0;
  }
  
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

// 平板：769px - 1024px 之间
@media (min-width: 769px) and (max-width: 1024px) {
  .layout-main {
    margin-left: 64px; // 平板显示折叠侧边栏
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