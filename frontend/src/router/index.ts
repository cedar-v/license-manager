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
    meta: { title: "仪表板", requiresAuth: true }
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

// 路由守卫
router.beforeEach((to, _from, next) => {
  // 设置页面标题
  document.title = typeof(to.meta.title) === "string" ? to.meta.title : "授权管理平台";
  
  // 检查是否需要认证
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