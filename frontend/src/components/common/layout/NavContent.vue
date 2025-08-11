<!--
/**
 * 顶部导航栏组件
 * 提供面包屑导航、页面标题、搜索框、通知、设置和用户信息
 * 支持响应式设计和移动端适配
 */
-->
<template>
  <!-- 顶部导航栏容器 -->
  <header class="nav-content" :class="{ 'nav-content--mobile': isMobile }">
    <!-- 移动端菜单按钮 -->
    <button
      v-if="showMenuButton"
      class="nav-content__menu-btn"
      @click="handleMenuToggle"
      :aria-label="'打开菜单'"
    >
      <i class="el-icon-s-unfold"></i>
    </button>

    <!-- 面包屑导航区域 -->
    <div class="nav-content__breadcrumb">
      <div class="breadcrumb-container">
        <slot name="breadcrumb">
          <nav class="breadcrumb" v-if="breadcrumbItems.length">
            <el-icon><SetUp /></el-icon>
            <span 
              v-for="(item, index) in breadcrumbItems" 
              :key="item.path || index"
              class="breadcrumb-item"
            >
              <a 
                v-if="item.path && index < breadcrumbItems.length - 1"
                :href="item.path"
                class="breadcrumb-link"
                @click="handleBreadcrumbClick(item, $event)"
              >
                {{ item.title }}
              </a>
              <span 
                v-else
                class="breadcrumb-current"
              >
                {{ item.title }}
              </span>
              <i 
                v-if="index < breadcrumbItems.length - 1"
                class="breadcrumb-separator"
              >
                /
              </i>
            </span>
          </nav>
        </slot>
      </div>
    </div>

    <!-- 操作按钮区域 -->
    <div class="nav-content__actions">
      <slot name="actions">
        <div class="actions-container">
          <!-- 搜索框 -->
          <div v-if="showSearch" class="search-box">
            <input 
              v-model="searchQuery"
              type="text"
              :placeholder="searchPlaceholder"
              class="search-input"
              @input="handleSearch"
            />
            <i class="search-icon el-icon-search"></i>
          </div>

          <!-- 通知按钮 -->
          <button 
            v-if="showNotifications"
            class="action-btn notification-btn"
            @click="handleNotificationClick"
            :aria-label="'通知'"
          >
            <i class="el-icon-bell"></i>
            <span v-if="notificationCount" class="notification-badge">
              {{ notificationCount > 99 ? '99+' : notificationCount }}
            </span>
          </button>

          <!-- 设置按钮 -->
          <button 
            v-if="showSettings"
            class="action-btn settings-btn"
            @click="handleSettingsClick"
            :aria-label="'设置'"
          >
            <i class="el-icon-setting"></i>
          </button>

          <!-- 用户头像 -->
          <div v-if="showUserAvatar" class="user-avatar" @click="handleUserClick">
            <img 
              v-if="userAvatar"
              :src="userAvatar"
              :alt="userName || '用户头像'"
              class="avatar-img"
            />
            <div v-else class="avatar-placeholder">
              {{ userInitials }}
            </div>
          </div>
        </div>
      </slot>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'

// 面包屑导航项接口定义
interface BreadcrumbItem {
  title: string // 显示标题
  path?: string // 链接路径（可选，最后一项通常不需要链接）
}

// 组件属性接口定义
interface Props {
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
  showMenuButton?: boolean // 是否显示菜单按钮（移动端使用）
}

// 定义组件属性和默认值
const props = withDefaults(defineProps<Props>(), {
  breadcrumbItems: () => [],
  pageTitle: '',
  showSearch: true,
  searchPlaceholder: '搜索...',
  showNotifications: true,
  notificationCount: 0,
  showSettings: true,
  showUserAvatar: true,
  userName: '',
  showMenuButton: false
})

// 定义组件事件
const emit = defineEmits<{
  search: [query: string]
  breadcrumbClick: [item: BreadcrumbItem, event: Event]
  notificationClick: []
  settingsClick: []
  userClick: []
  menuToggle: []
}>()

// 响应式状态管理
const searchQuery = ref('')
const isMobile = ref(false)

// 计算用户姓名缩写
const userInitials = computed(() => {
  if (!props.userName) return 'U'
  return props.userName
    .split(' ')
    .map(name => name.charAt(0))
    .join('')
    .toUpperCase()
    .slice(0, 2)
})

