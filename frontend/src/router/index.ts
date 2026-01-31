/*
 * @Author: 13895237362 2205451508@qq.com
 * @Date: 2025-08-01 09:32:42
 * @LastEditors: 13895237362 2205451508@qq.com
 * @LastEditTime: 2025-10-21 16:59:51
 * @FilePath: /frontend/src/router/index.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { createRouter, createWebHistory, RouteRecordRaw } from "vue-router";

// 路由配置
const routes: Array<RouteRecordRaw> = [
  {
    path: "/",
    redirect: "/login"
  },
  {
    path: "/login",
    component: () => import("@/views/Login.vue"),
    meta: { title: "登录" }
  },
  {
    path: "/dashboard",
    component: () => import("@/views/Dashboard.vue"),
    meta: { title: "仪表盘", requiresAuth: true }
  },
  {
    path: "/customers",
    component: () => import("@/views/Customers/index.vue"),
    meta: { title: "客户管理", requiresAuth: true }
  },
  {
    path: "/licenses",
    name: "licenses",
    component: () => import("@/views/Licenses/index.vue"),
    meta: { title: "授权管理", requiresAuth: true },
    children: [
      {
        path: "",
        name: "licenses-search",
        component: () => import("@/views/Licenses/LicenseSearch.vue"),
        meta: { title: "授权搜索", requiresAuth: true }
      },
      {
        path: "list",
        name: "licenses-list",
        component: () => import("@/views/Licenses/LinenseList.vue"),
        meta: { title: "授权列表", requiresAuth: true }
      },
      {
        path: "create",
        name: "licenses-create",
        component: () => import("@/views/Licenses/LicenseForm.vue"),
        meta: { title: "创建授权", requiresAuth: true }
      },
      {
        path: "view/:id",
        name: "licenses-view",
        component: () => import("@/views/Licenses/LicensesView/index.vue"),
        meta: { title: "查看授权", requiresAuth: true }
      }
    ]
  },
  {
    path: "/enterprise-leads",
    component: () => import("@/views/EnterpriseLeads/index.vue"),
    meta: { title: "企业线索", requiresAuth: true }
  },
  {
    path: "/invoices",
    component: () => import("@/views/Invoices/index.vue"),
    meta: { title: "发票管理", requiresAuth: true }
  },
  {
    path: "/packages",
    component: () => import("@/views/Packages/index.vue"),
    meta: { title: "套餐管理", requiresAuth: true }
  },
  {
    path: "/invoices/detail/:id",
    name: "invoice-detail",
    component: () => import("@/views/Invoices/InvoiceDetail.vue"),
    meta: { title: "发票申请详情", requiresAuth: true }
  },
  {
    path: "/roles",
    component: () => import("@/views/Roles.vue"),
    meta: { title: "角色管理", requiresAuth: true }
  },
  {
    path: "/users",
    component: () => import("@/views/Users.vue"),
    meta: { title: "用户管理", requiresAuth: true }
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

// 路由守卫
router.beforeEach((to, _from, next) => {
  // 设置页面标题
  document.title = typeof (to.meta.title) === "string" ? to.meta.title : "授权管理平台";

  // 检查是否需要认�?
  if (to.meta.requiresAuth) {
    const token = localStorage.getItem("token");
    if (!token) {
      next("/login");
      return;
    }
  }

  next();
});

export default router;

