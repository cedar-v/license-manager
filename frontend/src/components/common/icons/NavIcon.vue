<!--
 * @Author: 13895237362 2205451508@qq.com
 * @Date: 2025-08-12 00:00:00
 * @LastEditors: 13895237362 2205451508@qq.com
 * @LastEditTime: 2025-08-12 00:00:00
 * @FilePath: /frontend/src/components/common/icons/NavIcon.vue
 * @Description: 导航栏图标组件
-->
<template>
  <div class="nav-icon" :class="iconClass">
    <component :is="iconComponent" />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import SidebarToggleIcon from './svg/SidebarToggleIcon.vue'
import SearchIcon from './svg/SearchIcon.vue'
import NotificationIcon from './svg/NotificationIcon.vue'
import LanguageIcon from './svg/LanguageIcon.vue'
import DarkModeIcon from './svg/DarkModeIcon.vue'
import UserIcon from './svg/UserIcon.vue'

interface Props {
  name: string
  size?: 'small' | 'medium' | 'large'
}

const props = withDefaults(defineProps<Props>(), {
  size: 'medium'
})

const iconMap = {
  'sidebar-toggle': SidebarToggleIcon,
  'search': SearchIcon,
  'notification': NotificationIcon,
  'language': LanguageIcon,
  'dark-mode': DarkModeIcon,
  'user': UserIcon
}

const iconComponent = computed(() => {
  return iconMap[props.name as keyof typeof iconMap] || SearchIcon
})

const iconClass = computed(() => ({
  [`nav-icon--${props.size}`]: true
}))
</script>

<style scoped>
.nav-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.nav-icon--small {
  width: 16px;
  height: 16px;
}

.nav-icon--medium {
  width: 20px;
  height: 20px;
}

.nav-icon--large {
  width: 24px;
  height: 24px;
}

.nav-icon :deep(svg) {
  width: 100%;
  height: 100%;
}

.nav-icon :deep(path) {
  fill: #1D1D1D;
  transition: fill 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}
</style>