<!--
 * @Author: 13895237362 2205451508@qq.com
 * @Date: 2025-08-12 00:00:00
 * @LastEditors: 13895237362 2205451508@qq.com
 * @LastEditTime: 2025-08-12 00:00:00
 * @FilePath: /frontend/src/views/Users.vue
 * @Description: 系统用户页面
-->
<template>
  <Layout app-name="Cedar-V" page-title="系统用户">
    <!-- 页面内容 -->
    <div class="content-container">
      <h1>欢迎你来到系统用户页面</h1>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import Layout from '@/components/common/layout/Layout.vue';
</script>

<style scoped>
.content-container {
  min-height: calc(100vh - 80px);
  padding: 24px;
  width: 100%;
  margin: 0;
  box-sizing: border-box;
}

@media (max-width: 768px) {
  .content-container {
    padding: 16px;
  }
}

@media (max-width: 480px) {
  .content-container {
    padding: 12px;
  }
}

@media (min-width: 1025px) {
  .content-container {
    height: calc(100vh - 4.17vw);
    padding: 1.25vw;
    width: 100%;
    margin: 0;
    box-sizing: border-box;
    display: flex;
    flex-direction: column;
  }
}
</style>