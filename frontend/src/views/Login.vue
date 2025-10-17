<template>
  <div class="login-page">
    <!-- 背景层 -->
    <div class="background-layer" :style="backgroundStyle"></div>

    <!-- Logo -->
    <div class="logo-section">
      <div class="logo-icon">
        <svg
          width="54"
          height="52"
          viewBox="0 0 54 52"
          fill="none"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            d="M26.3125 11.4814L22.25 19.5947V22.7148L27.1191 13.0576L29.7393 18.1777L18.7988 40H0L13.5938 22.8037H14L13.8125 23.1201L7.46875 33.541H17.8438V16.9111L7.1875 25.6475L18.0312 10.1406L11.625 14.4463V14.2588L20.4375 0L26.3125 11.4814Z"
            fill="#019C7C"
          />
          <path
            d="M34.5498 39.9996H28.75L24.5938 32.8864L27.125 27.6246L34.5498 39.9996ZM41 39.9996H36.2705L27.9346 25.941L30.7188 20.1559L41 39.9996Z"
            fill="#146B59"
          />
        </svg>
      </div>
      <span class="logo-text">Cedar-V</span>
    </div>

    <!-- 语言切换器 -->
    <div class="language-switcher">
      <el-select v-model="currentLanguage" @change="handleLanguageChange" size="default">
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
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'
import { Login, type LoginRequest } from '@/api/user'
import { changeLanguage, type SupportedLocale } from '@/utils/language'
import { useAppStore } from '@/store/modules/app'
import { useUserStore } from '@/store/modules/user'
// 导入背景图片
import lightBgImg from '@/assets/images/login-background.png'
import darkBgImg from '@/assets/images/login-background-dark.png'

const router = useRouter()
const appStore = useAppStore()
const userStore = useUserStore()

// 动态背景图样式
const backgroundStyle = computed(() => {
  const isDark = appStore.theme === 'dark'
  const backgroundImage = isDark ? `url(${darkBgImg})` : `url(${lightBgImg})`

  return {
    backgroundImage,
    backgroundSize: 'cover',
    backgroundPosition: 'center',
    backgroundRepeat: 'no-repeat'
  }
})

const { t, locale } = useI18n()
const loginFormRef = ref<FormInstance>()

// 当前语言
const currentLanguage = ref(locale.value)
const loading = ref(false)

// 登录表单数据
const loginForm = reactive({
  username: '',
  password: ''
})

const rememberMe = ref(false)

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
}

// 切换语言
function handleLanguageChange(lang: string) {
  // 使用统一的语言管理器切换语言
  changeLanguage(lang as SupportedLocale)
  currentLanguage.value = lang
}

onMounted(() => {
  // 确保主题正确初始化
  appStore.initTheme()

  // 记住密码功能
  const saved = localStorage.getItem('loginInfo')
  if (saved) {
    const info = JSON.parse(saved)
    loginForm.username = info.username || ''
    loginForm.password = info.password || ''
    rememberMe.value = true
  }

  // 同步当前语言状态
  currentLanguage.value = locale.value
})

function handleForgotPassword() {
  ElMessage.info(t('login.forgotPasswordAlert'))
}

