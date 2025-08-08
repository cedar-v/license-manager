<!--
 * @Author: 13895237362 2205451508@qq.com
 * @Date: 2025-08-01 09:32:42
 * @LastEditors: 13895237362 2205451508@qq.com
 * @LastEditTime: 2025-08-08 14:08:31
 * @FilePath: /frontend/src/views/Dashboard.vue
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
-->
<template>
  <Layout app-name="License Manager" :nav-items="navigationItems" :breadcrumb-items="breadcrumbs" page-title="仪表盘"
    @nav-click="handleNavigation">
    <!-- 页面内容 -->
    <div>
      <h1>欢迎来到授权管理平台</h1>
      <p>您已成功登录！</p>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router';
import Layout from '@/views/components/layout/Layout.vue';
import type { NavItem, BreadcrumbItem } from '@/views/components/layout/index';
const router = useRouter();

const navigationItems = ref<NavItem[]>([
  { id: "dashboard", label: "仪表盘", href: "/dashboard", icon: "Dashboard", active: true },
  { id: "customers", label: "客户管理", href: "/customers", icon: "User" },
  { id: "licenses", label: "授权管理", href: "/licenses", icon: "Key" },
  { id: "roles", label: "角色权限", href: "/roles", icon: "UserFilled" },
  { id: "users", label: "系统用户", href: "/users", icon: "Avatar" }
]);
const breadcrumbs = ref<BreadcrumbItem[]>([
  { title: "仪表盘", path: "/dashboard" }
]);
const handleNavigation = (item: NavItem) => {
  router.push(item.href);
};
</script>

<style scoped></style>