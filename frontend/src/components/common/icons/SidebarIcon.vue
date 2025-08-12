<!--
 * @Author: 13895237362 2205451508@qq.com
 * @Date: 2025-08-12 00:00:00
 * @LastEditors: 13895237362 2205451508@qq.com
 * @LastEditTime: 2025-08-12 00:00:00
 * @FilePath: /frontend/src/components/common/icons/SidebarIcon.vue
 * @Description: 侧边栏图标组件
-->
<template>
  <div class="sidebar-icon" :class="iconClass">
    <component :is="iconComponent" />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import DashboardIcon from './svg/DashboardIcon.vue'
import CustomersIcon from './svg/CustomersIcon.vue'
import LicensesIcon from './svg/LicensesIcon.vue'
import RolesIcon from './svg/RolesIcon.vue'
import UsersIcon from './svg/UsersIcon.vue'

interface Props {
  name: string
  active?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  active: false
})

const iconMap = {
  'dashboard': DashboardIcon,
  'customers': CustomersIcon,
  'licenses': LicensesIcon,
  'roles': RolesIcon,
  'users': UsersIcon
}

const iconComponent = computed(() => {
  return iconMap[props.name as keyof typeof iconMap] || DashboardIcon
})

const iconClass = computed(() => ({
  'sidebar-icon--active': props.active
}))
</script>

<style scoped>
.sidebar-icon {
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.sidebar-icon :deep(svg) {
  width: 100%;
  height: 100%;
}

.sidebar-icon :deep(path) {
  fill: #B2B8C2;
  transition: fill 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.sidebar-icon--active :deep(path) {
  fill: #019C7C;
}
</style>