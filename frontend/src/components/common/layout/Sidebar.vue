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
      <div class="sidebar__logo">
        <div class="logo-container">
          <div class="logo-icon">
            <svg width="41" height="40" viewBox="0 0 41 40" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M26.3125 11.4814L22.25 19.5947V22.7148L27.1191 13.0576L29.7393 18.1777L18.7988 40H0L13.5938 22.8037H14L13.8125 23.1201L7.46875 33.541H17.8438V16.9111L7.1875 25.6475L18.0312 10.1406L11.625 14.4463V14.2588L20.4375 0L26.3125 11.4814Z" fill="#019C7C"/>
              <path d="M34.5498 39.9996H28.75L24.5938 32.8864L27.125 27.6246L34.5498 39.9996ZM41 39.9996H36.2705L27.9346 25.941L30.7188 20.1559L41 39.9996Z" fill="#146B59"/>
            </svg>
          </div>
          <span class="logo-text" v-show="!isCollapsed">Cedar-V</span>
        </div>
      </div>
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
  </aside>
</template>

<script setup lang="ts">
import { computed } from 'vue'
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

// 使用全局状态管理
const appStore = useAppStore()

// 从 store 获取折叠状态
const isCollapsed = computed(() => appStore.sidebarCollapsed)

// 处理导航项点击事件
const handleNavClick = (item: NavItem, event: Event) => {
  emit('navClick', item, event)
}
</script>

<style lang="scss" scoped>
// Variables and mixins are auto-injected via Vite configuration
@use 'sass:color';

/* 侧边栏 */
.sidebar {
  position: fixed;
  left: 0;
  top: 0;
  width: 280px;
  height: 100vh;
  background: var(--app-sidebar-bg);
  border-right: 1px solid var(--app-border-light);
  display: flex;
  flex-direction: column;
  transition: all 0.3s ease;
  z-index: 2000;
  
  @include mobile {
    transform: translateX(-100%);
    width: 100vw;
    max-width: 320px;
    
    &--mobile-open {
      transform: translateX(0);
    }
  }
}

.sidebar--collapsed {
  width: 64px;
  
  @include mobile {
    width: 100vw;
    max-width: 320px;
  }
}

/* Header 区域 */
.sidebar__header {
  height: 80px;
  @include flex-between;
  border-bottom: 1px solid var(--app-border-light);
}

.sidebar__logo {
  @include flex-center-vertical;
  flex: 1;
  justify-content: center;
}

.logo-container {
  @include flex-center-vertical;
  gap: $spacing-base;
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
  color: var(--app-text-primary);
  @include text-ellipsis;
  line-height: 1.2;

  @include mobile {
    font-size: 24px;
  }
}

.sidebar__toggle {
  // 所有CSS声明放在@include之前
  width: 36px;
  height: 36px;
  min-width: 36px;
  min-height: 36px;
  padding: 0;
  background: var(--app-content-bg);
  border: 1px solid var(--app-border-light);
  border-radius: 8px;
  color: var(--app-text-secondary);

  box-shadow: var(--app-shadow);
  
  @include button-base;
  
  @include non-touch-device {
    &:hover {
      background: var(--app-bg-color);
      color: var(--el-color-primary);
      border-color: var(--el-color-primary-light-8);
      transform: scale(1.02);
    }
  }
  
  &--collapsed {
    background: var(--el-color-primary);
    color: white;
    border-color: var(--el-color-primary);
    box-shadow: 0 2px 8px var(--el-color-primary-light-7);

    @include non-touch-device {
      &:hover {
        background: var(--el-color-primary-dark-2);
        box-shadow: 0 4px 12px var(--el-color-primary-light-5);
      }
    }
  }
}

.toggle-icon {
  width: 16px;
  height: 16px;
  font-size: $font-size-base;
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

@include non-touch-device {
  .sidebar__toggle:hover .toggle-icon {
    transform: translateX(1px);
  }
  
  .sidebar__toggle--collapsed:hover .toggle-icon {
    transform: translateX(-1px);
  }
}

/* 导航区域 */
.sidebar__nav {
  flex: 1;
  padding: 8px 0;
  overflow-y: auto;
  @include smooth-scroll;
}

.nav-section {
  padding: 8px 16px;
}

.nav-item {
  margin-bottom: 16px;
}

.nav-link {
  @include flex-center-vertical;
  gap: 16px;
  padding: 16px 24px;
  color: var(--app-text-primary);
  text-decoration: none;
  border-radius: 24px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  font-size: 14px;
  font-weight: 400;
  line-height: 1;
  
  @include mobile {
    min-height: 44px; // 触摸目标大小
    padding: 16px;
  }

  @include non-touch-device {
    &:hover {
      background: var(--el-color-primary-light-9);
      color: var(--el-color-primary);

      .nav-icon {
        color: var(--el-color-primary);
      }
    }
  }

  &--active {
    background: var(--el-color-primary-light-8);
    color: var(--el-color-primary);
    font-weight: 500;

    .nav-text {
      color: var(--el-color-primary);
    }

    .nav-icon {
      color: var(--el-color-primary);
    }
  }
}

.nav-icon-wrapper {
  width: 20px;
  height: 20px;
  @include flex-center;
  border-radius: 4px;
  flex-shrink: 0;
}

.nav-icon {
  width: 20px;
  height: 20px;
  font-size: 14px;
  color: var(--app-text-secondary);
}

.nav-text {
  @include text-ellipsis;
}

/* 子菜单 */
.nav-submenu {
  margin-left: 32px;
  margin-top: 4px;
}

.nav-sublink {
  display: block;
  padding: 8px 16px;
  color: var(--app-text-secondary);
  text-decoration: none;
  border-radius: 4px;
  font-size: 12px;
  transition: all 0.2s;
  margin-bottom: 4px;
  
  @include mobile {
    min-height: 40px;
  }
  
  @include non-touch-device {
    &:hover {
      background: var(--el-color-primary-light-9);
      color: var(--el-color-primary);
    }
  }

  &--active {
    background: var(--el-color-primary-light-9);
    color: var(--el-color-primary);
    font-weight: 500;
  }
}

.nav-subtext {
  @include text-ellipsis;
}

/* 底部区域 */
.sidebar__footer {
  padding: 16px;
  border-top: 1px solid var(--app-border-light);
}

/* 收起状态下的样式调整 */
.sidebar--collapsed {
  .sidebar__logo {
    @include flex-center;
  }
  
  .logo-container {
    @include flex-center;
  }
  
  .logo-text {
    display: none;
  }
  
  .nav-link {
    @include flex-center;
    padding: 4px 16px;
  }
  
  .nav-text {
    display: none;
  }
  
  .nav-submenu {
    display: none;
  }
}

/* 滚动条样式 */
.sidebar__nav {
  &::-webkit-scrollbar {
    width: 4px;
  }
  
  &::-webkit-scrollbar-track {
    background: transparent;
  }
  
  &::-webkit-scrollbar-thumb {
    background: rgba(0, 0, 0, 0.1);
    border-radius: 2px;
    
    @include non-touch-device {
      &:hover {
        background: rgba(0, 0, 0, 0.2);
      }
    }
  }
}
</style>