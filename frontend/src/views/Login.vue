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

<style lang="scss" scoped>
// Variables and mixins are auto-injected via Vite configuration
@use 'sass:color';

/* 基础页面样式 */
.login-page {
  position: relative;
  width: 100vw;
  height: 100vh;
  min-height: 100vh;
  display: flex;
  overflow: hidden;
  
  @include mobile {
    flex-direction: column;
  }
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
  
  @include mobile {
    background-image: url('@/assets/images/login-background.png');
    filter: blur(2px);
    position: fixed;
  }
}

/* Logo区域样式 */
.logo-section {
  position: absolute;
  top: $spacing-extra-large;
  left: $spacing-extra-large;
  @include flex-center-vertical;
  gap: $spacing-base;
  z-index: 100;
  
  @include tablet {
    position: fixed;
    top: $spacing-medium;
    left: $spacing-medium;
    z-index: 300;
  }
  
  @include mobile {
    position: fixed;
    top: $spacing-medium;
    left: $spacing-medium;
    z-index: 300;
  }
  
  @include mobile {
    transform: scale(0.8);
    transform-origin: left top;
  }
}

.logo-icon {
  position: relative;
  width: 54px;
  height: 52px;
  flex-shrink: 0;
}

.logo-text {
  font-family: 'Swis721 BlkCn BT', Arial, sans-serif;
  font-weight: 400;
  font-size: 26px;
  line-height: 1.2;
  color: $text-color-primary;
  @include text-ellipsis;
  
  @include mobile {
    font-size: 22px;
  }
}

/* 登录区域样式 */
.login-section {
  position: absolute;
  top: 0;
  right: 0;
  width: min(900px, 50vw);
  height: 100vh;
  background: rgba($background-color-white, 0.6);
  backdrop-filter: blur(30px);
  box-shadow: -8px 0px 40px 0px rgba(63, 139, 255, 0.08);
  @include flex-center;
  z-index: 50;
  padding: $spacing-large;
  box-sizing: border-box;
  
  @include desktop {
    width: 60vw;
  }
  
  @include tablet {
    width: 100vw;
    right: 0;
    background: rgba($background-color-white, 0.95);
    backdrop-filter: blur(20px);
    position: relative;
    padding: $spacing-medium;
    min-height: 100vh;
  }
  
  @include mobile {
    width: 100vw;
    right: 0;
    background: rgba($background-color-white, 0.95);
    backdrop-filter: blur(20px);
    position: relative;
    padding: $spacing-medium;
    min-height: 100vh;
  }
  
  @include mobile {
    padding: $spacing-medium $spacing-small;
  }
}

.login-container {
  width: 100%;
  max-width: 450px;
  margin: 0 auto;
  
  @include mobile {
    max-width: 90vw;
  }
  
  @include mobile {
    max-width: 100%;
  }
}

.login-card {
  padding: $spacing-large;
  
  @include mobile {
    padding: $spacing-medium;
  }
  
  @include mobile {
    padding: $spacing-base;
  }
}

.title {
  margin: 0 0 $spacing-medium 0;
  font-family: 'Source Han Sans CN', sans-serif;
  font-weight: 700;
  font-size: 30px;
  line-height: 1.2;
  color: $text-color-primary;
  text-align: center;
  
  @include mobile {
    font-size: 24px;
    margin-bottom: $spacing-base;
  }
  
  @include mobile {
    font-size: 20px;
    margin-bottom: $spacing-small;
  }
}

.subtitle {
  margin: 0 0 $spacing-large 0;
  font-family: 'Source Han Sans CN', sans-serif;
  font-weight: $font-weight-primary;
  font-size: $font-size-medium;
  line-height: 1.4;
  color: $text-color-secondary;
  text-align: center;
  
  @include mobile {
    font-size: $font-size-base;
    margin-bottom: $spacing-medium;
  }
  
  @include mobile {
    margin-bottom: $spacing-base;
  }
}

/* 表单样式 */
.login-form {
  width: 100%;
}

.login-input :deep(.el-input__wrapper) {
  @include input-style;
  height: 55px;
  border-radius: 10px;
  background-color: $background-color-white;
  box-shadow: none;
  transition: all 0.3s ease;
  padding: 0 $spacing-base;
  
  @include mobile {
    height: 48px;
    min-height: 44px; // 触摸目标大小
  }
  
  &:hover {
    border-color: $border-color-base;
  }
  
  &.is-focus {
    border-color: $primary-color !important;
    box-shadow: 0 0 0 2px rgba($primary-color, 0.2) !important;
  }
}

