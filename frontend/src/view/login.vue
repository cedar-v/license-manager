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

    <!-- 语言切换器 -->
    <div class="language-switcher">
      <button 
        @click="changeLanguage('en-US')" 
        :class="{ active: currentLanguage === 'en-US' }"
      >
        English
      </button>
      <button 
        @click="changeLanguage('zh-CN')" 
        :class="{ active: currentLanguage === 'zh-CN' }"
      >
        中文
      </button>
    </div>

    <!-- 登录区域 -->
    <div class="login-section">
      <div class="login-container">
        <h1 class="title">{{ t('login.title') }}</h1>
        <p class="subtitle">{{ t('login.subtitle') }}</p>
        
        <form @submit.prevent="handleLogin" class="login-form">
          <!-- 用户名输入框 -->
          <div class="form-group">
            <div class="input-container active">
              <div class="input-icon">
                <svg width="32" height="32" viewBox="0 0 24 24" fill="none">
                  <circle cx="12" cy="7" r="5" stroke="#AABCE6" stroke-width="1.5"/>
                  <path d="M3 20c0-5.523 4.477-10 9-10s9 4.477 9 10" stroke="#AABCE6" stroke-width="1.5"/>
                </svg>
              </div>
              <input
                v-model="username"
                type="text"
                :placeholder="t('login.usernamePlaceholder')"
                required
              />
            </div>
          </div>
          
          <!-- 密码输入框 -->
          <div class="form-group">
            <div class="input-container">
              <div class="input-icon">
                <svg width="32" height="32" viewBox="0 0 24 24" fill="none">
                  <path d="M3.1 11h17.84v9a2 2 0 0 1-2 2H5.1a2 2 0 0 1-2-2v-9Z" fill="#AABCE6"/>
                  <circle cx="12" cy="16" r="1" fill="#FFFFFF"/>
                  <path d="M7 11V7a5 5 0 0 1 10 0v4" stroke="#AABCE6" stroke-width="1.5" fill="none"/>
                </svg>
              </div>
              <input
                v-model="password"
                type="password"
                :placeholder="t('login.passwordPlaceholder')"
                required
              />
            </div>
          </div>
          
          <!-- 记住密码和忘记密码 -->
          <div class="form-options">
            <div class="remember-password">
              <div class="checkbox-container" @click="toggleRemember">
                <div class="checkbox" :class="{ checked: rememberMe }">
                  <svg v-if="rememberMe" width="14" height="14" viewBox="0 0 10 10" fill="none">
                    <path d="M8.5 2.5L3.5 7.5L1.5 5.5" stroke="white" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
                  </svg>
                </div>
                <span class="remember-text">{{ t('login.remember') }}</span>
              </div>
            </div>
            <div class="forgot-password">
              <a href="#" @click.prevent="handleForgotPassword">{{ t('login.forgotPassword') }}</a>
            </div>
          </div>
          
          <!-- 登录按钮 -->
          <button type="submit" class="login-button">
            {{ t('login.submit') }}
          </button>
        </form>
        
        <div v-if="errorMsg" class="error-message">{{ errorMsg }}</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";
import { useI18n } from 'vue-i18n';
import { Login, type LoginRequest } from "@/api/user";
import { setI18nLanguage, SupportedLocale } from '@/i18n'

const router = useRouter();
const { t, locale } = useI18n();

// 当前语言
const currentLanguage = ref(locale.value);

// 登录表单数据
const username = ref("");
const password = ref("");
const rememberMe = ref(false);
const errorMsg = ref("");

// 切换语言
function changeLanguage(lang: string) {
  locale.value = lang;
  currentLanguage.value = lang;
  localStorage.setItem('userLanguage', lang);
  setI18nLanguage(lang as SupportedLocale);
}

onMounted(() => {
  // 记住密码功能
  const saved = localStorage.getItem("loginInfo");
  if (saved) {
    const info = JSON.parse(saved);
    username.value = info.username || "";
    password.value = info.password || "";
    rememberMe.value = true;
  }
  
  // 从本地存储获取用户语言偏好
  const savedLanguage = localStorage.getItem('userLanguage');
  if (savedLanguage) {
    locale.value = savedLanguage;
    currentLanguage.value = savedLanguage;
  }
});

function toggleRemember() {
  rememberMe.value = !rememberMe.value;
}

function handleForgotPassword() {
  alert(t('login.forgotPasswordAlert'));
}

