<template>
  <div class="app-layout" :class="layoutClasses">
    <!-- 移动端遮罩层 -->
    <div 
      v-if="appStore.isMobile && !appStore.sidebarCollapsed" 
      class="layout-overlay"
      @click="handleMobileOverlay"
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
        <slot />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
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
    // 对于有子路由的项目，匹配路径前缀
    active: route.path === item.href || route.path.startsWith(item.href + '/')
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

// 移动端自动关闭侧边栏
const handleMobileOverlay = () => {
  if (appStore.isMobile) {
    appStore.setSidebarCollapsed(true)
  }
}
</script>

<style lang="scss" scoped>
// Variables and mixins are auto-injected via Vite configuration

.app-layout {
  position: relative;
  width: 100vw;
  height: 100vh;
  background: $background-color-base;
  overflow: hidden;
  display: flex;
}

// 遮罩层
.layout-overlay {
  @include full-overlay;
  background: rgba(0, 0, 0, 0.5);
  z-index: 1998;
  backdrop-filter: blur(2px);
}

// 主内容区域
.layout-main {
  margin-left: 280px;
  height: 100vh;
  display: flex;
  flex-direction: column;
  transition: margin-left 0.3s ease;
  flex: 1;
  min-width: 0; // 防止 flex 溢出
  
  &--sidebar-collapsed {
    margin-left: 64px;
  }
  
  &--mobile {
    margin-left: 0;
  }
  
  // 桌面端响应式
  @include desktop-up {
    margin-left: 280px;
    
    &--sidebar-collapsed {
      margin-left: 64px;
    }
  }
  
  // 平板端
  @include tablet {
    margin-left: 64px;
  }
  
  // 移动端
  @include mobile {
    margin-left: 0;
  }
}

// 页面内容
.layout-content {
  flex: 1;
  padding-top: 80px; // 为固定导航栏预留空间
  overflow: auto; 
  position: relative;
  min-height: 0; // 防止 flex 溢出
  display: flex; // 添加flex布局传递给子组件
  flex-direction: column;
  
  // 移动端优化
  @include mobile {
    padding: $spacing-medium;
    padding-top: calc($spacing-medium + 80px); // 移动端也要为导航栏预留空间
  }
  
  
  // 滚动条样式
  &::-webkit-scrollbar {
    width: 6px;
  }
  
  &::-webkit-scrollbar-track {
    background: $background-color-base;
    border-radius: 3px;
  }
  
  &::-webkit-scrollbar-thumb {
    background: rgba(0, 0, 0, 0.2);
    border-radius: 3px;
    
    @include non-touch-device {
      &:hover {
        background: rgba(0, 0, 0, 0.3);
      }
    }
  }
}

// 响应式布局类
.app-layout {
  &--mobile {
    .layout-main {
      margin-left: 0;
    }
  }
  
  &--sidebar-collapsed {
    .layout-main {
      margin-left: 64px;
      
      @include mobile {
        margin-left: 0;
      }
    }
  }
}

// 动画效果
.layout-main,
.layout-overlay {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

// 打印样式
@media print {
  .app-layout {
    background: white;
  }
  
  .layout-main {
    margin-left: 0 !important;
  }
  
  .layout-content {
    padding: 0;
    overflow: visible;
  }
}

// 无障碍访问
@media (prefers-contrast: high) {
  .app-layout {
    background: white;
  }
  
  .layout-overlay {
    background: rgba(0, 0, 0, 0.8);
  }
}

// 减少动画模式
@media (prefers-reduced-motion: reduce) {
  .layout-main,
  .layout-overlay {
    transition: none !important;
  }
}
</style>