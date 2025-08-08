<!--
/**
 * 主布局组件
 * 提供完整的后台管理系统布局结构，包括侧边栏、顶部导航和主内容区域
 * 支持响应式设计和侧边栏收起/展开功能
 */
-->
<template>
  <!-- 主布局容器 -->
  <div class="app-layout" :class="layoutClasses">
    <!-- 遮罩层 - 移动端使用 -->
    <div 
      v-if="isMobile && showSidebar"
      class="layout-overlay"
      @click="closeSidebar"
    ></div>

    <!-- 侧边栏 -->
    <Sidebar
      :app-name="appName"
      :nav-items="navItems"
      :collapsible="sidebarCollapsible"
      :default-collapsed="sidebarCollapsed"
      :class="{ 
        'sidebar--mobile-open': isMobile && showSidebar,
        'sidebar--collapsed': sidebarCollapsed 
      }"
      @nav-click="handleNavClick"
      @toggle="handleSidebarToggle"
    >
      <template #logo>
        <slot name="sidebar-logo">
          <span class="logo-text">{{ appName }}</span>
        </slot>
      </template>
      
      <template #nav-items>
        <slot name="sidebar-nav"></slot>
      </template>
      
      <template #footer>
        <slot name="sidebar-footer"></slot>
      </template>
    </Sidebar>

    <!-- 主内容区域 -->
    <div class="layout-main" :class="mainContentClasses">
      <!-- 顶部导航 -->
      <NavContent
        :breadcrumb-items="breadcrumbItems"
        :page-title="pageTitle"
        :show-search="showSearch"
        :search-placeholder="searchPlaceholder"
        :show-notifications="showNotifications"
        :notification-count="notificationCount"
        :show-settings="showSettings"
        :show-user-avatar="showUserAvatar"
        :user-avatar="userAvatar"
        :user-name="userName"
        :show-menu-button="isMobile"
        :class="{ 'nav-content--sidebar-collapsed': sidebarCollapsed }"
        @search="handleSearch"
        @breadcrumb-click="handleBreadcrumbClick"
        @notification-click="handleNotificationClick"
        @settings-click="handleSettingsClick"
        @user-click="handleUserClick"
        @menu-toggle="toggleSidebar"
      >
        <template #breadcrumb>
          <slot name="nav-breadcrumb"></slot>
        </template>
        
        <template #actions>
          <slot name="nav-actions"></slot>
        </template>
      </NavContent>
      <!-- 页面内容 -->
      <main class="layout-content">
        <div class="content-container">
          <slot name="default">
            <router-view />
          </slot>
        </div>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import Sidebar from './Sidebar.vue'
import NavContent from './NavContent.vue'

// 导航项接口定义
interface NavItem {
  id: string // 导航项唯一标识
  label: string // 显示文本
  href: string // 链接地址
  icon?: string // 图标类名（可选）
  active?: boolean // 是否为当前活跃项（可选）
  children?: NavItem[] // 子菜单项（可选）
}

// 面包屑导航项接口定义
interface BreadcrumbItem {
  title: string // 显示标题
  path?: string // 链接路径（可选，最后一项通常不需要链接）
}

// 组件属性接口定义
interface Props {
  appName?: string // 应用名称
  navItems?: NavItem[] // 导航菜单项列表
  sidebarCollapsible?: boolean // 侧边栏是否可折叠
  defaultSidebarCollapsed?: boolean // 侧边栏默认是否折叠
  breadcrumbItems?: BreadcrumbItem[] // 面包屑导航项列表
  pageTitle?: string // 当前页面标题
  showSearch?: boolean // 是否显示搜索框
  searchPlaceholder?: string // 搜索框占位符文本
  showNotifications?: boolean // 是否显示通知按钮
  notificationCount?: number // 未读通知数量
  showSettings?: boolean // 是否显示设置按钮
  showUserAvatar?: boolean // 是否显示用户头像
  userAvatar?: string // 用户头像URL
  userName?: string // 用户名称
}

// 定义组件属性和默认值
const props = withDefaults(defineProps<Props>(), {
  appName: 'Cedar',
  navItems: () => [],
  sidebarCollapsible: true,
  defaultSidebarCollapsed: false,
  breadcrumbItems: () => [],
  pageTitle: '',
  showSearch: true,
  searchPlaceholder: '搜索...',
  showNotifications: true,
  notificationCount: 0,
  showSettings: true,
  showUserAvatar: true,
  userName: ''
})

