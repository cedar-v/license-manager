<!--
 * @Author: 13895237362 2205451508@qq.com
 * @Date: 2025-08-12 00:00:00
 * @LastEditors: 13895237362 2205451508@qq.com
 * @LastEditTime: 2025-10-21 17:30:00
 * @FilePath: /frontend/src/views/Licenses/index.vue
 * @Description: 授权管理公共路由容器
-->
<template>
  <Layout app-name="Cedar-V" :page-title="t('pages.licenses.title')">
    <router-view />
  </Layout>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import Layout from '@/components/common/layout/Layout.vue'

const { t } = useI18n()
</script>

<style scoped lang="scss">
// 公共路由容器，不需要额外样式
</style>
