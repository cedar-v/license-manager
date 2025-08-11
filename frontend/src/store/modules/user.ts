import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

interface UserInfo {
  id: string
  username: string
  name: string
  email: string
  avatar?: string
  role: 'admin' | 'user' | 'guest'
  permissions: string[]
}

export const useUserStore = defineStore('user', () => {
  // 状态
  const token = ref<string | null>(localStorage.getItem('token'))
  const userInfo = ref<UserInfo | null>(null)
  const isLoggedIn = ref(false)

  // 计算属性
  const isAdmin = computed(() => userInfo.value?.role === 'admin')
  const hasPermission = computed(() => (permission: string) => {
    return userInfo.value?.permissions?.includes(permission) || false
  })

  // 操作
  const setToken = (newToken: string | null) => {
    token.value = newToken
    if (newToken) {
      localStorage.setItem('token', newToken)
    } else {
      localStorage.removeItem('token')
    }
  }

  const setUserInfo = (info: UserInfo | null) => {
    userInfo.value = info
    isLoggedIn.value = !!info
  }

  const login = async (credentials: { username: string; password: string }) => {
    // 这里应该调用登录API
    // const response = await loginApi(credentials)
    // setToken(response.data.token)
    // setUserInfo(response.data.user)
    
    // 模拟登录
    setToken('mock-token-123')
    setUserInfo({
      id: '1',
      username: credentials.username,
      name: '管理员',
      email: 'admin@example.com',
      role: 'admin',
      permissions: ['*']
    })
  }

  const logout = () => {
    setToken(null)
    setUserInfo(null)
  }

  const updateProfile = (profile: Partial<UserInfo>) => {
    if (userInfo.value) {
      userInfo.value = { ...userInfo.value, ...profile }
    }
  }

  // 初始化时检查登录状态
  const initAuth = () => {
    if (token.value) {
      // 这里应该调用API验证token有效性
      // 暂时模拟已登录状态
      isLoggedIn.value = true
    }
  }

  return {
    // 状态
    token,
    userInfo,
    isLoggedIn,
    // 计算属性
    isAdmin,
    hasPermission,
    // 操作
    setToken,
    setUserInfo,
    login,
    logout,
    updateProfile,
    initAuth
  }
})