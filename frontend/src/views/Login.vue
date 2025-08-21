<template>
  <div class="login-page">
    <!-- 背景层 -->
    <div class="background-layer"></div>

    <!-- Logo -->
    <div class="logo-section">
      <div class="logo-icon">
        <svg width="54" height="52" viewBox="0 0 54 52" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path
            d="M26.3125 11.4814L22.25 19.5947V22.7148L27.1191 13.0576L29.7393 18.1777L18.7988 40H0L13.5938 22.8037H14L13.8125 23.1201L7.46875 33.541H17.8438V16.9111L7.1875 25.6475L18.0312 10.1406L11.625 14.4463V14.2588L20.4375 0L26.3125 11.4814Z"
            fill="#019C7C" />
          <path
            d="M34.5498 39.9996H28.75L24.5938 32.8864L27.125 27.6246L34.5498 39.9996ZM41 39.9996H36.2705L27.9346 25.941L30.7188 20.1559L41 39.9996Z"
            fill="#146B59" />
        </svg>
      </div>
      <span class="logo-text">Cedar-V</span>
    </div>

    <!-- 语言切换器 -->
    <div class="language-switcher">
      <el-select v-model="currentLanguage" @change="handleLanguageChange" size="default" popper-class="language-popper">
        <el-option label="English" value="en" />
        <el-option label="中文" value="zh" />
        <el-option label="日本語" value="ja" />
      </el-select>
    </div>

    <!-- 登录区域 -->
    <div class="login-section">
      <div class="login-container">
        <div class="login-card">
          <h1 class="title">{{ t('login.title') }}</h1>
          <p class="subtitle">{{ t('login.subtitle') }}</p>
          <el-form ref="loginFormRef" :model="loginForm" :rules="loginRules" @submit.prevent="handleLogin"
            class="login-form" size="large">
            <!-- 用户名输入框 -->
            <el-form-item prop="username">
              <el-input v-model="loginForm.username" :placeholder="t('login.usernamePlaceholder')" :prefix-icon="User"
                clearable class="login-input" />
            </el-form-item>

            <!-- 密码输入框 -->
            <el-form-item prop="password">
              <el-input v-model="loginForm.password" type="password" :placeholder="t('login.passwordPlaceholder')"
                :prefix-icon="Lock" show-password clearable class="login-input" @keyup.enter="handleLogin" />
            </el-form-item>

            <!-- 记住密码和忘记密码 -->
            <div class="form-options">
              <el-checkbox v-model="rememberMe" size="default">
                {{ t('login.remember') }}
              </el-checkbox>
              <el-link type="info" @click="handleForgotPassword" :underline="false" class="forgot-link">
                {{ t('login.forgotPassword') }}
              </el-link>
            </div>

            <!-- 登录按钮 -->
            <el-form-item>
              <el-button type="primary" @click="handleLogin" :loading="loading" class="login-button" size="large">
                {{ t('login.submit') }}
              </el-button>
            </el-form-item>
          </el-form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from "vue";
import { useRouter } from "vue-router";
import { useI18n } from 'vue-i18n';
import { ElMessage, type FormInstance, type FormRules } from 'element-plus';
import { User, Lock } from '@element-plus/icons-vue';
import { Login, type LoginRequest } from "@/api/user";
import { changeLanguage, type SupportedLocale } from '@/utils/language'

const router = useRouter();
const { t, locale } = useI18n();
const loginFormRef = ref<FormInstance>();

// 当前语言
const currentLanguage = ref(locale.value);
const loading = ref(false);

// 登录表单数据
const loginForm = reactive({
  username: "",
  password: "",
});

const rememberMe = ref(false);

// 表单验证规则
const loginRules: FormRules = {
  username: [
    { required: true, message: () => t('login.error.usernameRequired'), trigger: 'blur' },
    { min: 3, message: () => t('login.error.usernameMinLength'), trigger: 'blur' }
  ],
  password: [
    { required: true, message: () => t('login.error.passwordRequired'), trigger: 'blur' },
    { min: 6, message: () => t('login.error.passwordMinLength'), trigger: 'blur' }
  ]
};