.login-input :deep(.el-input__inner) {
  height: 100%;
  font-size: $font-size-base;
  color: $text-color-primary;
  background-color: transparent;
  border: none;
  box-shadow: none;
  
  @include mobile {
    font-size: $font-size-small;
  }
}

.login-input :deep(.el-input__prefix) {
  color: rgba($primary-color, 0.6);
}

.login-button {
  @include button-primary;
  width: 100%;
  height: 60px;
  font-size: $font-size-medium;
  font-weight: 600;
  border-radius: 10px;
  background: linear-gradient(135deg, $primary-color 0%, color.adjust($primary-color, $lightness: -5%) 100%);
  border: none;
  margin-top: $spacing-medium;
  
  @include mobile {
    height: 50px;
    min-height: 44px; // 触摸目标大小
    font-size: $font-size-base;
  }
  
  @include non-touch-device {
    &:hover {
      background: linear-gradient(135deg, color.adjust($primary-color, $lightness: -5%) 0%, $primary-color 100%);
      @include hover-effect;
      box-shadow: 0 8px 25px rgba($primary-color, 0.3);
    }
  }
}

.form-options {
  @include flex-between;
  margin: $spacing-medium 0;
  font-size: $font-size-base;
  
  @include mobile {
    flex-direction: column;
    gap: $spacing-base;
    align-items: flex-start;
  }
}

.forgot-link {
  font-size: $font-size-base;
  color: $text-color-secondary;
  
  @include non-touch-device {
    &:hover {
      color: $primary-color;
    }
  }
  
  @include mobile {
    font-size: $font-size-small;
  }
}

/* 语言切换器样式 */
.language-switcher {
  position: absolute;
  top: $spacing-extra-large;
  right: $spacing-extra-large;
  z-index: 200;
  
  @include tablet {
    position: fixed;
    top: $spacing-medium;
    right: $spacing-medium;
    z-index: 300;
  }
  
  @include mobile {
    position: fixed;
    top: $spacing-medium;
    right: $spacing-medium;
    z-index: 300;
  }
  
  @include mobile {
    transform: scale(0.8);
    transform-origin: right top;
  }
}

.language-switcher :deep(.el-select) {
  width: 150px;
  
  @include mobile {
    width: 100px;
  }
}

.language-switcher :deep(.el-input__inner) {
  background: rgba($background-color-white, 0.9);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(156, 167, 199, 0.3);
  border-radius: $border-radius-base;
  font-size: $font-size-base;
  transition: all 0.3s ease;
  
  @include mobile {
    font-size: $font-size-small;
  }
  
  @include non-touch-device {
    &:hover {
      border-color: $primary-color;
      box-shadow: 0 0 0 2px rgba($primary-color, 0.2);
    }
  }
}

/* 高度适配 */
@media (max-height: 700px) {
  .login-section {
    padding: $spacing-medium;
    align-items: flex-start;
    padding-top: 80px;
  }
}

@media (max-height: 600px) {
  .login-section {
    padding-top: 48px;
  }
  
  .title {
    margin-bottom: $spacing-small;
  }
  
  .subtitle {
    margin-bottom: $spacing-base;
  }
  
  .form-options {
    margin: $spacing-base 0;
  }
}

/* Element Plus 组件样式自定义 */
:global(.language-popper) {
  border-radius: $border-radius-base;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  
  .el-select-dropdown__item {
    font-size: $font-size-base;
    padding: $spacing-small $spacing-medium;
    
    @include non-touch-device {
      &:hover {
        background-color: rgba($primary-color, 0.1);
        color: $primary-color;
      }
    }
    
    &.selected {
      background-color: $primary-color;
      color: white;
    }
  }
}

/* Element Plus checkbox 自定义样式 */
.login-form :deep(.el-checkbox__input.is-checked + .el-checkbox__label) {
  color: $primary-color;
}

.login-form :deep(.el-checkbox__label) {
  display: inline-block;
  font-size: $font-size-base;
  line-height: 1;
  padding-left: $spacing-small;
}
</style>