import axios from 'axios'

// http状态码
import { errorCodeType } from './errorCodeType'
//element-plus提示消息组件
import { ElMessage } from "element-plus"
const envUrl = import.meta.env.VITE_API_BASE_URL

//使用create方法创建axios实例
const Axios = axios.create({
  timeout: 10000, // 请求超时时间
  baseURL: envUrl,
  headers: {
    'Content-Type': 'application/json;charset=UTF-8',
  },
  // 跨域时候允许携带凭证
  withCredentials: true
})
// 添加请求拦截器
Axios.interceptors.request.use(config => {
  const token = localStorage.getItem('blockchain_special_token')
  // 自定义请求头携带token
  if(token){ config.headers['token'] = token}
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
    ElMessage.error("非本系统的接口")
    return Promise.reject(new Error("非本系统的接口"))
  } else {
    switch (code) {
      case 0:
        // code === 0 代表没有错误
        return apiData
      default:
        // 不是正确的 Code
        ElMessage.error(apiData.message || "Error")
        return Promise.reject(new Error("Error"))
    }
  }
}, error => {
  console.log('Response的error', error)
  showErrMessage(errorCodeType(error.response.status))
  return Promise.reject(error)
})

/**
 * @description 显示错误消息
 * opt 传入参数
 * err 错误信息
 * type 消息类型
 * duration 消息持续时间
 */
function showErrMessage (errMessage: string, type:any= 'error', duration:number = 5000){
  ElMessage({
      message: errMessage,
      type:type,
      duration: duration
  })
}

export default Axios;