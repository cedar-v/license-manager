<template>
  <!-- 切换路由 -->
  <router-link to="/login">登录</router-link>
  <router-link to="/register">注册</router-link>
  <!-- 测试路径别名配置 -->
  <img src="@assets/image/1.jpg" />
</template>
   
<script setup lang="ts">
console.log('测试生产环境移出console')
console.log(import.meta.env)
</script>
   
<style scoped></style>