async function handleLogin() {
  if (!username.value || !password.value) {
    errorMsg.value = t("login.error.required");
    return;
  }

  try {
    errorMsg.value = "";
    
    const loginData: LoginRequest = {
      username: username.value,
      password: password.value
    };
    
    const response = await Login(loginData);
    
    if (response.code === 200) {
      // 登录成功
      if (response.data && response.data.token) {
        localStorage.setItem("token", response.data.token);
      }
      
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
      errorMsg.value = response.message || t("login.error.invalid");
    }
  } catch (error: any) {
    console.error("登录错误:", error);
    errorMsg.value = error.response?.data?.message || t("login.error.general");
  }
}
</script>

<style scoped>
/* 基础页面样式 */
.login-page {
  position: relative;
  width: 1920px;
  height: 1080px;
  display: flex;
  overflow: hidden;
}

/* 背景层样式 */
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

/* Logo区域样式 */
.logo-section {
  position: absolute;
  top: 80px;
  left: 100px;
  display: flex;
  align-items: center;
  gap: 20px;
  z-index: 100;
}

.logo-icon {
  position: relative;
  width: 72px;
  height: 70px;
}

.logo-shape-1 {
  position: absolute;
  top: 0;
  left: 0;
  width: 52px;
  height: 70px;
  background: #019C7C;
  clip-path: polygon(0 0, 100% 0, 70% 100%, 0 100%);
}

.logo-shape-2 {
  position: absolute;
  top: 36px;
  right: 0;
  width: 30px;
  height: 34px;
  background: #146B59;
  clip-path: polygon(30% 0, 100% 0, 100% 100%, 0 100%);
}

.logo-text {
  font-family: 'Swis721 BlkCn BT', Arial, sans-serif;
  font-weight: 400;
  font-size: 42px;
  line-height: 1.2;
  color: #333333;
}

/* 登录区域样式 */
.login-section {
  position: absolute;
  top: 0;
  right: 0;
  width: 900px;
  height: 1080px;
  background: rgba(255, 255, 255, 0.6);
  backdrop-filter: blur(50px);
  box-shadow: -8px 0px 40px 0px rgba(63, 139, 255, 0.08);
  display: flex;
  align-items: flex-start;
  justify-content: center;
  z-index: 50;
  padding-top: 280px;
  padding-left: 160px;
  padding-right: 160px;
  box-sizing: border-box;
}

.login-container {
  width: 750px;
  max-width: 750px;
}

.title {
  margin: 0 0 80px 0;
  font-family: 'Source Han Sans CN', sans-serif;
  font-weight: 700;
  font-size: 56px;
  line-height: 1.5;
  color: #333333;
}

.subtitle {
  margin: 0 0 48px 0;
  font-family: 'Source Han Sans CN', sans-serif;
  font-weight: 500;
  font-size: 32px;
  line-height: 1.5;
  color: #333333;
}

/* 输入框样式 */
.input-container {
  position: relative;
  width: 100%;
  height: 80px;
  border: 2px solid #9CA7C7;
  border-radius: 10px;
  background: #FFFFFF;
  display: flex;
  align-items: center;
  padding: 0 20px;
  transition: border-color 0.3s ease;
}

.input-container.active {
  border-color: #2E69F7;
}

.input-icon {
  width: 32px;
  height: 32px;
  margin-right: 16px;
}

input {
  font-size: 24px;
}

/* 登录按钮样式 */
.login-button {
  width: 100%;
  height: 80px;
  font-size: 26px;
}

/* 错误信息样式 */
.error-message {
  margin-top: 30px;
  padding: 20px;
  font-size: 20px;
}
/*新增语言切换器样式 */
.language-switcher {
  position: absolute;
  top: 80px;
  right: 100px;
  display: flex;
  gap: 12px;
  z-index: 200;
}

.language-switcher button {
  padding: 8px 16px;
  border: 1px solid #9CA7C7;
  border-radius: 6px;
  background: rgba(255, 255, 255, 0.8);
  cursor: pointer;
  font-size: 16px;
  font-family: 'Source Han Sans CN', sans-serif;
  transition: all 0.3s ease;
}

.language-switcher button.active {
  background: #2E69F7;
  color: white;
  border-color: #2E69F7;
}

.language-switcher button:hover {
  background: rgba(46, 105, 247, 0.1);
  border-color: #2E69F7;
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