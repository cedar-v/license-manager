<!--
/**
 * 侧边栏组件
 * 提供应用导航菜单，支持折叠/展开功能和响应式设计
 * 包含Logo区域、导航菜单和底部区域
 */
-->
<template>
  <!-- 侧边栏容器 -->
  <aside class="sidebar" :class="{
    'sidebar--collapsed': isCollapsed,
    'sidebar--mobile-open': appStore.isMobile && !isCollapsed
  }">
    <!-- Logo区域 -->
    <div class="sidebar__header">
      <div class="sidebar__logo" v-show="!isCollapsed">
        <div class="logo-container">
          <div class="logo-icon">
            <svg width="41" height="40" viewBox="0 0 41 40" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M26.3125 11.4814L22.25 19.5947V22.7148L27.1191 13.0576L29.7393 18.1777L18.7988 40H0L13.5938 22.8037H14L13.8125 23.1201L7.46875 33.541H17.8438V16.9111L7.1875 25.6475L18.0312 10.1406L11.625 14.4463V14.2588L20.4375 0L26.3125 11.4814Z" fill="#019C7C"/>
              <path d="M34.5498 39.9996H28.75L24.5938 32.8864L27.125 27.6246L34.5498 39.9996ZM41 39.9996H36.2705L27.9346 25.941L30.7188 20.1559L41 39.9996Z" fill="#146B59"/>
            </svg>
          </div>
          <span class="logo-text">Cedar-V</span>
        </div>
      </div>
      <button v-if="collapsible" class="sidebar__toggle" :class="{ 'sidebar__toggle--collapsed': isCollapsed }"
        @click="toggleSidebar" :aria-label="isCollapsed ? t('sidebar.expand') : t('sidebar.collapse')">
        <el-icon class="toggle-icon">
          <component :is="isCollapsed ? 'ArrowRight' : 'ArrowLeft'" />
        </el-icon>
      </button>
    </div>

    <!-- 导航菜单 -->
    <nav class="sidebar__nav">
      <div class="nav-section">
        <slot name="nav-items">
          <div class="nav-item" v-for="item in navItems" :key="item.id">
            <a :href="item.href" class="nav-link" :class="{ 'nav-link--active': item.active }"
              @click="handleNavClick(item, $event)">
              <div class="nav-icon-wrapper">
                <SidebarIcon v-if="item.icon" :name="item.icon" :active="item.active" />
              </div>
              <span v-show="!isCollapsed" class="nav-text">{{ item.label }}</span>
            </a>

            <!-- 子菜单 -->
            <div v-if="item.children && !isCollapsed" class="nav-submenu">
              <a v-for="child in item.children" :key="child.id" :href="child.href" class="nav-sublink"
                :class="{ 'nav-sublink--active': child.active }" @click="handleNavClick(child, $event)">
                <span class="nav-subtext">{{ child.label }}</span>
              </a>
            </div>
          </div>
        </slot>
      </div>
    </nav>

    <!-- 底部区域 -->
    <div class="sidebar__footer">
      <slot name="footer">
        <!-- 用户信息等 -->
      </slot>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/store/modules/app'
import SidebarIcon from '@/components/common/icons/SidebarIcon.vue'

// 导航项接口定义
interface NavItem {
  id: string // 导航项唯一标识
  label: string // 显示文本
  href: string // 链接地址
  icon?: string // 图标类名（可选）
  active?: boolean // 是否为当前活跃项（可选）
  children?: NavItem[] // 子菜单项（可选）
}

// 组件属性接口定义
interface Props {
  appName?: string // 应用名称，显示在Logo区域
  navItems?: NavItem[] // 导航菜单项列表
  collapsible?: boolean // 是否允许折叠侧边栏
  defaultCollapsed?: boolean // 默认是否为折叠状态
}

// 定义组件属性和默认值  
withDefaults(defineProps<Props>(), {
  appName: 'Cedar',
  navItems: () => [],
  collapsible: true,
  defaultCollapsed: false
})

// 定义组件事件
const emit = defineEmits<{
  navClick: [item: NavItem, event: Event]
  toggle: [collapsed: boolean]
}>()

// 使用国际化
const { t } = useI18n()

// 使用全局状态管理
const appStore = useAppStore()

// 从 store 获取折叠状态
const isCollapsed = computed(() => appStore.sidebarCollapsed)

// 切换侧边栏折叠状态
const toggleSidebar = () => {
  appStore.toggleSidebar()
  emit('toggle', !appStore.sidebarCollapsed)
}

// 处理导航项点击事件
const handleNavClick = (item: NavItem, event: Event) => {
  emit('navClick', item, event)
}
</script>

