<!--
 * @Author: 13895237362 2205451508@qq.com
 * @Date: 2025-08-12 00:00:00
 * @LastEditors: 13895237362 2205451508@qq.com
 * @LastEditTime: 2025-08-12 00:00:00
 * @FilePath: /frontend/src/views/Roles.vue
 * @Description: 角色权限页面
-->
<template>
  <Layout app-name="License Manager" :nav-items="navigationItems" :breadcrumb-items="breadcrumbs" page-title="角色权限"
    @nav-click="handleNavigation">
    <!-- 页面内容 -->
    <div>
      <h1>欢迎你来到角色权限页面</h1>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router';
import Layout from '@/components/common/layout/Layout.vue';
import type { NavItem, BreadcrumbItem } from '@/components/common/layout/index';
const router = useRouter();

const navigationItems = ref<NavItem[]>([
  { id: "dashboard", label: "仪表盘", href: "/dashboard", icon: "dashboard" },
  { id: "customers", label: "客户管理", href: "/customers", icon: "customers" },
  { id: "licenses", label: "授权管理", href: "/licenses", icon: "licenses" },
  { id: "roles", label: "角色权限", href: "/roles", icon: "roles", active: true },
  { id: "users", label: "系统用户", href: "/users", icon: "users" }
]);
const breadcrumbs = ref<BreadcrumbItem[]>([
  { title: "角色权限", path: "/roles" }
]);
const handleNavigation = (item: NavItem) => {
  router.push(item.href);
};
</script>

<style scoped></style>