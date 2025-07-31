<template>
  <div class="dashboard">
    <h1>欢迎来到授权管理平台</h1>
    <p>您已成功登录！</p>
    <button @click="logout" class="logout-btn">退出登录</button>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from "vue-router";

const router = useRouter();

function logout() {
  localStorage.removeItem("token");
  router.push("/login");
}
</script>

<style scoped>
.dashboard {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
  background: #f5f5f5;
}

h1 {
  color: #333;
  margin-bottom: 20px;
}

p {
  color: #666;
  margin-bottom: 30px;
}

.logout-btn {
  padding: 10px 20px;
  background: #2E69F7;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
}

.logout-btn:hover {
  background: #1557E5;
}
</style>