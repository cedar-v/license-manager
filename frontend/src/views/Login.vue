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
      <el-select 
        v-model="currentLanguage" 
        @change="changeLanguage"
        size="default"
        popper-class="language-popper"
      >
        <el-option label="English" value="en" />
        <el-option label="中文" value="zh" />
        <el-option label="日本語" value="ja" />
      </el-select>
    </div>

    <!-- 登录区域 -->
    <div class="login-section">
      <div class="login-container">
        <el-card class="login-card" :body-style="{ padding: '2rem' }">
          <h1 class="title">{{ t('login.title') }}</h1>
          <p class="subtitle">{{ t('login.subtitle') }}</p>
          
          <el-form 
            ref="loginFormRef"
            :model="loginForm" 
            :rules="loginRules" 
            @submit.prevent="handleLogin"
            class="login-form"
            size="large"
          >
            <!-- 用户名输入框 -->
            <el-form-item prop="username">
              <el-input
                v-model="loginForm.username"
                :placeholder="t('login.usernamePlaceholder')"
                :prefix-icon="User"
                clearable
                class="login-input"
              />
            </el-form-item>
            
            <!-- 密码输入框 -->
            <el-form-item prop="password">
              <el-input
                v-model="loginForm.password"
                type="password"
                :placeholder="t('login.passwordPlaceholder')"
                :prefix-icon="Lock"
                show-password
                clearable
                class="login-input"
                @keyup.enter="handleLogin"
              />
            </el-form-item>
            
            <!-- 记住密码和忘记密码 -->
            <div class="form-options">
              <el-checkbox v-model="rememberMe" size="default">
                {{ t('login.remember') }}
              </el-checkbox>
              <el-link 
                type="info" 
                @click="handleForgotPassword"
                :underline="false"
                class="forgot-link"
              >
                {{ t('login.forgotPassword') }}
              </el-link>
            </div>
            
            <!-- 登录按钮 -->
            <el-form-item>
              <el-button 
                type="primary" 
                @click="handleLogin"
                :loading="loading"
                class="login-button"
                size="large"
              >
                {{ t('login.submit') }}
              </el-button>
            </el-form-item>
          </el-form>
        </el-card>
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
import { setI18nLanguage, SupportedLocale } from '@/i18n'

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
    
    if (response.code === 200) {
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
      
      ElMessage.success(t('login.success'));
      router.push("/dashboard");
    } else {
      ElMessage.error(response.message || t("login.error.invalid"));
    }
  } catch (error: any) {
    console.error("登录错误:", error);
    ElMessage.error(error.response?.data?.message || t("login.error.general"));
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
  width: clamp(50px, 6vw, 72px);
  height: clamp(48px, 6vw, 70px);
}

.logo-shape-1 {
  position: absolute;
  top: 0;
  left: 0;
  width: 72%;
  height: 100%;
  background: #019C7C;
  clip-path: polygon(0 0, 100% 0, 70% 100%, 0 100%);
}

.logo-shape-2 {
  position: absolute;
  top: 51%;
  right: 0;
  width: 42%;
  height: 49%;
  background: #146B59;
  clip-path: polygon(30% 0, 100% 0, 100% 100%, 0 100%);
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
  border: none;
  border-radius: 20px;
  box-shadow: 0 8px 32px rgba(31, 38, 135, 0.37);
  backdrop-filter: blur(4px);
  background: rgba(255, 255, 255, 0.95);
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
  border: 1px solid #019C7C;
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

.language-switcher :deep(.el-input__inner):focus {
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
  
  .login-card :deep(.el-card__body) {
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
  
  .login-card :deep(.el-card__body) {
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
  
  .login-card :deep(.el-card__body) {
    padding: 1.5rem !important;
  }
}

@media (max-height: 600px) {
  .login-section {
    padding-top: 3rem;
  }
  
  .login-card :deep(.el-card__body) {
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