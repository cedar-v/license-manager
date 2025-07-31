import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

// 全局样式重置
import './assets/styles/global.css'

import i18n from './i18n' // 国际化

const app = createApp(App)
app.use(router)
app.use(i18n)
app.mount('#app')