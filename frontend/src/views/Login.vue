<template>
  <div class="login-page">
    <!-- 背景层 -->
    <div class="background-layer"></div>
    
    <!-- Logo -->
    <div class="logo-section">
      <div class="logo-icon">
        <div class="logo-shape-1"></div>
        <div class="logo-shape-2"></div>
      </div>
      <span class="logo-text">Cedar-V</span>
    </div>

    <!-- 登录区域 -->
    <div class="login-section">
      <div class="login-container">
        <h1 class="title">授权管理平台</h1>
        <p class="subtitle">登录您的账户</p>
        
        <form @submit.prevent="handleLogin" class="login-form">
          <!-- 用户名输入框 -->
          <div class="form-group">
            <div class="input-container active">
              <div class="input-icon">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
                  <circle cx="12" cy="7" r="5" stroke="#AABCE6" stroke-width="1.5"/>
                  <path d="M3 20c0-5.523 4.477-10 9-10s9 4.477 9 10" stroke="#AABCE6" stroke-width="1.5"/>
                </svg>
              </div>
              <input
                v-model="username"
                type="text"
                placeholder="kuikuiya 1020666"
                required
              />
            </div>
          </div>
          
          <!-- 密码输入框 -->
          <div class="form-group">
            <div class="input-container">
              <div class="input-icon">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
                  <path d="M3.1 11h17.84v9a2 2 0 0 1-2 2H5.1a2 2 0 0 1-2-2v-9Z" fill="#AABCE6"/>
                  <circle cx="12" cy="16" r="1" fill="#FFFFFF"/>
                  <path d="M7 11V7a5 5 0 0 1 10 0v4" stroke="#AABCE6" stroke-width="1.5" fill="none"/>
                </svg>
              </div>
              <input
                v-model="password"
                type="password"
                placeholder="请输入密码"
                required
              />
            </div>
          </div>
          
          <!-- 记住密码和忘记密码 -->
          <div class="form-options">
            <div class="remember-password">
              <div class="checkbox-container" @click="toggleRemember">
                <div class="checkbox" :class="{ checked: rememberMe }">
                  <svg v-if="rememberMe" width="10" height="10" viewBox="0 0 10 10" fill="none">
                    <path d="M8.5 2.5L3.5 7.5L1.5 5.5" stroke="white" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
                  </svg>
                </div>
                <span class="remember-text">记住密码</span>
              </div>
            </div>
            <div class="forgot-password">
              <a href="#" @click.prevent="handleForgotPassword">忘记密码？</a>
            </div>
          </div>
          
          <!-- 登录按钮 -->
          <button type="submit" class="login-button">
            开始使用
          </button>
        </form>
        
        <div v-if="errorMsg" class="error-message">{{ errorMsg }}</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";

const router = useRouter();

const username = ref("");
const password = ref("");
const rememberMe = ref(false);
const errorMsg = ref("");

function toggleRemember() {
  rememberMe.value = !rememberMe.value;
}

function handleForgotPassword() {
  alert("忘记密码功能暂未实现");
}

function handleLogin() {
  if (!username.value || !password.value) {
    errorMsg.value = "请输入用户名和密码";
    return;
  }

  if (username.value === "admin" && password.value === "123456") {
    errorMsg.value = "";
    
    localStorage.setItem("token", "mock-jwt-token-12345");
    
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
    
    router.push("/dashboard");
  } else {
    errorMsg.value = "用户名或密码错误";
  }
}
</script>

<style scoped>
.login-page {
  position: relative;
  width: 100vw;
  height: 100vh;
  display: flex;
  overflow: hidden;
}

/* 背景层 */
.background-layer {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-image: url('@/assets/images/login-background.png');
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  z-index: 1;
}

/* Logo区域 */
.logo-section {
  position: absolute;
  top: 50px;
  left: 70px;
  display: flex;
  align-items: center;
  gap: 14px;
  z-index: 100;
}

.logo-icon {
  position: relative;
  width: 54px;
  height: 53px;
}

.logo-shape-1 {
  position: absolute;
  top: 0;
  left: 0;
  width: 39px;
  height: 53px;
  background: #019C7C;
  clip-path: polygon(0 0, 100% 0, 70% 100%, 0 100%);
}

.logo-shape-2 {
  position: absolute;
  top: 27px;
  right: 0;
  width: 22px;
  height: 26px;
  background: #146B59;
  clip-path: polygon(30% 0, 100% 0, 100% 100%, 0 100%);
}

.logo-text {
  font-family: 'Swis721 BlkCn BT', Arial, sans-serif;
  font-weight: 400;
  font-size: 32px;
  line-height: 1.2;
  color: #333333;
}