async function handleLogin() {
  if (!loginFormRef.value) return

  try {
    const valid = await loginFormRef.value.validate()
    if (!valid) return

    loading.value = true

    const loginData: LoginRequest = {
      username: loginForm.username,
      password: loginForm.password
    }

    const response = await Login(loginData)
    console.log('登录响应:', response)

    // 检查响应是否成功 (支持 code 000000)
    if (response.code === '000000') {
      // 登录成功，存储用户信息到store
      if (response.data && response.data.token && response.data.user_info) {
        // 使用store的方法存储登录数据
        userStore.setLoginData(response.data.token, {
          username: response.data.user_info.username,
          role: response.data.user_info.role
        })
      }

      if (rememberMe.value) {
        localStorage.setItem(
          'loginInfo',
          JSON.stringify({
            username: loginForm.username,
            password: loginForm.password
          })
        )
      } else {
        localStorage.removeItem('loginInfo')
      }

      ElMessage.success(response.message)
      router.push('/dashboard')
    } else {
      ElMessage.error(response.message)
    }
  } catch (error: any) {
    // 使用后端返回的错误信息
    let errorMessage = error.backendMessage || error.response?.data?.message || error.message
    if (errorMessage) {
      ElMessage.error(errorMessage)
    }
  } finally {
    loading.value = false
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
  background-color: var(--app-bg-color);
  transition: background-color 0.3s ease;

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
  z-index: 1;
  transition: all 0.3s ease;

  @include mobile {
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
    transform: scale(0.8);
    transform-origin: left top;
  }
}

.logo-icon {
  position: relative;
  width: 56px; // 改为56px符合8px栅格
  height: 56px; // 改为56px符合8px栅格
  flex-shrink: 0;
  transition: all 0.3s ease;

  // 暗模式下为Logo添加发光效果
  :global([data-theme='dark']) & {
    filter: drop-shadow(0 0 8px rgba(1, 156, 124, 0.4));

    svg {
      path {
        &:first-child {
          filter: drop-shadow(0 0 2px rgba(1, 156, 124, 0.6));
        }
      }
    }
  }
}

.logo-text {
  font-family: 'Swis721 BlkCn BT', Arial, sans-serif;
  font-weight: 400;
  font-size: 24px; // 改为24px符合8px栅格
  line-height: 1.2;
  color: var(--app-text-primary);
  transition: color 0.3s ease;
  @include text-ellipsis;

  @include mobile {
    font-size: 24px; // 改为24px符合8px栅格
  }
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
  box-shadow: var(--app-shadow);
  @include flex-center;
  z-index: 50;
  padding: $spacing-large;
  box-sizing: border-box;
  transition: all 0.3s ease;

  // 暗模式样式 - 使用30%透明度的黑色背景
  :global([data-theme='dark']) & {
    background: rgba(0, 0, 0, 0.3) !important;
    backdrop-filter: blur(40px) saturate(1.2) !important;
    box-shadow:
      -8px 0px 32px 0px rgba(0, 0, 0, 0.4),
      0 4px 12px 0px rgba(0, 0, 0, 0.2) !important;
    border-left: 1px solid rgba(255, 255, 255, 0.08) !important;
  }

  @include desktop {
    width: 60vw;
  }

  @include tablet {
    width: 100vw;
    right: 0;
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(20px);
    position: relative;
    padding: $spacing-medium;
    min-height: 100vh;

    :global([data-theme='dark']) & {
      background: rgba(0, 0, 0, 0.3) !important;
      backdrop-filter: blur(30px) saturate(1.2) !important;
    }
  }

  @include mobile {
    width: 100vw;
    right: 0;
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(20px);
    position: relative;
    padding: $spacing-medium $spacing-small;
    min-height: 100vh;

    :global([data-theme='dark']) & {
      background: rgba(0, 0, 0, 0.3) !important;
      backdrop-filter: blur(30px) saturate(1.2) !important;
    }
  }
}

.login-container {
  width: 100%;
  max-width: 448px; // 改为448px (56*8)
  margin: 0 auto;

  @include mobile {
    max-width: 100%;
  }
}

.login-card {
  padding: $spacing-large;

  @include mobile {
    padding: $spacing-base;
  }
}

.title {
  margin: 0 0 $spacing-medium 0;
  font-family: 'Source Han Sans CN', sans-serif;
  font-weight: 700;
  font-size: 32px; // 改为32px符合8px栅格
  line-height: 1.2;
  color: var(--app-text-primary);
  text-align: center;
  transition: color 0.3s ease;

  @include mobile {
    font-size: 24px; // 改为24px符合8px栅格
    margin-bottom: $spacing-small;
  }
}

.subtitle {
  margin: 0 0 $spacing-large 0;
  font-family: 'Source Han Sans CN', sans-serif;
  font-weight: $font-weight-primary;
  font-size: $font-size-medium;
  line-height: 1.4;
  color: var(--app-text-secondary);
  text-align: center;
  transition: color 0.3s ease;

  @include mobile {
    font-size: $font-size-base;
    margin-bottom: $spacing-base;
  }
}

/* 表单样式 */
.login-form {
  width: 100%;
}

.login-input :deep(.el-input__wrapper) {
  height: 56px; // 改为56px符合8px栅格
  border-radius: 8px; // 改为8px符合8px栅格
  background-color: var(--app-content-bg);
  border: 1px solid var(--app-border-color);
  box-shadow: none;
  transition: all 0.3s ease;
  padding: 0 $spacing-base;

  @include input-style;

  // 暗模式下的特殊样式
  :global([data-theme='dark']) & {
    background-color: rgba(45, 45, 45, 0.6);
    border: 1px solid rgba(255, 255, 255, 0.1);
    backdrop-filter: blur(8px);

    &:hover {
      border-color: rgba(255, 255, 255, 0.2);
      background-color: rgba(45, 45, 45, 0.7);
    }

    &.is-focus {
      border-color: $primary-color !important;
      background-color: rgba(45, 45, 45, 0.8);
      box-shadow: 0 0 0 2px rgba($primary-color, 0.3) !important;
    }
  }

  @include mobile {
    height: 48px;
    min-height: 44px; // 触摸目标大小
  }

  &:hover {
    border-color: var(--app-border-light);
  }

  &.is-focus {
    border-color: $primary-color !important;
    box-shadow: 0 0 0 2px rgba($primary-color, 0.2) !important;
  }
}

.login-input :deep(.el-input__inner) {
  height: 100%;
  font-size: $font-size-base;
  color: var(--app-text-primary);
  background-color: transparent;
  border: none;
  box-shadow: none;
  transition: color 0.3s ease;

  &::placeholder {
    color: var(--app-text-secondary);
  }

  @include mobile {
    font-size: $font-size-small;
  }
}

.login-input :deep(.el-input__prefix) {
  color: rgba($primary-color, 0.6);
  transition: color 0.3s ease;

  :global([data-theme='dark']) & {
    color: rgba($primary-color, 0.8);
  }
}

.login-button {
  width: 100%;
  height: 64px; // 改为64px符合8px栅格
  font-size: $font-size-medium;
  font-weight: 600;
  border-radius: 8px; // 改为8px符合8px栅格
  background: linear-gradient(
    135deg,
    $primary-color 0%,
    color.adjust($primary-color, $lightness: -5%) 100%
  );
  border: none;
  margin-top: $spacing-medium;
  position: relative;
  overflow: hidden;

  @include button-primary;

  // 暗模式下的按钮增强效果
  :global([data-theme='dark']) & {
    background: linear-gradient(
      135deg,
      $primary-color 0%,
      color.adjust($primary-color, $lightness: -8%) 100%
    );
    box-shadow:
      0 4px 20px rgba($primary-color, 0.4),
      0 2px 8px rgba(0, 0, 0, 0.3);

    &::before {
      content: '';
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background: linear-gradient(135deg, rgba(255, 255, 255, 0.1) 0%, transparent 50%);
      opacity: 0;
      transition: opacity 0.3s ease;
    }

    &:hover::before {
      opacity: 1;
    }
  }

  @include mobile {
    height: 48px; // 改为48px符合8px栅格
    min-height: 44px; // 触摸目标大小
    font-size: $font-size-base;
  }

  @include non-touch-device {
    &:hover {
      background: linear-gradient(
        135deg,
        color.adjust($primary-color, $lightness: -5%) 0%,
        $primary-color 100%
      );
      @include hover-effect;
      box-shadow: 0 8px 25px rgba($primary-color, 0.3);

      :global([data-theme='dark']) & {
        background: linear-gradient(
          135deg,
          color.adjust($primary-color, $lightness: -3%) 0%,
          color.adjust($primary-color, $lightness: -10%) 100%
        );
        box-shadow:
          0 8px 32px rgba($primary-color, 0.5),
          0 4px 16px rgba(0, 0, 0, 0.4);
        transform: translateY(-1px);
      }
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
  color: var(--app-text-secondary);
  transition: color 0.3s ease;

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
    transform: scale(0.8);
    transform-origin: right top;
  }
}

.language-switcher :deep(.el-select) {
  width: 152px; // 改为152px (19*8)

  @include mobile {
    width: 96px; // 改为96px (12*8)
  }
}

.language-switcher :deep(.el-input__inner) {
  background: var(--app-content-bg);
  backdrop-filter: blur(10px);
  border: 1px solid var(--app-border-color);
  border-radius: $border-radius-base;
  font-size: $font-size-base;
  color: var(--app-text-primary);
  transition: all 0.3s ease;

  // 暗模式下的特殊样式
  :global([data-theme='dark']) & {
    background: rgba(45, 45, 45, 0.7);
    border: 1px solid rgba(255, 255, 255, 0.15);
    backdrop-filter: blur(15px);

    &:hover {
      border-color: rgba(255, 255, 255, 0.3);
      background: rgba(45, 45, 45, 0.8);
    }
  }

  @include mobile {
    font-size: $font-size-small;
  }

  @include non-touch-device {
    &:hover {
      border-color: $primary-color;
      box-shadow: 0 0 0 2px rgba($primary-color, 0.2);

      :global([data-theme='dark']) & {
        border-color: $primary-color;
        box-shadow: 0 0 0 2px rgba($primary-color, 0.4);
      }
    }
  }
}

/* 高度适配 */
@media (max-height: 700px) {
  .login-section {
    padding: $spacing-medium;
    align-items: flex-start;
    padding-top: 80px; // 符合8px栅格
  }
}

@media (max-height: 600px) {
  .login-section {
    padding-top: 48px; // 符合8px栅格
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

/* Element Plus checkbox 自定义样式 */
.login-form :deep(.el-checkbox__input.is-checked + .el-checkbox__label) {
  color: $primary-color;
}

.login-form :deep(.el-checkbox__label) {
  display: inline-block;
  font-size: $font-size-base;
  line-height: 1;
  padding-left: $spacing-small;
  color: var(--app-text-regular);
  transition: color 0.3s ease;
}

/* 语言切换器下拉菜单暗模式适配 */
.language-switcher :deep(.el-select-dropdown) {
  background-color: var(--app-content-bg) !important;
  border: 1px solid var(--app-border-color) !important;
  box-shadow: var(--app-shadow) !important;
  transition: all 0.3s ease;

  :global([data-theme='dark']) & {
    background: rgba(45, 45, 45, 0.95) !important;
    backdrop-filter: blur(20px) !important;
    border: 1px solid rgba(255, 255, 255, 0.1) !important;
    box-shadow:
      0 8px 32px rgba(0, 0, 0, 0.4),
      0 4px 16px rgba(0, 0, 0, 0.2) !important;
  }
}

.language-switcher :deep(.el-popper) {
  background-color: var(--app-content-bg) !important;

  :global([data-theme='dark']) & {
    background: rgba(45, 45, 45, 0.95) !important;
  }
}

.language-switcher :deep(.el-select-dropdown__item) {
  color: var(--app-text-primary) !important;
  transition: all 0.2s ease;

  &:hover {
    background-color: var(--app-action-btn-bg) !important;

    :global([data-theme='dark']) & {
      background-color: rgba(1, 156, 124, 0.15) !important;
      color: rgba(1, 156, 124, 1) !important;
    }
  }

  &.selected {
    background-color: rgba(1, 156, 124, 0.1) !important;
    color: $primary-color !important;

    :global([data-theme='dark']) & {
      background-color: rgba(1, 156, 124, 0.2) !important;
      color: rgba(1, 156, 124, 1) !important;
    }
  }

  :global([data-theme='dark']) & {
    color: #e5eaf3 !important;
  }
}
</style>

<style lang="scss">
/* 暗模式登录区域样式 - 完美还原截图设计 */
[data-theme='dark'] .login-section {
  background: rgba(0, 0, 0, 0.3) !important;
  opacity: 0.8;
  backdrop-filter: blur(40px) !important;
  box-shadow: -6px 0px 30px 0px rgba(63, 139, 255, 0.05) !important;
}

/* 暗模式下登录卡片 - 简洁设计 */
[data-theme='dark'] .login-card {
  background: transparent !important;
  backdrop-filter: none !important;
  border-radius: 0 !important;
  padding: 40px !important;
  border: none !important;
  box-shadow: none !important;
}

/* 暗模式下文字样式 - 还原截图 */
[data-theme='dark'] .title {
  color: #ffffff !important;
  text-shadow: none !important;
  font-weight: 400 !important;
}

[data-theme='dark'] .subtitle {
  color: #9ca3af !important;
  text-shadow: none !important;
}

/* 暗模式下输入框样式 - 还原截图 */
[data-theme='dark'] .login-input .el-input__wrapper {
  background-color: transparent !important;
  border: none !important;
  border-bottom: 1px solid #4b5563 !important;
  border-radius: 0 !important;
  backdrop-filter: none !important;
  box-shadow: none !important;
  padding: 8px 0 !important;
}

[data-theme='dark'] .login-input .el-input__wrapper:hover {
  border-bottom-color: #6b7280 !important;
}

[data-theme='dark'] .login-input .el-input__wrapper.is-focus {
  border-bottom-color: #00c896 !important;
  box-shadow: none !important;
}

[data-theme='dark'] .login-input .el-input__inner {
  color: #ffffff !important;
  background: transparent !important;
}

[data-theme='dark'] .login-input .el-input__inner::placeholder {
  color: #6b7280 !important;
}

[data-theme='dark'] .login-input .el-input__prefix {
  color: #6b7280 !important;
}

/* 暗模式下记住密码和忘记密码链接 - 还原截图 */
[data-theme='dark'] .form-options .el-checkbox__label {
  color: #9ca3af !important;
}

[data-theme='dark'] .forgot-link {
  color: #9ca3af !important;
}

[data-theme='dark'] .forgot-link:hover {
  color: #00c896 !important;
}

/* 暗模式下登录按钮 - 还原截图样式 */
[data-theme='dark'] .login-button {
  background: #00c896 !important;
  box-shadow: none !important;
  border: none !important;
  border-radius: 6px !important;
  color: #ffffff !important;
  font-weight: 500 !important;
}

[data-theme='dark'] .login-button:hover {
  background: #00b085 !important;
  box-shadow: none !important;
  transform: none !important;
}

[data-theme='dark'] .login-button:active {
  background: #009973 !important;
}

/* 暗模式下语言选择器 */
[data-theme='dark'] .language-switcher .el-input__inner {
  background: rgba(25, 25, 25, 0.8) !important;
  border: 1px solid rgba(255, 255, 255, 0.15) !important;
  color: #e5eaf3 !important;
  backdrop-filter: blur(15px) !important;
}

[data-theme='dark'] .language-switcher .el-input__inner:hover {
  border-color: rgba(255, 255, 255, 0.3) !important;
  background: rgba(25, 25, 25, 0.9) !important;
}

/* 暗模式下复选框 - 还原截图样式 */
[data-theme='dark'] .el-checkbox__input.is-checked .el-checkbox__inner {
  background-color: #00c896 !important;
  border-color: #00c896 !important;
}

[data-theme='dark'] .el-checkbox__inner {
  background-color: transparent !important;
  border-color: #6b7280 !important;
}

[data-theme='dark'] .el-checkbox__inner:hover {
  border-color: #00c896 !important;
}

/* 响应式暗模式样式 - 保持一致性 */
@media (max-width: 1024px) {
  [data-theme='dark'] .login-section {
    background: #1a2332 !important;
  }
}

@media (max-width: 768px) {
  [data-theme='dark'] .login-section {
    background: #1a2332 !important;
  }

  [data-theme='dark'] .login-card {
    padding: 32px 24px !important;
  }
}
</style>
