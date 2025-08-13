<!--
/**
 * 顶部导航栏组件
 * 提供面包屑导航、页面标题、搜索、通知、语言切换、主题切换和用户信息
 * 支持响应式设计和移动端适配
 */
-->
<template>
  <!-- 顶部导航栏容器 -->
  <header class="nav-content" :class="navClasses">
    <!-- 左侧：侧边栏切换按钮和面包屑 -->
    <div class="nav-content__left">
      <!-- 侧边栏切换按钮 -->
      <button 
        class="sidebar-toggle-btn"
        @click="handleSidebarToggle"
        :aria-label="'切换侧边栏'"
      >
        <NavIcon name="sidebar-toggle" size="large" />
      </button>

      <!-- 面包屑导航 -->
      <div class="breadcrumb-section">
        <nav class="breadcrumb" v-if="breadcrumbs.length > 0">
          <span 
            v-for="(item, index) in breadcrumbs" 
            :key="item.path || index"
            class="breadcrumb-item"
          >
            <button 
              v-if="item.path && index < breadcrumbs.length - 1"
              class="breadcrumb-link"
              @click="navigateTo(item)"
            >
              {{ item.title }}
            </button>
            <span 
              v-else
              class="breadcrumb-current page-title"
            >
              {{ item.title }}
            </span>
            <span 
              v-if="index < breadcrumbs.length - 1"
              class="breadcrumb-separator"
            >
              /
            </span>
          </span>
        </nav>
      </div>
    </div>

    <!-- 右侧：用户头像和操作按钮 -->
    <div class="nav-content__right">
      <!-- 用户头像 -->
      <div class="user-avatar" @click="handleUserClick">
        <NavIcon name="user" />
      </div>

      <!-- 分割线 -->
      <div class="divider"></div>

      <!-- 操作按钮组 -->
      <div class="action-buttons">
        <!-- 搜索按钮 -->
        <button class="action-btn" @click="handleSearchClick">
          <NavIcon name="search" />
        </button>

        <!-- 通知按钮 -->
        <button class="action-btn notification-btn" @click="handleNotificationClick">
          <NavIcon name="notification" />
          <span v-if="notificationCount" class="notification-badge">
            {{ notificationCount > 99 ? '99+' : notificationCount }}
          </span>
        </button>

        <!-- 语言切换按钮 -->
        <button class="action-btn" @click="handleLanguageClick">
          <NavIcon name="language" />
        </button>

        <!-- 主题切换按钮 -->
        <button class="action-btn" @click="handleThemeClick">
          <NavIcon name="dark-mode" />
        </button>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAppStore } from '@/store/modules/app'
import NavIcon from '@/components/common/icons/NavIcon.vue'
import { useBreadcrumb } from '@/utils/breadcrumb'

// 组件属性接口定义
interface Props {
  notificationCount?: number // 未读通知数量
}

// 定义组件属性和默认值
const { notificationCount = 24 } = defineProps<Props>()

// 定义组件事件
const emit = defineEmits<{
  sidebarToggle: []
  searchClick: []
  notificationClick: []
  languageClick: []
  themeClick: []
  userClick: []
}>()

// 使用store和组合函数
const appStore = useAppStore()
const { breadcrumbs, navigateTo } = useBreadcrumb()

// 计算导航栏样式类
const navClasses = computed(() => ({
  'nav-content--sidebar-collapsed': appStore.sidebarCollapsed,
  'nav-content--mobile': appStore.isMobile
}))

// 处理侧边栏切换
const handleSidebarToggle = () => {
  appStore.toggleSidebar()
  emit('sidebarToggle')
}

// 处理搜索点击
const handleSearchClick = () => {
  emit('searchClick')
}

// 处理通知按钮点击
const handleNotificationClick = () => {
  emit('notificationClick')
}

// 处理语言切换点击
const handleLanguageClick = () => {
  emit('languageClick')
}

// 处理主题切换点击
const handleThemeClick = () => {
  emit('themeClick')
}

// 处理用户头像点击
const handleUserClick = () => {
  emit('userClick')
}
</script>

