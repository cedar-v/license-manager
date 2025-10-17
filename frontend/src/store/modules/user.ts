/*
 * @Author: 13895237362 2205451508@qq.com
 * @Date: 2025-08-26 13:23:45
 * @LastEditors: 13895237362 2205451508@qq.com
 * @LastEditTime: 2025-10-17 10:35:22
 * @FilePath: \license-manager\frontend\src\store\modules\user.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

interface UserInfo {
  id?: string
  username: string
  email?: string
  avatar?: string
  role: string
  permissions?: string[]
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
    
    // 同时保存到localStorage
    if (info) {
      localStorage.setItem('userInfo', JSON.stringify(info))
    } else {
      localStorage.removeItem('userInfo')
    }
  }

  const setLoginData = (token: string, userInfo: UserInfo) => {
    setToken(token)
    setUserInfo(userInfo)
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
      // 从localStorage恢复用户信息
      const savedUserInfo = localStorage.getItem('userInfo')
      if (savedUserInfo) {
        try {
          const userInfoData = JSON.parse(savedUserInfo)
          setUserInfo(userInfoData)
        } catch (error) {
          console.error('恢复用户信息失败:', error)
          // 如果解析失败，清除无效数据
          localStorage.removeItem('userInfo')
          setToken(null)
        }
      } else {
        // 如果没有保存的用户信息，但有token，设置基本登录状态
        isLoggedIn.value = true
      }
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
    setLoginData,
    logout,
    updateProfile,
    initAuth
  }
})