// 切换语言
function handleLanguageChange(lang: string) {
  // 使用统一的语言管理器切换语言
  changeLanguage(lang as SupportedLocale);
  currentLanguage.value = lang;
}

onMounted(() => {
  // 记住密码功能
  const saved = localStorage.getItem("loginInfo");
  if (saved) {
    const info = JSON.parse(saved);
    loginForm.username = info.username || "";
    loginForm.password = info.password || "";
    rememberMe.value = true;
  }

  // 同步当前语言状态
  currentLanguage.value = locale.value;
});

function handleForgotPassword() {
  ElMessage.info(t('login.forgotPasswordAlert'));
}

async function handleLogin() {
  if (!loginFormRef.value) return;

  try {
    const valid = await loginFormRef.value.validate();
    if (!valid) return;

    loading.value = true;

    const loginData: LoginRequest = {
      username: loginForm.username,
      password: loginForm.password
    };

    const response = await Login(loginData);
    console.log('登录响应:', response);

    // 检查响应是否成功 (支持 code 000000)
    if (response.code === '000000') {
      // 登录成功
      if (response.data && response.data.token) {
        localStorage.setItem("token", response.data.token);
      }

      if (rememberMe.value) {
        localStorage.setItem(
          "loginInfo",
          JSON.stringify({
            username: loginForm.username,
            password: loginForm.password,
          })
        );
      } else {
        localStorage.removeItem("loginInfo");
      }

      ElMessage.success(response.message);
      router.push("/dashboard");
    } else {
      ElMessage.error(response.message);
    }
  } catch (error: any) {
    // 使用后端返回的错误信息
    let errorMessage = error.backendMessage || error.response?.data?.message || error.message;
    if (errorMessage) {
      ElMessage.error(errorMessage);
    }
  } finally {
    loading.value = false;
  }
}
</script>

<style scoped>
/* 基础页面样式 */
.login-page {
  position: relative;
  width: 100vw;
  height: 100vh;
  min-height: 100vh;
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
  top: clamp(2rem, 8vh, 5rem);
  left: clamp(1.5rem, 5vw, 6rem);
  display: flex;
  align-items: center;
  gap: 1.25rem;
  z-index: 100;
}

.logo-icon {
  position: relative;
  width: clamp(54px, 6vw, 54px);
  height: clamp(52px, 6vw, 52px);
}


.logo-text {
  font-family: 'Swis721 BlkCn BT', Arial, sans-serif;
  font-weight: 400;
  font-size: clamp(1.8rem, 4vw, 2.6rem);
  line-height: 1.2;
  color: #333333;
}

/* 登录区域样式 */
.login-section {
  position: absolute;
  top: 0;
  right: 0;
  width: min(900px, 50vw);
  height: 100vh;
  background: rgba(255, 255, 255, 0.6);
  backdrop-filter: blur(30px);
  box-shadow: -8px 0px 40px 0px rgba(63, 139, 255, 0.08);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 50;
  padding: 2rem;
  box-sizing: border-box;
}

.login-container {
  width: 100%;
  max-width: 450px;
  margin: 0 auto;
}

.login-card {
  padding: 2rem;
}

.title {
  margin: 0 0 1.5rem 0;
  font-family: 'Source Han Sans CN', sans-serif;
  font-weight: 700;
  font-size: clamp(2rem, 4vw, 3rem);
  line-height: 1.2;
  color: #333333;
  text-align: center;
}

.subtitle {
  margin: 0 0 2rem 0;
  font-family: 'Source Han Sans CN', sans-serif;
  font-weight: 500;
  font-size: clamp(1rem, 2.5vw, 1.5rem);
  line-height: 1.4;
  color: #666666;
  text-align: center;
}

/* Element Plus 输入框组件自定义样式 */
.login-form {
  width: 100%;
}

.login-input :deep(.el-input__wrapper) {
  height: clamp(45px, 8vh, 55px);
  border-radius: 10px;
  /* border: 1px solid #019C7C; */
  background-color: #ffffff;
  box-shadow: none;
  transition: all 0.3s ease;
  padding: 0 12px;
}

.login-input :deep(.el-input__wrapper):hover {
  border-color: #c0c4cc;
}

