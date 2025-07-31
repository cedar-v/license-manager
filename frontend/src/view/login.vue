<template>
  <div class="login-container">
    <h2>{{ t('login.title') }}</h2>
    <form @submit.prevent="handleLogin">
      <div class="form-group">
        <label for="username">{{ t('login.username') }}</label>
        <input
          id="username"
          v-model="username"
          type="text"
          :placeholder="t('login.usernamePlaceholder')"
          required
        />
      </div>
      <div class="form-group">
        <label for="password">{{ t('login.password') }}</label>
        <input
          id="password"
          v-model="password"
          type="password"
          :placeholder="t('login.passwordPlaceholder')"
          required
        />
      </div>
      <div class="form-group remember-row">
        <input
          id="remember"
          v-model="rememberMe"
          type="checkbox"
        />
        <label for="remember">{{ t('login.remember') }}</label>
      </div>
      <button type="submit">{{ t('login.submit') }}</button>
    </form>
    <p
      v-if="errorMsg"
      class="error"
    >{{ errorMsg }}</p>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useI18n } from "vue-i18n";
import { useRouter } from "vue-router";

const router = useRouter();
const { t } = useI18n();

// 登录表单数据
const username = ref("");
const password = ref("");
const rememberMe = ref(false);
const errorMsg = ref("");

onMounted(() => {
  // 记住密码功能
  const saved = localStorage.getItem("loginInfo");
  if (saved) {
    const info = JSON.parse(saved);
    username.value = info.username || "";
    password.value = info.password || "";
    rememberMe.value = true;
  }
});

function handleLogin() {
  // 基本验证
  if (!username.value || !password.value) {
    errorMsg.value = t("login.error.required");
    return;
  }

  // 模拟登录逻辑
  if (username.value === "admin" && password.value === "123456") {
    errorMsg.value = "";
    if (rememberMe.value) {
      localStorage.setItem(
        "loginInfo",
        JSON.stringify({
          username: username.value,
          password: password.value,
        })
      );
    } else {
      localStorage.removeItem("loginInfo");
    }
    // 这里可以添加登录成功后的跳转逻辑
    router.push("/");
  } else {
    errorMsg.value = t("login.error.invalid");
  }
}
</script>

<style scoped>
.login-container {
  max-width: 350px;
  margin: 60px auto;
  padding: 32px 24px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  display: flex;
  flex-direction: column;
  align-items: center;
  position: relative;
}

h2 {
  margin-bottom: 24px;
  color: #333;
}

.form-group {
  width: 100%;
  margin-bottom: 18px;
  display: flex;
  flex-direction: column;
}

.remember-row {
  flex-direction: row;
  align-items: center;
  margin-bottom: 10px;
}

.remember-row input[type="checkbox"] {
  margin-right: 6px;
}

label {
  margin-bottom: 6px;
  color: #555;
}

input[type="text"],
input[type="password"] {
  width: 100%;
  padding: 8px 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 15px;
}

button[type="submit"] {
  width: 100%;
  padding: 10px;
  background: #42b983;
  color: #fff;
  border: none;
  border-radius: 4px;
  font-size: 16px;
  cursor: pointer;
  margin-top: 10px;
}

button[type="submit"]:hover {
  background: #369870;
}

.error {
  color: #e74c3c;
  margin-top: 12px;
}
</style>