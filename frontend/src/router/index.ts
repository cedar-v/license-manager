import { createRouter, createWebHistory, RouteRecordRaw } from "vue-router";
 
// 路由类型:RouteRecordRaw
const routes: Array<RouteRecordRaw> = [
  {
    path: "/",
    component: () => import("@/view/Home.vue"),
    meta: { title: '首页' }
  },
  {
    path: "/login",
    component: () => import("@/view/login.vue"),
    meta: { title: '登录' }
  },
  {
    path: "/register",
    component: () => import("@/view/register.vue"),
    meta: { title: '注册' }
  },
];
 
const router = createRouter({
  // 路由模式
  history: createWebHistory(),
  routes,
});

router.beforeEach(async(to) => {
  document.title = typeof(to.meta.title) === "string" ? to.meta.title : ""
})
 
export default router;