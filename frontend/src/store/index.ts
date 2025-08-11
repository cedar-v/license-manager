import { createPinia } from 'pinia'

const pinia = createPinia()

export default pinia

// 导出所有store
export * from './modules/app'
export * from './modules/user'