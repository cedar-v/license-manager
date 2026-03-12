<template>
  <div class="login-page">
    <!-- 动画背景层 -->
    <div class="animated-background">
      <!-- 动态渐变背景 -->
      <div class="gradient-bg"></div>
      
      <!-- 渐变网格 -->
      <div class="gradient-mesh"></div>
      
      <!-- 连接线画布 -->
      <canvas ref="canvasRef" class="connection-canvas"></canvas>
      
      <!-- 浮动粒子 -->
      <div class="particles-container">
        <div v-for="n in 50" :key="n" class="particle" :style="getParticleStyle(n)"></div>
      </div>
      
      <!-- 大粒子光球 -->
      <div class="glow-orb orb-1"></div>
      <div class="glow-orb orb-2"></div>
      <div class="glow-orb orb-3"></div>
      
      <!-- 波浪效果 -->
      <div class="wave wave-1"></div>
      <div class="wave wave-2"></div>
      <div class="wave wave-3"></div>
      
      <!-- 网格线 -->
      <div class="grid-lines"></div>
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
    <div class="login-wrapper">
      <div class="login-card">
        <!-- Logo -->
        <div class="logo-section">
          <div class="logo-icon">
            <svg
              width="48"
              height="46"
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

    <!-- 底部版权 -->
    <div class="footer">
      <p>© 2025 Cedar-V. All rights reserved.</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'
import { Login, type LoginRequest } from '@/api/user'
import { changeLanguage, type SupportedLocale } from '@/utils/language'
import { useAppStore } from '@/store/modules/app'
import { useUserStore } from '@/store/modules/user'

const router = useRouter()
const appStore = useAppStore()
const userStore = useUserStore()

const { t, locale } = useI18n()
const loginFormRef = ref<FormInstance>()
const canvasRef = ref<HTMLCanvasElement>()

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

// 生成粒子样式
function getParticleStyle(index: number) {
  const size = Math.random() * 8 + 4
  const left = Math.random() * 100
  const delay = Math.random() * 15
  const duration = Math.random() * 8 + 10
  const opacity = Math.random() * 0.5 + 0.3

  return {
    width: `${size}px`,
    height: `${size}px`,
    left: `${left}%`,
    animationDelay: `${delay}s`,
    animationDuration: `${duration}s`,
    opacity: opacity
  }
}

// Canvas 动画
let animationId: number
let particles: Array<{
  x: number
  y: number
  vx: number
  vy: number
  radius: number
}> = []

function initCanvas() {
  const canvas = canvasRef.value
  if (!canvas) return

  const ctx = canvas.getContext('2d')
  if (!ctx) return

  const resizeCanvas = () => {
    canvas.width = window.innerWidth
    canvas.height = window.innerHeight
  }
  resizeCanvas()
  window.addEventListener('resize', resizeCanvas)

  // 初始化粒子
  const particleCount = 25
  particles = []
  for (let i = 0; i < particleCount; i++) {
    particles.push({
      x: Math.random() * canvas.width,
      y: Math.random() * canvas.height,
      vx: (Math.random() - 0.5) * 0.5,
      vy: (Math.random() - 0.5) * 0.5,
      radius: Math.random() * 2 + 1
    })
  }

  const isDark = document.documentElement.getAttribute('data-theme') === 'dark'
  const primaryColor = isDark ? '0, 200, 150' : '1, 156, 124'
  const secondaryColor = isDark ? '0, 150, 200' : '20, 107, 89'

  function animate() {
    if (!ctx || !canvas) return
    
    ctx.clearRect(0, 0, canvas.width, canvas.height)

    // 更新和绘制粒子
    particles.forEach((particle, i) => {
      // 更新位置
      particle.x += particle.vx
      particle.y += particle.vy

      // 边界检测
      if (particle.x < 0 || particle.x > canvas.width) particle.vx *= -1
      if (particle.y < 0 || particle.y > canvas.height) particle.vy *= -1

      // 绘制粒子
      ctx.beginPath()
      ctx.arc(particle.x, particle.y, particle.radius, 0, Math.PI * 2)
      ctx.fillStyle = `rgba(${primaryColor}, 0.6)`
      ctx.fill()

      // 绘制连接线
      for (let j = i + 1; j < particles.length; j++) {
        const dx = particles[j].x - particle.x
        const dy = particles[j].y - particle.y
        const distance = Math.sqrt(dx * dx + dy * dy)

        if (distance < 150) {
          ctx.beginPath()
          ctx.moveTo(particle.x, particle.y)
          ctx.lineTo(particles[j].x, particles[j].y)
          const opacity = (1 - distance / 150) * 0.3
          ctx.strokeStyle = `rgba(${primaryColor}, ${opacity})`
          ctx.lineWidth = 1
          ctx.stroke()
        }
      }
    })

    animationId = requestAnimationFrame(animate)
  }

  animate()

  return () => {
    window.removeEventListener('resize', resizeCanvas)
    cancelAnimationFrame(animationId)
  }
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

  // 初始化 Canvas 动画
  const cleanup = initCanvas()

  onUnmounted(() => {
    cleanup?.()
  })
})

