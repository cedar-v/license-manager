import axios from 'axios'

// http状态码
import { errorCodeType } from './errorCodeType'
// 语言管理器
import { generateAcceptLanguageHeader } from '@/utils/language'
// 路由
import router from '@/router'

const envUrl = import.meta.env.VITE_API_BASE_URL

//使用create方法创建axios实例
const Axios = axios.create({
  timeout: 10000, // 请求超时时间
  baseURL: envUrl,
  headers: {
    'Content-Type': 'application/json;charset=UTF-8',
  },
  // 跨域时候允许携带凭证 (开发环境不启用，避免CORS问题)
  withCredentials: false
})
// 添加请求拦截器
Axios.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  // 自定义请求头携带token (Bearer认证)
  if(token){ 
    config.headers['Authorization'] = `Bearer ${token}`
  }
  
  // 设置 Accept-Language 头，符合 RFC 7231 规范
  // 优先级：用户选择;q=1.0 > 备用语言;q=0.8 > 通配符;q=0.5
  config.headers['Accept-Language'] = generateAcceptLanguageHeader()
  
  return config
})
// 添加响应拦截器
Axios.interceptors.response.use(response => {
  // console.log(response)
  const apiData = response.data as any
  // 这个 Code 是和后端约定的业务 Code
  const code = apiData.code
  // 如果没有 Code, 代表这不是项目后端开发的 API
  if (code === undefined) {
    // 不在拦截器中显示错误，让组件自己处理
    return Promise.reject(new Error("非本系统的接口"))
  } else {
    switch (code) {
      case 200:
        // code === 200 代表成功
        return apiData
      case 0:
        // code === 0 也代表成功 (兼容)
        return apiData
      case 401:
        // code === 401 代表未授权，清除token并跳转到登录页
        localStorage.removeItem('token')
        router.push('/login')
        const authError = new Error(apiData.message || "未授权访问")
        ;(authError as any).response = { data: apiData }
        return Promise.reject(authError)
      default:
        // 不在拦截器中显示错误，让组件自己处理
        // 将错误信息附加到 Error 对象上，供组件使用
        const error = new Error(apiData.message || "Error")
        ;(error as any).response = { data: apiData }
        return Promise.reject(error)
    }
  }
}, error => {
  console.log('Response的error', error)
  // 不在拦截器中显示错误消息，让具体的页面组件处理
  // showErrMessage(errorCodeType(error.response?.status))
  
  // 为错误对象添加更多信息，便于组件处理
  if (error.response) {
    const status = error.response.status
    const data = error.response.data
    
    // 处理HTTP 401状态码
    if (status === 401) {
      localStorage.removeItem('token')
      router.push('/login')
      ;(error as any).backendMessage = '登录已过期，请重新登录'
      return Promise.reject(error)
    }
    
    // 如果后端返回了错误消息，优先使用后端消息
    if (data && data.message) {
      ;(error as any).backendMessage = data.message
    } else {
      // 否则使用状态码对应的消息
      ;(error as any).backendMessage = errorCodeType(status)
    }
  }
  
  return Promise.reject(error)
})

// 错误处理函数已移除，由页面组件统一处理错误显示

export default Axios;