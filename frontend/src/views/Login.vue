<template>
  <div class="login-page">
    <!-- 左侧视觉区域 - 黄金比例 61.8% -->
    <div class="visual-section">
      <!-- 动画背景层 -->
      <div class="animated-background">
        <div class="gradient-bg"></div>
        <div class="gradient-mesh"></div>
        <canvas ref="canvasRef" class="connection-canvas"></canvas>
        <div class="particles-container">
          <div v-for="n in 50" :key="n" class="particle" :style="getParticleStyle(n)"></div>
        </div>
        <div class="glow-orb orb-1"></div>
        <div class="glow-orb orb-2"></div>
        <div class="glow-orb orb-3"></div>
        <div class="wave wave-1"></div>
        <div class="wave wave-2"></div>
        <div class="grid-lines"></div>
      </div>

      <!-- 品牌标识 -->
      <div class="brand-header">
        <div class="brand-logo">
          <svg width="32" height="31" viewBox="0 0 54 52" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M26.3125 11.4814L22.25 19.5947V22.7148L27.1191 13.0576L29.7393 18.1777L18.7988 40H0L13.5938 22.8037H14L13.8125 23.1201L7.46875 33.541H17.8438V16.9111L7.1875 25.6475L18.0312 10.1406L11.625 14.4463V14.2588L20.4375 0L26.3125 11.4814Z" fill="#ffffff"/>
            <path d="M34.5498 39.9996H28.75L24.5938 32.8864L27.125 27.6246L34.5498 39.9996ZM41 39.9996H36.2705L27.9346 25.941L30.7188 20.1559L41 39.9996Z" fill="rgba(255,255,255,0.7)"/>
          </svg>
        </div>
        <span class="brand-name">Cedar-V</span>
      </div>

      <!-- 左侧主文案 -->
      <div class="hero-content">
        <h1 class="hero-title">{{ t('login.heroTitleLine1') }}{{ t('login.heroTitleLine2') }}</h1>
        <p class="hero-subtitle">{{ t('login.heroSubtitle') }}</p>
        <div class="hero-stats">
          <div class="stat-item">
            <span class="stat-label">{{ t('login.stats.totalAuthCodes') }}</span>
            <span class="stat-number">10K+</span>
            <span class="stat-sublabel">{{ t('login.stats.totalAuthCodesSubLabel') }}</span>
          </div>
          <div class="stat-divider"></div>
          <div class="stat-item">
            <span class="stat-label">{{ t('login.stats.activeLicenses') }}</span>
            <span class="stat-number">8,234</span>
            <span class="stat-sublabel">{{ t('login.stats.activeLicensesSubLabel') }}</span>
          </div>
          <div class="stat-divider"></div>
          <div class="stat-item">
            <span class="stat-label">{{ t('login.stats.systemUptime') }}</span>
            <span class="stat-number">99.9%</span>
            <span class="stat-sublabel">{{ t('login.stats.systemUptimeSubLabel') }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 右侧登录区域 - 黄金比例 38.2% -->
    <div class="login-section">
      <!-- 语言切换器 -->
      <div class="language-switcher">
        <el-select v-model="currentLanguage" @change="handleLanguageChange" size="default">
          <el-option label="English" value="en" />
          <el-option label="中文" value="zh" />
          <el-option label="日本語" value="ja" />
        </el-select>
      </div>

      <div class="login-container">
        <div class="login-card">
          <div class="login-header">
            <h2 class="title">{{ t('login.welcomeTitle') }}</h2>
            <p class="subtitle">{{ t('login.welcomeSubtitle') }}</p>
          </div>

          <el-form
            ref="loginFormRef"
            :model="loginForm"
            :rules="loginRules"
            @submit.prevent="handleLogin"
            class="login-form"
            size="large"
          >
            <el-form-item prop="username">
              <el-input
                v-model="loginForm.username"
                :placeholder="t('login.usernamePlaceholder')"
                :prefix-icon="User"
                clearable
                class="login-input"
              />
            </el-form-item>

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

          <div class="login-footer">
            <p>{{ t('login.copyright') }}</p>
          </div>
        </div>
      </div>
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

const currentLanguage = ref(locale.value)
const loading = ref(false)

const loginForm = reactive({
  username: '',
  password: ''
})

const rememberMe = ref(false)

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

function getParticleStyle(_index: number) {
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

let animationId: number
let particles: Array<{ x: number; y: number; vx: number; vy: number; radius: number }> = []

function initCanvas() {
  const canvas = canvasRef.value
  if (!canvas) return

  const ctx = canvas.getContext('2d')
  if (!ctx) return

  const resizeCanvas = () => {
    const parent = canvas.parentElement
    if (parent) {
      canvas.width = parent.offsetWidth
      canvas.height = parent.offsetHeight
    }
  }
  resizeCanvas()
  window.addEventListener('resize', resizeCanvas)

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

  function animate() {
    if (!ctx || !canvas) return
    ctx.clearRect(0, 0, canvas.width, canvas.height)

    particles.forEach((particle, i) => {
      particle.x += particle.vx
      particle.y += particle.vy
      if (particle.x < 0 || particle.x > canvas.width) particle.vx *= -1
      if (particle.y < 0 || particle.y > canvas.height) particle.vy *= -1

      ctx.beginPath()
      ctx.arc(particle.x, particle.y, particle.radius, 0, Math.PI * 2)
      ctx.fillStyle = `rgba(${primaryColor}, 0.6)`
      ctx.fill()

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

function handleLanguageChange(lang: string) {
  changeLanguage(lang as SupportedLocale)
  currentLanguage.value = lang
}

onMounted(() => {
  appStore.initTheme()

  const saved = localStorage.getItem('loginInfo')
  if (saved) {
    const info = JSON.parse(saved)
    loginForm.username = info.username || ''
    loginForm.password = info.password || ''
    rememberMe.value = true
  }

  currentLanguage.value = locale.value

  const cleanup = initCanvas()
  onUnmounted(() => {
    cleanup?.()
  })
})

onUnmounted(() => {
  if (animationId) cancelAnimationFrame(animationId)
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
    if (response.code === '000000') {
      if (response.data?.token && response.data?.user_info) {
        userStore.setLoginData(response.data.token, {
          username: response.data.user_info.username,
          role: response.data.user_info.role
        })
      }

      if (rememberMe.value) {
        localStorage.setItem('loginInfo', JSON.stringify({
          username: loginForm.username,
          password: loginForm.password
        }))
      } else {
        localStorage.removeItem('loginInfo')
      }

      ElMessage.success(response.message)
      router.push('/dashboard')
    } else {
      ElMessage.error(response.message)
    }
  } catch (error: any) {
    const errorMessage = error.backendMessage || error.response?.data?.message || error.message
    if (errorMessage) ElMessage.error(errorMessage)
  } finally {
    loading.value = false
  }
}
</script>

<style lang="scss" scoped>
@use 'sass:color';

/* 黄金比例: 61.8% / 38.2% */
$phi: 61.8%;
$phi-inv: 38.2%;

.login-page {
  position: relative;
  width: 100vw;
  height: 100vh;
  min-height: 100vh;
  overflow: hidden;
  display: flex;

  @include mobile {
    flex-direction: column;
  }
}

/* 左侧视觉区域 - 61.8% */
.visual-section {
  position: relative;
  width: $phi;
  height: 100%;
  overflow: hidden;
  background: linear-gradient(145deg, #0a3d2e 0%, #0d5c48 40%, #0f766e 100%);

  :global([data-theme='dark']) & {
    background: linear-gradient(145deg, #052e22 0%, #064335 40%, #065a52 100%);
  }

  @media (max-width: 1024px) {
    display: none;
  }
}

/* 动画背景层 */
.animated-background {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;
}

.gradient-bg {
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: 
    radial-gradient(ellipse at 25% 25%, rgba(255, 255, 255, 0.12) 0%, transparent 45%),
    radial-gradient(ellipse at 75% 75%, rgba(255, 255, 255, 0.08) 0%, transparent 40%),
    radial-gradient(ellipse at 50% 50%, rgba(255, 255, 255, 0.05) 0%, transparent 50%);
  animation: gradientMove 25s ease-in-out infinite;

  :global([data-theme='dark']) & {
    background: 
      radial-gradient(ellipse at 25% 25%, rgba(0, 255, 200, 0.12) 0%, transparent 45%),
      radial-gradient(ellipse at 75% 75%, rgba(0, 200, 255, 0.08) 0%, transparent 40%),
      radial-gradient(ellipse at 50% 50%, rgba(0, 255, 200, 0.05) 0%, transparent 50%);
    animation: gradientMoveDark 18s ease-in-out infinite;
  }
}

@keyframes gradientMove {
  0%, 100% { transform: translate(0, 0) rotate(0deg) scale(1); }
  33% { transform: translate(3%, -2%) rotate(2deg) scale(1.05); }
  66% { transform: translate(-2%, 3%) rotate(-1deg) scale(1.02); }
}

@keyframes gradientMoveDark {
  0%, 100% { transform: translate(0, 0) rotate(0deg) scale(1); filter: hue-rotate(0deg); }
  33% { transform: translate(4%, -3%) rotate(3deg) scale(1.08); filter: hue-rotate(8deg); }
  66% { transform: translate(-3%, 4%) rotate(-2deg) scale(1.04); filter: hue-rotate(-5deg); }
}

.gradient-mesh {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: 
    linear-gradient(90deg, rgba(255, 255, 255, 0.04) 1px, transparent 1px),
    linear-gradient(rgba(255, 255, 255, 0.04) 1px, transparent 1px);
  background-size: 60px 60px;
  animation: meshFloat 12s linear infinite;

  :global([data-theme='dark']) & {
    background: 
      linear-gradient(90deg, rgba(0, 255, 200, 0.03) 1px, transparent 1px),
      linear-gradient(rgba(0, 255, 200, 0.03) 1px, transparent 1px);
  }
}

@keyframes meshFloat {
  0% { background-position: 0 0, 0 0; }
  100% { background-position: 60px 60px, 60px 60px; }
}

.connection-canvas {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 1;
}

.particles-container {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;
  z-index: 2;
}

.particle {
  position: absolute;
  bottom: -20px;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.5), rgba(255, 255, 255, 0.2));
  border-radius: 50%;
  animation: floatUp linear infinite;
  box-shadow: 0 0 12px rgba(255, 255, 255, 0.4);

  :global([data-theme='dark']) & {
    background: linear-gradient(135deg, rgba(0, 255, 200, 0.6), rgba(0, 200, 150, 0.3));
    box-shadow: 0 0 16px rgba(0, 255, 200, 0.5);
  }
}

@keyframes floatUp {
  0% { transform: translateY(0) scale(1) rotate(0deg); opacity: 0; }
  5% { opacity: 0.8; }
  95% { opacity: 0.4; }
  100% { transform: translateY(-120vh) scale(0.3) rotate(360deg); opacity: 0; }
}

.glow-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  animation: orbFloat 25s ease-in-out infinite;
  z-index: 0;
}

.orb-1 {
  width: 400px;
  height: 400px;
  top: 15%;
  left: 15%;
  background: radial-gradient(circle, rgba(255, 255, 255, 0.2) 0%, transparent 70%);
  animation-delay: 0s;

  :global([data-theme='dark']) & {
    background: radial-gradient(circle, rgba(0, 255, 200, 0.15) 0%, transparent 70%);
  }
}

.orb-2 {
  width: 280px;
  height: 280px;
  top: 55%;
  right: 20%;
  background: radial-gradient(circle, rgba(255, 255, 255, 0.15) 0%, transparent 70%);
  animation-delay: -8s;

  :global([data-theme='dark']) & {
    background: radial-gradient(circle, rgba(0, 200, 255, 0.12) 0%, transparent 70%);
  }
}

.orb-3 {
  width: 200px;
  height: 200px;
  bottom: 25%;
  left: 35%;
  background: radial-gradient(circle, rgba(255, 255, 255, 0.12) 0%, transparent 70%);
  animation-delay: -16s;

  :global([data-theme='dark']) & {
    background: radial-gradient(circle, rgba(100, 255, 200, 0.1) 0%, transparent 70%);
  }
}

@keyframes orbFloat {
  0%, 100% { transform: translate(0, 0) scale(1); }
  33% { transform: translate(40px, -30px) scale(1.08); }
  66% { transform: translate(-20px, 25px) scale(0.95); }
}

.wave {
  position: absolute;
  bottom: 0;
  left: 0;
  width: 200%;
  height: 200px;
  background-repeat: repeat-x;
  background-size: 50% 100%;
  opacity: 0.35;

  :global([data-theme='dark']) & { opacity: 0.45; }
}

.wave-1 {
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 1440 320'%3E%3Cpath fill='%23ffffff' fill-opacity='0.15' d='M0,160L48,176C96,192,192,224,288,213.3C384,203,480,149,576,138.7C672,128,768,160,864,181.3C960,203,1056,213,1152,192C1248,171,1344,117,1392,90.7L1440,64L1440,320L1392,320C1344,320,1248,320,1152,320C1056,320,960,320,864,320C768,320,672,320,576,320C480,320,384,320,288,320C192,320,96,320,48,320L0,320Z'%3E%3C/path%3E%3C/svg%3E");
  animation: waveMove 18s linear infinite;
  height: 240px;

  :global([data-theme='dark']) & {
    background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 1440 320'%3E%3Cpath fill='%2300ffc8' fill-opacity='0.12' d='M0,160L48,176C96,192,192,224,288,213.3C384,203,480,149,576,138.7C672,128,768,160,864,181.3C960,203,1056,213,1152,192C1248,171,1344,117,1392,90.7L1440,64L1440,320L1392,320C1344,320,1248,320,1152,320C1056,320,960,320,864,320C768,320,672,320,576,320C480,320,384,320,288,320C192,320,96,320,48,320L0,320Z'%3E%3C/path%3E%3C/svg%3E");
  }
}

.wave-2 {
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 1440 320'%3E%3Cpath fill='%23ffffff' fill-opacity='0.08' d='M0,224L48,213.3C96,203,192,181,288,181.3C384,181,480,203,576,224C672,245,768,267,864,250.7C960,235,1056,181,1152,165.3C1248,149,1344,171,1392,181.3L1440,192L1440,320L1392,320C1344,320,1248,320,1152,320C1056,320,960,320,864,320C768,320,672,320,576,320C480,320,384,320,288,320C192,320,96,320,48,320L0,320Z'%3E%3C/path%3E%3C/svg%3E");
  animation: waveMove 14s linear infinite reverse;
  height: 200px;
  bottom: 30px;
  opacity: 0.2;

  :global([data-theme='dark']) & {
    background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 1440 320'%3E%3Cpath fill='%2300c8ff' fill-opacity='0.08' d='M0,224L48,213.3C96,203,192,181,288,181.3C384,181,480,203,576,224C672,245,768,267,864,250.7C960,235,1056,181,1152,165.3C1248,149,1344,171,1392,181.3L1440,192L1440,320L1392,320C1344,320,1248,320,1152,320C1056,320,960,320,864,320C768,320,672,320,576,320C480,320,384,320,288,320C192,320,96,320,48,320L0,320Z'%3E%3C/path%3E%3C/svg%3E");
    opacity: 0.3;
  }
}

@keyframes waveMove {
  0% { transform: translateX(0); }
  100% { transform: translateX(-50%); }
}

.grid-lines {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-image: 
    linear-gradient(rgba(255, 255, 255, 0.025) 1px, transparent 1px),
    linear-gradient(90deg, rgba(255, 255, 255, 0.025) 1px, transparent 1px);
  background-size: 100px 100px;

  :global([data-theme='dark']) & {
    background-image: 
      linear-gradient(rgba(0, 200, 150, 0.02) 1px, transparent 1px),
      linear-gradient(90deg, rgba(0, 200, 150, 0.02) 1px, transparent 1px);
  }
}

/* 品牌头部 */
.brand-header {
  position: absolute;
  top: 40px;
  left: 48px;
  display: flex;
  align-items: center;
  gap: 12px;
  z-index: 10;
}

.brand-name {
  font-family: 'Swis721 BlkCn BT', Arial, sans-serif;
  font-size: 22px;
  font-weight: 600;
  color: #ffffff;
  letter-spacing: 0.5px;
}

/* Hero 内容区 */
.hero-content {
  position: absolute;
  bottom: 80px;
  left: 48px;
  right: 48px;
  z-index: 10;
}

.hero-title {
  font-size: 42px;
  font-weight: 700;
  color: #ffffff;
  line-height: 1.2;
  margin: 0 0 16px 0;
  letter-spacing: -0.5px;
}

.hero-subtitle {
  font-size: 16px;
  color: rgba(255, 255, 255, 0.7);
  line-height: 1.6;
  margin: 0 0 40px 0;
  white-space: nowrap;
}

.hero-stats {
  display: flex;
  align-items: center;
  gap: 24px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 2px;
}

.stat-label {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.7);
  font-weight: 400;
  margin-bottom: 2px;
}

.stat-number {
  font-size: 28px;
  font-weight: 600;
  color: #ffffff;
  line-height: 1.2;
  margin-bottom: 2px;
}

.stat-sublabel {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.75);
  font-weight: 400;
}

.stat-divider {
  width: 1px;
  height: 56px;
  background: rgba(255, 255, 255, 0.15);
}

/* 右侧登录区域 - 38.2% */
.login-section {
  position: relative;
  width: $phi-inv;
  height: 100%;
  background: #fafafa;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;

  :global([data-theme='dark']) & {
    background: #0a0f1a;
  }

  @media (max-width: 1024px) {
    width: 100%;
  }
}

.language-switcher {
  position: absolute;
  top: 32px;
  right: 40px;
  z-index: 100;
}

.language-switcher :deep(.el-select) {
  width: 100px;
}

.language-switcher :deep(.el-input__wrapper) {
  background: transparent;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  box-shadow: none;
  transition: all 0.3s ease;

  :global([data-theme='dark']) & {
    border-color: #2d3748;
  }

  &:hover {
    border-color: #019c7c;
  }
}

.language-switcher :deep(.el-input__inner) {
  font-size: 13px;
  color: #374151;

  :global([data-theme='dark']) & {
    color: #e5e7eb;
  }
}

.login-container {
  width: 100%;
  max-width: 340px;
  padding: 0 40px;
}

.login-card {
  width: 100%;
  animation: cardFadeIn 0.6s ease-out;
}

@keyframes cardFadeIn {
  0% { opacity: 0; transform: translateY(20px); }
  100% { opacity: 1; transform: translateY(0); }
}

.login-header {
  margin-bottom: 32px;
}

.title {
  margin: 0 0 8px 0;
  font-weight: 600;
  font-size: 26px;
  line-height: 1.3;
  color: #111827;

  :global([data-theme='dark']) & {
    color: #f9fafb;
  }
}

.subtitle {
  margin: 0;
  font-size: 14px;
  color: #6b7280;

  :global([data-theme='dark']) & {
    color: #9ca3af;
  }
}

.login-form {
  width: 100%;
}

.login-input :deep(.el-input__wrapper) {
  height: 46px;
  border-radius: 10px;
  background-color: #ffffff;
  border: 1px solid #e5e7eb;
  box-shadow: none;
  transition: all 0.25s ease;
  padding: 0 14px;

  :global([data-theme='dark']) & {
    background-color: #111827;
    border-color: #374151;
  }

  &:hover {
    border-color: #d1d5db;

    :global([data-theme='dark']) & {
      border-color: #4b5563;
    }
  }

  &.is-focus {
    border-color: #019c7c !important;
    box-shadow: 0 0 0 4px rgba(1, 156, 124, 0.08) !important;

    :global([data-theme='dark']) & {
      border-color: #00c896 !important;
      box-shadow: 0 0 0 4px rgba(0, 200, 150, 0.1) !important;
    }
  }
}

.login-input :deep(.el-input__inner) {
  height: 100%;
  font-size: 14px;
  color: #111827;
  background-color: transparent;
  border: none;
  box-shadow: none;

  &::placeholder {
    color: #9ca3af;
  }

  :global([data-theme='dark']) & {
    color: #f9fafb;

    &::placeholder {
      color: #6b7280;
    }
  }
}

.login-input :deep(.el-input__prefix) {
  color: #9ca3af;
  margin-right: 8px;

  :global([data-theme='dark']) & {
    color: #6b7280;
  }
}

.login-button {
  width: 100%;
  height: 46px;
  font-size: 14px;
  font-weight: 500;
  border-radius: 10px;
  background: #111827;
  border: none;
  margin-top: 8px;
  transition: all 0.25s ease;

  :global([data-theme='dark']) & {
    background: #00c896;
    color: #000000;
  }

  &:hover {
    background: #374151;
    transform: translateY(-1px);

    :global([data-theme='dark']) & {
      background: color.adjust(#00c896, $lightness: -5%);
    }
  }

  &:active {
    transform: translateY(0);
  }
}

.form-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin: 16px 0 24px 0;
  font-size: 13px;
}

.forgot-link {
  font-size: 13px;
  color: #6b7280;
  transition: color 0.25s ease;

  :global([data-theme='dark']) & {
    color: #9ca3af;
  }

  &:hover {
    color: #019c7c;

    :global([data-theme='dark']) & {
      color: #00c896;
    }
  }
}

.login-form :deep(.el-checkbox__label) {
  font-size: 13px;
  color: #374151;
  padding-left: 6px;

  :global([data-theme='dark']) & {
    color: #d1d5db;
  }
}

.login-form :deep(.el-checkbox__inner) {
  border-radius: 5px;
  border-color: #d1d5db;
  width: 16px;
  height: 16px;

  :global([data-theme='dark']) & {
    border-color: #4b5563;
    background-color: transparent;
  }
}

.login-form :deep(.el-checkbox__input.is-checked .el-checkbox__inner) {
  background-color: #019c7c;
  border-color: #019c7c;

  :global([data-theme='dark']) & {
    background-color: #00c896;
    border-color: #00c896;
  }
}

.login-footer {
  margin-top: 32px;
  padding-top: 24px;
  border-top: 1px solid #e5e7eb;
  text-align: center;

  :global([data-theme='dark']) & {
    border-color: #374151;
  }

  p {
    font-size: 12px;
    color: #9ca3af;
    margin: 0;

    :global([data-theme='dark']) & {
      color: #6b7280;
    }
  }
}

.language-switcher :deep(.el-select-dropdown) {
  background-color: #ffffff !important;
  border: 1px solid #e5e7eb !important;
  border-radius: 8px !important;
  box-shadow: 0 10px 40px -10px rgba(0, 0, 0, 0.15) !important;

  :global([data-theme='dark']) & {
    background: #1f2937 !important;
    border-color: #374151 !important;
    box-shadow: 0 10px 40px -10px rgba(0, 0, 0, 0.4) !important;
  }
}

.language-switcher :deep(.el-select-dropdown__item) {
  color: #374151 !important;
  font-size: 13px;
  padding: 10px 14px;

  :global([data-theme='dark']) & {
    color: #e5e7eb !important;
  }

  &:hover {
    background-color: #f3f4f6 !important;

    :global([data-theme='dark']) & {
      background-color: #374151 !important;
    }
  }

  &.selected {
    color: #019c7c !important;
    font-weight: 500;

    :global([data-theme='dark']) & {
      color: #00c896 !important;
    }
  }
}
</style>
