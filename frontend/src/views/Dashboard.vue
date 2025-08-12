<!--
 * @Author: 13895237362 2205451508@qq.com
 * @Date: 2025-08-01 09:32:42
 * @LastEditors: 13895237362 2205451508@qq.com
 * @LastEditTime: 2025-08-12 09:44:47
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
import Layout from '@/components/common/layout/Layout.vue';
import type { NavItem, BreadcrumbItem } from '@/components/common/layout/index';
const router = useRouter();

const navigationItems = ref<NavItem[]>([
  { id: "dashboard", label: "仪表盘", href: "/dashboard", icon: "dashboard", active: true },
  { id: "customers", label: "客户管理", href: "/customers", icon: "customers" },
  { id: "licenses", label: "授权管理", href: "/licenses", icon: "licenses" },
  { id: "roles", label: "角色权限", href: "/roles", icon: "roles" },
  { id: "users", label: "系统用户", href: "/users", icon: "users" }
]);
const breadcrumbs = ref<BreadcrumbItem[]>([
  { title: "仪表盘", path: "/dashboard" }
]);
const handleNavigation = (item: NavItem) => {
  router.push(item.href);
};
</script>

<style scoped></style>