<style scoped>
.sidebar {
  position: fixed;
  left: 0;
  top: 0;
  width: 280px;
  height: 100vh;
  background: linear-gradient(180deg, #F5F7FA 0%, #FFFFFF 100%);
  border-right: 1px solid rgba(29, 29, 29, 0.12);
  display: flex;
  flex-direction: column;
  transition: all 0.3s ease;
  z-index: 1000;
}

.sidebar--collapsed {
  width: 64px;
}

/* Header 区域 */
.sidebar__header {
  height: 80px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
  border-bottom: 1px solid rgba(29, 29, 29, 0.06);
}

.sidebar__logo {
  display: flex;
  align-items: center;
  flex: 1;
}

.logo-container {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo-icon {
  width: 41px;
  height: 40px;
  flex-shrink: 0;
}

.logo-text {
  font-family: 'Swis721 BlkCn BT', sans-serif;
  font-size: 30px;
  font-weight: 400;
  color: #333333;
  white-space: nowrap;
  overflow: hidden;
  line-height: 1.2;
}

.sidebar__toggle {
  width: 36px;
  height: 36px;
  border: none;
  background: #FFFFFF;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #B2B8C2;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0px 1px 3px 0px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(29, 29, 29, 0.08);
}

/* 未选中状态的hover */
.sidebar__toggle:hover {
  background: #F8F9FA;
  color: #019C7C;
  border-color: rgba(1, 156, 124, 0.2);
  transform: scale(1.02);
}

/* 选中状态（折叠时的样式） */
.sidebar__toggle--collapsed {
  background: linear-gradient(135deg, #019C7C 0%, #0BB68A 100%);
  color: #FFFFFF;
  border-color: #019C7C;
  box-shadow: 0px 2px 8px 0px rgba(1, 156, 124, 0.24);
}

/* 选中状态的hover */
.sidebar__toggle--collapsed:hover {
  background: linear-gradient(135deg, #028970 0%, #0AA67D 100%);
  transform: scale(1.02);
  box-shadow: 0px 4px 12px 0px rgba(1, 156, 124, 0.32);
}

.toggle-icon {
  width: 16px;
  height: 16px;
  font-size: 14px;
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.sidebar__toggle:hover .toggle-icon {
  transform: translateX(1px);
}

.sidebar__toggle--collapsed:hover .toggle-icon {
  transform: translateX(-1px);
}

/* 导航区域 */
.sidebar__nav {
  flex: 1;
  padding: 8px 0;
  overflow-y: auto;
}

.nav-section {
  padding: 8px 16px;
}

.nav-item {
  margin-bottom: 10px;
  padding: 8px 16px;
}

.nav-link {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 20px;
  color: #1D1D1D;
  text-decoration: none;
  border-radius: 28px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  font-size: 16px;
  font-family: 'OPPOSans', sans-serif;
  font-weight: 400;
  line-height: 1;
}

.nav-link:hover {
  background: rgba(1, 156, 124, 0.08);
  color: #019C7C;
}

.nav-link--active {
  background: rgba(0, 194, 124, 0.12);
  color: #019C7C;
  font-weight: 700;
  box-shadow: 0px 2px 32px 0px rgba(0, 0, 0, 0.02);
}

.nav-link--active .nav-text {
  color: #019C7C;
}

.nav-icon-wrapper {
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  flex-shrink: 0;
}

.nav-icon {
  width: 20px;
  height: 20px;
  font-size: 16px;
  color: #B2B8C2;
}

.nav-link:hover .nav-icon {
  color: #019C7C;
}

.nav-link--active .nav-icon {
  color: #019C7C;
}

.nav-text {
  white-space: nowrap;
  overflow: hidden;
}

/* 子菜单 */
.nav-submenu {
  margin-left: 32px;
  margin-top: 4px;
}

.nav-sublink {
  display: block;
  padding: 8px 16px;
  color: #888;
  text-decoration: none;
  border-radius: 6px;
  font-size: 13px;
  transition: all 0.2s;
  margin-bottom: 2px;
}

.nav-sublink:hover {
  background: rgba(1, 156, 124, 0.06);
  color: #019C7C;
}

.nav-sublink--active {
  background: rgba(1, 156, 124, 0.08);
  color: #019C7C;
  font-weight: 500;
}

.nav-subtext {
  white-space: nowrap;
  overflow: hidden;
}

/* 底部区域 */
.sidebar__footer {
  padding: 16px;
  border-top: 1px solid rgba(29, 29, 29, 0.06);
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .sidebar {
    transform: translateX(-100%);
  }

  .sidebar--mobile-open {
    transform: translateX(0);
  }

  .sidebar--collapsed {
    width: 280px;
  }
}

@media (max-width: 768px) {
  .sidebar {
    width: 100vw;
    max-width: 320px;
  }

  .sidebar--collapsed {
    width: 100vw;
    max-width: 320px;
  }
}

/* 收起状态下的样式调整 */
.sidebar--collapsed .sidebar__logo {
  justify-content: center;
}

.sidebar--collapsed .logo-container {
  justify-content: center;
}

.sidebar--collapsed .logo-text {
  display: none;
}

.sidebar--collapsed .nav-link {
  justify-content: center;
  padding: 3px 12px;
}

.sidebar--collapsed .nav-text {
  display: none;
}

.sidebar--collapsed .nav-submenu {
  display: none;
}

/* 滚动条样式 */
.sidebar__nav::-webkit-scrollbar {
  width: 4px;
}

.sidebar__nav::-webkit-scrollbar-track {
  background: transparent;
}

.sidebar__nav::-webkit-scrollbar-thumb {
  background: rgba(0, 0, 0, 0.1);
  border-radius: 2px;
}

.sidebar__nav::-webkit-scrollbar-thumb:hover {
  background: rgba(0, 0, 0, 0.2);
}
</style>