.login-input :deep(.el-input__wrapper.is-focus) {
  border-color: #019C7C !important;
  box-shadow: 0 0 0 2px rgba(1, 156, 124, 0.2) !important;
}

.login-input :deep(.el-input__inner) {
  height: 100%;
  font-size: clamp(14px, 2vw, 16px);
  color: #333333;
  background-color: transparent;
  border: none;
  box-shadow: none;
}

.login-input :deep(.el-input__prefix) {
  color: rgba(1, 156, 124, 0.6);
}

.login-button {
  width: 100%;
  height: clamp(50px, 8vh, 60px);
  font-size: clamp(16px, 2.5vw, 18px);
  font-weight: 600;
  border-radius: 10px;
  background: linear-gradient(135deg, #019C7C 0%, #0e8a71 100%);
  border: none;
  margin-top: 1rem;
}

.login-button:hover {
  background: linear-gradient(135deg, #0e8a71 0%, #019C7C 100%);
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(1, 156, 124, 0.3);
}

.form-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin: 1.5rem 0;
  font-size: clamp(13px, 1.8vw, 14px);
}

.forgot-link {
  font-size: clamp(13px, 1.8vw, 14px);
  color: #999999;
}

.forgot-link:hover {
  color: #019C7C;
}

/* 语言切换器样式 */
.language-switcher {
  position: absolute;
  top: clamp(2rem, 8vh, 5rem);
  right: clamp(1.5rem, 5vw, 6rem);
  z-index: 200;
}

.language-switcher :deep(.el-select) {
  width: clamp(120px, 15vw, 150px);
}

.language-switcher :deep(.el-input__inner) {
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(156, 167, 199, 0.3);
  border-radius: 8px;
  font-size: clamp(13px, 1.5vw, 14px);
  transition: all 0.3s ease;
}

.language-switcher :deep(.el-input__inner):hover {
  border-color: #019C7C;
  box-shadow: 0 0 0 2px rgba(1, 156, 124, 0.2);
}

/* 响应式布局 */
@media (max-width: 1200px) {
  .login-section {
    width: 60vw;
  }
}

@media (max-width: 1024px) {
  .login-section {
    width: 100vw;
    right: 0;
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(20px);
    position: relative;
  }

  .login-page {
    flex-direction: column;
  }

  .background-layer {
    filter: blur(2px);
    position: fixed;
  }

  .logo-section {
    position: fixed;
    top: 1.5rem;
    left: 1.5rem;
    z-index: 300;
  }

  .language-switcher {
    position: fixed;
    top: 1.5rem;
    right: 1.5rem;
    z-index: 300;
  }
}

@media (max-width: 768px) {
  .login-section {
    padding: 1rem;
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .login-container {
    max-width: 90vw;
  }

  .login-card {
    padding: 1.5rem !important;
  }

  .form-options {
    flex-direction: column;
    gap: 1rem;
    align-items: flex-start;
  }

  .language-switcher :deep(.el-select) {
    width: 100px;
  }
}

@media (max-width: 480px) {
  .login-section {
    padding: 1rem 0.5rem;
  }

  .login-container {
    max-width: 100%;
  }

  .login-card {
    padding: 1rem !important;
  }

  .logo-section {
    transform: scale(0.8);
    transform-origin: left top;
  }

  .language-switcher {
    transform: scale(0.8);
    transform-origin: right top;
  }
}

/* 桌面端vw适配 - 2K/4K屏幕优化 */
@media (min-width: 1025px) {
  /* Logo区域 - 使用vw单位 */
  .logo-section {
    top: 2.6vw; /* 50px/1920 = 2.6vw */
    left: 3.13vw; /* 60px/1920 = 3.13vw */
    gap: 0.65vw; /* 12.5px/1920 = 0.65vw */
  }
  
  .logo-icon {
    width: 2.81vw; /* 54px/1920 = 2.81vw */
    height: 2.71vw; /* 52px/1920 = 2.71vw */
  }
  
  .logo-text {
    font-size: 1.35vw; /* 26px/1920 = 1.35vw */
  }
  
  /* 登录区域 - 使用vw单位 */
  .login-section {
    width: 46.875vw; /* 900px/1920 = 46.875vw */
    padding: 1.04vw; /* 20px/1920 = 1.04vw */
  }
  
  .login-container {
    max-width: 23.44vw; /* 450px/1920 = 23.44vw */
  }
  
  .login-card {
    padding: 1.04vw; /* 20px/1920 = 1.04vw */
  }
  
  .title {
    font-size: 1.56vw; /* 30px/1920 = 1.56vw */
    margin: 0 0 0.78vw 0; /* 15px/1920 = 0.78vw */
  }
  
  .subtitle {
    font-size: 0.78vw; /* 15px/1920 = 0.78vw */
    margin: 0 0 1.04vw 0; /* 20px/1920 = 1.04vw */
  }
  
  /* 输入框适配 */
  .login-input :deep(.el-input__wrapper) {
    height: 2.86vw; /* 55px/1920 = 2.86vw */
    border-radius: 0.52vw; /* 10px/1920 = 0.52vw */
    padding: 0 0.625vw; /* 12px/1920 = 0.625vw */
  }
  
  .login-input :deep(.el-input__inner) {
    font-size: 0.83vw; /* 16px/1920 = 0.83vw */
  }
  
  /* 登录按钮适配 */
  .login-button {
    height: 3.125vw; /* 60px/1920 = 3.125vw */
    font-size: 0.94vw; /* 18px/1920 = 0.94vw */
    border-radius: 0.52vw; /* 10px/1920 = 0.52vw */
    margin-top: 0.52vw; /* 10px/1920 = 0.52vw */
  }
  
  /* 表单选项 */
  .form-options {
    margin: 0.78vw 0; /* 15px/1920 = 0.78vw */
    font-size: 0.73vw; /* 14px/1920 = 0.73vw */
  }
  
  .forgot-link {
    font-size: 0.73vw; /* 14px/1920 = 0.73vw */
  }
  
  /* 语言切换器 */
  .language-switcher {
    top: 2.6vw; /* 50px/1920 = 2.6vw */
    right: 3.13vw; /* 60px/1920 = 3.13vw */
  }
  
  .language-switcher :deep(.el-select) {
    width: 7.81vw; /* 150px/1920 = 7.81vw */
  }
  
  .language-switcher :deep(.el-input__inner) {
    font-size: 0.73vw; /* 14px/1920 = 0.73vw */
    border-radius: 0.42vw; /* 8px/1920 = 0.42vw */
  }
}

/* 高度适配 */
@media (max-height: 700px) {
  .login-section {
    padding: 1rem;
    align-items: flex-start;
    padding-top: 5rem;
  }

  .title {
    margin-bottom: 1rem;
  }

  .subtitle {
    margin-bottom: 1.5rem;
  }

  .login-card {
    padding: 1.5rem !important;
  }
}

@media (max-height: 600px) {
  .login-section {
    padding-top: 3rem;
  }

  .login-card {
    padding: 1rem !important;
  }

  .title {
    margin-bottom: 0.5rem;
  }

  .subtitle {
    margin-bottom: 1rem;
  }

  .form-options {
    margin: 1rem 0;
  }
}

/* Element Plus 弹窗样式自定义 */
:global(.language-popper) {
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

:global(.language-popper .el-select-dropdown__item) {
  font-size: 14px;
  padding: 8px 16px;
}

:global(.language-popper .el-select-dropdown__item:hover) {
  background-color: rgba(1, 156, 124, 0.1);
  color: #019C7C;
}

:global(.language-popper .el-select-dropdown__item.selected) {
  background-color: #019C7C;
  color: white;
}

/* Element Plus checkbox 自定义样式 */
.login-form :deep(.el-checkbox__input.is-checked + .el-checkbox__label) {
  color: var(--el-color-primary, #019C7C);
}

.login-form :deep(.el-checkbox__label) {
  display: inline-block;
  font-size: var(--el-color-primary, #019C7C);
  line-height: 1;
  padding-left: 8px;
}

/* 确保一屏显示 */
@media (min-height: 600px) {
  .login-section {
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
  }
}
</style>