/* 登录区域 */
.login-section {
  position: absolute;
  top: 0;
  right: 0;
  width: 828px;
  height: 100vh;
  background: rgba(255, 255, 255, 0.6);
  backdrop-filter: blur(40px);
  box-shadow: -6px 0px 30px 0px rgba(63, 139, 255, 0.05);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 50;
  padding: 40px 30px;
  box-sizing: border-box;
}

.login-container {
  width: 581px;
  max-width: 581px;
}

.title {
  margin: 0 0 29px 0;
  font-family: 'Source Han Sans CN', sans-serif;
  font-weight: 700;
  font-size: 44px;
  line-height: 1.5;
  color: #333333;
}

.subtitle {
  margin: 0 0 72px 0;
  font-family: 'Source Han Sans CN', sans-serif;
  font-weight: 500;
  font-size: 24px;
  line-height: 1.5;
  color: #333333;
}

.login-form {
  width: 100%;
}

.form-group {
  margin-bottom: 40px;
}

.input-container {
  position: relative;
  width: 100%;
  height: 62px;
  border: 1px solid #9CA7C7;
  border-radius: 8px;
  background: #FFFFFF;
  display: flex;
  align-items: center;
  padding: 0 16px;
  transition: border-color 0.3s ease;
}

.input-container.active {
  border-color: #2E69F7;
}

.input-container:focus-within {
  border-color: #2E69F7;
  box-shadow: 0 0 0 2px rgba(46, 105, 247, 0.1);
}

.input-icon {
  width: 24px;
  height: 24px;
  margin-right: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

input {
  flex: 1;
  border: none;
  outline: none;
  background: transparent;
  font-family: 'Source Han Sans CN', sans-serif;
  font-weight: 500;
  font-size: 18px;
  line-height: 1.5;
  color: #333333;
}

input::placeholder {
  color: #999999;
  font-weight: 400;
}

/* 记住密码和忘记密码 */
.form-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 96px;
  height: 20px;
}

.remember-password {
  display: flex;
  align-items: center;
}

.checkbox-container {
  display: flex;
  align-items: center;
  cursor: pointer;
  user-select: none;
}

.checkbox {
  width: 14px;
  height: 14px;
  border-radius: 2px;
  background: #4B7FF9;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 10px;
  transition: all 0.2s ease;
}

.checkbox:not(.checked) {
  background: transparent;
  border: 1px solid #9CA7C7;
}

.remember-text {
  font-family: 'Alibaba PuHuiTi 2.0', sans-serif;
  font-weight: 400;
  font-size: 14px;
  line-height: 1.4;
  color: #999999;
}

.forgot-password a {
  font-family: 'Alibaba PuHuiTi 2.0', sans-serif;
  font-weight: 400;
  font-size: 14px;
  line-height: 1.4;
  color: #999999;
  text-decoration: none;
  transition: color 0.3s ease;
}

.forgot-password a:hover {
  color: #2E69F7;
}

.login-button {
  width: 100%;
  height: 62px;
  background: #2E69F7;
  border: none;
  border-radius: 8px;
  font-family: 'Source Han Sans CN', sans-serif;
  font-weight: 500;
  font-size: 20px;
  letter-spacing: 0.15em;
  color: #FFFFFF;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

.login-button:hover {
  background: #1557E5;
  transform: translateY(-1px);
  box-shadow: 0 6px 20px rgba(46, 105, 247, 0.3);
}

.login-button:active {
  transform: translateY(0);
}

.error-message {
  margin-top: 20px;
  padding: 15px;
  background: #fee;
  border: 1px solid #f99;
  border-radius: 8px;
  color: #d33;
  text-align: center;
  font-size: 16px;
}

/* 响应式 */
@media (max-width: 1024px) {
  .login-section {
    width: 100vw;
    right: 0;
    background: rgba(255, 255, 255, 0.95);
  }
  
  .background-layer {
    filter: blur(2px);
  }
}

@media (max-width: 768px) {
  .login-section {
    padding: 20px;
  }
  
  .login-container {
    width: 100%;
    max-width: 400px;
  }
  
  .logo-section {
    top: 20px;
    left: 20px;
  }
  
  .logo-text {
    font-size: 24px;
  }
  
  .title {
    font-size: 32px;
    margin-bottom: 20px;
  }
  
  .subtitle {
    font-size: 18px;
    margin-bottom: 40px;
  }
  
  .form-group {
    margin-bottom: 32px;
  }
  
  .form-options {
    margin-bottom: 60px;
  }
  
  .input-container {
    height: 55px;
  }
  
  .login-button {
    height: 55px;
    font-size: 18px;
  }
}

@media (max-width: 480px) {
  .login-section {
    width: 100vw;
    padding: 15px;
  }
  
  .login-container {
    width: 100%;
    max-width: 100%;
  }
  
  .title {
    font-size: 24px;
  }
  
  .input-container {
    height: 50px;
  }
  
  .login-button {
    height: 50px;
    font-size: 16px;
  }
}
</style>