onUnmounted(() => {
  if (animationId) {
    cancelAnimationFrame(animationId)
  }
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
@use 'sass:color';

/* 基础页面样式 */
.login-page {
  position: relative;
  width: 100vw;
  height: 100vh;
  min-height: 100vh;
  overflow: hidden;
  background: linear-gradient(135deg, #f5f7fa 0%, #e4e8ec 100%);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;

  :global([data-theme='dark']) & {
    background: linear-gradient(135deg, #0a0f14 0%, #141925 50%, #0a0e14 100%);
  }
}

/* 动画背景层 */
.animated-background {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;
  z-index: 0;
}

/* 动态渐变背景 */
.gradient-bg {
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: 
    radial-gradient(ellipse at 30% 20%, rgba(1, 156, 124, 0.15) 0%, transparent 40%),
    radial-gradient(ellipse at 70% 80%, rgba(20, 107, 89, 0.12) 0%, transparent 40%),
    radial-gradient(ellipse at 50% 50%, rgba(1, 156, 124, 0.08) 0%, transparent 50%);
  animation: gradientMove 20s ease-in-out infinite;

  :global([data-theme='dark']) & {
    background: 
      radial-gradient(ellipse at 30% 20%, rgba(0, 200, 150, 0.25) 0%, transparent 40%),
      radial-gradient(ellipse at 70% 80%, rgba(0, 180, 200, 0.18) 0%, transparent 40%),
      radial-gradient(ellipse at 50% 50%, rgba(0, 200, 150, 0.12) 0%, transparent 50%);
    animation: gradientMoveDark 15s ease-in-out infinite;
  }
}

@keyframes gradientMove {
  0%, 100% { 
    transform: translate(0, 0) rotate(0deg) scale(1);
  }
  25% { 
    transform: translate(2%, 2%) rotate(2deg) scale(1.05);
  }
  50% { 
    transform: translate(-1%, 3%) rotate(-1deg) scale(1.02);
  }
  75% { 
    transform: translate(3%, -2%) rotate(1deg) scale(1.03);
  }
}

@keyframes gradientMoveDark {
  0%, 100% { 
    transform: translate(0, 0) rotate(0deg) scale(1);
    filter: hue-rotate(0deg);
  }
  25% { 
    transform: translate(3%, 3%) rotate(3deg) scale(1.1);
    filter: hue-rotate(10deg);
  }
  50% { 
    transform: translate(-2%, 4%) rotate(-2deg) scale(1.05);
    filter: hue-rotate(-5deg);
  }
  75% { 
    transform: translate(4%, -3%) rotate(2deg) scale(1.08);
    filter: hue-rotate(5deg);
  }
}

/* 渐变网格 */
.gradient-mesh {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: 
    linear-gradient(90deg, rgba(1, 156, 124, 0.03) 1px, transparent 1px),
    linear-gradient(rgba(1, 156, 124, 0.03) 1px, transparent 1px);
  background-size: 50px 50px;
  animation: meshFloat 10s linear infinite;

  :global([data-theme='dark']) & {
    background: 
      linear-gradient(90deg, rgba(0, 200, 150, 0.05) 1px, transparent 1px),
      linear-gradient(rgba(0, 200, 150, 0.05) 1px, transparent 1px);
    background-size: 40px 40px;
    animation: meshFloatDark 8s linear infinite;
  }
}

@keyframes meshFloat {
  0% { background-position: 0 0, 0 0; }
  100% { background-position: 50px 50px, 50px 50px; }
}

@keyframes meshFloatDark {
  0% { background-position: 0 0, 0 0; }
  100% { background-position: 40px 40px, 40px 40px; }
}

/* Canvas 连接线 */
.connection-canvas {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 1;
}

/* 粒子容器 */
.particles-container {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;
  z-index: 2;
}

/* 浮动粒子 */
.particle {
  position: absolute;
  bottom: -20px;
  background: linear-gradient(135deg, rgba(1, 156, 124, 0.6), rgba(20, 107, 89, 0.3));
  border-radius: 50%;
  animation: floatUp linear infinite;
  filter: blur(0.5px);
  box-shadow: 0 0 10px rgba(1, 156, 124, 0.4);

  :global([data-theme='dark']) & {
    background: linear-gradient(135deg, rgba(0, 255, 200, 0.7), rgba(0, 200, 150, 0.4));
    box-shadow: 
      0 0 15px rgba(0, 255, 200, 0.5),
      0 0 30px rgba(0, 200, 150, 0.3);
    filter: blur(0px);
  }
}

@keyframes floatUp {
  0% {
    transform: translateY(0) scale(1) rotate(0deg);
    opacity: 0;
  }
  5% {
    opacity: 0.8;
  }
  95% {
    opacity: 0.5;
  }
  100% {
    transform: translateY(-120vh) scale(0.3) rotate(360deg);
    opacity: 0;
  }
}

/* 大光球 */
.glow-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(60px);
  animation: orbFloat 20s ease-in-out infinite;
  z-index: 0;

  :global([data-theme='dark']) & {
    filter: blur(40px);
  }
}

.orb-1 {
  width: 400px;
  height: 400px;
  top: 10%;
  left: 10%;
  background: radial-gradient(circle, rgba(1, 156, 124, 0.3) 0%, transparent 70%);
  animation-delay: 0s;

  :global([data-theme='dark']) & {
    width: 500px;
    height: 500px;
    background: radial-gradient(circle, rgba(0, 255, 200, 0.25) 0%, transparent 70%);
    box-shadow: 0 0 100px rgba(0, 255, 200, 0.2);
  }
}

.orb-2 {
  width: 300px;
  height: 300px;
  top: 60%;
  right: 15%;
  background: radial-gradient(circle, rgba(20, 107, 89, 0.25) 0%, transparent 70%);
  animation-delay: -7s;

  :global([data-theme='dark']) & {
    width: 400px;
    height: 400px;
    background: radial-gradient(circle, rgba(0, 200, 255, 0.2) 0%, transparent 70%);
    box-shadow: 0 0 80px rgba(0, 200, 255, 0.15);
  }
}

.orb-3 {
  width: 250px;
  height: 250px;
  bottom: 20%;
  left: 30%;
  background: radial-gradient(circle, rgba(1, 156, 124, 0.2) 0%, transparent 70%);
  animation-delay: -14s;

  :global([data-theme='dark']) & {
    width: 350px;
    height: 350px;
    background: radial-gradient(circle, rgba(100, 255, 200, 0.15) 0%, transparent 70%);
    box-shadow: 0 0 60px rgba(100, 255, 200, 0.1);
  }
}

@keyframes orbFloat {
  0%, 100% { 
    transform: translate(0, 0) scale(1);
  }
  25% { 
    transform: translate(30px, -30px) scale(1.1);
  }
  50% { 
    transform: translate(-20px, 20px) scale(0.95);
  }
  75% { 
    transform: translate(20px, 30px) scale(1.05);
  }
}

/* 波浪效果 - 更明显 */
.wave {
  position: absolute;
  bottom: 0;
  left: 0;
  width: 200%;
  height: 250px;
  background-repeat: repeat-x;
  background-size: 50% 100%;
  opacity: 0.6;

  :global([data-theme='dark']) & {
    opacity: 0.8;
  }
}

.wave-1 {
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 1440 320'%3E%3Cpath fill='%23019C7C' fill-opacity='0.12' d='M0,160L48,176C96,192,192,224,288,213.3C384,203,480,149,576,138.7C672,128,768,160,864,181.3C960,203,1056,213,1152,192C1248,171,1344,117,1392,90.7L1440,64L1440,320L1392,320C1344,320,1248,320,1152,320C1056,320,960,320,864,320C768,320,672,320,576,320C480,320,384,320,288,320C192,320,96,320,48,320L0,320Z'%3E%3C/path%3E%3C/svg%3E");
  animation: waveMove 15s linear infinite;
  height: 280px;

  :global([data-theme='dark']) & {
    background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 1440 320'%3E%3Cpath fill='%2300ffc8' fill-opacity='0.15' d='M0,160L48,176C96,192,192,224,288,213.3C384,203,480,149,576,138.7C672,128,768,160,864,181.3C960,203,1056,213,1152,192C1248,171,1344,117,1392,90.7L1440,64L1440,320L1392,320C1344,320,1248,320,1152,320C1056,320,960,320,864,320C768,320,672,320,576,320C480,320,384,320,288,320C192,320,96,320,48,320L0,320Z'%3E%3C/path%3E%3C/svg%3E");
  }
}

.wave-2 {
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 1440 320'%3E%3Cpath fill='%23146B59' fill-opacity='0.08' d='M0,224L48,213.3C96,203,192,181,288,181.3C384,181,480,203,576,224C672,245,768,267,864,250.7C960,235,1056,181,1152,165.3C1248,149,1344,171,1392,181.3L1440,192L1440,320L1392,320C1344,320,1248,320,1152,320C1056,320,960,320,864,320C768,320,672,320,576,320C480,320,384,320,288,320C192,320,96,320,48,320L0,320Z'%3E%3C/path%3E%3C/svg%3E");
  animation: waveMove 12s linear infinite reverse;
  height: 220px;
  bottom: 20px;
  opacity: 0.4;

  :global([data-theme='dark']) & {
    background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 1440 320'%3E%3Cpath fill='%2300c8ff' fill-opacity='0.1' d='M0,224L48,213.3C96,203,192,181,288,181.3C384,181,480,203,576,224C672,245,768,267,864,250.7C960,235,1056,181,1152,165.3C1248,149,1344,171,1392,181.3L1440,192L1440,320L1392,320C1344,320,1248,320,1152,320C1056,320,960,320,864,320C768,320,672,320,576,320C480,320,384,320,288,320C192,320,96,320,48,320L0,320Z'%3E%3C/path%3E%3C/svg%3E");
    opacity: 0.6;
  }
}

.wave-3 {
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 1440 320'%3E%3Cpath fill='%23019C7C' fill-opacity='0.05' d='M0,96L48,112C96,128,192,160,288,160C384,160,480,128,576,112C672,96,768,96,864,112C960,128,1056,160,1152,160C1248,160,1344,128,1392,112L1440,96L1440,320L1392,320C1344,320,1248,320,1152,320C1056,320,960,320,864,320C768,320,672,320,576,320C480,320,384,320,288,320C192,320,96,320,48,320L0,320Z'%3E%3C/path%3E%3C/svg%3E");
  animation: waveMove 18s linear infinite;
  height: 180px;
  bottom: 40px;
  opacity: 0.3;

  :global([data-theme='dark']) & {
    background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 1440 320'%3E%3Cpath fill='%2364ffc8' fill-opacity='0.08' d='M0,96L48,112C96,128,192,160,288,160C384,160,480,128,576,112C672,96,768,96,864,112C960,128,1056,160,1152,160C1248,160,1344,128,1392,112L1440,96L1440,320L1392,320C1344,320,1248,320,1152,320C1056,320,960,320,864,320C768,320,672,320,576,320C480,320,384,320,288,320C192,320,96,320,48,320L0,320Z'%3E%3C/path%3E%3C/svg%3E");
    opacity: 0.5;
  }
}

@keyframes waveMove {
  0% { transform: translateX(0); }
  100% { transform: translateX(-50%); }
}

/* 网格线 */
.grid-lines {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-image: 
    linear-gradient(rgba(1, 156, 124, 0.05) 1px, transparent 1px),
    linear-gradient(90deg, rgba(1, 156, 124, 0.05) 1px, transparent 1px);
  background-size: 100px 100px;
  animation: gridMove 20s linear infinite;

  :global([data-theme='dark']) & {
    background-image: 
      linear-gradient(rgba(0, 200, 150, 0.08) 1px, transparent 1px),
      linear-gradient(90deg, rgba(0, 200, 150, 0.08) 1px, transparent 1px);
    background-size: 80px 80px;
  }
}

@keyframes gridMove {
  0% { background-position: 0 0; }
  100% { background-position: 100px 100px; }
}

/* 登录包装器 */
.login-wrapper {
  position: relative;
  z-index: 10;
  width: 100%;
  max-width: 480px;
  padding: $spacing-base;
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
}

/* 登录卡片 */
.login-card {
  width: 100%;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(20px);
  border-radius: 24px;
  padding: 48px 40px;
  box-shadow: 
    0 4px 6px -1px rgba(0, 0, 0, 0.05),
    0 10px 15px -3px rgba(0, 0, 0, 0.08),
    0 20px 25px -5px rgba(0, 0, 0, 0.05),
    inset 0 1px 0 rgba(255, 255, 255, 0.6);
  border: 1px solid rgba(255, 255, 255, 0.5);
  animation: cardSlideUp 0.6s ease-out;

  :global([data-theme='dark']) & {
    background: rgba(20, 25, 35, 0.75);
    border: 1px solid rgba(255, 255, 255, 0.1);
    box-shadow: 
      0 4px 6px -1px rgba(0, 0, 0, 0.3),
      0 10px 15px -3px rgba(0, 0, 0, 0.4),
      0 20px 25px -5px rgba(0, 0, 0, 0.3),
      inset 0 1px 0 rgba(255, 255, 255, 0.08),
      0 0 60px rgba(0, 200, 150, 0.1);
  }

  @include mobile {
    padding: 36px 24px;
    border-radius: 20px;
    margin: $spacing-base;
  }
}

@keyframes cardSlideUp {
  0% {
    opacity: 0;
    transform: translateY(30px);
  }
  100% {
    opacity: 1;
    transform: translateY(0);
  }
}

/* Logo区域 */
.logo-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: $spacing-base;
  margin-bottom: 32px;
}

.logo-icon {
  width: 64px;
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, rgba(1, 156, 124, 0.1), rgba(1, 156, 124, 0.05));
  border-radius: 16px;
  transition: all 0.3s ease;

  :global([data-theme='dark']) & {
    background: linear-gradient(135deg, rgba(0, 200, 150, 0.15), rgba(1, 156, 124, 0.08));
    box-shadow: 0 0 20px rgba(0, 200, 150, 0.2);
  }

  svg {
    filter: drop-shadow(0 2px 4px rgba(1, 156, 124, 0.2));
  }
}

.logo-text {
  font-family: 'Swis721 BlkCn BT', Arial, sans-serif;
  font-weight: 600;
  font-size: 24px;
  color: var(--app-text-primary);
  letter-spacing: 0.5px;

  :global([data-theme='dark']) & {
    color: #ffffff;
    text-shadow: 0 0 20px rgba(0, 200, 150, 0.3);
  }
}

/* 标题样式 */
.title {
  margin: 0 0 $spacing-small 0;
  font-weight: 700;
  font-size: 28px;
  line-height: 1.3;
  color: var(--app-text-primary);
  text-align: center;

  :global([data-theme='dark']) & {
    color: #ffffff;
  }

  @include mobile {
    font-size: 24px;
  }
}

.subtitle {
  margin: 0 0 $spacing-large 0;
  font-weight: 400;
  font-size: $font-size-base;
  line-height: 1.5;
  color: var(--app-text-secondary);
  text-align: center;

  :global([data-theme='dark']) & {
    color: #9ca3af;
  }

  @include mobile {
    font-size: $font-size-small;
    margin-bottom: $spacing-medium;
  }
}

/* 表单样式 */
.login-form {
  width: 100%;
}

.login-input :deep(.el-input__wrapper) {
  height: 52px;
  border-radius: 12px;
  background-color: rgba(255, 255, 255, 0.8);
  border: 1px solid rgba(0, 0, 0, 0.08);
  box-shadow: none;
  transition: all 0.3s ease;
  padding: 0 $spacing-base;

  :global([data-theme='dark']) & {
    background-color: rgba(15, 20, 30, 0.6);
    border: 1px solid rgba(255, 255, 255, 0.1);
  }

  @include mobile {
    height: 48px;
  }

  &:hover {
    border-color: rgba(1, 156, 124, 0.4);
    background-color: rgba(255, 255, 255, 0.95);

    :global([data-theme='dark']) & {
      border-color: rgba(0, 200, 150, 0.3);
      background-color: rgba(20, 25, 35, 0.7);
    }
  }

  &.is-focus {
    border-color: $primary-color !important;
    box-shadow: 0 0 0 3px rgba(1, 156, 124, 0.15) !important;
    background-color: #ffffff;

    :global([data-theme='dark']) & {
      border-color: #00c896 !important;
      box-shadow: 0 0 0 3px rgba(0, 200, 150, 0.2) !important;
      background-color: rgba(20, 25, 35, 0.8);
    }
  }
}

.login-input :deep(.el-input__inner) {
  height: 100%;
  font-size: $font-size-base;
  color: var(--app-text-primary);
  background-color: transparent;
  border: none;
  box-shadow: none;

  &::placeholder {
    color: var(--app-text-secondary);
  }

  :global([data-theme='dark']) & {
    color: #ffffff;

    &::placeholder {
      color: #6b7280;
    }
  }
}

.login-input :deep(.el-input__prefix) {
  color: rgba(1, 156, 124, 0.6);
  margin-right: 8px;

  :global([data-theme='dark']) & {
    color: rgba(0, 200, 150, 0.7);
  }
}

/* 登录按钮 */
.login-button {
  width: 100%;
  height: 52px;
  font-size: $font-size-medium;
  font-weight: 600;
  border-radius: 12px;
  background: linear-gradient(
    135deg,
    $primary-color 0%,
    color.adjust($primary-color, $lightness: -8%) 100%
  );
  border: none;
  margin-top: $spacing-base;
  position: relative;
  overflow: hidden;
  transition: all 0.3s ease;

  :global([data-theme='dark']) & {
    background: linear-gradient(
      135deg,
      #00c896 0%,
      #019c7c 100%
    );
    box-shadow: 0 4px 15px rgba(0, 200, 150, 0.3);
  }

  @include mobile {
    height: 48px;
  }

  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: -100%;
    width: 100%;
    height: 100%;
    background: linear-gradient(
      90deg,
      transparent,
      rgba(255, 255, 255, 0.2),
      transparent
    );
    transition: left 0.5s ease;
  }

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 25px rgba(1, 156, 124, 0.35);

    :global([data-theme='dark']) & {
      box-shadow: 0 8px 25px rgba(0, 200, 150, 0.4);
    }

    &::before {
      left: 100%;
    }
  }

  &:active {
    transform: translateY(0);
  }
}

