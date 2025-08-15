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
        :aria-label="t('navigation.toggleSidebar')"
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
import { useI18n } from 'vue-i18n'
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

// 使用国际化
const { t } = useI18n()

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
/* 顶部导航栏 - 基于1920*1080设计的vw适配，确保4K下正确显示 */
.nav-content {
  position: fixed;
  top: 0;
  left: 14.58vw; /* 280px/1920 = 14.58vw */
  right: 0;
  height: 4.17vw; /* 80px/1920 = 4.17vw */
  background: #FFFFFF;
  border-bottom: 1px solid rgba(29, 29, 29, 0.12);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1.04vw 1.25vw; /* 20px 24px/1920 = 1.04vw 1.25vw */
  z-index: 2001;
  transition: all 0.3s ease;
}

/* 左侧区域 */
.nav-content__left {
  display: flex;
  align-items: center;
  gap: 0.63vw; /* 12px/1920 = 0.63vw */
}

/* 侧边栏切换按钮 */
.sidebar-toggle-btn {
  width: 1.67vw; /* 32px/1920 = 1.67vw */
  height: 1.67vw; /* 32px/1920 = 1.67vw */
  border: none;
  background: transparent;
  border-radius: 0.42vw; /* 8px/1920 = 0.42vw */
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0.21vw; /* 4px/1920 = 0.21vw */
  transition: all 0.2s;
}

.sidebar-toggle-btn:hover {
  background: rgba(1, 156, 124, 0.1);
}

/* 面包屑区域 */
.breadcrumb-section {
  display: flex;
  flex-direction: column;
  gap: 0.10vw; /* 2px/1920 = 0.10vw */
}

/* 面包屑导航 */
.breadcrumb {
  display: flex;
  align-items: center;
  gap: 0.21vw; /* 4px/1920 = 0.21vw */
  font-size: 0.63vw; /* 12px/1920 = 0.63vw */
  color: #888;
  flex-wrap: wrap;
}

.breadcrumb-item {
  display: flex;
  align-items: center;
  gap: 0.21vw; /* 4px/1920 = 0.21vw */
  white-space: nowrap;
}

/* 面包屑当前页面样式，使用页面标题样式 */
.breadcrumb-current.page-title {
  font-family: 'OPPOSans', sans-serif;
  font-size: 1.04vw; /* 20px/1920 = 1.04vw */
  font-weight: 400;
  color: #1D1D1D;
  line-height: 1.32;
}

.breadcrumb-link {
  border: none;
  background: none;
  color: #019C7C;
  font-size: 0.63vw; /* 12px/1920 = 0.63vw */
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
  font-size: 0.63vw; /* 12px/1920 = 0.63vw */
}

.breadcrumb-separator {
  color: #ccc;
  margin: 0 0.21vw; /* 4px/1920 = 0.21vw */
  font-size: 0.63vw; /* 12px/1920 = 0.63vw */
}

.page-title {
  font-family: 'OPPOSans', sans-serif;
  font-size: 1.04vw; /* 20px/1920 = 1.04vw */
  font-weight: 400;
  color: #1D1D1D;
  line-height: 1.32;
}

/* 右侧区域 */
.nav-content__right {
  display: flex;
  align-items: center;
  gap: 1.25vw; /* 24px/1920 = 1.25vw */
}

/* 用户头像 */
.user-avatar {
  width: 1.04vw; /* 20px/1920 = 1.04vw */
  height: 1.04vw; /* 20px/1920 = 1.04vw */
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
  height: 1.25vw; /* 24px/1920 = 1.25vw */
  background: rgba(29, 29, 29, 0.12);
  border-radius: 0.10vw; /* 2px/1920 = 0.10vw */
}

/* 操作按钮组 */
.action-buttons {
  display: flex;
  gap: 0.63vw; /* 12px/1920 = 0.63vw */
}

/* 操作按钮 */
.action-btn {
  position: relative;
  width: 2.08vw; /* 40px/1920 = 2.08vw */
  height: 2.08vw; /* 40px/1920 = 2.08vw */
  border: none;
  background: #F7F8FA;
  border-radius: 1.04vw; /* 20px/1920 = 1.04vw */
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
  top: -0.16vw; /* -3px/1920 = -0.16vw */
  right: -0.05vw; /* -1px/1920 = -0.05vw */
  width: 1.15vw; /* 22px/1920 = 1.15vw */
  height: 0.83vw; /* 16px/1920 = 0.83vw */
  background: #00C27C;
  color: white;
  font-size: 0.52vw; /* 10px/1920 = 0.52vw */
  font-weight: 700;
  border-radius: 4.69vw; /* 90px/1920 = 4.69vw */
  border: 1px solid #FFFFFF;
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: 'Roboto', sans-serif;
  line-height: 1.3;
}

/* 侧边栏状态配合 */
.nav-content--sidebar-collapsed {
  left: 3.33vw; /* 64px/1920 = 3.33vw */
}

.nav-content--mobile {
  left: 0;
}

/* 响应式设计 - 移动端切换回px单位 */
@media (max-width: 1024px) {
  .nav-content {
    left: 0;
    padding: 16px 20px;
    height: 80px; /* 移动端使用固定高度 */
  }
  
  /* 覆盖侧边栏收起状态，移动端和平板端始终从左侧边缘开始 */
  .nav-content--sidebar-collapsed {
    left: 0;
  }
  
  /* 移动端使用固定像素值 */
  .nav-content__left {
    gap: 12px;
  }
  
  .sidebar-toggle-btn {
    width: 32px;
    height: 32px;
    padding: 4px;
    border-radius: 8px;
  }
  
  .breadcrumb {
    font-size: 12px;
    gap: 4px;
  }
  
  .breadcrumb-current.page-title,
  .page-title {
    font-size: 20px;
  }
  
  .nav-content__right {
    gap: 24px;
  }
  
  .user-avatar {
    width: 20px;
    height: 20px;
  }
  
  .divider {
    height: 24px;
  }
  
  .action-buttons {
    gap: 12px;
  }
  
  .action-btn {
    width: 40px;
    height: 40px;
    border-radius: 20px;
  }
  
  .notification-badge {
    width: 22px;
    height: 16px;
    font-size: 10px;
    top: -3px;
    right: -1px;
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
    border-radius: 18px;
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
    border-radius: 16px;
  }
}

/* 桌面端vw单位在基础样式中已设置，移动端通过媒体查询覆盖为px单位 */
</style>