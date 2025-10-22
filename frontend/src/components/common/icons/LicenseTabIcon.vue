<!--
 * @Description: License Tab Icon Component
 * 用于授权详情页面左侧Tab导航的图标组件
 -->
<template>
  <div class="license-tab-icon" :class="iconClass">
    <component :is="iconComponent" />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import BasicInfoIcon from './svg/BasicInfoIcon.vue'
import AuthorizationIcon from './svg/AuthorizationIcon.vue'
import LicenseCertIcon from './svg/LicenseCertIcon.vue'
import HistoryIcon from './svg/HistoryIcon.vue'

interface Props {
  name: string
  active?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  active: false
})

const iconMap = {
  'basic': BasicInfoIcon,
  'authorization': AuthorizationIcon,
  'license': LicenseCertIcon,
  'history': HistoryIcon
}

const iconComponent = computed(() => {
  return iconMap[props.name as keyof typeof iconMap] || BasicInfoIcon
})

const iconClass = computed(() => ({
  'license-tab-icon--active': props.active
}))
</script>

<style scoped>
.license-tab-icon {
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.license-tab-icon :deep(svg) {
  width: 100%;
  height: 100%;
}

.license-tab-icon :deep(path) {
  fill: #B2B8C2;
  transition: fill 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.license-tab-icon--active :deep(path) {
  fill: #019C7C;
}
</style>