/* 表单选项 */
.form-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
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

  :global([data-theme='dark']) & {
    color: #9ca3af;
  }

  &:hover {
    color: $primary-color;

    :global([data-theme='dark']) & {
      color: #00c896;
    }
  }
}

/* 语言切换器 */
.language-switcher {
  position: fixed;
  top: $spacing-large;
  right: $spacing-large;
  z-index: 100;

  @include mobile {
    top: $spacing-medium;
    right: $spacing-medium;
  }
}

.language-switcher :deep(.el-select) {
  width: 120px;
}

.language-switcher :deep(.el-input__wrapper) {
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(0, 0, 0, 0.08);
  border-radius: 10px;
  box-shadow: none;
  transition: all 0.3s ease;

  :global([data-theme='dark']) & {
    background: rgba(20, 25, 35, 0.7);
    border: 1px solid rgba(255, 255, 255, 0.1);
  }

  &:hover {
    border-color: rgba(1, 156, 124, 0.4);
    box-shadow: 0 2px 8px rgba(1, 156, 124, 0.1);

    :global([data-theme='dark']) & {
      border-color: rgba(0, 200, 150, 0.3);
      box-shadow: 0 2px 8px rgba(0, 200, 150, 0.15);
    }
  }
}

.language-switcher :deep(.el-input__inner) {
  font-size: $font-size-base;
  color: var(--app-text-primary);

  :global([data-theme='dark']) & {
    color: #e5eaf3;
  }
}