// 定义组件事件
const emit = defineEmits<{
  navClick: [item: NavItem, event: Event]
  search: [query: string]
  breadcrumbClick: [item: BreadcrumbItem, event: Event]
  notificationClick: []
  settingsClick: []
  userClick: []
}>()

// 响应式状态管理
const isMobile = ref(false)
const showSidebar = ref(false)
const sidebarCollapsed = ref(props.defaultSidebarCollapsed)

// 计算布局样式类
const layoutClasses = computed(() => ({
  'app-layout--mobile': isMobile.value,
  'app-layout--sidebar-collapsed': sidebarCollapsed.value
}))

// 计算主内容区样式类
const mainContentClasses = computed(() => ({
  'layout-main--sidebar-collapsed': sidebarCollapsed.value,
  'layout-main--mobile': isMobile.value
}))

// 处理导航点击事件
const handleNavClick = (item: NavItem, event: Event) => {
  if (isMobile.value) {
    showSidebar.value = false
  }
  emit('navClick', item, event)
}

// 处理侧边栏收起/展开
const handleSidebarToggle = (collapsed: boolean) => {
  sidebarCollapsed.value = collapsed
}

// 切换侧边栏显示状态
const toggleSidebar = () => {
  if (isMobile.value) {
    showSidebar.value = !showSidebar.value
  } else {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }
}

// 关闭侧边栏（移动端）
const closeSidebar = () => {
  showSidebar.value = false
}

// 处理搜索事件
const handleSearch = (query: string) => {
  emit('search', query)
}

// 处理面包屑导航点击
const handleBreadcrumbClick = (item: BreadcrumbItem, event: Event) => {
  emit('breadcrumbClick', item, event)
}

// 处理通知按钮点击
const handleNotificationClick = () => {
  emit('notificationClick')
}

// 处理设置按钮点击
const handleSettingsClick = () => {
  emit('settingsClick')
}

// 处理用户头像点击
const handleUserClick = () => {
  emit('userClick')
}

// 检测是否为移动设备
const checkMobile = () => {
  isMobile.value = window.innerWidth <= 1024
  if (!isMobile.value) {
    showSidebar.value = false
  }
}

// 组件挂载时初始化
onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
})

// 组件卸载时清理
onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})
</script>

<style scoped>
.app-layout {
  position: relative;
  width: 100vw;
  height: 100vh;
  background: #F7F8FA;
  overflow: hidden;
}

/* 遮罩层 */
.layout-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 999;
  backdrop-filter: blur(2px);
}

/* 主内容区域 */
.layout-main {
  margin-left: 280px;
  height: 100vh;
  display: flex;
  flex-direction: column;
  transition: all 0.3s ease;
}

.layout-main--sidebar-collapsed {
  margin-left: 64px;
}

.layout-main--mobile {
  margin-left: 0;
}

/* 页面内容 */
.layout-content {
  flex: 1;
  padding-top: 80px;
  overflow-y: auto;
  position: relative;
}

.content-container {
  min-height: calc(100vh - 80px);
  padding: 24px;
  max-width: 1600px;
  margin: 0 auto;
}

/* 移动端适配 */
@media (max-width: 1024px) {
  .layout-main {
    margin-left: 0;
  }
}

@media (max-width: 768px) {
  .content-container {
    padding: 16px;
  }
}

@media (max-width: 480px) {
  .content-container {
    padding: 12px;
  }
}

/* 滚动条样式 */
.layout-content::-webkit-scrollbar {
  width: 6px;
}

.layout-content::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

.layout-content::-webkit-scrollbar-thumb {
  background: rgba(0, 0, 0, 0.2);
  border-radius: 3px;
}

.layout-content::-webkit-scrollbar-thumb:hover {
  background: rgba(0, 0, 0, 0.3);
}

/* 动画效果 */
.app-layout * {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

/* 打印样式 */
@media print {
  .app-layout {
    background: white;
  }
  
  .layout-main {
    margin-left: 0;
  }
  
  .layout-content {
    padding-top: 0;
  }
  
  .content-container {
    max-width: none;
    padding: 0;
  }
}

/* 高对比度模式支持 */
@media (prefers-contrast: high) {
  .app-layout {
    background: white;
  }
}

/* 减少动画模式 */
@media (prefers-reduced-motion: reduce) {
  .app-layout *,
  .layout-main,
  .layout-content {
    transition: none !important;
  }
}
</style>