<style scoped>
.nav-content {
  position: fixed;
  top: 0;
  left: 280px;
  right: 0;
  height: 80px;
  background: #FFFFFF;
  border-bottom: 1px solid rgba(29, 29, 29, 0.12);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  z-index: 999;
  transition: all 0.3s ease;
}

/* 左侧区域 */
.nav-content__left {
  display: flex;
  align-items: center;
  gap: 12px;
}

/* 侧边栏切换按钮 */
.sidebar-toggle-btn {
  width: 32px;
  height: 32px;
  border: none;
  background: transparent;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 4px;
  transition: all 0.2s;
}

.sidebar-toggle-btn:hover {
  background: rgba(1, 156, 124, 0.1);
}

/* 面包屑区域 */
.breadcrumb-section {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

/* 面包屑导航 */
.breadcrumb {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #888;
  flex-wrap: wrap;
}

.breadcrumb-item {
  display: flex;
  align-items: center;
  gap: 4px;
  white-space: nowrap;
}

/* 面包屑当前页面样式，使用页面标题样式 */
.breadcrumb-current.page-title {
  font-family: 'OPPOSans', sans-serif;
  font-size: 20px;
  font-weight: 400;
  color: #1D1D1D;
  line-height: 1.32;
}

.breadcrumb-link {
  border: none;
  background: none;
  color: #019C7C;
  font-size: 12px;
  cursor: pointer;
  padding: 0;
  text-decoration: none;
  transition: all 0.2s;
}

.breadcrumb-link:hover {
  color: #017A63;
  text-decoration: underline;
}

.breadcrumb-current {
  color: #666;
  font-size: 12px;
}

.breadcrumb-separator {
  color: #ccc;
  margin: 0 4px;
  font-size: 12px;
}

.page-title {
  font-family: 'OPPOSans', sans-serif;
  font-size: 20px;
  font-weight: 400;
  color: #1D1D1D;
  line-height: 1.32;
}

/* 右侧区域 */
.nav-content__right {
  display: flex;
  align-items: center;
  gap: 24px;
}

/* 用户头像 */
.user-avatar {
  width: 20px;
  height: 20px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.user-avatar:hover {
  opacity: 0.7;
}

/* 分割线 */
.divider {
  width: 1px;
  height: 24px;
  background: rgba(29, 29, 29, 0.12);
  border-radius: 2px;
}

/* 操作按钮组 */
.action-buttons {
  display: flex;
  gap: 12px;
}

/* 操作按钮 */
.action-btn {
  position: relative;
  width: 40px;
  height: 40px;
  border: none;
  background: #F7F8FA;
  border-radius: 20px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.action-btn:hover {
  background: rgba(1, 156, 124, 0.1);
}

/* 通知徽章 */
.notification-badge {
  position: absolute;
  top: -3px;
  right: -1px;
  width: 22px;
  height: 16px;
  background: #00C27C;
  color: white;
  font-size: 10px;
  font-weight: 700;
  border-radius: 90px;
  border: 1px solid #FFFFFF;
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: 'Roboto', sans-serif;
  line-height: 1.3;
}

/* 侧边栏状态配合 */
.nav-content--sidebar-collapsed {
  left: 64px;
}

.nav-content--mobile {
  left: 0;
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .nav-content {
    left: 0;
    padding: 16px 20px;
  }
  
  /* 覆盖侧边栏收起状态，移动端和平板端始终从左侧边缘开始 */
  .nav-content--sidebar-collapsed {
    left: 0;
  }
}

@media (max-width: 768px) {
  .nav-content {
    padding: 12px 16px;
  }
  
  .nav-content__right {
    gap: 16px;
  }
  
  .action-buttons {
    gap: 8px;
  }
  
  .action-btn {
    width: 36px;
    height: 36px;
  }
  
  .page-title {
    font-size: 18px;
  }
}

@media (max-width: 480px) {
  .nav-content {
    padding: 10px 12px;
  }
  
  .page-title {
    font-size: 16px;
  }
  
  .action-buttons {
    gap: 6px;
  }
  
  .action-btn {
    width: 32px;
    height: 32px;
  }
}

/* 桌面端特定样式 */
@media (min-width: 1025px) {
  .nav-content--sidebar-collapsed {
    left: 64px;
  }
}
</style>