/* 底部版权 */
.footer {
  position: fixed;
  bottom: $spacing-medium;
  left: 0;
  right: 0;
  text-align: center;
  z-index: 10;

  p {
    font-size: 12px;
    color: var(--app-text-secondary);
    opacity: 0.6;

    :global([data-theme='dark']) & {
      color: #6b7280;
    }
  }

  @include mobile {
    bottom: $spacing-small;
  }
}

/* 复选框样式 */
.login-form :deep(.el-checkbox__input.is-checked + .el-checkbox__label) {
  color: $primary-color;

  :global([data-theme='dark']) & {
    color: #00c896;
  }
}

.login-form :deep(.el-checkbox__label) {
  font-size: $font-size-base;
  color: var(--app-text-regular);
  padding-left: $spacing-small;

  :global([data-theme='dark']) & {
    color: #9ca3af;
  }
}

.login-form :deep(.el-checkbox__inner) {
  border-radius: 6px;
  border-color: var(--app-border-color);

  :global([data-theme='dark']) & {
    border-color: #4b5563;
    background-color: transparent;
  }
}

.login-form :deep(.el-checkbox__input.is-checked .el-checkbox__inner) {
  background-color: $primary-color;
  border-color: $primary-color;

  :global([data-theme='dark']) & {
    background-color: #00c896;
    border-color: #00c896;
  }
}