// 处理搜索事件
const handleSearch = () => {
  emit('search', searchQuery.value)
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

// 处理菜单按钮点击（移动端）
const handleMenuToggle = () => {
  emit('menuToggle')
}

// 检测是否为移动设备
const checkMobile = () => {
  isMobile.value = window.innerWidth <= 1024
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

.nav-content--mobile {
  left: 0;
  padding: 16px 20px;
}

/* 移动端菜单按钮 */
.nav-content__menu-btn {
  display: none;
  width: 40px;
  height: 40px;
  border: none;
  background: transparent;
  border-radius: 6px;
  cursor: pointer;
  color: #666;
  font-size: 18px;
  align-items: center;
  justify-content: center;
  margin-right: 12px;
  transition: all 0.2s;
}

.nav-content__menu-btn:hover {
  background: rgba(1, 156, 124, 0.1);
  color: #019C7C;
}

/* 面包屑区域 */
.nav-content__breadcrumb {
  display: flex;
  align-items: center;
  flex: 1;
  min-width: 0;
}

.breadcrumb-container {
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
}

.breadcrumb {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  color: #888;
  flex-wrap: wrap;
}

.breadcrumb-item {
  display: flex;
  align-items: center;
  gap: 8px;
  white-space: nowrap;
}

.breadcrumb-link {
  color: #019C7C;
  text-decoration: none;
  transition: all 0.2s;
}

.breadcrumb-link:hover {
  color: #017A63;
}

.breadcrumb-current {
  color: #333;
  font-weight: 500;
}

.breadcrumb-separator {
  color: #ccc;
  font-style: normal;
  margin: 0 4px;
}

.page-title {
  font-size: 18px;
  font-weight: 600;
  color: #333;
  margin: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* 操作按钮区域 */
.nav-content__actions {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-shrink: 0;
}

.actions-container {
  display: flex;
  align-items: center;
  gap: 16px;
}

/* 搜索框 */
.search-box {
  position: relative;
  min-width: 200px;
}

.search-input {
  width: 100%;
  height: 36px;
  padding: 0 40px 0 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
  background: #fafafa;
  transition: all 0.2s;
}

.search-input:focus {
  outline: none;
  border-color: #019C7C;
  background: #ffffff;
  box-shadow: 0 0 0 2px rgba(1, 156, 124, 0.1);
}

.search-icon {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  color: #888;
  font-size: 16px;
  pointer-events: none;
}

/* 操作按钮 */
.action-btn {
  position: relative;
  width: 40px;
  height: 40px;
  border: none;
  background: transparent;
  border-radius: 6px;
  cursor: pointer;
  color: #666;
  font-size: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.action-btn:hover {
  background: rgba(1, 156, 124, 0.1);
  color: #019C7C;
}

.notification-badge {
  position: absolute;
  top: 6px;
  right: 6px;
  min-width: 16px;
  height: 16px;
  background: #ff4757;
  color: white;
  font-size: 10px;
  font-weight: 600;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 4px;
  line-height: 1;
}

/* 用户头像 */
.user-avatar {
  width: 40px;
  height: 40px;
  border-radius: 20px;
  overflow: hidden;
  cursor: pointer;
  border: 2px solid transparent;
  transition: all 0.2s;
}

.user-avatar:hover {
  border-color: #019C7C;
}

.avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-placeholder {
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #019C7C, #017A63);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 600;
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .nav-content {
    left: 0;
    padding: 16px 20px;
  }
  
  .nav-content__menu-btn {
    display: flex;
  }
  
  .search-box {
    min-width: 150px;
  }
}

@media (max-width: 768px) {
  .nav-content {
    padding: 12px 16px;
    gap: 12px;
  }
  
  .actions-container {
    gap: 12px;
  }
  
  .search-box {
    min-width: 120px;
  }
  
  .search-input {
    height: 32px;
    font-size: 13px;
  }
  
  .action-btn {
    width: 36px;
    height: 36px;
    font-size: 16px;
  }
  
  .user-avatar {
    width: 36px;
    height: 36px;
    border-radius: 18px;
  }
  
  .breadcrumb {
    font-size: 11px;
  }
  
  .page-title {
    font-size: 16px;
  }
}

@media (max-width: 480px) {
  .nav-content {
    padding: 10px 12px;
  }
  
  .search-box {
    display: none;
  }
  
  .breadcrumb {
    display: none;
  }
  
  .page-title {
    font-size: 14px;
  }
}

/* 与侧边栏收起状态的配合 */
.nav-content--sidebar-collapsed {
  left: 64px;
}
</style>