/* 语言切换器下拉菜单 */
.language-switcher :deep(.el-select-dropdown) {
  background-color: rgba(255, 255, 255, 0.95) !important;
  backdrop-filter: blur(10px) !important;
  border: 1px solid rgba(0, 0, 0, 0.08) !important;
  border-radius: 10px !important;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1) !important;

  :global([data-theme='dark']) & {
    background: rgba(20, 25, 35, 0.95) !important;
    border: 1px solid rgba(255, 255, 255, 0.1) !important;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3) !important;
  }
}

.language-switcher :deep(.el-select-dropdown__item) {
  color: var(--app-text-primary) !important;
  border-radius: 6px;
  margin: 2px 8px;
  padding: 8px 12px;

  :global([data-theme='dark']) & {
    color: #e5eaf3 !important;
  }

  &:hover {
    background-color: rgba(1, 156, 124, 0.08) !important;

    :global([data-theme='dark']) & {
      background-color: rgba(0, 200, 150, 0.15) !important;
    }
  }

  &.selected {
    background-color: rgba(1, 156, 124, 0.12) !important;
    color: $primary-color !important;
    font-weight: 500;

    :global([data-theme='dark']) & {
      background-color: rgba(0, 200, 150, 0.2) !important;
      color: #00c896 !important;
    }
  }
}

/* 高度适配 */
@media (max-height: 700px) {
  .login-wrapper {
    align-items: flex-start;
    padding-top: 80px;
  }

  .footer {
    position: relative;
    margin-top: auto;
    padding-bottom: $spacing-medium;
  }
}

@media (max-height: 600px) {
  .login-wrapper {
    padding-top: 60px;
  }

  .login-card {
    padding: 32px 28px;
  }

  .logo-section {
    margin-bottom: 24px;
  }
}
</style>
